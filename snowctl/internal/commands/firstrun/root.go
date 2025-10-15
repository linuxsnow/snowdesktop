package firstrun

import (
	"fmt"
	"os"
	"os/exec"
)

func FirstBootRoot() error {
	var err error

	err = checkRoot()
	if err != nil {
		return err
	}
	fmt.Println("Starting SNOW first-run configuration...")
	err = firstUser()
	if err != nil {
		return err
	}
	fmt.Println("First user created successfully.")
	// Enable GDM
	err = enableGDM()
	if err != nil {
		return err
	}
	fmt.Println("Graphical Login enabled successfully.")
	err = installBazaar()
	if err != nil {
		return err
	}
	fmt.Println("Bazaar installed successfully.")
	return nil
}

func checkRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run as root")
	}
	return nil
}

func enableGDM() error {
	cmd := exec.Command("systemctl", "enable", "gdm.service")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func firstUser() error {
	// Check if any users exist
	// If not, prompt to create the first user
	fmt.Println("Checking for existing users...")
	fmt.Println("No users found. Please create the first user.")
	cmd := exec.Command(
		"homectl",
		"firstboot",
		"--prompt-new-user",
		"--storage=subvolume",
		"--prompt-shell=no",
		"--prompt-groups=no",
		"--mute-console=yes",
		"--uid=1000",
		"--member-of=adm,systemd-journal,sudo,docker,incus-admin",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()

}

func installBazaar() error {
	// flatpak install --or-update --noninteractive --assumeyes flathub "$flatpak"

	cmd := exec.Command(
		"flatpak",
		"install",
		"--or-update",
		"--noninteractive",
		"--assumeyes",
		"flathub",
		"io.github.kolunmi.Bazaar",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()

}
