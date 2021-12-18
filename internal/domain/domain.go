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

package domain

import (
	"encoding/binary"
	"github.com/mdhender/fhcms/internal/dat32"
	"github.com/mdhender/fhcms/internal/models"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/mdhender/fhcms/internal/repos/games"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"sort"
)

type Store struct {
	Accounts *accounts.AccountList
	Games    *games.GameList
	Species  *games.SpecieList
}

func New(opts ...func(*Store) error) (*Store, error) {
	s := &Store{}

	// apply the list of options to Store
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) FetchGame(uid, gid string) *models.Game {
	log.Printf("domain: FetchGame %q %q\n", uid, gid)
	if s.Games == nil {
		return &models.Game{}
	}
	for _, game := range s.Games.ByPlayer[uid] {
		if game.Id == gid {
			return &models.Game{
				Id:          game.Id,
				Name:        game.Name,
				Files:       game.Files,
				PlayerCount: len(game.Players),
			}
		}
	}
	return &models.Game{}
}

func (s *Store) FetchGames(uid string) models.Games {
	log.Printf("domain: FetchGames %q\n", uid)
	var set models.Games
	if s.Games == nil {
		return set
	}
	for _, game := range s.Games.ByPlayer[uid] {
		mg := &models.Game{
			Id:          game.Id,
			Name:        game.Name,
			Files:       game.Files,
			PlayerCount: len(game.Players),
		}
		mg.Specie.No = "00"
		mg.Specie.Name = "(missing name)"
		if spNo, ok := game.Players[uid]; ok {
			for _, file := range game.Turns.Files {
				if file.Turn == game.Turns.Current {
					if sp, err := s.loadSpecie(file.Files, spNo); err != nil {
						log.Printf("[domain] FetchGames %q %q %+v\n", uid, spNo, err)
					} else {
						mg.Specie.No = spNo
						mg.Specie.Name = sp.Name
					}
					break
				}
			}
		}
		mg.Turns.Current = game.Turns.Current
		set = append(set, mg)
	}
	sort.Sort(set)
	log.Printf("domain: FetchGames %q returning %d games\n", uid, len(set))
	return set
}

func (s *Store) FetchSpecie(uid, gid string) *models.Specie {
	log.Printf("domain: FetchSpecie %q %q\n", uid, gid)
	bigEndian := viper.GetBool("files.big_endian")
	var bo binary.ByteOrder
	if bigEndian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}

	u, ok := s.Accounts.ById[uid]
	if !ok || u == nil {
		log.Printf("domain: FetchSpecie %q %q: no such user\n", uid, gid)
		return &models.Specie{}
	}
	g, ok := s.Games.ById[gid]
	if !ok || u == nil {
		log.Printf("domain: FetchSpecie %q %q: no such game\n", uid, gid)
		return &models.Specie{}
	}
	spid, ok := g.Players[uid]
	if !ok || u == nil {
		log.Printf("domain: FetchSpecie %q %q: no such player\n", uid, gid)
		return &models.Specie{}
	}
	sp, err := dat32.ReadSpecies(filepath.Join(g.Files, "sp"+spid+".dat"), 0, bo)
	if err != nil {
		log.Printf("[domain] FetchSpecies %q %+v\n", spid, err)
		return &models.Specie{}
	}
	return &models.Specie{
		Id:   spid,
		Name: sp.Name,
	}
}

func (s *Store) FetchUser(uid string) *models.User {
	if s.Accounts == nil {
		return &models.User{}
	}
	u, ok := s.Accounts.ById[uid]
	if !ok {
		return &models.User{}
	}
	return &models.User{
		Id:      u.Id,
		Name:    u.Username,
		IsAdmin: u.Username == "mdhender",
	}
}
