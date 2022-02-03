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

package jsonapi

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Details []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func Error(w http.ResponseWriter, status, detail string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ErrorResponse{Details: []ErrorDetail{{Status: status, Detail: detail}}})
}
