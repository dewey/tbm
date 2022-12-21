package cmd

import (
	"errors"
	"github.com/dewey/tbm/proc"
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
		cfg := make(proc.Configuration)
		cfg["ping"] = proc.Service{
			Command:     "ping google.com",
			Environment: "prod",
			Enable:      true,
		}
		b, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		// Create config file if it doesn't exist yet
		if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(configFilePath)
			if err != nil {
				return err
			}
			_, err = f.Write(b)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully initialized. An example config file created in: %s. Use `tbm start` to give it a try based on the example config file.\n", configFilePath)
			return nil
		} else {
			cmd.Printf("Config file already exists in %s. Manually delete it to recreate the example config file.\n", configFilePath)
			return nil
		}
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
