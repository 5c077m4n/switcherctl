// Package main
package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net"
	"os"
	"switcherctl/connections"
	"switcherctl/consts"
)

func main() {
	ip := flag.String("ip", consts.DefaultIP.String(), "The local Switcher device's IP address")
	port := flag.Int("port", consts.UDPPortType1New, "The local Switcher device's port")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	parsedIP := net.ParseIP(*ip)
	if parsedIP == nil {
		panic(consts.ErrInvalidIP)
	}
	if port == nil || *port < 100 || *port >= 65_000 {
		panic(consts.ErrInvalidPort)
	}

	conn, err := connections.TryNewListener(parsedIP, *port)
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			panic(err)
		}
	}()

	data, err := conn.Read()
	if err != nil {
		panic(err)
	}

	results, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	slog.Info(
		"switcher device data",
		"value", string(results),
	)
}
