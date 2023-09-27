package main

import (
	"example.com/m/v2/NetworkScanner"
	"fmt"
	"os"
	"strconv"
)

const (
	MySQLProtocolVersion = 10
)

func main() {

	hostIP := os.Getenv("HOST_IP")
	if hostIP == "" {
		hostIP = "localhost"
	}

	portStr := os.Getenv("HOST_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || portStr == "" {
		port = 3306
	}

	scanner := NetworkScanner.NewNetworkScanner(hostIP, port)

	instance := scanner.FindMySqlInstance()

	fmt.Println("MySQL instance found: ", instance.MySQLActive)
	if instance.MySQLActive {
		fmt.Println("MySQL version: ", instance.MySQLServerVersion)
		fmt.Println("MySQL banner: ", instance.MySQLBanner)
		fmt.Println("MySQL protocol: ", instance.MySQLProtocolVersion)
		fmt.Println("MySQL target port: ", instance.OriginalPort)
		fmt.Println("MySQL target Host: ", instance.OriginalHost)
	}
}
