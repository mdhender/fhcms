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

// Package responders implements the presentation logic for the application.
// It accepts a context from the action (which is the domain payload),
// sets the HTTP status and headers, and writes the body of the response.
package responders

//import (
//	"bytes"
//	"github.com/mdhender/fhcms/cms/config"
//	"github.com/mdhender/fhcms/internal/view"
//	"log"
//)
//
//// Responder is pretty simple for our application.
//type Responder struct {
//	id         string
//	components string
//	statusCode int
//	headers    map[string]string
//	buf        *bytes.Buffer
//	body       []byte
//	view       *view.View
//}
//
//// New returns a new Responder.
//func New(id string, components string) (*Responder, error) {
//	var v *view.View
//	var err error
//
//	cfg := config.Get()
//	if v, err = view.New("layout", cfg.Templates+"/layout.gohtml", components, cfg.Templates); err != nil {
//		log.Printf("[responder] new: %+v\n", err)
//		return nil, err
//	}
//
//	return &Responder{
//		id:         id,
//		components: components,
//		headers:    make(map[string]string),
//		buf:        bytes.NewBuffer(nil),
//		view:       v,
//	}, nil
//}
//
//// Render is
//// Body returns a slice of bytes built using the payload, or nil if the response body should be empty.
//func (r *Responder) Render(payload interface{}) (body []byte, err error) {
//	log.Printf("[responder] render\n")
//	t, err := r.view.Load()
//	if err != nil {
//		log.Printf("[responder] render load: %+v\n", err)
//		return nil, err
//	}
//	buf := bytes.NewBuffer(nil)
//	if err = t.ExecuteTemplate(r.buf, r.view.Name, payload); err != nil {
//		log.Printf("[responder] render execute: %+v\n", err)
//		return nil, err
//	}
//	log.Printf("[responder] render %d bytes\n", buf.Len())
//	return buf.Bytes(), nil
//}
//
//// ID returns the identifier for this responder.
//func (r *Responder) ID() string {
//	return r.id
//}
