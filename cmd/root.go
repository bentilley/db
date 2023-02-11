package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		uris, err := config.URIs()
		if err != nil {
			fmt.Printf("could not get uris: %v\n", err)
			return
		}
		list := strings.Join(uris, "\n")

		fzfCmd := exec.Command("fzf", "--height", "100%")
		fzfCmd.Stdin = strings.NewReader(list)
		fzfCmd.Stdout = os.Stdout
		fzfCmd.Stderr = os.Stderr
		if err := fzfCmd.Run(); err != nil {
			fmt.Printf("could not run fzf: %v", err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
