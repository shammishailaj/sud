package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/crazyprograms/callpull"
	"github.com/crazyprograms/sud/client"
	"github.com/crazyprograms/sud/core"
	"github.com/crazyprograms/sud/corebase"
)

type Session struct {
	lock          sync.RWMutex
	sessionID     string
	Param         interface{}
	ticCount      int
	remoteAddress []string
	client        client.IClient
	auth          bool
	lockResult    sync.Mutex
	result        map[string]chan callpull.Result
}

func (s *Session) setAddress(address string) {
	i := strings.LastIndex(address, ":")
	s.remoteAddress = []string{address[0:i], address[i+1:]}
}
func (s *Session) CheckAddress(address string) bool {
	i := strings.LastIndex(address, ":")
	return address[0:i] == s.remoteAddress[0]
}
func (s *Session) goTic(server *Server) {
	for {
		time.Sleep(server.ticTimeOut)
		s.ticCount++
		if s.ticCount >= server.ticCount {
			server.stopSession(s.sessionID)
			return
		}
	}
}

type fHandler func(w http.ResponseWriter, request *http.Request, session *Session) error
type Server struct {
	address           string
	core              *core.Core
	sessionCookieName string
	ticTimeOut        time.Duration
	ticCount          int
	sessionsLock      sync.RWMutex
	sessions          map[string]*Session
	handlers          map[string]fHandler
}

func (server *Server) getHandler(URL string) (fHandler, string) {
	var ok bool
	var h fHandler = nil
	for {
		if h, ok = server.handlers[URL+"/"]; ok {
			return h, URL + "/"
		}
		if h, ok = server.handlers[URL]; ok {
			return h, URL
		}
		i := strings.LastIndex(URL, "/")
		if URL == "" || URL == "/" || i == -1 {
			return nil, ""
		}
		URL = URL[0:i]
	}
}
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RequestURI := strings.SplitN(r.RequestURI, "?", 2)
	if h, URL := server.getHandler(RequestURI[0]); h != nil {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, URL)
		session := server.getSession(w, r)
		err := h(w, r, session)
		if err != nil {
			server.errorHandler(w, r, err)
		}
	} else {
		server.errorHandler(w, r, errors.New("http404"))
	}
}
func (server *Server) stopSession(SessionID string) {
	server.sessionsLock.Lock()
	delete(server.sessions, SessionID)
	server.sessionsLock.Unlock()
}
func (server *Server) newSession(w http.ResponseWriter, r *http.Request) *Session {
	SessionID := corebase.NewUUID().String()
	c := &http.Cookie{}
	c.Name = server.sessionCookieName
	c.Value = SessionID
	c.Path = "/"
	http.SetCookie(w, c)
	s := &Session{sessionID: SessionID, result: make(map[string]chan callpull.Result)}
	s.setAddress(r.RemoteAddr)
	server.sessionsLock.Lock()
	server.sessions[SessionID] = s
	server.sessionsLock.Unlock()
	go s.goTic(server)
	return s
}
func (server *Server) getSession(w http.ResponseWriter, r *http.Request) *Session {
	var ok bool
	var session *Session
	if c, err := r.Cookie(server.sessionCookieName); err == nil {
		SessionID := c.Value
		server.sessionsLock.RLock()
		session, ok = server.sessions[SessionID]
		server.sessionsLock.RUnlock()
		if !ok {
			return nil
		}
		if !session.CheckAddress(r.RemoteAddr) {
			return nil
		}
		session.ticCount = 0
		return session
	}
	return server.newSession(w, r)
}
func (server *Server) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Write(([]byte)(err.Error()))
}
func (server *Server) httpTest(w http.ResponseWriter, r *http.Request, session *Session) error {
	w.Write(([]byte)("OK"))
	return nil
}

type JsonHandler func(w http.ResponseWriter, r *http.Request, Param interface{}, session *Session) (interface{}, error)
type JsonHandlerError func(w http.ResponseWriter, r *http.Request, err error) interface{}

