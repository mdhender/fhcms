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

func (e *Engine) do_germ_warfare(attacking_species, defending_species, defender_index int, bat *battle_data, act *action_data) {
	attacker_BI := e.c_species[attacking_species].tech_level[BI]
	defender_BI := e.c_species[defending_species].tech_level[BI]
	attacked_nampla, ok := act.fighting_unit[defender_index].(*nampla_data)
	if !ok {
		panic("act.fighting_unit[defender_index].(*nampla_data); !ok")
	}

	success_chance := 50 + (2 * (attacker_BI - defender_BI))
	success := FALSE
	num_bombs := e.germ_bombs_used[attacking_species][defending_species]

	for i := 0; i < num_bombs; i++ {
		if e.rnd(100) <= success_chance {
			success = TRUE
			break
		}
	}

	if success != FALSE {
		e.log_string("        Unfortunately")
	} else {
		e.log_string("        Fortunately")
	}
	e.log_printf(" for the %s defenders of PL %s, the ", e.c_species[defending_species].name, attacked_nampla.name)
	i := bat.spec_num[attacking_species]
	if e.field_distorted[attacking_species] != FALSE {
		e.log_int(e.distorted(i))
	} else {
		e.log_string(e.c_species[attacking_species].name)
	}
	e.log_string(" attackers ")

	if success == FALSE {
		e.log_string("failed")
		if num_bombs <= 0 {
			e.log_string(" because they didn't have any germ warfare bombs")
		}
		e.log_string("!\n")
		return
	}

	e.log_string("succeeded, using ")
	e.log_int(num_bombs)
	e.log_string(" germ warfare bombs. The defenders were wiped out!\n")

	/* Take care of looting. */
	econ_units_from_looting := attacked_nampla.mi_base + attacked_nampla.ma_base

	if (attacked_nampla.status & HOME_PLANET) != 0 {
		if e.c_species[defending_species].hp_original_base < econ_units_from_looting {
			e.c_species[defending_species].hp_original_base = econ_units_from_looting
		}
		econ_units_from_looting *= 5
	}

	if econ_units_from_looting > 0 {
		/* Check if there's enough memory for a new interspecies transaction. */
		if e.num_transactions == MAX_TRANSACTIONS {
			fprintf(e.stderr, "\nRan out of memory! MAX_TRANSACTIONS is too small!\n\n")
			panic("\"\\nRan out of memory! MAX_TRANSACTIONS is too small!\\n\\n\"")
		}
		i = e.num_transactions
		e.num_transactions++

		/* Define this transaction. */
		e.transaction[i]._type = LOOTING_EU_TRANSFER
		e.transaction[i].donor = bat.spec_num[defending_species]
		e.transaction[i].recipient = bat.spec_num[attacking_species]
		e.transaction[i].value = econ_units_from_looting
		e.transaction[i].name1 = e.c_species[defending_species].name
		e.transaction[i].name2 = e.c_species[attacking_species].name
		e.transaction[i].name3 = attacked_nampla.name
	}

	/* Finish off defenders. */
	attacked_nampla.mi_base = 0
	attacked_nampla.ma_base = 0
	attacked_nampla.IUs_to_install = 0
	attacked_nampla.AUs_to_install = 0
	attacked_nampla.pop_units = 0
	attacked_nampla.siege_eff = 0
	attacked_nampla.shipyards = 0
	attacked_nampla.hiding = 0
	attacked_nampla.hidden = 0
	attacked_nampla.use_on_ambush = 0

	for i := 0; i < MAX_ITEMS; i++ {
		attacked_nampla.item_quantity[i] = 0
	}

	/* Reset status word. */
	if (attacked_nampla.status & HOME_PLANET) != 0 {
		attacked_nampla.status = HOME_PLANET
	} else {
		attacked_nampla.status = COLONY
	}

	/* Delete any ships that were under construction on the planet. */
	//sh := e.c_ship[defending_species] - 1;
	for i := 0; i < e.c_species[defending_species].num_ships; i++ {
		sh := e.c_ship[defending_species][i]

		if sh.x != attacked_nampla.x {
			continue
		} else if sh.y != attacked_nampla.y {
			continue
		} else if sh.z != attacked_nampla.z {
			continue
		} else if sh.pn != attacked_nampla.pn {
			continue
		}

		e.delete_ship(sh)
	}
}
