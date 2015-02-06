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
	if !sort.	(ipAddresses) {log.Fatal("IP addresses not sorted!")}

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

func get4thIpByte(ipAddress string) int {
	loIndex := strings.LastIndex(ipAddress, ".")
	hiIndex := strings.Index(ipAddress, ":")

	parsedInt, err := strconv.ParseInt(ipAddress[loIndex+1:hiIndex], 10, 0)
	if err != nil {log.Fatal(err)}

	return int(parsedInt)
}

func main() {
	rawIps := []string{"129.241.187.16:20014",
					   "129.01.187.57:9354",
					   "129.241.187.225:3687",
					   "129.2.187.79:10784",
					   "129.37.187.0:41634"}

	for i := 0; i < len(rawIps); i++ {
		fmt.Println(get4thIpByte(rawIps[i]))
	}

}