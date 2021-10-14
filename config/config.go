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
		BigEndian bool   // byte order of v1 binary data files
		Files     string // name of files data file
		Logs      string // path to create log files
		Orders    string // path to orders files
		Path      string // path to v1 binary data files
		Players   string // name of player data file
		Reports   string // path to turn reports
		Sessions  string // name of session data file
		Site      string // name of site data file
		Stats     string // name of stats data file
		Turn      string // name of turn data file
		TurnFiles string // path to turn files (orders and reports)
	}
	Server struct {
		Host string
		Port int
	}
	PIDFile bool // create pid file if set
	Log     struct {
		Flags   int // use as log.SetFlags(cfg.Log.Flags)
		Verbose bool
	}
}

// DefaultConfig returns a default configuration.
// These are the values without loading the environment, configuration file, or command line.
func DefaultConfig() *Config {
	var cfg Config
	root := "D:\\FarHorizons\\testdata"
	// byte order is the order in the data file, not the computer we're running on!
	cfg.Data.BigEndian = false // leave this to document the choice
	cfg.Data.Files = filepath.Join(root, "files.json")
	cfg.Data.Logs = root
	cfg.Data.Orders = root
	cfg.Data.Path = root
	cfg.Data.Players = filepath.Join(root, "players.json")
	cfg.Data.Reports = filepath.Join(root, "reports")
	cfg.Data.Sessions = filepath.Join(root, "sessions.json")
	cfg.Data.Site = filepath.Join(root, "site.json")
	cfg.Data.Stats = filepath.Join(root, "stats.json")
	cfg.Data.Turn = filepath.Join(root, "turn.json")
	cfg.Data.TurnFiles = filepath.Join(root, "reports")
	cfg.Log.Flags = log.Ldate | log.Ltime | log.LUTC // force logs to be UTC
	cfg.Log.Verbose = true
	cfg.Server.Port = 8080
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
	dataFiles := fs.String("files", cfg.Data.Files, "name of files data json file")
	dataLogs := fs.String("logs", cfg.Data.Logs, "path to log files")
	dataOrders := fs.String("orders", cfg.Data.Orders, "path to orders files")
	dataPath := fs.String("data", cfg.Data.Path, "path to v1 binary data files")
	dataPIDFile := fs.Bool("pid-file", false, "create pid file")
	dataPlayers := fs.String("players", cfg.Data.Players, "name of players data json file")
	dataReports := fs.String("reports", cfg.Data.Reports, "path to turn report files")
	dataSessions := fs.String("sessions", cfg.Data.Sessions, "name of sessions data json file")
	dataSite := fs.String("site", cfg.Data.Site, "name of sites data json file")
	dataStats := fs.String("stats", cfg.Data.Stats, "name of stats data json file")
	dataTurn := fs.String("turn", cfg.Data.Turn, "name of turn data json file")
	dataTurnFiles := fs.String("turn-files", cfg.Data.TurnFiles, "path to turn orders and report files")
	logVerbose := fs.Bool("verbose", cfg.Log.Verbose, "log extra information to the console")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("FH"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return err
	}

	cfg.Debug = *debug
	cfg.Data.Files = filepath.Clean(*dataFiles)
	cfg.Data.Logs = filepath.Clean(*dataLogs)
	cfg.Data.Orders = filepath.Clean(*dataOrders)
	cfg.Data.Path = filepath.Clean(*dataPath)
	cfg.Data.Players = filepath.Clean(*dataPlayers)
	cfg.Data.Reports = filepath.Clean(*dataReports)
	cfg.Data.Sessions = filepath.Clean(*dataSessions)
	cfg.Data.Site = filepath.Clean(*dataSite)
	cfg.Data.Stats = filepath.Clean(*dataStats)
	cfg.Data.Turn = filepath.Clean(*dataTurn)
	cfg.Data.TurnFiles = filepath.Clean(*dataTurnFiles)
	cfg.Log.Verbose = *logVerbose
	cfg.PIDFile = *dataPIDFile

	log.Printf("config: %-30s == %q\n", "files", cfg.Data.Files)
	log.Printf("config: %-30s == %q\n", "logs", cfg.Data.Logs)
	log.Printf("config: %-30s == %q\n", "orders", cfg.Data.Orders)
	log.Printf("config: %-30s == %q\n", "path", cfg.Data.Path)
	log.Printf("config: %-30s == %v\n", "pid-file", cfg.PIDFile)
	log.Printf("config: %-30s == %q\n", "players", cfg.Data.Players)
	log.Printf("config: %-30s == %q\n", "reports", cfg.Data.Reports)
	log.Printf("config: %-30s == %q\n", "sessions", cfg.Data.Sessions)
	log.Printf("config: %-30s == %q\n", "site", cfg.Data.Site)
	log.Printf("config: %-30s == %q\n", "stats", cfg.Data.Stats)
	log.Printf("config: %-30s == %q\n", "turn", cfg.Data.Turn)
	log.Printf("config: %-30s == %q\n", "turn-files", cfg.Data.TurnFiles)
	log.Printf("config: %-30s == %v\n", "verbose", cfg.Log.Verbose)

	return nil
}
