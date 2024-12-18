// Package main
package main

import (
	"errors"
	"flag"
	"log/slog"
	"net"
	"os"
	"strconv"
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
			"[schedule] connection data",
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
			"[schedule] switcher device",
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

	ip := consts.DefaultIP
	flag.Func("ip", "The local Switcher device's IP address", func(maybeIP string) error {
		if maybeIP == "" {
			return consts.ErrInvalidIP
		}

		parsedIP := net.ParseIP(maybeIP)
		if parsedIP == nil {
			return consts.ErrInvalidIP
		}

		ip = parsedIP
		return nil
	})

	port := consts.UDPPortType1New
	flag.Func("port", "The local Switcher device's port", func(maybePort string) error {
		p, err := strconv.ParseUint(maybePort, 10, 32)
		if err != nil {
			return errors.Join(consts.ErrInvalidPort, err)
		}
		if p < 100 || p >= 65_000 {
			return consts.ErrInvalidPort
		}
		port = uint(p)

		return nil
	})

	shouldGetSchedule := flag.Bool("schedule", false, "Get you Switcher device's work schedule")
	flag.Parse()

	if err := Start(ip, port, *shouldGetSchedule); err != nil {
		panic(err)
	}
}
