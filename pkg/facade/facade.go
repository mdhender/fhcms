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

// Package facade implements Mat Ryer's Facade from Go Programming Blueprints.
package facade

// Facade requires a Public() method to expose an object's desires.
type Facade interface {
	Public() interface{}
}

// Public will call an object's Public method if it has one.
// Otherwise, it just returns the object.
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
