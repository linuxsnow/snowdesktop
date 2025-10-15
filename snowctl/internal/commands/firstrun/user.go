package firstrun

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
)

var brew = false
var brew_description = "Install Homebrew"
var flatpaks = false
var flatpaks_description = "Install Basic Gnome Flatpaks"
var shellexts = false
var shellexts_description = "Install Gnome Shell Extensions"
var gdm = false
var gdm_description = "Enable Gnome Display Manager on startup"
var rebootrequired = true

func Configure() error {
	var err error
	err = enableTask(&brew, brew_description)
	if err != nil {
		return err
	}
	err = enableTask(&flatpaks, flatpaks_description)
	if err != nil {
		return err
	}
	err = enableTask(&shellexts, shellexts_description)
	if err != nil {
		return err
	}
	err = enableTask(&gdm, gdm_description)
	if err != nil {
		return err
	}

	fmt.Println("\n\nConfiguration Summary:")
	if brew {
		fmt.Println("  -", brew_description)
	}
	if flatpaks {
		fmt.Println("  -", flatpaks_description)
	}
	if shellexts {
		fmt.Println("  -", shellexts_description)
	}
	if gdm {
		fmt.Println("  -", gdm_description)
	}
	fmt.Println()

	if !confirm("Apply these changes now?") {
		fmt.Println("No changes applied.")
		return nil
	}
	if brew {
		err = installBrew()
		if err != nil {
			return err
		}
	}
	if flatpaks {
		err = installFlatpaks()
		if err != nil {
			return err
		}
	}
	if gdm {
		err = enableGdm()
		if err != nil {
			return err
		}
		rebootrequired = true
	}

	if rebootrequired {
		fmt.Println("\nA reboot is required to apply all changes.")
		if confirm("Reboot now?") {
			fmt.Println("Restarting SNOW...")
			cmd := exec.Command("systemctl", "reboot")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			err = cmd.Run()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Please remember to reboot later to apply all changes.")
		}
	}
	return nil
}
func enableTask(defaultValue *bool, desc string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Value(defaultValue).
				Title(desc + "?"),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	return nil
}
func confirm(message string) bool {
	var happy bool

	confirm := huh.NewConfirm().
		Title("Are you sure? ").
		Description(message).
		Affirmative("Yes!").
		Negative("No.").
		Value(&happy)

	huh.NewForm(huh.NewGroup(confirm)).Run()
	return happy
}

func installBrew() error {
	// todo: check if already installed
	cmd := exec.Command(
		"/usr/bin/bash",
		"-c",
		"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func installFlatpaks() error {
	// todo: check if already installed
	// or maybe just run the script that is already in place
	// since it already checks for existing installs
	cmd := exec.Command(
		"sudo",
		"systemctl",
		"enable",
		"gdm.service",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func enableGdm() error {
	cmd := exec.Command(
		"sudo",
		"/usr/bin/bash",
		"-c",
		"/usr/local/bin/_enable-gdm.sh",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
