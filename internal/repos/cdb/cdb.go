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

// Package cdb implements a CockroachDB repository.
package cdb

// https://stackoverflow.com/questions/69404758/transaction-in-golang-with-pgx
// https://stackoverflow.com/questions/64357313/postgres-table-batch-updates-using-golang-pgxpool-not-reflected-in-database
// b.Queue("INSERT INTO colonies VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", id, sp.Id, planetId, isHomeworld)

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mdhender/fhcms/internal/dat32"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type DBConfig struct {
	Cockroachdb struct {
		ConnectString  string `json:"connect_string"`
		Host           string `json:"host"`
		Port           int    `json:"port"`
		Authentication struct {
			User     string `json:"user"`
			Password string `json:"password"`
		} `json:"authentication"`
		Cluster  string `json:"cluster"`
		Database string `json:"database"`
		Ssl      struct {
			Mode string `json:"mode"`
			Cert string `json:"cert"`
		} `json:"ssl"`
	} `json:"cockroachdb"`
	Postgres struct {
		ConnectString  string `json:"connect_string"`
		Host           string `json:"host"`
		Port           int    `json:"port"`
		Authentication struct {
			User     string `json:"user"`
			Password string `json:"password"`
		} `json:"authentication"`
		Database string `json:"database"`
	} `json:"postgres"`
}

type DB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *DBConfig) (*DB, error) {
	log.Printf("[cdb] new: connect %q\n", cfg.Postgres.ConnectString)
	if pool, err := pgxpool.Connect(context.Background(), cfg.Postgres.ConnectString); err != nil {
		return nil, err
	} else {
		return &DB{pool: pool}, nil
	}
}

func (db *DB) Close() {
	db.pool.Close()
}

type GalaxySetup struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	NumSpecies int    `json:"-"`
	InitGalaxy bool   `json:"init_galaxy"`
	Root       string `json:"root"`
	Files      []struct {
		Root   string `json:"root"`
		Loaded bool   `json:"loaded"`
	} `json:"files"`
	Players []struct {
		AccountId int    `json:"account_id"`
		UserName  string `json:"username"`
		Species   int    `json:"species"`
	} `json:"players"`
	// Systems is a map of system_no to system_id
	Systems map[int]int `json:"systems"`
	// Planets is a map of planet_no to planet_id
	Planets map[int]int `json:"planets"`
	// Species is a map of species_no to species_id
	Species map[int]int `json:"species"`
	Endian  binary.ByteOrder
}

func (gs *GalaxySetup) GasCode(i int) string {
	switch i {
	case 1:
		return "H2"
	case 2:
		return "CH4"
	case 3:
		return "He"
	case 4:
		return "NH3"
	case 5:
		return "N2"
	case 6:
		return "CO2"
	case 7:
		return "O2"
	case 8:
		return "HCl"
	case 9:
		return "Cl2"
	case 10:
		return "F2"
	case 11:
		return "H2O"
	case 12:
		return "SO2"
	case 13:
		return "H2S"
	}
	return " "
}

func (gs *GalaxySetup) ItemCode(i int) (string, string) {
	switch i {
	case 0:
		return "RM", "Raw Material Unit"
	case 1:
		return "PD", "Planetary Defense Unit"
	case 2:
		return "SU", "Starbase Unit"
	case 3:
		return "DR", "Damage Repair Unit"
	case 4:
		return "CU", "Colonist Unit"
	case 5:
		return "IU", "Colonial Mining Unit"
	case 6:
		return "AU", "Colonial Manufacturing Unit"
	case 7:
		return "FS", "Fail-Safe Jump Unit"
	case 8:
		return "JP", "Jump Portal Unit"
	case 9:
		return "FM", "Forced Misjump Unit"
	case 10:
		return "FJ", "Forced Jump Unit"
	case 11:
		return "GT", "Gravitic Telescope Unit"
	case 12:
		return "FD", "Field Distortion Unit"
	case 13:
		return "TP", "Terraforming Plant"
	case 14:
		return "GW", "Germ Warfare Bomb"
	case 15:
		return "SG1", "Mark-1 Shield Generator"
	case 16:
		return "SG2", "Mark-2 Shield Generator"
	case 17:
		return "SG3", "Mark-3 Shield Generator"
	case 18:
		return "SG4", "Mark-4 Shield Generator"
	case 19:
		return "SG5", "Mark-5 Shield Generator"
	case 20:
		return "SG6", "Mark-6 Shield Generator"
	case 21:
		return "SG7", "Mark-7 Shield Generator"
	case 22:
		return "SG8", "Mark-8 Shield Generator"
	case 23:
		return "SG9", "Mark-9 Shield Generator"
	case 24:
		return "GU1", "Mark-1 Gun Unit"
	case 25:
		return "GU2", "Mark-2 Gun Unit"
	case 26:
		return "GU3", "Mark-3 Gun Unit"
	case 27:
		return "GU4", "Mark-4 Gun Unit"
	case 28:
		return "GU5", "Mark-5 Gun Unit"
	case 29:
		return "GU6", "Mark-6 Gun Unit"
	case 30:
		return "GU7", "Mark-7 Gun Unit"
	case 31:
		return "GU8", "Mark-8 Gun Unit"
	case 32:
		return "GU9", "Mark-9 Gun Unit"
	case 33:
		return "X1", "X1 Unit"
	case 34:
		return "X2", "X2 Unit"
	case 35:
		return "X3", "X3 Unit"
	case 36:
		return "X4", "X4 Unit"
	case 37:
		return "X5", "X5 Unit"
	}
	panic(fmt.Sprintf("assert(item != %d)", i))
}

