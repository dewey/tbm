package cmd

import (
	"context"
	"github.com/dewey/tbm/proc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all services",
	Long: `Start all services that are enabled in the configuration file. Only services with a valid configuration will
be started`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := proc.NotifyCh()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Read configuration file
		hd, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFilePath := path.Join(hd, ".tbm.yaml")
		b, err := os.ReadFile(configFilePath)
		if err != nil {
			return err
		}

		var configuration proc.Configuration
		if err := yaml.Unmarshal(b, &configuration); err != nil {
			return err
		}

		svc := proc.NewServicesService(configuration)
		err = svc.ReadProcfile(configuration)
		if err != nil {
			return err
		}
		return svc.StartProcs(c, true, true)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	startCmd.PersistentFlags().String("config", "~/tbm.yaml", "Location of the configuration file.")
	startCmd.PersistentFlags().String("exit-on-stop", "true", "Exit tbm if all services stop")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
