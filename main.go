package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

	log.Printf("The UDP server is connected @ %s\n", conn.RemoteAddress().String())

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		data := []byte(text + "\n")
		_, err = conn.Write(data)
		if err != nil {
			log.Fatalln(err)
		}

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		resp, err := conn.ReadFromUDP()
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Reply: %s\n", resp)
	}
}
