package cmd

import (
	"github.com/dewey/tbm/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize tbm by creating an example configuration file",
	Long:  `Initializing tbm by creating an example configuration file in the home directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Configuration

		// Fetch remote configuration file if provided by user, otherwise create default example configuration
		configURL := cmd.Flag("config-url")
		if configURL.Value.String() != "" {
			resp, err := http.Get(configURL.Value.String())
			if err != nil {
				return err
			}
			if err := yaml.NewDecoder(resp.Body).Decode(&cfg); err != nil {
				return err
			}
			b, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configFilePath := path.Join(hd, ".tbm.yaml")
			existed, err := config.Create(configFilePath, b)
			if err != nil {
				return err
			}
			if existed {
				cmd.Printf("Config file already exists in %s. Manually delete it to recreate the example config file.\n", configFilePath)
			} else {
				cmd.Printf("Successfully initialized based on remote configuration file downloaded to: %s. Use `tbm start` to give it a try based on the config file.\n", configFilePath)
			}
		} else {
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configFilePath := path.Join(hd, ".tbm.yaml")
			services := make(map[string]config.Service)
			services["ping-one"] = config.Service{
				Command:     "ping google.com",
				Environment: "prod",
				Enable:      true,
			}
			services["ping-two"] = config.Service{
				Command:     "ping duckduckgo.com",
				Environment: "stage",
				Enable:      true,
			}
			cfg = config.Configuration{Services: services}
			b, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			existed, err := config.Create(configFilePath, b)
			if err != nil {
				return err
			}
			if existed {
				cmd.Printf("Config file already exists in %s. Manually delete it to recreate the example config file.\n", configFilePath)
			} else {
				cmd.Printf("Successfully initialized. An example config file created in: %s. Use `tbm start` to give it a try based on the example config file.\n", configFilePath)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("config-url", "", "Provide a URL hosting a configuration file. This will be stored at the default configuration location.")
}
