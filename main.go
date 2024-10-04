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

func Start(ip net.IP, port uint, shouldGetSchedule bool) error {
	listener, err := connections.TryNewListener(ip, int(port))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			panic(err)
		}
	}()

	data, err := listener.Read()
	if err != nil {
		return err
	}
	baseDeviceData, err := data.ToJSON()
	if err != nil {
		return err
	}

	slog.Debug(
		"switcher device data",
		"value", baseDeviceData,
	)

	if shouldGetSchedule {
		slog.Debug(
			"connection data",
			"ip", ip,
			"port", port,
			"device ID", baseDeviceData.ID,
		)

		biConn, err := connections.TryNewBidirectionalConn(ip, int(port), baseDeviceData.ID)
		if err != nil {
			return err
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

	return nil
}

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)
	slog.SetDefault(logger)

	var ip net.IP
	flag.Func("ip", "The local Switcher device's IP address", func(maybeIP string) error {
		ip = net.ParseIP(maybeIP)
		if ip == nil {
			return consts.ErrInvalidIP
		}
		return nil
	})
	port := flag.Uint("port", consts.UDPPortType1New, "The local Switcher device's port")
	shouldGetSchedule := flag.Bool("schedule", false, "Get you Switcher device's work schedule")
	flag.Parse()

	if port == nil || *port < 100 || *port >= 65_000 {
		panic(consts.ErrInvalidPort)
	}

	if err := Start(ip, *port, *shouldGetSchedule); err != nil {
		panic(err)
	}
}
