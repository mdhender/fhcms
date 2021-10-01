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
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet

// templateHandler implements an HTTP handler that renders a single template
type templateHandler struct {
	once     sync.Once
	prod     bool               // if true, parse the template only once
	filename string             // path to template to render
	t        *template.Template // a single template
}

// ServeHTTP handles the HTTP request by rendering the template
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := t.Render(nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// Render renders the template with the given data
func (t *templateHandler) Render(data interface{}) ([]byte, error) {
	if t.prod {
		t.once.Do(func() {
			t.t = template.Must(template.ParseFiles(filepath.Join("templates", t.filename+".html")))
		})
	}
	tmplt := t.t
	if tmplt == nil {
		var err error
		tmplt, err = template.ParseFiles(filepath.Join("templates", t.filename+".html"))
		if err != nil {
			return nil, err
		}
	}
	var br bytes.Buffer
	if err := tmplt.Execute(&br, data); err != nil {
		return nil, err
	}
	return br.Bytes(), nil
}

func (s *Server) render(name string, data interface{}) ([]byte, error) {
	t, err := template.ParseFiles(filepath.Join("templates", name+".html"))
	if err != nil {
		return nil, err
	}
	var br bytes.Buffer
	if err = t.Execute(&br, data); err != nil {
		return nil, err
	}
	return br.Bytes(), nil
}
