package main

import ( 
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
    "github.com/farmlab/devsim"
)

var (

	// alias for show
	helloCmd = &cobra.Command{
		Use:   "hello",
		Short: "Display a file from the hoarder storage",
		Long:  ``,

		Run: hello,
	}
)

// show utilizes the api to show data associated to key
func hello(ccmd *cobra.Command, args []string) {
	devsim.Hello()
}
