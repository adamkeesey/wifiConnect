package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/akamensky/argparse"
	uuid "github.com/satori/go.uuid"
)

func createFile(ssid string, password string) {

	uuid := uuid.Must(uuid.NewV4())

	fileContents := fmt.Sprintf(`[connection]
id=%s
uuid=%s
type=wifi
permissions=

[wifi]
mac-address=18:5E:0F:87:CC:33
mac-address-blacklist=
mode=infrastructure
ssid=%s

[wifi-security]
auth-alg=open
key-mgmt=wpa-psk
psk=%s

[ipv4]
dns-search=
method=auto

[ipv6]
addr-gen-mode=stable-privacy
dns-search=
method=auto
`, ssid, uuid, ssid, password)

	err := ioutil.WriteFile("/etc/NetworkManager/system-connections/"+ssid, []byte(fileContents), 0600)
	if err != nil {
		fmt.Print(err)
	}
}

func args() (string, string) {
	parser := argparse.NewParser("Wifi Connection Creation", "Create wifi connections from cli")
	ssid := parser.String("s", "ssid", &argparse.Options{Required: true, Help: "SSID to connect to"})
	password := parser.String("p", "pass", &argparse.Options{Required: true, Help: "Password to SSID"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	return *ssid, *password
}

func main() {
	ssid, password := args()
	createFile(ssid, password)
	cmd := exec.Command("service", "network-manager", "restart")
	err := cmd.Run()
	if err != nil {
		fmt.Print(err)
	}
}
