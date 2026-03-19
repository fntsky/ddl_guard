package ddlcmd

import (
	"github.com/fntsky/ddl_guard/internal/cli"
	"github.com/spf13/cobra"
)

var (
	dataDir string
	rootCmd = &cobra.Command{
		Use:   "ddl_guard",
		Short: "DDL Guard CLI",
		Long:  "DDL Guard CLI",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "run server",
		Long:  "run server",
		Run: func(_ *cobra.Command, _ []string) {
			runApp()
		},
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "init config",
		Long:  "init config",
		Run: func(_ *cobra.Command, _ []string) {
			cli.Install(dataDir)
		},
	}
)

func init() {
	initCmd.Flags().StringVarP(&dataDir, "dir", "d", "./data", "data directory path")
	rootCmd.AddCommand(runCmd, initCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
