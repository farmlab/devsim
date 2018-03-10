
package main

import (
	"os"
	"path/filepath"
	"strings"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  string //
	daemon  bool   //
	version bool   //

	// devsCmd ...
	DevsCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		// parse the config if one is provided, or use the defaults. Set the backend
		// driver to be used
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// if --config is passed, attempt to parse the config file
			if config != "" {

				// get the filepath
				abs, err := filepath.Abs(config)
				if err != nil {
					//lumber.Error("Error reading filepath: ", err.Error())
				}

				// get the config name
				base := filepath.Base(abs)

				// get the path
				path := filepath.Dir(abs)

				//
				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.AddConfigPath(path)

				// Find and read the config file; Handle errors reading the config file
				if err := viper.ReadInConfig(); err != nil {
					//lumber.Fatal("Failed to read config file: ", err.Error())
					os.Exit(1)
				}
			}
		},

		// either run devs as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(ccmd *cobra.Command, args []string) {

			// if --server is passed start the devs server
			if daemon {
				// do server stuff...
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {
// local flags;
	DevsCmd.Flags().StringVar(&config, "config", "", "./config.yml")
	DevsCmd.Flags().BoolVar(&daemon, "server", false, "Run devs as a server")
	DevsCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	// commands
	DevsCmd.AddCommand(helloCmd)
}
