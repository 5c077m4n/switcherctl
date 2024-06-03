// Package main
package main

import (
	"encoding/json"
	"log"
	"switcherctl/connection"
	"switcherctl/consts"
)

func main() {
	port, ok := consts.DeviceCategoryToUDPPort[consts.DeviceCategoryWaterHeater]
	if !ok {
		log.Fatalln("Could not find port for this device")
	}

	conn, err := connection.TryNew(consts.DefaultIP, port)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Fatalln(closeErr)
		}
	}()

	data, err := conn.Read()
	if err != nil {
		log.Fatalln(err)
	}

	results, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(results))
}
