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
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/mdhender/fhcms/internal/repos/site"
	"path/filepath"
)

type Option func(*Server) error

// Options turns a list of Option instances into an Option.
func Options(opts ...Option) Option {
	return func(s *Server) error {
		for _, opt := range opts {
			if err := opt(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithAccounts(accts *accounts.AccountList) Option {
	return func(s *Server) (err error) {
		s.accts = accts
		return nil
	}
}

func WithDomain(ds *domain.Store) Option {
	return func(s *Server) (err error) {
		s.ds = ds
		return nil
	}
}

func WithJotFactory(jf *jot.Factory) Option {
	return func(s *Server) (err error) {
		s.jf = jf
		return nil
	}
}

func WithSite(ds *site.Store) Option {
	return func(s *Server) (err error) {
		s.site = ds
		return nil
	}
}

func WithTemplates(root string) Option {
	return func(s *Server) (err error) {
		s.templates = filepath.Clean(root)
		return nil
	}
}
