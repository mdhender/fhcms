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

import "github.com/mdhender/fhcms/cms/prng"

func New(promptGM bool) *Engine {
	return &Engine{
		correct_spelling_required: FALSE,
		defaultPRNG:               prng.New(0xBADC0FFEE),
		input_abbr:                make([]byte, 256, 256),
		input_line:                make([]byte, 256, 256),
		log_line:                  make([]byte, 1024, 1024),
		original_line:             make([]byte, 256, 256),
		original_name:             make([]byte, 32, 32),
		prompt_gm:                 promptGM,
		stderr:                    fopen("*stderr*", nil),
		upper_name:                make([]byte, 32, 32),
	}
}
