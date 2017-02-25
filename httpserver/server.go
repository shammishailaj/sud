package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/crazyprograms/sud/core"
)

type Session struct {
	sessionID     string
	Param         interface{}
	ticCount      int
	remoteAddress []string
	client        core.IClient
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
	s := &Session{sessionID: SessionID}
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

type jsonLogin struct {
	Login             string `json:"login"`
	Password          string `json:"password"`
	ConfigurationName string `json:"configurationName"`
}

func (server *Server) httpJsonSend(w http.ResponseWriter, result interface{}) error {
	var b []byte
	var err error
	if b, err = json.Marshal(result); err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}
func (server *Server) httpJsonLogin(w http.ResponseWriter, r *http.Request, session *Session) error {
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}

	var b bytes.Buffer
	var result struct {
		Login bool   `json:"result,omitempty"`
		Error string `json:"error,omitempty"`
	}
	if _, err = b.ReadFrom(r.Body); err == nil {
		var j struct {
			Login             string `json:"login"`
			Password          string `json:"password"`
			ConfigurationName string `json:"configurationName"`
		}
		if err = json.Unmarshal(b.Bytes(), &j); err == nil {
			fmt.Println("Login: ", (string)(b.Bytes()))
			if j.Login == "" || j.Password == "" || j.ConfigurationName == "" {
				err = errors.New("param error")
			} else if session.client == nil {
				session.client = server.core.NewClient(j.Login, j.Password, j.ConfigurationName)
				if session.client == nil {
					err = errors.New("Error login")
				} else {
					result.Login = true
				}
			} else {
				err = errors.New("Already login")
			}
		}
	}
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}
func (server *Server) httpJsonBeginTransaction(w http.ResponseWriter, r *http.Request, session *Session) error {
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}
	if session.client == nil {
		err = errors.New("not login")
		http.Error(w, err.Error(), 403)
		return err
	}
	var result struct {
		TransactionUID string `json:"transactionUID,omitempty"`
		Error          string `json:"error,omitempty"`
	}
	var TransactionUID string
	TransactionUID, err = session.client.BeginTransaction()
	result.TransactionUID = TransactionUID
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}
func (server *Server) httpJsonCommitTransaction(w http.ResponseWriter, r *http.Request, session *Session) error {
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}
	if session.client == nil {
		err = errors.New("not login")
		http.Error(w, err.Error(), 403)
		return err
	}
	var b bytes.Buffer
	var result struct {
		Commit bool   `json:"commit"`
		Error  string `json:"error,omitempty"`
	}
	if _, err = b.ReadFrom(r.Body); err == nil {
		var j struct {
			transactionUID string `json:"transactionUID"`
		}
		if err = json.Unmarshal(b.Bytes(), &j); err == nil {
			fmt.Println("CommitTransaction: ", (string)(b.Bytes()))
			if j.transactionUID == "" {
				err = errors.New("param error")
			} else {
				err = session.client.CommitTransaction(j.transactionUID)
				result.Commit = true
			}
		}
	}
	if err != nil {
		result.Error = err.Error()
	}
	return server.httpJsonSend(w, result)
}
func (server *Server) httpJsonRollbackTransaction(w http.ResponseWriter, r *http.Request, session *Session) error {
	var err error
	if session == nil {
		err = errors.New("session not found")
		http.Error(w, err.Error(), 403)
		return err
	}
	if session.client == nil {
		err = errors.New("not login")
		http.Error(w, err.Error(), 403)
		return err
	}
	var b bytes.Buffer
	var result struct {
		Rollback bool   `json:"rollback"`
		Error    string `json:"error,omitempty"`
	}
	if _, err = b.ReadFrom(r.Body); err == nil {
		var j struct {
			transactionUID string `json:"transactionUID"`
		}
		if err = json.Unmarshal(b.Bytes(), &j); err == nil {
			fmt.Println("RollbackTransaction: ", (string)(b.Bytes()))
			if j.transactionUID == "" {
				err = errors.New("param error")
			} else {
				err = session.client.RollbackTransaction(j.transactionUID)
				result.Rollback = true
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
	s.handlers["/json/begintransaction"] = s.httpJsonBeginTransaction
	s.handlers["/json/committransaction"] = s.httpJsonCommitTransaction
	s.handlers["/json/rollbacktransaction"] = s.httpJsonRollbackTransaction

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
