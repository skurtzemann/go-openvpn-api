package vpn

import (
	"bufio"
	"os"
	"strings"
)

type VpnUser struct {
	Name      string `json:"name"`
	Enable    bool   `json:"enable"`
	IpAddress string `json:"ip_address"`
	NetMask   string `json:"netmask"`
}

// set attributes functions
func (user *VpnUser) setEnable(status bool) {
	user.Enable = status
}
func (user *VpnUser) setIpAddress(ipAdrress string) {
	user.IpAddress = ipAdrress
}
func (user *VpnUser) setNetMask(netmask string) {
	user.NetMask = netmask
}

// Parse the client config file
func (user *VpnUser) ParseConfigFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		user.parseLine(scanner.Text())
	}

	return nil
}

// Parse the config file line and set the appropriate attributes
// We can find the below instructions
// - disable : disable the client
// - ifconfig-push  : push virtual endpoint to the client tunnel
func (user *VpnUser) parseLine(line string) (err error) {
	line = strings.TrimSpace(line)

	// we ignored empty line an comments (starting with a #)
	if (line == "") || (line[0] == '#') {
		return nil
	}

	fields := strings.Fields(line)
	if len(fields) == 1 {
		if line == "disable" {
			user.setEnable(false)
			return nil
		}
	} else {
		if fields[0] == "ifconfig-push" {
			user.setIpAddress(fields[1])
			user.setNetMask(fields[2])
			return nil
		}
	}

	return nil
}
