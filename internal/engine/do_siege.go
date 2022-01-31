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

package engine

func (e *Engine) do_siege(bat *battle_data, act *action_data) {
	for defender_index := 0; defender_index < act.num_units_fighting; defender_index++ {
		if act.unit_type[defender_index] == BESIEGED_NAMPLA {
			defending_nampla, ok := act.fighting_unit[defender_index].(*nampla_data)
			if !ok {
				panic("act.fighting_unit[defender_index].(*nampla_data); !ok")
			}
			defending_nampla.siege_eff = TRUE

			d := act.fighting_species_index[defender_index]
			defending_species := e.c_species[d]
			defending_species_number := bat.spec_num[d]

			for attacker_index := 0; attacker_index < act.num_units_fighting; attacker_index++ {
				if act.unit_type[attacker_index] != SHIP {
					continue
				}
				attacking_ship, ok := act.fighting_unit[attacker_index].(*ship_data)
				if !ok {
					panic("act.fighting_unit[attacker_index].(*ship_data); !ok")
				}
				a := act.fighting_species_index[attacker_index]
				if e.x_attacked_y[a][d] == FALSE {
					continue
				}
				attacking_species := e.c_species[a]
				attacking_species_number := bat.spec_num[a]

				/* Check if there's enough memory for a new interspecies transaction. */
				if e.num_transactions == MAX_TRANSACTIONS {
					fprintf(e.stderr, "\nRan out of memory! MAX_TRANSACTIONS is too small!\n\n")
					panic("\nRan out of memory! MAX_TRANSACTIONS is too small!\n\n")
				}
				i := e.num_transactions
				e.num_transactions++

				/* Define this transaction. */
				e.transaction[i]._type = BESIEGE_PLANET
				e.transaction[i].x = defending_nampla.x
				e.transaction[i].y = defending_nampla.y
				e.transaction[i].z = defending_nampla.z
				e.transaction[i].pn = defending_nampla.pn
				e.transaction[i].number1 = attacking_species_number
				e.transaction[i].name1 = attacking_species.name
				e.transaction[i].number2 = defending_species_number
				e.transaction[i].name2 = defending_species.name
				e.transaction[i].name3 = attacking_ship.name
			}
		}
	}

	e.log_string("      Only those ships that actually remain in the system will take part in the siege.\n")
}
