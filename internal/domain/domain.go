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
	"github.com/mdhender/fhcms/internal/models"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/mdhender/fhcms/internal/repos/games"
	"log"
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
	log.Printf("[domain] FetchGame %q %q\n", uid, gid)
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
				CurrentTurn: game.Turns.Current,
			}
		}
	}
	return &models.Game{}
}

func (s *Store) FetchGames(uid string) models.Games {
	log.Printf("[domain] FetchGames %q\n", uid)
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
		mg.Specie.Government.Name = "(missing name)"
		if spNo, ok := game.Players[uid]; ok {
			for _, file := range game.Turns.Files {
				if file.Turn == game.Turns.Current {
					if sp, err := s.loadSpecie(file.Files, spNo); err != nil {
						log.Printf("[domain] FetchGames %q %q %+v\n", uid, spNo, err)
					} else {
						mg.Specie.No = spNo
						mg.Specie.Name = sp.Name
						mg.Specie.Government.Name = sp.Government.Name
					}
					break
				}
			}
		}
		mg.Turns.Current = game.Turns.Current
		set = append(set, mg)
	}
	sort.Sort(set)
	log.Printf("[domain] FetchGames %q returning %d games\n", uid, len(set))
	return set
}

func (s *Store) FetchSpecie(uid, gid, spNo string, turnNo int) *models.Specie {
	log.Printf("[domain] FetchSpecie %q %q %q %d\n", uid, gid, spNo, turnNo)

	u, ok := s.Accounts.ById[uid]
	if !ok || u == nil {
		log.Printf("[domain] FetchSpecie %q %q %q %d: no such user\n", uid, gid, spNo, turnNo)
		return &models.Specie{}
	}
	g, ok := s.Games.ById[gid]
	if !ok || u == nil {
		log.Printf("[domain] FetchSpecie %q %q %q %d: no such game\n", uid, gid, spNo, turnNo)
		return &models.Specie{}
	}
	if spid, ok := g.Players[uid]; !ok || u == nil {
		log.Printf("[domain] FetchSpecie %q %q %q %d: no such player\n", uid, gid, spNo, turnNo)
		return &models.Specie{}
	} else if spNo != spid {
		log.Printf("[domain] FetchSpecie %q %q %q %d: player spoofing species!\n", uid, gid, spNo, turnNo)
		return &models.Specie{}
	}
	var gtf *games.GameTurnFile
	for _, file := range g.Turns.Files {
		if file.Turn == turnNo {
			gtf = file
			break
		}
	}
	if gtf == nil {
		log.Printf("[domain] FetchSpecie %q %q %q %d: no such turn\n", uid, gid, spNo, turnNo)
		return &models.Specie{}
	}
	sp, err := s.loadSpecie(gtf.Files, spNo)
	if err != nil {
		log.Printf("[domain] FetchSpecie %q %q %q %d: %+v\n", uid, gid, spNo, turnNo, err)
		return &models.Specie{}
	}
	o := &models.Specie{
		Id:   spNo,
		No:   spNo,
		Name: sp.Name,
	}
	o.Government.Name = sp.Government.Name
	for i := 0; i < 6; i++ {
		o.Technology[i].Code = sp.Technology[i].Code
		o.Technology[i].Name = sp.Technology[i].Name
		o.Technology[i].CurrentLevel = sp.Technology[i].CurrentLevel
		o.Technology[i].InitialLevel = sp.Technology[i].InitialLevel
		o.Technology[i].KnowledgeLevel = sp.Technology[i].KnowledgeLevel
		o.Technology[i].ExperiencePoints = sp.Technology[i].ExperiencePoints
	}
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Production", Units: "EU"})
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Banked", Units: "EU"})
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Colonies"})
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Ships"})
	for _, stat := range sp.Stats {
		o.Statistics = append(o.Statistics, &models.SpecieStatistic{
			Label: stat.Label,
			Value: stat.Value,
			Units: stat.Units,
		})
	}
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Offensive Power"})
	o.Statistics = append(o.Statistics, &models.SpecieStatistic{Label: "Defensive Power"})
	return o
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
