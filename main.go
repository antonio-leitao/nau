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
	"log"
	"os"
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

func rootCmd(config lib.Config, version string) *cobra.Command {
	var versionFlag bool
	//add root command
	rootCmd := &cobra.Command{
		Use:   "nau",
        Short: `|\| /\ |_|: command line project manager`,
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
		Long: `Manage nau's configuration.

If it is the first time using now start by running "nau config" to set all configuration parameters.
You can also set them individually by running "nau config field value" or print current values with
"nau config field". Configurations are stored at "~/.config/naurc"`,
        Example: `  nau config                # Set up all configuration values in NAU
  nau config author           # Print current value of "author"
  nau config author John Doe  # Set author field to "John Doe"`,
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
		Long: `Choose a project template and create a new instance.

Collapse an existing template from the "templates" folder. The user is prompted with necessary
information such as project_name and description. Upon initialization of the template, nau runs
"make init" command on the directory. Use this for adding extra features to the template,
such as enviroment initialization.`,
		Example: `  nau new template            # Start a new project using template "template"
  nau new  # Choose the template before starting`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
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
		Long: `Open a project in your preferred application.

This command opens the specified project. If no project is provided, it opens the default project.
You can specify the project by providing its name as an argument.`,
		Example: `  nau open myproject            # Open the project named "myproject"
  nau open myproj  # Open the project that best matches "myproj"`,
		Short: "Open a project",
		Run: func(cmd *cobra.Command, args []string) {
            if len(args)>0{
                open.Execute(config, args[0])
            }else{
                cmd.Help()
            }
		},
	}

	return cmd
}

func archiveCmd(config lib.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive [project]",
		Short: "Archive a project and move it to the `archive` folder.",
		Long: `Archive a project and move it to the "archive" folder.

The project considered is the best fuzzy match. Use "nau" for more control. Before .tar and .zip
the command "make archive" is run on the directory. Define it in your project's file to enable extra
features, such as deleting git and node dependedncies.`,
		Example: `  nau aarchive myproject            # Archive the project named "myproject"
  nau archive myproj  # Archive the project that best matches "myproj"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				archive.Execute(config, args[0])
			} else {
                cmd.Help()
			}
		},
	}

	return cmd
}
