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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Server http.Handler

// Router returns a router that will serve REST-ish data for the EmberJS front end.
func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/turn", apiGetTurn)
	r.Get("/widgets", apiGetWidgets)
	r.Post("/widgets", apiCreateWidget)
	r.Post("/widgets/{slug}", apiUpdateWidget)
	r.Post("/widgets/{slug}/parts", apiCreateWidgetPart)
	r.Post("/widgets/{slug}/parts/{id:[0-9]+}/update", apiUpdateWidgetPart)
	r.Post("/widgets/{slug}/parts/{id:[0-9]+}/delete", apiDeleteWidgetPart)

	// admin protected routes
	r.Group(func(r chi.Router) {
		r.Use(adminOnly)
		r.Get("/stats", apiGetStats)
	})

	var s Server
	s = r
	return s
}

// adminOnly is middleware that forwards only if the request has admin rights.
func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		isAdmin, ok := claims["admin"].(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func apiGetStats(w http.ResponseWriter, r *http.Request) {
	type species struct {
		Species             int    `json:"species"`
		Name                string `json:"name"`
		MI                  int    `json:"mi"`
		MA                  int    `json:"ma"`
		ML                  int    `json:"ml"`
		GV                  int    `json:"gv"`
		LS                  int    `json:"ls"`
		BI                  int    `json:"bi"`
		TotalProduction     int    `json:"total_production"`
		NumberOfPlanets     int    `json:"number_of_planets"`
		NumberOfShips       int    `json:"number_of_ships"`
		NumberOfShipYards   int    `json:"number_of_ship_yards"`
		OffensivePower      int    `json:"offensive_power"`
		DefensivePower      int    `json:"defensive_power"`
		BankedEconomicUnits int    `json:"banked_economic_units"`
	}
	type tech struct {
		Min     int `json:"min"`
		Average int `json:"average"`
		Max     int `json:"max"`
	}
	type tonnage struct {
		MinQty     string `json:"min_qty"`
		AvgQty     string `json:"avg_qty"`
		MaxQty     string `json:"max_qty"`
		AvgSize    string `json:"avg_size"`
		AvgTonnage string `json:"avg_tonnage"`
	}

	type response struct {
		Species    []species `json:"species"`
		MI         tech      `json:"mi"`
		MA         tech      `json:"ma"`
		ML         tech      `json:"ml"`
		GV         tech      `json:"gv"`
		LS         tech      `json:"ls"`
		BI         tech      `json:"bi"`
		Starbases  tonnage   `json:"starbases"`
		Transports tonnage   `json:"transports"`
		Warships   tonnage   `json:"warships"`
	}

	ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
	if err != nil {
		log.Printf("error: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var rsp response

	game := ds.Stats()
	for _, sp := range game.Stats {
		rsp.Species = append(rsp.Species, species{
			Species:             sp.No,
			Name:                sp.Name,
			MI:                  sp.MI,
			MA:                  sp.MA,
			ML:                  sp.ML,
			GV:                  sp.GV,
			LS:                  sp.LS,
			BI:                  sp.BI,
			TotalProduction:     int(sp.Production),
			NumberOfPlanets:     int(sp.PopulatedPlanets),
			NumberOfShips:       int(sp.Ships),
			NumberOfShipYards:   int(sp.Shipyards),
			OffensivePower:      int(sp.OffensivePower),
			DefensivePower:      int(sp.DefensivePower),
			BankedEconomicUnits: int(sp.BankedEconUnits),
		})
	}

	if game.TotalSpecies > 0 {
		// m attempts to compensate for rounding integer values down
		m := game.TotalSpecies / 2

		rsp.MI = tech{Min: game.MI.Min, Average: (game.MI.Total + m) / game.TotalSpecies, Max: game.MI.Max}
		rsp.MA = tech{Min: game.MA.Min, Average: (game.MA.Total + m) / game.TotalSpecies, Max: game.MA.Max}
		rsp.ML = tech{Min: game.ML.Min, Average: (game.ML.Total + m) / game.TotalSpecies, Max: game.ML.Max}
		rsp.GV = tech{Min: game.GV.Min, Average: (game.GV.Total + m) / game.TotalSpecies, Max: game.GV.Max}
		rsp.LS = tech{Min: game.LS.Min, Average: (game.LS.Total + m) / game.TotalSpecies, Max: game.LS.Max}
		rsp.BI = tech{Min: game.BI.Min, Average: (game.BI.Total + m) / game.TotalSpecies, Max: game.BI.Max}

		rsp.Starbases = tonnage{
			MinQty:     fmt.Sprintf("%.0f", game.Starbases.Min),
			AvgQty:     fmt.Sprintf("%.1f", game.Starbases.Average),
			MaxQty:     fmt.Sprintf("%.0f", game.Starbases.Max),
			AvgSize:    "0",
			AvgTonnage: commas(int(game.StarbaseTonnage.Average)),
		}
		if game.Starbases.Total > 0 {
			rsp.Starbases.AvgSize = commas(int(game.StarbaseTonnage.Total / game.Starbases.Total))
		}
		rsp.Transports = tonnage{
			MinQty:     fmt.Sprintf("%.0f", game.Transports.Min),
			AvgQty:     fmt.Sprintf("%.1f", game.Transports.Average),
			MaxQty:     fmt.Sprintf("%.0f", game.Transports.Max),
			AvgSize:    "0",
			AvgTonnage: commas(int(game.TransportTonnage.Average)),
		}
		if game.Transports.Total > 0 {
			rsp.Transports.AvgSize = commas(int(game.TransportTonnage.Total / game.Transports.Total))
		}
		rsp.Warships = tonnage{
			MinQty:     fmt.Sprintf("%.0f", game.Warships.Min),
			AvgQty:     fmt.Sprintf("%.1f", game.Warships.Average),
			MaxQty:     fmt.Sprintf("%.0f", game.Warships.Max),
			AvgSize:    "0",
			AvgTonnage: commas(int(game.WarshipTonnage.Average)),
		}
		if game.Warships.Total > 0 {
			rsp.Warships.AvgSize = commas(int(game.WarshipTonnage.Total / game.Warships.Total))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func apiGetTurn(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Turn int `json:"turn"`
	}

	ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
	if err != nil {
		log.Printf("error: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rsp := response{Turn: ds.Turn}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	_, _ = fmt.Fprintf(w, "apiGetWidgets: claims %v\n", claims["species"])
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	_, _ = fmt.Fprintf(w, "apiCreateWidget: claims %v\n", claims["species"])
}

func apiUpdateWidget(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "apiUpdateWidget: claims %v: slug %s\n", claims["species"], slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "apiCreateWidgetPart: claims %v: slug %s\n", claims["species"], slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, _ = fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, _ = fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

/* This routine is intended to take a long argument and return a pointer
 * to a string that has embedded commas to make the string more readable. */
func commas(value int) string {
	if value == 0 {
		return "0"
	} else if value < 0 {
		return "-" + commas(-1*value)
	}
	s := fmt.Sprintf("%024d", value)
	return strings.TrimLeft(s[0:3]+","+s[3:6]+","+s[6:9]+","+s[9:12]+","+s[12:15]+","+s[15:18]+","+s[18:21]+","+s[21:], "0,")
}

func loader(path string, bigEndian bool) (*cluster.Store, error) {
	if path = filepath.Clean(path); path == "." {
		if cwd, err := os.Getwd(); err != nil {
			return nil, err
		} else if path = filepath.Clean(cwd); path == "." {
			return nil, fmt.Errorf("unable to determine path to data files")
		}
	}
	return cluster.FromDat32(path, bigEndian)
}
