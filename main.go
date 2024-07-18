// Package main
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"switcherctl/connections"
	"switcherctl/consts"
)

func main() {
	ip := flag.String("ip", consts.DefaultIP.String(), "The local Switcher device's IP address")
	port := flag.Int("port", consts.UDPPortType1New, "The local Switcher device's port")
	flag.Parse()

	parsedIP := net.ParseIP(*ip)
	if parsedIP == nil {
		log.Fatalln(consts.ErrInvalidIP)
		return
	}
	if port == nil || *port < 100 || *port >= 65_000 {
		log.Fatalln(consts.ErrInvalidPort)
		return
	}

	conn, err := connections.TryNewListener(parsedIP, *port)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Fatalln(closeErr)
		}
	}()

	data, err := conn.Read()
	if err != nil {
		log.Fatalln(err)
		return
	}

	results, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(string(results))
}
