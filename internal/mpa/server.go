/*******************************************************************************
Far Horizons Engine
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

package mpa

import (
	"github.com/mdhender/fhcms/internal/domain"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/models"
	"github.com/mdhender/fhcms/internal/way"
	"net"
	"net/http"
	"time"
)

func New(host, port string, opts ...func(*Server) error) (*Server, error) {
	s := &Server{
		router: way.NewRouter(),
	}
	s.Addr = net.JoinHostPort(host, port)
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20 // 1mb?

	// apply the list of options to Store
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	s.routes("reports", "uploads")

	return s, nil
}

type Server struct {
	http.Server
	router    *way.Router
	ds        *domain.Store
	jf        *jot.Factory
	auth      AuthStore
	games     GamesStore
	profiles  ProfileStore
	site      SiteStore
	templates string // path to templates directory
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) currentUser(r *http.Request) models.Account {
	return models.CurrentUser(r)
}
