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
	"bytes"
	"github.com/mdhender/fhcms/internal/models"
	"html/template"
	"log"
	"path/filepath"
)

type homePayload struct {
	Site    models.Site
	Account models.Account
}

func homeRenderer(templates string) (func(p homePayload) ([]byte, error), error) {
	return func(p homePayload) ([]byte, error) {
		t, err := template.ParseFiles(filepath.Join(templates, "site.layout.gohtml"), filepath.Join(templates, "fragments", "navbar.gohtml"), filepath.Join(templates, "fragments", "footer.gohtml"), filepath.Join(templates, "home.index.gohtml"))
		if err != nil {
			log.Printf("[app] homeRenderer: %+v\n", err)
			return nil, err
		}

		b := &bytes.Buffer{}
		if err = t.ExecuteTemplate(b, "layout", p); err != nil {
			log.Printf("[app] homeRenderer: %+v\n", err)
			return nil, err
		}

		return b.Bytes(), nil
	}, nil
}
