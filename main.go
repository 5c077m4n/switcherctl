// Package main
package main

import (
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

	for {
		data, err := conn.Read()
		if err != nil {
			log.Fatalln(err)
		}
		ip, err := data.GetIPType1()
		if err != nil {
			log.Fatalln(err)
		}
		autoShutdown, err := data.GetTimeToShutdown()
		if err != nil {
			log.Fatalln(err)
		}
		remaining, err := data.GetRemainingTime()
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf(`
Received: "%s"
> From: "%s"
> Device ID: "%s"
> Key: "%s"
> IP: %s
> MAC: %s
> On: %v
> Auto shutdown in: %s
> Ramaingin time: %s

`,
			data,
			data.GetDeviceName(),
			data.GetDeviceID(),
			data.GetDeviceKey(),
			ip,
			data.GetDeviceMAC(),
			data.IsPoweredOn(),
			autoShutdown,
			remaining,
		)
	}
}
