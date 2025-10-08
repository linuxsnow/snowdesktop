package install

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

func Install() error {
	// Check /proc/cmdline for "root=tmpfs"
	// If found, proceed with installation
	// If not found, return an error indicating that installation can only be run from the live environment
	err := checkLiveEnvironment()
	if err != nil {
		return err
	}
	// Perform installation steps:
	// find target disks
	fmt.Println("Finding target disks...")
	disks, err := findTargetDisks()
	if err != nil {
		return err
	}
	if len(disks) == 0 {
		return fmt.Errorf("no suitable target disks found")
	}
	// prompt user to select a disk
	disk, err := chooseDisk(disks)
	if err != nil {
		return err
	}
	fmt.Println("Selected disk:", disk)
	// confirm with user
	if !confirm(fmt.Sprintf("All data on %s will be lost. Are you sure you want to continue?", disk)) {
		fmt.Println("Installation cancelled.")
		return fmt.Errorf("installation cancelled by user")
	}
	err = runRepart(disk)
	if err != nil {
		return err
	}
	return nil
}

func checkLiveEnvironment() error {
	// Read /proc/cmdline and check for "root=tmpfs"
	// If found, return nil
	// If not found, return an error
	cmpl, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return fmt.Errorf("failed to read /proc/cmdline: %w", err)
	}
	if strings.Contains(string(cmpl), "root=tmpfs") {
		return nil
	}
	return fmt.Errorf("installation can only be run from the live environment")
}

func findTargetDisks() ([]string, error) {
	// List disks in /sys/block
	// Filter out loop devices, ram devices, and the live environment disk (e.g., /dev/sda)
	// Return a list of candidate disks for installation
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, fmt.Errorf("failed to read /sys/block: %w", err)
	}
	var disks []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") || strings.HasPrefix(name, "zram") || strings.HasPrefix(name, "dm") {
			continue
		}
		// Add additional filtering logic as needed
		// find disks that are mounted
		mounted, err := os.ReadFile("/proc/mounts")
		if err != nil {
			return nil, fmt.Errorf("failed to read /proc/mounts: %w", err)
		}
		if strings.Contains(string(mounted), "/dev/"+name) {
			continue
		}
		// If the disk is not filtered out, add it to the list
		disks = append(disks, "/dev/"+name)
	}
	return disks, nil
}

func chooseDisk(disks []string) (string, error) {
	var disk string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(disks...)...).
				Value(&disk).
				Title("Choose target disk"),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return disk, nil
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

func runRepart(targetDisk string) error {
	cmd := exec.Command(
		"/usr/bin/systemd-repart",
		"--dry-run=no",
		"--empty=force",
		"--defer-partitions=root,swap",
		targetDisk,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
