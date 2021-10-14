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
	"crypto/sha1"
	"encoding/base64"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/mdhender/fhcms/internal/way"
	"log"
	"net/http"
)

type Server struct {
	http.Server
	router *way.Router
	data   struct {
		Store   *cluster.Store
		Engine  *Engine
		Site    *Site
		Players []*PlayerData
		Files   map[string][]*FileData // key is species id
		Stats   map[string]*StatsData  // key is species id
		Turn    struct {
			Due      string `json:"due"`
			By       string `json:"by"`
			TimeZone string `json:"tmz"`
		}
		TurnFiles string // path to files named sp??.t??.[orders|report].txt
	}
	sessions *SessionManager
}

// serverContextKey is the context key type for storing parameters in context.Context.
type serverContextKey string

func (s *Server) notImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: not implemented\n", r.Method, r.URL.Path)
	}
}

func mketag(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return "fhapp:" + base64.URLEncoding.EncodeToString(h.Sum(nil))
}
