package main

import (
	"fmt"
	"strings"

	"context"
	"os"
	"time"

	"github.com/bketelsen/snow/snowctl/internal/commands/firstrun"
	"github.com/bketelsen/snow/snowctl/internal/commands/install"
	"github.com/bketelsen/snow/snowctl/internal/features"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {

	var now bool

	cmd := &cobra.Command{
		Use:   "snowctl [args]",
		Short: "Manage your SNOW system with snowctl",
		Long:  `Manage your SNOW system with snowctl.`,

		RunE: func(c *cobra.Command, _ []string) error {

			c.Println("You ran the root command. Now try --help.")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&now, "now", "n", false, "Apply changes now")

	cmd.AddGroup(&cobra.Group{
		ID:    "features",
		Title: "Manage Features",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "extensions",
		Title: "Manage Extensions",
	})
	installer := &cobra.Command{
		Use:   "install",
		Short: "Install SNOW on this system",
		Run: func(c *cobra.Command, _ []string) {
			if err := install.Install(); err != nil {
				c.Println("Error installing SNOW:", err)
			} else {
				c.Println("SNOW installed successfully.")
			}
		},
	}
	cmd.AddCommand(installer)

	firstboot := &cobra.Command{
		Use:   "firstboot",
		Short: "Create the first user and enable graphical login (requires root)",
		Run: func(c *cobra.Command, _ []string) {
			if err := firstrun.FirstBootRoot(); err != nil {
				c.Println("Error configuring SNOW:", err)
			} else {
				c.Println("SNOW configured successfully. Reboot to continue.")
			}
		},
	}
	cmd.AddCommand(firstboot)
	firstuser := &cobra.Command{
		Use:   "setup",
		Short: "First-run configuration for SNOW user",
		Run: func(c *cobra.Command, _ []string) {
			if err := firstrun.Configure(); err != nil {
				c.Println("Error configuring SNOW:", err)
			} else {
				c.Println("SNOW configured successfully. ")
			}
		},
	}
	cmd.AddCommand(firstuser)
	feature := &cobra.Command{
		Use:     "feature [command] [flags] [args]",
		Aliases: []string{"f"},
		Short:   "Manage features",
		GroupID: "features",
		Run: func(c *cobra.Command, _ []string) {
			other()
		},
	}
	cmd.AddCommand(feature)
	feature.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all features",
		Example: `
snowctl feature list
`,
		RunE: func(c *cobra.Command, _ []string) error {
			cmd.Println("Working...")
			select {
			case <-time.After(time.Second * 5):
				cmd.Println("Done!")
			case <-c.Context().Done():
				return c.Context().Err()
			}
			return nil
		},
	})
	ext := &cobra.Command{
		Use:     "ext [command] [flags] [args]",
		Aliases: []string{"e"},
		Short:   "Manage extensions",
		GroupID: "extensions",
		Run: func(c *cobra.Command, _ []string) {
			other()
		},
	}
	cmd.AddCommand(ext)
	ext.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all extensions",
		Example: `
snowctl ext list
`,
		RunE: func(c *cobra.Command, _ []string) error {
			cmd.Println("Working...")
			select {
			case <-time.After(time.Second * 5):
				cmd.Println("Done!")
			case <-c.Context().Done():
				return c.Context().Err()
			}
			return nil
		},
	})

	// This is where the magic happens.
	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}

func other() {

	ff, err := features.ListFeatures()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listing features:", err)
		os.Exit(1)
	}
	fmt.Println("Result from calling ListFeatures function on org.freedesktop.sysupdate1.Target.ListFeatures interface:")
	for _, feature := range ff {
		desc, err := features.DescribeFeature(feature)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting feature information:", err)
			os.Exit(1)
		}
		fmt.Println("   Name:", desc.Name)
		fmt.Println("   Enabled:", desc.Enabled)
		fmt.Println("   Description:", desc.Description)
		fmt.Println("   Documentation:", desc.Documentation)
		fmt.Println("   Transfers:")
		// if desc.Transfers is empty, print "    (none)"
		if len(desc.Transfers) == 0 {
			fmt.Println("    (none)")
		} else {
			for _, transfer := range desc.Transfers {
				fmt.Println("    -", transfer)
			}
		}

		fmt.Println("   Associated files in /var/lib/extensions:")
		// get a list of files in /var/lib/extensions that end in feature.raw
		files, err := os.ReadDir("/var/lib/extensions")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read extensions directory:", err)
			os.Exit(1)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			suffix := feature + ".raw"
			// check if the file ends with suffix
			if strings.HasSuffix(file.Name(), suffix) {
				fmt.Println("   -", file.Name())
			}
		}

	}
}
