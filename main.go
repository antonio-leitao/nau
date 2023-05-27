package main

import (
	"fmt"
	root "github.com/antonio-leitao/nau/cmd"
	archive "github.com/antonio-leitao/nau/cmd/archive"
	configure "github.com/antonio-leitao/nau/cmd/configure"
	new "github.com/antonio-leitao/nau/cmd/new"
	open "github.com/antonio-leitao/nau/cmd/open"
	lib "github.com/antonio-leitao/nau/lib"
	"github.com/spf13/cobra"
	"os"
    "log"
)


func main() {
    //load config
	config, err := lib.LoadConfig()
	if err != nil {
		log.Printf("NAU error: %v", err)
		os.Exit(1)
	}
	// Add the --version flag
    app := rootCmd(config, config.Version) 
	if err := app.Execute(); err != nil {
		log.Printf("NAU error: %v", err)
        os.Exit(1)
	}
}

func rootCmd(config lib.Config, version string)*cobra.Command{
	var versionFlag bool
    //add root command
	rootCmd := &cobra.Command{
		Use:   "nau",
		Short: "A CLI application called nau",
		Run: func(cmd *cobra.Command, args []string) {
			if versionFlag {
				fmt.Println("nau version:", version)
				return
			}
			root.Execute(config)
		},
	}

	rootCmd.AddCommand(configCmd(), newCmd(config), openCmd(config), archiveCmd(config))
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print the version number")
    return rootCmd

}
func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [field] [value]",
		Short: "Set or get configuration values",
		Run: func(cmd *cobra.Command, args []string) {
            configure.Execute(args)
        },
    }
	return cmd
}

func newCmd(config lib.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [template]",
		Short: "Create a new project from a template",
		Run: func(cmd *cobra.Command, args []string) {
            if len(args)>0{
                new.Execute(config, args[0])
            } else {
                new.Execute(config, "")
            }
		},
	}

	return cmd
}

func openCmd(config lib.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "open [project]",
		Long: `Open a project.

This command opens the specified project. If no project is provided, it opens the default project.

You can specify the project by providing its name as an argument.`,
		Example: `  nau open            # Open the default project
  nau open myproject  # Open the project named "myproject"`,
		Short: "Open a project",
		Run: func(cmd *cobra.Command, args []string) {
            if len(args) >0{
                open.Execute(config, args[0])
            }else{
                fmt.Println("TODO: pass empty query to open")
            }
		},
	}

	return cmd
}

func archiveCmd(config lib.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive [project]",
		Short: "Archive a project",
		Run: func(cmd *cobra.Command, args []string) {
            if len(args) >0{
                archive.Execute(config, args[0])
            }else{
                fmt.Println("TODO: pass empty query to archive")
            }
		},
	}

	return cmd
}
