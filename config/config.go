/*****************************************************************************
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
****************************************************************************/

package config

import (
	"flag"
	ff "github.com/peterbourgon/ff/v3"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Debug bool
	Data  struct {
		JDB    string // path to json data file
		Log    string // path to create log files
		Orders string // path to orders files
	}
	Server struct {
		Host    string
		Port    int
		PidFile string // if set, save the PID to this file
		App     struct {
			Cluster  string // path to cluster data file
			Data     string // path to v1 data files
			Players  string // path to player data file
			Sessions string // path to session data file
			Site     string // path to site data file
		}
	}
	Log struct {
		Flags   int // use as log.SetFlags(cfg.Log.Flags)
		Verbose bool
	}
}

// DefaultConfig returns a default configuration.
// These are the values without loading the environment, configuration file, or command line.
func DefaultConfig() *Config {
	var cfg Config
	cfg.Data.JDB = "D:\\FarHorizons\\testdata\\t19"
	cfg.Data.Log = "D:\\FarHorizons\\testdata\\t19"
	cfg.Data.Orders = "D:\\FarHorizons\\testdata\\t19"
	cfg.Log.Flags = log.Ldate | log.Ltime | log.LUTC // force logs to be UTC
	cfg.Log.Verbose = true
	cfg.Server.Port = 8080
	cfg.Server.App.Data = "D:\\FarHorizons\\testdata\\t19"
	cfg.Server.App.Cluster = filepath.Join(cfg.Server.App.Data, "cluster.json")
	cfg.Server.App.Players = filepath.Join(cfg.Server.App.Data, "players.json")
	cfg.Server.App.Sessions = filepath.Join(cfg.Server.App.Data, "sessions.json")
	cfg.Server.App.Site = filepath.Join(cfg.Server.App.Data, "site.json")
	return &cfg
}

// Load updates the values in a Config in this order:
//   1. It will load a configuration file if one is given on the
//      command line via the `-config` flag. If provided, the file
//      must contain a valid JSON object.
//   2. Environment variables, using the prefix `CONDUIT_RYER_SERVER`
//   3. Command line flags
func (cfg *Config) Load() error {
	fs := flag.NewFlagSet("Server", flag.ExitOnError)
	debug := fs.Bool("debug", cfg.Debug, "log debug information (optional)")
	dataJDB := fs.String("jdb-path", cfg.Data.JDB, "path to read json data")
	dataLog := fs.String("data", cfg.Data.Log, "path to create log files")
	dataOrders := fs.String("orders", cfg.Data.Orders, "path to read orders files")
	logVerbose := fs.Bool("verbose", cfg.Log.Verbose, "log extra information to the console")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("FH"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return err
	}

	cfg.Debug = *debug
	cfg.Data.JDB = filepath.Clean(*dataJDB)
	cfg.Data.Log = filepath.Clean(*dataLog)
	cfg.Data.Orders = filepath.Clean(*dataOrders)
	cfg.Log.Verbose = *logVerbose

	return nil
}
