//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

//TODO: This time duration is to long time to time we could have error in test
const CheckSessionDuration = 5 * time.Second

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	mu       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data         map[string]interface{}
	lastActivity time.Time
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}
	go sessionCleaner(m)
	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	m.sessions[sessionID] = Session{
		Data: make(map[string]interface{}),
	}
	m.mu.Unlock()
	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mu.Lock()
	session, ok := m.sessions[sessionID]
	m.mu.Unlock()
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	m.mu.Lock()
	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:         data,
		lastActivity: time.Now(),
	}
	m.mu.Unlock()
	return nil
}

func (m *SessionManager) DeleteSession(sessionId string) {
	delete(m.sessions, sessionId)
}

func (m *SessionManager) DeleteOldSession(timeTicker time.Time) {
	m.mu.Lock()
	for i, s := range m.sessions {
		fmt.Println(timeTicker.Sub(s.lastActivity))
		if timeTicker.Sub(s.lastActivity) > CheckSessionDuration {
			m.DeleteSession(i)
		}
	}
	m.mu.Unlock()
}

func sessionCleaner(m *SessionManager) {
	c := time.Tick(CheckSessionDuration)
	for timeTicker := range c {
		m.DeleteOldSession(timeTicker)
	}
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()

	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	if err = m.UpdateSessionData(sID, data); err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}

// MakeSessionID is used to generate a random dummy sessionID
func MakeSessionID() (string, error) {
	buf := make([]byte, 26)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}
