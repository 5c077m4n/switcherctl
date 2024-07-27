// Package main
package main

import (
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
	shouldGetSchedule := flag.Bool("schedule", false, "Get you Switcher device's work schedule")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	parsedIP := net.ParseIP(*ip)
	if parsedIP == nil {
		panic(consts.ErrInvalidIP)
	}
	if port == nil || *port < 100 || *port >= 65_000 {
		panic(consts.ErrInvalidPort)
	}

	listener, err := connections.TryNewListener(parsedIP, *port)
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			panic(err)
		}
	}()

	data, err := listener.Read()
	if err != nil {
		panic(err)
	}
	baseDeviceData, err := data.ToJSON()
	if err != nil {
		panic(err)
	}

	slog.Debug(
		"switcher device data",
		"value", baseDeviceData,
	)

	if *shouldGetSchedule {
		biConn, err := connections.TryNewBidirectionalConn(parsedIP, *port, baseDeviceData.ID)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := biConn.Close(); err != nil {
				panic(err)
			}
		}()

		slog.Debug(
			"switcher device schedule",
			"value", biConn.GetSchedules(),
		)
	}
}
