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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *Server) authenticatedOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAuthenticated {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *Server) staticFiles(root string, h http.HandlerFunc) http.HandlerFunc {
	type cache struct {
		b []byte
		etag string
		contentType string
	}
	files := make(map[string]*cache)
	for _, file := range []string {
		filepath.Clean(filepath.Join(root, "browserconfig.xml")),
		filepath.Clean(filepath.Join(root, "favicon.ico")),
		filepath.Clean(filepath.Join(root, "humans.txt")),
		filepath.Clean(filepath.Join(root, "icon.png")),
		filepath.Clean(filepath.Join(root, "robots.txt")),
		filepath.Clean(filepath.Join(root, "rules.html")),
		filepath.Clean(filepath.Join(root, "site.webmanifest")),
		filepath.Clean(filepath.Join(root, "tile-wide.png")),
		filepath.Clean(filepath.Join(root, "tile.png")),
		filepath.Clean(filepath.Join(root, "css", "daleri-mega.css")),
		filepath.Clean(filepath.Join(root, "css", "main.css")),
		filepath.Clean(filepath.Join(root, "css", "normalize.css")),
		filepath.Clean(filepath.Join(root, "css", "img", "bodybg.gif")),
		filepath.Clean(filepath.Join(root, "js", "vendor", "modernizr-3.11.2.min.js")),
	}{
		log.Printf("[static] cacheing %q\n", file)
		if b, err := ioutil.ReadFile(file); err != nil {
			panic(err)
		} else {
			var contentType string
			switch filepath.Ext(file) {
			case ".css":
				contentType = "text/css; charset=utf-8"
			case ".gif":
				contentType = "image/gif"
			case ".html":
				contentType = "text/html"
			case ".ico":
				contentType = "image/x-icon"
			case ".js":
				contentType = "text/javascript"
			case ".png":
				contentType = "image/png"
			case ".txt":
				contentType = "text/plain; charset=utf-8"
			case ".webmanifest":
				contentType = "text/plain; charset=utf-8"
			case ".xml":
				contentType = "text/xml"
			}
			h := sha1.New()
			h.Write(b)
			files[filepath.Base(file)]= &cache{contentType: contentType, b: b, etag: "fhapp:" + base64.URLEncoding.EncodeToString(h.Sum(nil))}
		}
	}
	maxAge := fmt.Sprintf("max-age=%d", 4*7*24*60*60) // 4 weeks
	return func(w http.ResponseWriter, r *http.Request) {
		c, ok := files[filepath.Base(r.URL.Path)]
		if !ok {
			log.Printf("server: %s %q: pass from handleStaticFile\n", r.Method, r.URL.Path)
			h(w, r)
			return
		}
		log.Printf("server: %s %q: handleStaticFile\n", r.Method, r.URL.Path)
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, c.etag) {
				log.Printf("server: %s %q: handleStaticFile: cached\n", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
		w.Header().Set("Content-Type", c.contentType)
		w.Header().Set("Etag", c.etag)
		w.Header().Set("Cache-Control", maxAge)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(c.b)
	}
}
