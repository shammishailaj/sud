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
	fmt.Println("in:", RequestURI[0])
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
	fmt.Println("out:", RequestURI[0])
	//	fmt.Fprintln(w, r.Host, r.Method, r.RequestURI)

	/*//var err error


	if session != nil {
		err = errors.New("session end")
	}*/

}
func (server *Server) stopSession(SessionID string) {
	server.sessionsLock.Lock()
	delete(server.sessions, SessionID)
	server.sessionsLock.Unlock()
}
func (server *Server) newSession(w http.ResponseWriter, r *http.Request) *Session {
	SessionID := core.NewUUID().String()
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

func httpJson(InParamType reflect.Type, Handler JsonHandler, HandlerError JsonHandlerError) fHandler {
	return func(w http.ResponseWriter, request *http.Request, session *Session) error {
		var err error
		var inBuff bytes.Buffer
		var outBuff []byte
		if _, err = inBuff.ReadFrom(request.Body); err == nil {
			var RecvParam interface{}
			RecvParam = reflect.New(InParamType).Interface()
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
func (server *Server) httpJsonSend(w http.ResponseWriter, result interface{}) error {
	var b []byte
	var err error
	if b, err = json.Marshal(result); err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(string(b))
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}
func (server *Server) httpJsonLogin(w http.ResponseWriter, r *http.Request, session *Session) error {
	session.lock.Lock()
	defer session.lock.Unlock()
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}

	var b bytes.Buffer
	var result jsonLoginResult
	if _, err = b.ReadFrom(r.Body); err == nil {
		var j jsonLogin
		if err = json.Unmarshal(b.Bytes(), &j); err == nil {
			if j.Login == "" || j.Password == "" || j.ConfigurationName == "" {
				err = errors.New("param error")
			} else if session.client == nil {
				if session.client, err = server.core.NewClient(j.Login, j.Password, j.ConfigurationName); err == nil {
					session.auth = true
					result.Login = true
				} else {
					session.auth = false
				}
			} else {
				result.Login = true
				err = errors.New("Already login")
			}
		}
	}
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}
func (session *Session) setResult(result chan callpull.Result) string {
	session.lockResult.Lock()
	defer session.lockResult.Unlock()
	uid := core.NewUUID().String()
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

func (server *Server) httpJsonBeginTransaction(w http.ResponseWriter, r *http.Request, session *Session) error {
	session.lock.RLock()
	defer session.lock.RUnlock()
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}
	if !session.auth {
		err = errors.New("not login")
		http.Error(w, err.Error(), 403)
		return err
	}
	var result jsonBeginTransactionResult
	var TransactionUID string
	TransactionUID, err = session.client.BeginTransaction()
	result.TransactionUID = TransactionUID
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}

func (server *Server) httpJsonCall(w http.ResponseWriter, r *http.Request, session *Session) error {
	session.lock.RLock()
	defer session.lock.RUnlock()
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}
	if !session.auth {
		err = errors.New("not login")
		http.Error(w, err.Error(), 403)
		return err
	}
	var b bytes.Buffer
	var result jsonCallResult
	if _, err = b.ReadFrom(r.Body); err == nil {
		var j jsonCall
		if err = json.Unmarshal(b.Bytes(), &j); err == nil {
			if Params, err := jsonUnPackMap(*j.Params); err == nil {
				var r1 callpull.Result
				if r1, err = session.client.Call(j.Name, Params, time.Duration(j.TimeoutWait)*time.Millisecond); err == nil {
					if result.Result, err = jsonPack(r1.Result); err == nil {
						result.Error = r1.Error.Error()
					}
				}
			}
		}
	}
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}
func (server *Server) Start() {
	http.ListenAndServe(server.address, server)
}
func (server *Server) jsonError(w http.ResponseWriter, r *http.Request, err error) interface{} {
	return &jsonErrorReturn{Error: err.Error()}
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
func (server *Server) jsonCommit(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
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
func (server *Server) jsonRollback(w http.ResponseWriter, r *http.Request, In interface{}, session *Session) (interface{}, error) {
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
func NewServer(c *core.Core, Address string) *Server {
	s := &Server{
		core:              c,
		address:           Address,
		sessions:          make(map[string]*Session),
		handlers:          make(map[string]fHandler),
		sessionCookieName: "SessionID",
		ticTimeOut:        time.Millisecond * 1000 * 30,
		ticCount:          100}
	s.handlers["/json/login"] = s.httpJsonLogin
	//s.handlers["/json/listen"] = s.httpJsonListen
	s.handlers["/json/listen"] = httpJson(reflect.TypeOf(jsonListen{}), s.jsonListen, s.jsonError)
	s.handlers["/json/listenreturn"] = httpJson(reflect.TypeOf(jsonListenReturn{}), s.jsonListenResult, s.jsonError)
	//s.handlers["/json/listenreturn"] = s.httpJsonListenReturn
	s.handlers["/json/begin"] = s.httpJsonBeginTransaction
	s.handlers["/json/commit"] = httpJson(reflect.TypeOf(jsonCommitTransaction{}), s.jsonCommit, s.jsonError)
	s.handlers["/json/rollback"] = httpJson(reflect.TypeOf(jsonRollbackTransaction{}), s.jsonRollback, s.jsonError)

	//s.handlers["/json/commit"] = s.httpJsonCommitTransaction
	//s.handlers["/json/rollback"] = s.httpJsonRollbackTransaction
	s.handlers["/json/call"] = s.httpJsonCall

	staticFile := http.FileServer(http.Dir("./static/"))
	s.handlers["/"] = func(w http.ResponseWriter, request *http.Request, session *Session) error {
		/*if request.URL.Path == "/" || request.URL.Path == "" {
			request.URL.Path = "/index.html"
		}*/
		staticFile.ServeHTTP(w, request)
		return nil
	}
	return s
}