func (gs *GalaxySetup) power(tonnage int) int {
	if tonnage <= 0 {
		return 0
	} else if tonnage <= 100 {
		return [101]int{
			0, // zeroth element required but not used
			100, 230, 374, 528, 690, 859, 1033, 1213, 1397, 1585,
			1777, 1973, 2171, 2373, 2578, 2786, 2996, 3209, 3424, 3641,
			3861, 4082, 4306, 4532, 4759, 4988, 5220, 5452, 5687, 5923,
			6161, 6400, 6641, 6883, 7127, 7372, 7618, 7866, 8115, 8365,
			8617, 8870, 9124, 9379, 9635, 9893, 10151, 10411, 10672, 10934,
			11197, 11461, 11725, 11991, 12258, 12526, 12795, 13065, 13336, 13608,
			13880, 14154, 14428, 14703, 14979, 15256, 15534, 15813, 16092, 16373,
			16654, 16936, 17218, 17502, 17786, 18071, 18356, 18643, 18930, 19218,
			19507, 19796, 20086, 20377, 20668, 20960, 21253, 21547, 21841, 22136,
			22431, 22727, 23024, 23321, 23619, 23918, 24217, 24517, 24818, 25119,
		}[tonnage]
	} else if tonnage <= 4068 {
		// Break it up into two halves and use recursion to get approximate result = 1.149 * (x1 + x2).
		return 1149 * (gs.power(tonnage/2) + gs.power(tonnage-tonnage/2)) / 1000
	}
	log.Printf("[cdb] integer overflow will occur in call to 'power(%d)'!\n", tonnage)
	return 0
}

