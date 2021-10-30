/*
Copyright © 2021 Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string
var homeFolder string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fh",
	Short: "Far Horizons engine",
	Long: `FH is the game engine for Far Horizons. This application creates
new games, executes orders, and generates reports for each player.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("env: %-30s == %q\n", "HOME", homeFolder)
		log.Printf("env: %-30s == %q\n", "FH_CONFIG", viper.ConfigFileUsed())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fh.yaml)")

	// Cobra also supports local flags, which will only run when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	// Find home directory.
	homeFolder, err = homedir.Dir()
	cobra.CheckErr(err)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".fh" (without extension).
		viper.AddConfigPath(homeFolder)
		viper.SetConfigType("json")
		viper.SetConfigName(".fh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
		log.Printf("viper: using config file: %q\n", viper.ConfigFileUsed())
	}

	err = viper.WriteConfigAs("D:\\GoLand\\fhcms\\testdata\\viper.json")
	cobra.CheckErr(err)
}
