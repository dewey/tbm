package cmd

import (
	"context"
	"errors"
	"github.com/dewey/tbm/config"
	"github.com/dewey/tbm/proc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all enabled services",
	Long: `Start all services that are enabled in the configuration file. Only services with a valid configuration will
be started`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := proc.NotifyCh()
		ctx := context.Background()
		_, cancel := context.WithCancel(ctx)
		defer cancel()

		// If user provided a custom config file location, we read it from there. Otherwise, we are using the default
		// location in the user's home directory.
		configFlag := cmd.Flag("config")
		var configFilePath string
		if configFlag.Value.String() != configFlag.DefValue {
			// Replace tilde in user provided string, otherwise we can't resolve it
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configFilePath = strings.Replace(configFlag.Value.String(), "~", hd, -1)
		} else {
			hd, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configFilePath = path.Join(hd, ".tbm.yaml")
		}
		// Check if configuration file already exists, otherwise we direct the user to use `tbm init`
		if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
			return errors.New("configuration file doesn't exist. Use `tbm init` to create one or use --config to pass a path")
		}
		b, err := os.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		var configuration config.Configuration
		if err := yaml.Unmarshal(b, &configuration); err != nil {
			return err
		}

		if !configuration.Valid() {
			return errors.New("invalid configuration file, make sure ports are unique across services")
		}

		svc := proc.NewServicesService(configuration)
		err = svc.ReadProcfile(configuration)
		if err != nil {
			return err
		}

		exitOnError := true
		exitOnErrorVal, err := cmd.PersistentFlags().GetBool("exit-on-error")
		if err == nil {
			exitOnError = exitOnErrorVal
		} else {
			return errors.New("couldn't parse exit-on-error flag")
		}

		exitOnStop := true
		exitOnStopVal, err := cmd.PersistentFlags().GetBool("exit-on-stop")
		if err == nil {
			exitOnStop = exitOnStopVal
		} else {
			return errors.New("couldn't parse exit-on-stop flag")
		}

		return svc.StartProcs(c, exitOnError, exitOnStop)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	var configFilePath string
	hd, err := os.UserHomeDir()
	if err == nil {
		configFilePath = path.Join(hd, ".tbm.yaml")
	} else {
		configFilePath = "~/.tbm.yaml"
	}

	startCmd.PersistentFlags().String("config", configFilePath, "Location of the configuration file.")
	startCmd.PersistentFlags().Bool("exit-on-stop", true, "Exit tbm if all services stop")
	startCmd.PersistentFlags().Bool("exit-on-error", true, "Exit tbm if one of the services encounters an error")
}
