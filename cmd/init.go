package cmd

import (
	"github.com/dewey/tbm/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize tbm by creating an example configuration file",
	Long:  `Initializing tbm by creating an example configuration file in the home directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hd, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFilePath := path.Join(hd, ".tbm.yaml")
		services := make(map[string]config.Service)
		services["ping"] = config.Service{
			Command:     "ping google.com",
			Environment: "prod",
			Enable:      true,
		}
		cfg := config.Configuration{Services: services}
		b, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		exists, err := config.Create(configFilePath, b)
		if exists {
			cmd.Printf("Config file already exists in %s. Manually delete it to recreate the example config file.\n", configFilePath)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
