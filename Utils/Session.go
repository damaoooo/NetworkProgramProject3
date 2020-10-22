package Utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type Session struct {
	Session  string `json:"session"`
	LastTime int64  `json:"last_time"`
	Username string `json:"username"`
}

//Session Database With a Mutex lock
type SessionManager struct {
	Sessions map[string]Session //{"session key": session struct}
	Lock     sync.Mutex
	Err      SessionError
}

type SessionError struct {
	errorDescription string //error description
}

func (s *SessionError) Error() error {
	return errors.New(s.errorDescription)
}

func (s *SessionManager) GenerateNew(username string) string {
	UUID := uuid.Must(uuid.NewV4()).String()
	data := []byte(username + UUID)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash[:16])
	t := time.Now().Unix()
	ret := Session{
		Session:  md5str,
		LastTime: t,
		Username: username,
	}
	s.Lock.Lock()
	s.Sessions[md5str] = ret
	s.Lock.Unlock()
	return ret.Session
}

func (s *SessionManager) IsValid(session string) bool {
	s.Lock.Lock()
	if _, ok := s.Sessions[session]; ok {
		s.Lock.Unlock()
		return true
	} else {
		s.Lock.Unlock()
		return false
	}
}

func (s *SessionManager) Update(session string) error {
	s.Lock.Lock()
	if s.IsValid(session) {
		ns := s.Sessions[session]
		ns.LastTime = time.Now().Unix()
		s.Sessions[session] = ns
		s.Lock.Unlock()
		return nil
	} else {
		s.Err.errorDescription = "no such session, update failed"
		s.Lock.Unlock()
		return s.Err.Error()
	}
}

func (s *SessionManager) Destroy(session string) error {
	s.Lock.Lock()
	if s.IsValid(session) {
		delete(s.Sessions, session)
		s.Lock.Unlock()
		return nil
	} else {
		s.Err.errorDescription = "no such session, destroy failed"
		s.Lock.Unlock()
		return s.Err.Error()
	}
}