func httpJsonServer(InParamType reflect.Type, Handler JsonHandler, HandlerError JsonHandlerError) fHandler {
	return func(w http.ResponseWriter, request *http.Request, session *Session) error {
		var err error
		var inBuff bytes.Buffer
		var outBuff []byte
		if _, err = inBuff.ReadFrom(request.Body); err == nil {
			var RecvParam interface{}
			if InParamType != nil {
				RecvParam = reflect.New(InParamType).Interface()
			} else {
				RecvParam = &struct{}{}
			}
			if err = json.Unmarshal(inBuff.Bytes(), RecvParam); err == nil {
				result, err := Handler(w, request, RecvParam, session)
				if err != nil {
					result = HandlerError(w, request, err)
				}
				if outBuff, err = json.Marshal(result); err != nil {
					http.Error(w, err.Error(), 500)
					return err
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err = w.Write(outBuff); err != nil {
					return err
				}
			}
		}
		return err
	}
}
func (session *Session) setResult(result chan callpull.Result) string {
	session.lockResult.Lock()
	defer session.lockResult.Unlock()
	uid := corebase.NewUUID().String()
	session.result[uid] = result
	return uid
}
func (session *Session) getResult(resultUID string) chan callpull.Result {
	session.lockResult.Lock()
	defer session.lockResult.Unlock()
	if r, ok := session.result[resultUID]; ok {
		delete(session.result, resultUID)
		return r
	}
	return nil
}
func (server *Server) Start() {
	http.ListenAndServe(server.address, server)
}
func (server *Server) jsonError(w http.ResponseWriter, r *http.Request, err error) interface{} {
	return &jsonErrorReturn{Error: err.Error()}
}
func (server *Server) jsonLogin(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	InParam := In.(*jsonLogin)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if InParam.Login == "" || InParam.Password == "" || InParam.ConfigurationName == "" {
		return nil, errors.New("param error")
	}
	if session.client != nil {
		return nil, errors.New("Already login")
	}
	if session.client, err = server.core.NewClient(InParam.Login, InParam.Password, InParam.ConfigurationName); err != nil {
		session.auth = false
		return nil, err
	}
	session.auth = true
	return jsonLoginResult{Login: true}, nil
}
func (server *Server) jsonGetConfiguration(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	return jsonGetConfigurationResult{Configuration: session.client.GetConfiguration()}, nil
}
func (server *Server) jsonCall(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	InParam := In.(*jsonCall)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	var Params map[string]interface{}
	if Params, err = jsonUnPackMap(*InParam.Params); err != nil {
		return nil, err
	}
	var r1 callpull.Result
	if r1, err = session.client.Call(InParam.Name, Params, time.Duration(InParam.TimeoutWait)*time.Millisecond); err != nil {
		return nil, err
	}
	var callPullError string
	if r1.Error != nil {
		callPullError = r1.Error.Error()
	}
	var Result *jsonParam
	if Result, err = jsonPack(r1.Result); err != nil {
		return nil, err
	}
	return jsonCallResult{Result: Result, CallPullError: callPullError}, nil
}
func (server *Server) jsonBeginTransaction(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	var TransactionUID string
	if TransactionUID, err = session.client.BeginTransaction(); err != nil {
		return nil, err
	}
	return jsonBeginTransactionResult{TransactionUID: TransactionUID}, nil
}
func (server *Server) jsonCommitTransaction(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	InParam := In.(*jsonCommitTransaction)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	if InParam.TransactionUID == "" {
		return nil, errors.New("param error")
	}
	return jsonCommitTransactionResult{}, session.client.CommitTransaction(InParam.TransactionUID)
}
func (server *Server) jsonRollbackTransaction(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	InParam := In.(*jsonRollbackTransaction)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	if InParam.TransactionUID == "" {
		return nil, errors.New("param error")
	}
	return jsonRollbackTransactionResult{}, session.client.RollbackTransaction(InParam.TransactionUID)
}
func (server *Server) jsonListen(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	InParam := In.(*jsonListen)

	session.lock.RLock()
	defer session.lock.RUnlock()
	var err error
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	if InParam.Name == "" || InParam.TimeoutWait < 0 {
		return nil, errors.New("param error")
	}
	var ResultChan chan callpull.Result
	var Params map[string]interface{}
	var PParams map[string]*jsonParam
	if Params, ResultChan, err = session.client.Listen(InParam.Name, time.Millisecond*time.Duration(InParam.TimeoutWait)); err != nil {
		return nil, err
	}
	if PParams, err = jsonPackMap(Params); err != nil {
		return nil, err
	}
	return &jsonListenResult{Param: &PParams, ResultUID: session.setResult(ResultChan)}, nil
}
func (server *Server) jsonListenResult(w http.ResponseWriter, request *http.Request, In interface{}, session *Session) (OutParam interface{}, resultErr error) {
	InParam := In.(*jsonListenReturn)
	OutParam = nil
	session.lock.RLock()
	defer session.lock.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			resultErr = errors.New(fmt.Sprintln(r))
		}
	}()
	var err error
	if session == nil {
		resultErr = errors.New("session not found")
		return
	}
	if !session.auth {
		resultErr = errors.New("not login")
		return
	}
	if InParam.ResultUID == "" {
		resultErr = errors.New("param error")
		return
	}
	ResultChan := session.getResult(InParam.ResultUID)
	if ResultChan == nil {
		resultErr = errors.New("ResultUID not found")
		return
	}
	var rPack interface{}
	if rPack, err = jsonUnPack(InParam.Result); err != nil {
		resultErr = err
		return
	}
	r := callpull.Result{Result: rPack}
	if InParam.Error != "" {
		r.Error = errors.New(InParam.Error)
	}
	ResultChan <- r
	return &jsonListenReturnResult{}, nil
}
func (server *Server) jsonSetDocumentPoles(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	InParam := In.(*jsonSetDocumentPoles)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	var Poles map[string]interface{}
	if Poles, err = jsonUnPackMap(*InParam.Poles); err != nil {
		return nil, err
	}
	if err = session.client.SetDocumentPoles(InParam.TransactionUID, InParam.DocumentUID, Poles); err != nil {
		return nil, err
	}
	return jsonSetDocumentPolesResult{}, nil
}
func (server *Server) jsonNewDocument(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	InParam := In.(*jsonNewDocument)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	var Poles map[string]interface{}
	if Poles, err = jsonUnPackMap(*InParam.Poles); err != nil {
		return nil, err
	}
	var DocumentUID string
	if DocumentUID, err = session.client.NewDocument(InParam.TransactionUID, InParam.DocumentType, Poles); err != nil {
		return nil, err
	}
	return jsonNewDocumentResult{DocumentUID: DocumentUID}, nil
}
func (server *Server) jsonGetDocumentPoles(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
	var err error
	InParam := In.(*jsonGetDocumentPoles)
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session == nil {
		return nil, errors.New("session not found")
	}
	if !session.auth {
		return nil, errors.New("not login")
	}
	var Wheres = make([]corebase.IDocumentWhere, len(InParam.Wheres), len(InParam.Wheres))
	for i, wp := range InParam.Wheres {
		var wparam map[string]interface{}
		if wparam, err = jsonUnPackMap(wp); err != nil {
			return nil, err
		}
		var WhereTypeI interface{}
		var WhereType string
		var ok bool
		if WhereTypeI, ok = wparam["whereType"]; !ok {
			return nil, errors.New("WhereType not found")
		}
		if WhereType, ok = WhereTypeI.(string); !ok {
			return nil, errors.New("WhereType not is string")
		}
		var Where corebase.IDocumentWhere
		if Where, err = corebase.NewDocumentWhere(WhereType); err != nil {
			return nil, err
		}
		params := map[string]interface{}{}
		for paramName, paramValue := range wparam {
			if len(paramName) > 0 {
				params["DocumentWhere."+WhereType+"."+strings.ToUpper(paramName[0:1])+paramName[1:]] = paramValue
			}
		}
		Where.Load(params)
		Wheres[i] = Where
	}
	var Documents map[string]map[string]interface{}
	if Documents, err = session.client.GetDocumentsPoles(InParam.TransactionUID, InParam.DocumentType, InParam.Poles, Wheres); err != nil {
		return nil, err
	}
	DocumentsPack := map[string]map[string]*jsonParam{}
	for DocumentUID, DocumentPoles := range Documents {
		if DocumentsPack[DocumentUID], err = jsonPackMap(DocumentPoles); err != nil {
			return nil, err
		}
	}
	return jsonGetDocumentPolesResult{Documents: &DocumentsPack}, nil
}

