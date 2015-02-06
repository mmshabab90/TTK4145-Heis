package main

import(
	"fmt"
	"sort"
	"strings"
	"strconv"
	"log"
	//"net"
)

func ConnectToSlavesByTcp(ipAddresses []string) {
	// ! validIpList(ipAddresses) -> error!
	
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

func sortIpsAscending(ipAddresses []string) []string {
	var ipByte4 int
	for i := 0; i < len(ipAddresses); i++ {
		ipByte4 = get4thIpByte(ipAddresses[i])
		fmt.Println(ipByte4)
	}

	return ipAddresses
}

func padIpAddress(ipAddress string) string {
	// check valid ip
	var paddedIp string
	var index int
	for i := 0; i < 3; i++ {
		index = strings.Index(ipAddress, ".")
		if index == 1 {
			paddedIp += "00" + ipAddress[:2]
			ipAddress = ipAddress[2:]
		} else if index == 2 {
			paddedIp += "0" + ipAddress[:3]
			ipAddress = ipAddress[3:]
		}
	}

	index = strings.Index(ipAddress, ":")
	if index == 1 {
			paddedIp += "00" + ipAddress[:1]
		} else if index == 2 {
			paddedIp += "0" + ipAddress[:2]
	}


	return paddedIp

	// for hver av de fire segmentene:
		// stapp inn riktig antall nuller (index-3)
}

func get4thIpByte(ipAddress string) int {
	loIndex := strings.LastIndex(ipAddress, ".")
	hiIndex := strings.Index(ipAddress, ":")

	parsedInt, err := strconv.ParseInt(ipAddress[loIndex+1:hiIndex], 10, 0)
	if err != nil {log.Fatal(err)}

	return int(parsedInt)
}

func main() {
	rawIps := []string{"5.41.7.16:20014",
					   "129.01.187.57:9354",
					   "129.241.187.225:3687",
					   "129.2.187.79:10784",
					   "29.37.187.0:41634"}
	for i := 
	fmt.Println(padIpAddress(rawIps[0]))

}