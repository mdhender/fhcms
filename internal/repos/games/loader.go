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

package games

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Load reads a JSON file containing game data.
func Load(filename string) (*GameList, error) {
	log.Printf("domain: games.Load: %q\n", filename)
	var input struct {
		Games map[int]*Game `json:"games"`
	}
	input.Games = make(map[int]*Game)
	if data, err := ioutil.ReadFile(filename); err != nil {
		log.Printf("domain: games.Load: %q %+v\n", filename, err)
		return nil, err
	} else if err = json.Unmarshal(data, &input); err != nil {
		log.Printf("domain: games.Load: %q %+v\n", filename, err)
		return nil, err
	}

	g := &GameList{ById: make(map[int]*Game), ByPlayer: make(map[int][]*Game)}

	for id, game := range input.Games {
		game.Id = id
		g.ById[game.Id] = game
		for userId := range game.Players {
			g.ByPlayer[userId] = append(g.ByPlayer[userId], game)
		}
	}

	// sort each player's games.
	//// the bubbly sort assumes that the strings it is using as keys are actually small integer values.
	//for _, list := range g.ByPlayer {
	//	for i := 0; i < len(list); i++ {
	//		for j := i + 1; j < len(list); j++ {
	//			if len(list[i].Id) < len(list[j].Id) {
	//				continue
	//			} else if len(list[i].Id) == len(list[j].Id) {
	//				if list[i].Id < list[j].Id {
	//					continue
	//				}
	//			}
	//			list[i], list[j] = list[j], list[i]
	//		}
	//	}
	//}
	for _, list := range g.ByPlayer {
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				if list[i].Id > list[j].Id {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	}

	for _, game := range g.ById {
		for player, species := range game.Players {
			log.Printf("domain: games.Load: %q player %q: species %q\n", filename, player, species)
		}
	}
	for _, games := range g.ByPlayer {
		for player, game := range games {
			log.Printf("domain: games.Load: %q player %d: game %v\n", filename, player, game)
		}
	}

	return g, nil
}