func NewServer(c *core.Core, Address string) *Server {
	s := &Server{
		core:              c,
		address:           Address,
		sessions:          make(map[string]*Session),
		handlers:          make(map[string]fHandler),
		sessionCookieName: "SessionID",
		ticTimeOut:        time.Millisecond * 1000 * 30,
		ticCount:          100}
	s.handlers["/json/login"] = httpJsonServer(reflect.TypeOf(jsonLogin{}), s.jsonLogin, s.jsonError)
	s.handlers["/json/getconfiguration"] = httpJsonServer(nil, s.jsonGetConfiguration, s.jsonError)
	s.handlers["/json/listen"] = httpJsonServer(reflect.TypeOf(jsonListen{}), s.jsonListen, s.jsonError)
	s.handlers["/json/listenreturn"] = httpJsonServer(reflect.TypeOf(jsonListenReturn{}), s.jsonListenResult, s.jsonError)
	s.handlers["/json/begintransaction"] = httpJsonServer(nil, s.jsonBeginTransaction, s.jsonError)
	s.handlers["/json/committransaction"] = httpJsonServer(reflect.TypeOf(jsonCommitTransaction{}), s.jsonCommitTransaction, s.jsonError)
	s.handlers["/json/rollbacktransaction"] = httpJsonServer(reflect.TypeOf(jsonRollbackTransaction{}), s.jsonRollbackTransaction, s.jsonError)
	s.handlers["/json/call"] = httpJsonServer(reflect.TypeOf(jsonCall{}), s.jsonCall, s.jsonError)
	s.handlers["/json/setdocumentpoles"] = httpJsonServer(reflect.TypeOf(jsonSetDocumentPoles{}), s.jsonSetDocumentPoles, s.jsonError)
	s.handlers["/json/getdocumentpoles"] = httpJsonServer(reflect.TypeOf(jsonGetDocumentPoles{}), s.jsonGetDocumentPoles, s.jsonError)
	s.handlers["/json/newdocument"] = httpJsonServer(reflect.TypeOf(jsonNewDocument{}), s.jsonNewDocument, s.jsonError)
	staticFile := http.FileServer(http.Dir("./static/"))
	s.handlers["/"] = func(w http.ResponseWriter, request *http.Request, session *Session) error {
		staticFile.ServeHTTP(w, request)
		return nil
	}
	return s
}