func (db *DB) LoadGalaxy(setupJson string, initGalaxy, forcePurge bool) error {
	ctx := context.Background()

	setup := &GalaxySetup{
		Systems: make(map[int]int),
		Planets: make(map[int]int),
		Species: make(map[int]int),
	}
	if bigEndian := viper.GetBool("files.big_endian"); bigEndian {
		setup.Endian = binary.BigEndian
	} else {
		setup.Endian = binary.LittleEndian
	}
	b, err := ioutil.ReadFile(setupJson)
	if err != nil {
		return err
	} else if err = json.Unmarshal(b, &setup); err != nil {
		return err
	} else if len(setup.Name) == 0 {
		return fmt.Errorf("missing galaxy name")
	} else if len(setup.Name) > 32 {
		return fmt.Errorf("exceeded galaxy name length of 32")
	} else if len(setup.Files) == 0 {
		return fmt.Errorf("no files to load")
	}
	setup.Id, setup.InitGalaxy = 0, initGalaxy

	// load the accounts and be whiny about it
	log.Printf("[load] fetching account information for galaxy %q\n", setup.Name)
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	for k, v := range setup.Players {
		if err := conn.QueryRow(ctx, "select id from accounts where upper(username) = upper($1)", v.UserName).Scan(&setup.Players[k].AccountId); err != nil {
			log.Printf("[load] user %q does not exist\n", v.UserName)
			return err
		}
	}

	log.Printf("[load] fetching information for galaxy %q\n", setup.Name)

	// see if the galaxy already exists
	if err = conn.QueryRow(ctx, "select id, num_species from galaxies where upper(name) = upper($1)", setup.Name).Scan(&setup.Id, &setup.NumSpecies); err != nil {
		// the galaxy doesn't exist
		if !initGalaxy {
			// we must return an error if we're not initializing it
			return fmt.Errorf("no such galaxy")
		}
		// fetch next available galaxy id to use for the initial load
		if err = conn.QueryRow(ctx, "select nextval('seq_galaxy')").Scan(&setup.Id); err != nil {
			return err
		}
	} else if initGalaxy && !forcePurge {
		return fmt.Errorf("init requested but purge not enabled")
	}
	if setup.Id == 0 {
		log.Printf("[load] internal error: galaxy id is still zero\n")
		return fmt.Errorf("unable to fetch galaxy id")
	}

	if setup.InitGalaxy {
		log.Printf("[load] initializing galaxy %d %q root %q\n", setup.Id, setup.Name, setup.Files[0].Root)
		if err := db.initGalaxy(setup, filepath.Join(setup.Root, setup.Files[0].Root), forcePurge, ctx); err != nil {
			return err
		}
		log.Printf("[load] initialized  galaxy %d %q root %q\n", setup.Id, setup.Name, setup.Files[0].Root)
		return nil
	}

	for i, t := range setup.Files {
		if t.Loaded {
			//log.Printf("[load] skipping galaxy %d %q root %q (%d)\n", setup.Id, setup.Name, t.Root, i)
			continue
		}
		log.Printf("[load] staging galaxy %d %q root %q (%d)\n", setup.Id, setup.Name, t.Root, i)
		if err := db.stageGalaxy(setup, filepath.Join(setup.Root, setup.Files[i].Root), ctx); err != nil {
			return err
		} else if err = db.stageSystems(setup, filepath.Join(setup.Root, setup.Files[i].Root), ctx); err != nil {
			return err
		}
		for no := 1; no <= setup.NumSpecies; no++ {
			if err = db.stageSpecies(setup, filepath.Join(setup.Root, setup.Files[i].Root), no, ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) initGalaxy(gs *GalaxySetup, root string, forcePurge bool, ctx context.Context) error {
	// read in all the data before starting any of the database updates
	galaxy, err := dat32.ReadGalaxy(filepath.Join(root, "galaxy.dat"), gs.Endian)
	if err != nil {
		return err
	}
	stars, err := dat32.ReadStars(filepath.Join(root, "stars.dat"), gs.Endian)
	if err != nil {
		return err
	}
	planets, err := dat32.ReadPlanets(filepath.Join(root, "planets.dat"), gs.Endian)
	if err != nil {
		return err
	}
	var species []*dat32.Species
	for no := 0; no < galaxy.NumSpecies; no++ {
		if sp, err := dat32.ReadSpecies(filepath.Join(root, fmt.Sprintf("sp%02d.dat", no+1)), no+1, gs.Endian); err != nil {
			return err
		} else {
			species = append(species, sp)
		}
	}

	log.Printf("[init] initializing galaxy %d %q purge %v\n", gs.Id, gs.Name, forcePurge)

	if forcePurge {
		log.Printf("[init] purging galaxy %d %q\n", gs.Id, gs.Name)
		if tx, err := db.pool.Begin(ctx); err != nil {
			return err
		} else {
			b := &pgx.Batch{}
			b.Queue("delete from colonies where planet_id in (select id from planets where galaxy_id = $1)", gs.Id)
			b.Queue("delete from species where galaxy_id = $1", gs.Id)
			b.Queue("delete from planet_turn where planet_id in (select id from planets where galaxy_id = $1)", gs.Id)
			b.Queue("delete from planets where galaxy_id = $1", gs.Id)
			b.Queue("delete from systems where galaxy_id = $1", gs.Id)
			b.Queue("delete from galaxies where id = $1", gs.Id)
			for batchResults := tx.SendBatch(ctx, b); ; {
				rows, err := batchResults.Query()
				rows.Close()
				if err != nil {
					log.Printf("[batch.loop] purge %+v\n", err)
					break
				}
			}
			if err = tx.Commit(ctx); err != nil {
				return err
			}
		}
	}

	if tx, err := db.pool.Begin(ctx); err != nil {
		return err
	} else {
		b := &pgx.Batch{}
		b.Queue("INSERT INTO galaxies VALUES ($1, $2, $3, $4, 0, true)", gs.Id, gs.Name, galaxy.Radius, galaxy.NumSpecies)
		for id, system := range stars.Stars {
			systemNo := id
			b.Queue("INSERT INTO systems VALUES(nextval('serial'), $1, $2, $3, $4, $5, $6)", gs.Id, system.X, system.Y, system.Z, system.NumPlanets, systemNo)
		}
		for batchResults := tx.SendBatch(ctx, b); ; {
			rows, err := batchResults.Query()
			rows.Close()
			if err != nil {
				log.Printf("[batch.loop] systems %+v\n", err)
				break
			}
		}
		if err = tx.Commit(ctx); err != nil {
			return err
		}
	}

	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// verify that we have keys for all the systems before loading the planet data
	if rows, err := conn.Query(ctx, "select system_no, id from systems where galaxy_id = $1", gs.Id); err != nil {
		return err
	} else {
		for rows.Next() {
			var no, id int
			if err = rows.Scan(&no, &id); err != nil {
				return err
			}
			gs.Systems[no] = id
		}
		rows.Close()
	}
	for id := range stars.Stars {
		systemNo := id
		if systemId, ok := gs.Systems[systemNo]; !ok {
			return fmt.Errorf("system_no %d has no system_id", systemNo)
		} else if systemId == 0 {
			return fmt.Errorf("system_no %d missing system_id", systemNo)
		}
	}

	if tx, err := db.pool.Begin(ctx); err != nil {
		return err
	} else {
		b := &pgx.Batch{}
		for id, system := range stars.Stars {
			systemNo := id
			systemId := gs.Systems[systemNo]
			for orbit, planet := range planets.Planets[system.PlanetIndex : system.PlanetIndex+system.NumPlanets] {
				planetNo := system.PlanetIndex + orbit
				b.Queue("INSERT INTO planets VALUES(nextval('serial'), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", gs.Id, systemId, system.X, system.Y, system.Z, orbit+1, planet.Diameter, planet.Gravity, systemNo, planetNo)
			}
		}

		for batchResults := tx.SendBatch(ctx, b); ; {
			rows, err := batchResults.Query()
			rows.Close()
			if err != nil {
				log.Printf("[batch.loop] planets %+v\n", err)
				break
			}
		}
		if err = tx.Commit(ctx); err != nil {
			return err
		}
	}

	// verify that we have keys for all the planets before loading the species data
	if rows, err := conn.Query(ctx, "select planet_no, id from planets where galaxy_id = $1", gs.Id); err != nil {
		return err
	} else {
		for rows.Next() {
			var no, id int
			if err = rows.Scan(&no, &id); err != nil {
				return err
			}
			gs.Planets[no] = id
		}
		rows.Close()
	}

	if tx, err := db.pool.Begin(ctx); err != nil {
		return err
	} else {
		b := &pgx.Batch{}

		for _, sp := range species {
			if sp != nil {
				accountId, homeWorldId := 0, 0
				for _, v := range gs.Players {
					if v.Species == sp.Id {
						accountId = v.AccountId
					}
				}
				// TODO: tech knowledge is not exactly the initial tech level for the species.
				// But, early enough in the game, it's close.
				b.Queue("INSERT INTO species VALUES (nextval('serial'), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)",
					gs.Id,
					sp.Name, sp.GovtName, sp.GovtType,
					homeWorldId,
					sp.TechKnowledge[0], sp.TechKnowledge[1], sp.TechKnowledge[2], sp.TechKnowledge[3], sp.TechKnowledge[4], sp.TechKnowledge[5],
					gs.GasCode(sp.RequiredGas), sp.RequiredGasMin, sp.RequiredGasMax,
					gs.GasCode(sp.NeutralGas[0]), gs.GasCode(sp.NeutralGas[1]), gs.GasCode(sp.NeutralGas[2]), gs.GasCode(sp.NeutralGas[3]), gs.GasCode(sp.NeutralGas[4]), gs.GasCode(sp.NeutralGas[5]),
					gs.GasCode(sp.PoisonGas[0]), gs.GasCode(sp.PoisonGas[1]), gs.GasCode(sp.PoisonGas[2]), gs.GasCode(sp.PoisonGas[3]), gs.GasCode(sp.PoisonGas[4]), gs.GasCode(sp.PoisonGas[5]),
					sp.Id, sp.NamplaBase[0].PlanetIndex,
					accountId)
			}
		}
		b.Queue("update species set homeworld_id = (select id from planets where planets.galaxy_id = species.galaxy_id and planets.planet_no = species.planet_no) where species.galaxy_id = $1 and species.homeworld_id = 0", gs.Id)

		for batchResults := tx.SendBatch(ctx, b); ; {
			rows, err := batchResults.Query()
			rows.Close()
			if err != nil {
				log.Printf("[batch.loop] species %+v\n", err)
				break
			}
		}
		if err = tx.Commit(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) stageGalaxy(gs *GalaxySetup, root string, ctx context.Context) error {
	log.Printf("[stage] galaxy %d %q: reading galaxy file\n", gs.Id, gs.Name)
	galaxy, err := dat32.ReadGalaxy(filepath.Join(root, "galaxy.dat"), gs.Endian)
	if err != nil {
		return err
	}
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	b := &pgx.Batch{}
	b.Queue("delete from stage_galaxies where galaxy_no = $1", gs.Id)
	b.Queue("insert into stage_galaxies (galaxy_no, name, radius, num_species, turn_number) values ($1, $2, $3, $4, $5)",
		gs.Id, gs.Name, galaxy.Radius, galaxy.NumSpecies, galaxy.TurnNumber)
	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[stage] batch.galaxies: %+v\n", err)
			}
			break
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (db *DB) stageSpecies(gs *GalaxySetup, root string, no int, ctx context.Context) error {
	log.Printf("[stage] galaxy %d %q: reading species %2d file\n", gs.Id, gs.Name, no)
	sp, err := dat32.ReadSpecies(filepath.Join(root, fmt.Sprintf("sp%02d.dat", no)), no, gs.Endian)
	if err != nil {
		return err
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	b := &pgx.Batch{}
	b.Queue("delete from stage_species where galaxy_no = $1 and species_no = $2", gs.Id, no)
	homeworldNo := 0
	if len(sp.NamplaBase) > 0 {
		homeworldNo = sp.NamplaBase[0].PlanetIndex
	}
	b.Queue("insert into stage_species (galaxy_no, species_no, auto_orders, display_name, econ_units, govt_name, govt_type, name, num_namplas, homeworld_no) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		gs.Id, sp.Id, sp.AutoOrders, "SP "+sp.Name, sp.EconUnits, sp.GovtName, sp.GovtType, strings.ToUpper(sp.Name), sp.NumNamplas, homeworldNo)
	b.Queue("delete from stage_species_contact where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for _, alienNo := range sp.Contact {
		b.Queue("insert into stage_species_contact (galaxy_no, species_no, alien_no) values ($1, $2, $3)",
			gs.Id, sp.Id, alienNo)
	}
	b.Queue("delete from stage_species_fleet where galaxy_no = $1 and species_no = $2", gs.Id, no)
	b.Queue("insert into stage_species_fleet (galaxy_no, species_no, fleet_cost, fleet_pct_cost, num_ships) values ($1, $2, $3, $4, $5)",
		gs.Id, sp.Id, sp.FleetCost, sp.FleetPercentCost, sp.NumShips)
	b.Queue("delete from stage_species_homeworld where galaxy_no = $1 and species_no = $2", gs.Id, no)
	b.Queue("insert into stage_species_homeworld (galaxy_no, species_no, homeworld_no, hp_original_base, neutral_gas_code_1, neutral_gas_code_2, neutral_gas_code_3, neutral_gas_code_4, neutral_gas_code_5, neutral_gas_code_6, poison_gas_code_1, poison_gas_code_2, poison_gas_code_3, poison_gas_code_4, poison_gas_code_5, poison_gas_code_6, required_gas_code, required_gas_min, required_gas_max) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)",
		gs.Id, sp.Id, homeworldNo, sp.HPOriginalBase,
		gs.GasCode(sp.NeutralGas[0]), gs.GasCode(sp.NeutralGas[1]), gs.GasCode(sp.NeutralGas[2]), gs.GasCode(sp.NeutralGas[3]), gs.GasCode(sp.NeutralGas[4]), gs.GasCode(sp.NeutralGas[5]),
		gs.GasCode(sp.PoisonGas[0]), gs.GasCode(sp.PoisonGas[1]), gs.GasCode(sp.PoisonGas[2]), gs.GasCode(sp.PoisonGas[3]), gs.GasCode(sp.PoisonGas[4]), gs.GasCode(sp.PoisonGas[5]),
		gs.GasCode(sp.RequiredGas), sp.RequiredGasMin, sp.RequiredGasMax)
	b.Queue("delete from stage_species_relations where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for _, alienNo := range sp.Ally {
		b.Queue("insert into stage_species_relations (galaxy_no, species_no, alien_no, relation) values ($1, $2, $3, 'ALLY')",
			gs.Id, sp.Id, alienNo)
	}
	for _, alienNo := range sp.Enemy {
		b.Queue("insert into stage_species_relations (galaxy_no, species_no, alien_no, relation) values ($1, $2, $3, 'ENEMY')",
			gs.Id, sp.Id, alienNo)
	}
	b.Queue("delete from stage_species_tech where galaxy_no = $1 and species_no = $2", gs.Id, no)
	b.Queue("insert into stage_species_tech (galaxy_no, species_no, mi, mi_init, mi_knowledge, mi_xp, ma, ma_init, ma_knowledge, ma_xp, ml, ml_init, ml_knowledge, ml_xp, gv, gv_init, gv_knowledge, gv_xp, ls, ls_init, ls_knowledge, ls_xp, bi, bi_init, bi_knowledge, bi_xp) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)",
		gs.Id, sp.Id,
		sp.TechLevel[0], sp.InitTechLevel[0], sp.TechKnowledge[0], sp.TechEps[0],
		sp.TechLevel[1], sp.InitTechLevel[1], sp.TechKnowledge[1], sp.TechEps[1],
		sp.TechLevel[2], sp.InitTechLevel[2], sp.TechKnowledge[2], sp.TechEps[2],
		sp.TechLevel[3], sp.InitTechLevel[3], sp.TechKnowledge[3], sp.TechEps[3],
		sp.TechLevel[4], sp.InitTechLevel[4], sp.TechKnowledge[4], sp.TechEps[4],
		sp.TechLevel[5], sp.InitTechLevel[5], sp.TechKnowledge[5], sp.TechEps[5])

	b.Queue("delete from stage_namplas where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for npNo, np := range sp.NamplaBase {
		if strings.ToUpper(np.Name) == "UNUSED" {
			continue
		}
		b.Queue("insert into stage_namplas (galaxy_no, species_no, nampla_no, name, aus_auto, aus_needed, aus_to_install, base_ma, base_mi, ius_auto, ius_needed, ius_to_install, message, planet_no, pop_units, shipyards, siege_effect, special, use_on_ambush) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)",
			gs.Id, sp.Id, npNo,
			np.Name,
			np.AutoAUs, np.AUsNeeded, np.AUsToInstall,
			np.MaBase, np.MiBase,
			np.AutoIUs, np.IUsNeeded, np.IUsToInstall,
			np.Message, np.PlanetIndex, np.PopUnits, np.Shipyards, np.SiegeEff, np.Special, np.UseOnAmbush)
	}
	b.Queue("delete from stage_namplas_inventory where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for npNo, np := range sp.NamplaBase {
		if strings.ToUpper(np.Name) == "UNUSED" {
			continue
		}
		for item, qty := range np.ItemQuantity {
			if qty > 0 {
				code, _ := gs.ItemCode(item)
				b.Queue("insert into stage_namplas_inventory (galaxy_no, species_no, nampla_no, item, quantity) values ($1, $2, $3, $4, $5)",
					gs.Id, sp.Id, npNo, code, qty)
			}
		}
	}
	b.Queue("delete from stage_namplas_status where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for npNo, np := range sp.NamplaBase {
		if strings.ToUpper(np.Name) == "UNUSED" {
			continue
		}
		isColony := (np.Status & 2) != 0
		isDisbandedColony := (np.Status & 64) != 0
		isHidden := np.Hidden
		isHiding := np.Hiding
		isHomePlanet := (np.Status & 1) != 0
		isMiningColony := (np.Status & 16) != 0
		isPopulated := (np.Status & 8) != 0
		isResortColony := (np.Status & 32) != 0
		b.Queue("insert into stage_namplas_status (galaxy_no, species_no, nampla_no, colony, disbanded_colony, hidden, hiding, home_planet, mining_colony, populated, resort_colony) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			gs.Id, sp.Id, npNo,
			isColony, isDisbandedColony, isHidden, isHiding, isHomePlanet, isMiningColony, isPopulated, isResortColony)
	}
	b.Queue("delete from stage_ships where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		var code string
		transport := false
		switch ship.Class {
		case 0:
			code = "PB"
		case 1:
			code = "CT"
		case 2:
			code = "ES"
		case 3:
			code = "DD"
		case 4:
			code = "FG"
		case 5:
			code = "CL"
		case 6:
			code = "CS"
		case 7:
			code = "CA"
		case 8:
			code = "CC"
		case 9:
			code = "BC"
		case 10:
			code = "BS"
		case 11:
			code = "DN"
		case 12:
			code = "SD"
		case 13:
			code = "BM"
		case 14:
			code = "BW"
		case 15:
			code = "BR"
		case 16:
			code = "BA"
		case 17: // TR
			code, transport = fmt.Sprintf("TR%d", ship.Tonnage), true
		default:
			code = "??"
		}
		var starbase, sublight bool
		switch ship.Type {
		case 0:
			starbase, sublight = false, false
		case 1:
			starbase, sublight = false, true
		case 2:
			starbase, sublight = true, true
		}
		if starbase || sublight {
			code += "S"
		}
		var class string
		if transport {
			class = "TRANSPORT"
		} else if starbase {
			class = "STARBASE"
		} else {
			class = "WARSHIP"
		}
		var carryingCapacity int
		if starbase {
			carryingCapacity = 10 * ship.Tonnage
		} else if transport {
			carryingCapacity = (10 + (ship.Tonnage / 2)) * ship.Tonnage
		} else {
			carryingCapacity = ship.Tonnage
		}
		b.Queue("insert into stage_ships (galaxy_no, species_no, ship_no, age, carrying_capacity, class, display_name, name, special, sublight, tonnage) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			gs.Id, sp.Id, shipNo,
			ship.Age, carryingCapacity, class, code+" "+ship.Name, strings.ToUpper(ship.Name), ship.Special, sublight, ship.Tonnage*10_000)
	}
	b.Queue("delete from stage_ships_combat where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		unitPower := gs.power(ship.Tonnage)
		defensivePower, offensivePower := unitPower, unitPower
		if ship.Class == 17 { // transports are not designed for combat
			defensivePower, offensivePower = defensivePower/10, offensivePower/10
		} else if ship.Class == 16 { // starbase has no modifiers?
			// no modifiers? that is odd
		} else {
			// add auxiliary gun and shield generator contributions
			auxGunOffset := 23          // index for gun quantity is GU1 - 1
			shieldGeneratorOffset := 15 // index for generator quantity is SG1 - 1
			for i := 1; i <= 9; i++ {
				if ship.ItemQuantity[i+auxGunOffset] > 0 {
					// each level of GU will have increasing power
					offensivePower += ship.ItemQuantity[i] * gs.power(i*5)
				}
				if ship.ItemQuantity[i+shieldGeneratorOffset] > 0 {
					// each level of SG will have increasing power
					defensivePower += ship.ItemQuantity[i] * gs.power(i*5)
				}
			}
		}
		// adjust for ship aging
		offensivePower -= (ship.Age * offensivePower) / 50
		defensivePower -= (ship.Age * defensivePower) / 50

		// adjust values for tech levels
		LS, ML := 4, 2
		offensivePower += (offensivePower * sp.TechLevel[ML]) / 50
		defensivePower += (defensivePower * sp.TechLevel[LS]) / 50
		if offensivePower < 0 {
			offensivePower = 0
		}
		if defensivePower < 0 {
			defensivePower = 0
		}
		if ship.Age < 0 {
			defensivePower, offensivePower = 0, 0
		}
		b.Queue("insert into stage_ships_combat (galaxy_no, species_no, ship_no, defense, offense) values ($1, $2, $3, $4, $5)",
			gs.Id, sp.Id, shipNo, defensivePower, offensivePower)
	}
	b.Queue("delete from stage_ships_cost where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		var starbase, sublight bool
		switch ship.Type {
		case 0:
			starbase, sublight = false, false
		case 1:
			starbase, sublight = false, true
		case 2:
			starbase, sublight = true, true
		}
		var cost int
		if starbase || !sublight {
			cost = ship.Tonnage * 100
		} else {
			cost = ship.Tonnage * 75
		}
		b.Queue("insert into stage_ships_cost (galaxy_no, species_no, ship_no, cost, remaining_cost) values ($1, $2, $3, $4, $5)",
			gs.Id, sp.Id, shipNo, cost, ship.RemainingCost)
	}
	b.Queue("delete from stage_ships_inventory where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		for item, qty := range ship.ItemQuantity {
			if qty > 0 {
				code, _ := gs.ItemCode(item)
				b.Queue("insert into stage_ships_inventory (galaxy_no, species_no, ship_no, item, quantity) values ($1, $2, $3, $4, $5)",
					gs.Id, sp.Id, shipNo, code, qty)
			}
		}
	}
	b.Queue("delete from stage_ships_location where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		b.Queue("insert into stage_ships_location (galaxy_no, species_no, ship_no, x, y, z, orbit, station) values ($1, $2, $3, $4, $5, $6, $7, $8)",
			gs.Id, sp.Id, shipNo,
			ship.X, ship.Y, ship.Z, ship.PN, "?")
	}
	b.Queue("delete from stage_ships_orders where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		b.Queue("insert into stage_ships_orders (galaxy_no, species_no, ship_no, dest_x, dest_y, dest_z) values ($1, $2, $3, $4, $5, $6)",
			gs.Id, sp.Id, shipNo, ship.DestX, ship.DestY, ship.DestZ)
	}
	b.Queue("delete from stage_ships_status where galaxy_no = $1 and species_no = $2", gs.Id, no)
	for shipNo, ship := range sp.ShipBase {
		if strings.ToUpper(ship.Name) == "UNUSED" {
			continue
		}
		b.Queue("insert into stage_ships_status (galaxy_no, species_no, ship_no, arrived_via_wormhole, just_jumped, remaining_cost, status) values ($1, $2, $3, $4, $5, $6, $7)",
			gs.Id, sp.Id, shipNo, ship.ArrivedViaWormhole, ship.JustJumped, ship.RemainingCost, ship.Status)
	}

	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[stage] batch.species: %+v\n", err)
			}
			break
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (db *DB) stageSystems(gs *GalaxySetup, root string, ctx context.Context) error {
	log.Printf("[stage] galaxy %d %q: reading stars file\n", gs.Id, gs.Name)
	stars, err := dat32.ReadStars(filepath.Join(root, "stars.dat"), gs.Endian)
	if err != nil {
		return err
	}
	log.Printf("[stage] galaxy %d %q: reading planets file\n", gs.Id, gs.Name)
	planets, err := dat32.ReadPlanets(filepath.Join(root, "planets.dat"), gs.Endian)
	if err != nil {
		return err
	}
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	b := &pgx.Batch{}
	b.Queue("delete from stage_stars where galaxy_no = $1", gs.Id)
	b.Queue("delete from stage_planets where galaxy_no = $1", gs.Id)
	for systemNo, system := range stars.Stars {
		b.Queue("insert into stage_stars (galaxy_no, system_no, x, y, z, num_planets) values ($1, $2, $3, $4, $5, $6)",
			gs.Id, systemNo, system.X, system.Y, system.Z, system.NumPlanets)
		for orbit, planet := range planets.Planets[system.PlanetIndex : system.PlanetIndex+system.NumPlanets] {
			planetNo := system.PlanetIndex + orbit
			b.Queue("insert into stage_planets (galaxy_no, planet_no, diameter, econ_efficiency, gas_code_1, gas_code_2, gas_code_3, gas_code_4, gas_percent_1, gas_percent_2, gas_percent_3, gas_percent_4, gravity, message, mining_difficulty, mining_difficulty_increase, pressure_class, special, temperature_class, system_no, orbit) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)",
				gs.Id, planetNo, planet.Diameter, planet.EconEfficiency,
				gs.GasCode(planet.Gas[0]), gs.GasCode(planet.Gas[1]), gs.GasCode(planet.Gas[2]), gs.GasCode(planet.Gas[3]),
				planet.GasPercent[0], planet.GasPercent[1], planet.GasPercent[2], planet.GasPercent[3],
				planet.Gravity, planet.Message, planet.MiningDifficulty, planet.MDIncrease, planet.PressureClass, planet.Special, planet.TemperatureClass,
				systemNo, orbit+1)
		}
	}
	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[stage] batch.stars: %+v\n", err)
			}
			break
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
