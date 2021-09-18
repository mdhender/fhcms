/*******************************************************************************
fhapp - a player aid for Far Horizons
Copyright (C) 2021  Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

******************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mdhender/fhcms/internal/cluster"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	Uuid      string `json:"uuid"`
	ExpiresAt time.Time
	User      UserData `json:"user"`
}

type SessionManager struct {
	lock       sync.Mutex // protects session
	cookieName string
	sessions   map[string]Session
	players    []*PlayerData
	species    map[string]*cluster.Species
}

func NewSessionManager(name string) *SessionManager {
	return &SessionManager{
		cookieName: "_" + name,
		sessions:   make(map[string]Session),
	}
}

func (s *SessionManager) SessionStart(u UserData) Session {
	sess := Session{
		Uuid:      uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
		User:      u,
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessions[sess.Uuid] = sess
	if b, err := json.MarshalIndent(s, "  ", "  "); err == nil {
		if err = ioutil.WriteFile("testdata/sessions.json", b, 0600); err != nil {
			log.Printf("error: %+v\n", err)
		}
	}

	return sess
}

func (s *SessionManager) SessionGet(uuid string) Session {
	s.lock.Lock()
	if s.sessions == nil {
		s.sessions = make(map[string]Session)
	}
	sess, ok := s.sessions[uuid]
	if ok {
		if !sess.ExpiresAt.After(time.Now()) {
			delete(s.sessions, uuid)
			sess = Session{Uuid: uuid}
		}
	}
	s.lock.Unlock()
	return sess
}

func (s *SessionManager) SessionDelete(uuid string) {
	s.lock.Lock()
	delete(s.sessions, uuid)
	s.lock.Unlock()
}

// sessionContextKey is the context key type for storing parameters in context.Context.
type sessionContextKey string

// SessionUserHandler extracts the session from the cookie and uses it to find the associated user.
// If there is no cookie, the session is invalid, or there is no associated user, then the "zero" user is added.
func (s *SessionManager) SessionUserHandler(h http.Handler) http.HandlerFunc {
	log.Printf("sessions: adding session user handler as middleware\n")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("adapter: %s %q\n", r.Method, r.URL.Path)
		var user UserData
		if cookie, err := r.Cookie(s.cookieName); err == nil {
			user = s.SessionGet(cookie.Value).User
		}
		h.ServeHTTP(w, r.WithContext(user.NewContext(r.Context())))
	}
}

func (s *SessionManager) MarshalJSON() ([]byte, error) {
	// don't bother saving sessions that are about to expire
	now := time.Now().Add(15 * time.Second)
	type export struct {
		Uuid      string `json:"uuid"`
		ExpiresAt string `json:"expires_at"`
		Player    string `json:"player"`
	}
	var data []export
	for _, sess := range s.sessions {
		if sess.ExpiresAt.After(now) {
			data = append(data, export{
				Uuid:      sess.Uuid,
				ExpiresAt: sess.ExpiresAt.Format(time.RFC3339),
				Player:    sess.User.Player,
			})
		}
	}
	return json.Marshal(data)
}

func (s *SessionManager) UnmarshalJSON(data []byte) error {
	var sessions []struct {
		Uuid      string `json:"uuid"`
		ExpiresAt string `json:"expires_at"`
		Player    string `json:"player"`
	}
	if err := json.Unmarshal(data, &sessions); err != nil {
		return err
	}
	if s.sessions == nil {
		s.sessions = make(map[string]Session)
	}
	for i, sess := range sessions {
		expiresAt, err := time.Parse(time.RFC3339, sess.ExpiresAt)
		if err != nil {
			return fmt.Errorf("error importing session %d: %+v", i+1, err)
		}
		var u UserData
		for _, p := range s.players {
			if p.User == sess.Player {
				if sp, ok := s.species[p.SpeciesId]; ok {
					u.Player = p.User // confusing, I know
					u.Species = sp
					u.SpeciesId = sp.Id
					u.IsAuthenticated = true
					break
				}
			}
		}
		s.sessions[sess.Uuid] = Session{
			Uuid:      sess.Uuid,
			ExpiresAt: expiresAt,
			User:      u,
		}
	}
	return nil
}
