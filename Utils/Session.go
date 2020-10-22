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
	errorDescription string
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
	defer s.Lock.Unlock()
	s.Sessions[md5str] = ret
	return ret.Session
}

// isValid func haven't lock to avoid duplicate lock operation,
//that means, you should add lock by your self.
func (s *SessionManager) isValid(session string) bool {
	if _, ok := s.Sessions[session]; ok {
		return true
	} else {
		return false
	}
}

func (s *SessionManager) Update(session string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if s.isValid(session) {
		ns := s.Sessions[session]
		ns.LastTime = time.Now().Unix()
		s.Sessions[session] = ns
		return nil
	} else {
		s.Err.errorDescription = "no such session, update failed"
		return s.Err.Error()
	}
}

func (s *SessionManager) Destroy(session string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if s.isValid(session) {
		delete(s.Sessions, session)
		return nil
	} else {
		s.Err.errorDescription = "no such session, destroy failed"
		return s.Err.Error()
	}
}
