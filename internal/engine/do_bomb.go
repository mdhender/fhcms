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

func (e *Engine) do_bombardment(unit_index int, act *action_data) {
	attacked_nampla, ok := act.fighting_unit[unit_index].(*nampla_data)
	if !ok {
		panic("act.fighting_unit[unit_index].(*nampla_data); !ok")
	}
	initial_base := attacked_nampla.mi_base + attacked_nampla.ma_base
	total_pop := initial_base
	if attacked_nampla.item_quantity[CU] > 0 {
		total_pop += 1
	}
	if total_pop < 1 {
		e.log_string("        The planet is completely uninhabited. There is nothing to bomb!\n")
		return
	}

	// total damage done by ten strike cruisers (ML = 50) in ten rounds is 100 x 4 x the power value for a single ship.
	// to eliminate the chance of overflow, the algorithm has been carefully chosen.
	CS_bomb_damage := 400 * power(ship_tonnage[CS]) // should be 400 * 4759 = 1,903,600

	total_bomb_damage := act.bomb_damage[unit_index]

	/* Keep about 2 significant digits. */
	for total_bomb_damage > 1000 {
		total_bomb_damage /= 10
		CS_bomb_damage /= 10
	}

	var percent_damage int
	if CS_bomb_damage == 0 {
		percent_damage = 101
	} else {
		percent_damage = ((total_bomb_damage * 250000) / CS_bomb_damage) / total_pop
	}

	if percent_damage > 100 {
		percent_damage = 101
	}

	new_mi := attacked_nampla.mi_base - (percent_damage*attacked_nampla.mi_base)/100
	new_ma := attacked_nampla.ma_base - (percent_damage*attacked_nampla.ma_base)/100
	new_pop := attacked_nampla.pop_units - (percent_damage*attacked_nampla.pop_units)/100

	if new_mi == attacked_nampla.mi_base && new_ma == attacked_nampla.ma_base && new_pop == attacked_nampla.pop_units {
		e.log_string("        Damage due to bombardment was insignificant.\n")
		return
	}

	defending_species := act.fighting_species_index[unit_index]
	if (attacked_nampla.status & HOME_PLANET) != 0 {
		n := attacked_nampla.mi_base + attacked_nampla.ma_base
		if e.c_species[defending_species].hp_original_base < n {
			e.c_species[defending_species].hp_original_base = n
		}
	}

	if new_mi <= 0 && new_ma <= 0 && new_pop <= 0 {
		e.log_string("        Everyone and everything was completely wiped out!\n")

		attacked_nampla.mi_base = 0
		attacked_nampla.ma_base = 0
		attacked_nampla.pop_units = 0
		attacked_nampla.siege_eff = 0
		attacked_nampla.shipyards = 0
		attacked_nampla.hiding = 0
		attacked_nampla.hidden = 0
		attacked_nampla.use_on_ambush = 0

		/* Reset status. */
		if (attacked_nampla.status & HOME_PLANET) != 0 {
			attacked_nampla.status = HOME_PLANET
		} else {
			attacked_nampla.status = COLONY
		}

		for i := 0; i < MAX_ITEMS; i++ {
			attacked_nampla.item_quantity[i] = 0
		}

		/* Delete any ships that were under construction on the planet. */
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

		return
	}

	e.log_printf("        Mining base of PL %s went from %d.%d to %d.%d.\n", attacked_nampla.name, attacked_nampla.mi_base/10, attacked_nampla.mi_base%10, new_mi/10, new_mi%10)
	e.log_printf("        Manufacturing base of PL %s went from %d.%d to %d.%d.\n", attacked_nampla.name, attacked_nampla.ma_base/10, attacked_nampla.ma_base%10, new_ma/10, new_ma%10)

	attacked_nampla.mi_base = new_mi
	attacked_nampla.ma_base = new_ma
	attacked_nampla.pop_units = new_pop

	for i := 0; i < MAX_ITEMS; i++ {
		if n := (percent_damage * attacked_nampla.item_quantity[i]) / 100; n > 0 {
			attacked_nampla.item_quantity[i] -= n
			e.log_printf("        %d %s", n, item_name[i])
			if n > 1 {
				e.log_string("s were")
			} else {
				e.log_string(" was")
			}
			e.log_string(" destroyed.\n")
		}
	}

	if n := (percent_damage * attacked_nampla.shipyards) / 100; n > 0 {
		attacked_nampla.shipyards -= n
		e.log_printf("        %d shipyard")
		if n > 1 {
			e.log_string("s were")
		} else {
			e.log_string(" was")
		}
		e.log_string(" also destroyed.\n")
	}

	e.check_population(attacked_nampla)
}
