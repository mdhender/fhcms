/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

package app

import (
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/models"
	"github.com/mdhender/fhcms/internal/way"
	"log"
	"net"
	"net/http"
	"time"
)

// New returns an initialized application server.
func New(opts ...Option) (s *Server, err error) {
	// create a server with default values
	s = &Server{}
	s.args.host, s.args.port = "", "8080"
	s.router = way.NewRouter()
	s.Addr = net.JoinHostPort(s.args.host, s.args.port)
	s.MaxHeaderBytes = 1 << 20 // 1mb?
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second

	// apply the list of options to the server
	for _, opt := range opts {
		if err = opt(s); err != nil {
			return nil, err
		}
	}

	if s.renderHome, err = homeRenderer(s.templates); err != nil {
		return nil, err
	}

	s.routes("reports", "uploads")

	return s, nil
}

type Server struct {
	http.Server
	router *way.Router
	args   struct {
		host string
		port string
	}
	data         string // path to data directory
	templates    string // path to templates directory
	accountStore models.AccountStore
	jf           *jot.Factory
	renderHome   func(p homePayload) ([]byte, error)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) currentUser(r *http.Request) models.Account {
	log.Printf("[app] server.currentUser\n")
	return models.CurrentUser(r)
}
