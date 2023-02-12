package cmd

import (
	"fmt"
	"os"

	"github.com/bentilley/db/pkg/db"
	"github.com/spf13/cobra"
)

var configFile string

func init() {
	rootCmd.PersistentFlags().
		StringVar(&configFile, "config", "$HOME/.config/db/config.yaml", "config file")
}

var rootCmd = &cobra.Command{
	Use:   "db",
	Short: "db is a tool for managing database connection information",
	Long: `A simple tool for managing database connection URIs.
	Print URIs or automatically connect to them with your favourite database cli.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := db.LoadConfig(os.ExpandEnv(configFile))
		if err != nil {
			fmt.Printf("could not load config: %v\n", err)
			return
		}

		index, err := db.FuzzyFind(config.SearchStrings())
		if err != nil {
			fmt.Printf("could not fuzzy find: %v", err)
			return
		}

		uri, err := config.Sessions[index-1].URI()
		if err != nil {
			fmt.Printf("could not compile uri: %v", err)
			return
		}
		fmt.Println(uri)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
