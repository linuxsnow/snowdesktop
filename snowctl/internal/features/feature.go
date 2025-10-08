package features

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

type FeatureDescription struct {
	Name          string   `json:"name"`
	Enabled       bool     `json:"enabled"`
	Description   string   `json:"description"`
	Documentation string   `json:"documentation"`
	Transfers     []string `json:"transfers"`
}

func ListFeatures() ([]string, error) {

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var s []string
	obj := conn.Object("org.freedesktop.sysupdate1", "/org/freedesktop/sysupdate1/target/host")
	err = obj.Call("org.freedesktop.sysupdate1.Target.ListFeatures", 0, uint64(0)).Store(&s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to call org.freedesktop.sysupdate1.Target.ListFeatures:", err)
		os.Exit(1)
	}
	return s, nil
}

func DescribeFeature(name string) (*FeatureDescription, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		os.Exit(1)
	}
	defer conn.Close()
	var res string
	obj := conn.Object("org.freedesktop.sysupdate1", "/org/freedesktop/sysupdate1/target/host")
	err = obj.Call("org.freedesktop.sysupdate1.Target.DescribeFeature", 0, name, uint64(0)).Store(&res)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to call DescribeFeature function (is the server snowctl running?):", err)
		os.Exit(1)
	}

	var desc FeatureDescription
	err = json.Unmarshal([]byte(res), &desc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to unmarshal JSON response:", err)
		os.Exit(1)
	}

	return &desc, nil
}
