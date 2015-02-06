package main

import(
	"fmt"
	"sort"
	"strings"
	"strconv"
	"log"
	"net"
)

var _ = net.ParseIP // For debugging; delete when done.
var _ = strconv.ParseInt // For debugging; delete when done.

func ConnectToSlavesByTcp(ipAddresses []string) {
	for i := 0; i < len(ipAddresses); i++ {
		if !isValidPaddedIp(ipAddresses[i]) {
			log.Fatalf("IP adress %s is not valid!", ipAddresses[i])
		}
	}
	
	sort.Strings(ipAddresses)
	if !sort.StringsAreSorted(ipAddresses) {log.Fatal("IP addresses not sorted!")}

	masterIp := ipAddresses[0]
	// masterIp != getOwnIp() -> error!
	slaveIps := ipAddresses[1:]

	fmt.Println(masterIp)

	for i := 0; i < len(slaveIps); i++ {
		fmt.Println(slaveIps[i])
		//raddr, err := net.ResolveTCPAddr("tcp", )
	}
}

func padIpAddress(ipAddress string) string {
	index := strings.Index(ipAddress, ":")
	ipAddress = ipAddress[:index]
	
	if net.ParseIP(ipAddress) == nil {log.Fatalf("IP address %s not valid!", ipAddress)}

	var paddedIp string
	
	for i := 0; i < 3; i++ {
		index = strings.Index(ipAddress, ".")
		if index == 1 {
			paddedIp += "00" + ipAddress[:2]
			ipAddress = ipAddress[2:]
		} else if index == 2 {
			paddedIp += "0" + ipAddress[:3]
			ipAddress = ipAddress[3:]
		} else if index == 3 {
			paddedIp += ipAddress[:4]
			ipAddress = ipAddress[4:]
		}
	}

	if len(ipAddress) == 1 {
		paddedIp += "00" + ipAddress[:1]
	} else if len(ipAddress) == 2 {
		paddedIp += "0" + ipAddress[:2]
	} else if len(ipAddress) == 3 {
		paddedIp += ipAddress[:3]
	}

	return paddedIp
}

func isValidPaddedIp(ipAddress string) bool {
	// use net.ParseIP to validate unpadded IPs
	if len(ipAddress) != 15 {
		return false
	}

	for i := 3; i < 12; i += 4 {
		if ipAddress[i] != '.' {
			return false
		}
	}

	for i := 0; i < 4; i++ {
		parsedInt, err := strconv.ParseInt(ipAddress[4*i+0:4*i+3], 10, 0)
		if err != nil {log.Fatal(err)}
		if parsedInt < 0 || parsedInt > 256 {
			return false
		}
	}
	return true
}

func main() {
	rawIps := []string{"5.41.7.16:20014",
					   "129.01.187.57:9354",
					   "129.241.187.225:3687",
					   "129.2.187.79:10784",
					   "29.37.187.0:41634"}

	var parsedIps []string

	for i := 0; i < len(rawIps); i++ {
		parsedIps = append(parsedIps, padIpAddress(rawIps[i]))
	}

	for i := 0; i < 1; i++ {
		fmt.Println(isValidPaddedIp(parsedIps[i]))
	}
}