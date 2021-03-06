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
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

var cfgFile string
var homeFolder string
var testFlag bool
var verboseFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fh",
	Short: "Far Horizons engine",
	Long: `FH is the game engine for Far Horizons. This application creates
new games, executes orders, and generates reports for each player.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// bind viper and cobra here since this hook runs early and always
		return bindConfig(cmd)
	},
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fh.yaml)")
	rootCmd.PersistentFlags().BoolVar(&testFlag, "test", false, "test mode")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "verbose mode")

	// Cobra also supports local flags, which will only run when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// bindConfig reads in config file and ENV variables if set.
// logic for binding viper and cobra taken from
// https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
func bindConfig(cmd *cobra.Command) error {
	var err error

	// Find home directory.
	homeFolder, err = homedir.Dir()
	if err != nil {
		return err
	}

	if cfgFile == "" { // use default location of ~/.fh
		viper.AddConfigPath(homeFolder)
		viper.SetConfigType("json")
		viper.SetConfigName(".fh")
	} else { // Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	// Try to read the config file. Ignore file-not-found errors.
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		log.Printf("viper: using config file: %q\n", viper.ConfigFileUsed())
		if err = viper.WriteConfigAs(filepath.Join(viper.Get("files.path").(string), "viper.json")); err != nil {
			return err
		}
	}

	// read in environment variables that match
	viper.SetEnvPrefix("FH")
	viper.AutomaticEnv()

	// bind the current command's flags to viper
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", "FH", envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}
