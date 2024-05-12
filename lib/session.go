package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

var sessions = struct {
	sync.RWMutex
	store map[string]*Session
}{
	store: make(map[string]*Session),
}

func createSessionId() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:]), nil
}

func NewSession(userID int) (*Session, error) {
	sessionId, err := createSessionId()
	if err != nil {
		return nil, err
	}
	newSession := &Session{
		ID:        sessionId,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessions.Lock()
	sessions.store[sessionId] = newSession
	sessions.Unlock()

	return newSession, nil
}

func IsValidSession(sessionCookie string) bool {
	sessions.RLock()
	session, exists := sessions.store[sessionCookie]
	sessions.RUnlock()
	return exists && time.Now().Before(session.ExpiresAt)
}

func DeleteSession(sessionCookie string) {
	sessions.Lock()
	delete(sessions.store, sessionCookie)
	sessions.Unlock()
}

func GetUserByCookie(sessionCookie string) *Session {
	return sessions.store[sessionCookie]
}
