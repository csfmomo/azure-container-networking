// Copyright Microsoft Corp.
// All rights reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Azure/Aqua/core"
	"github.com/Azure/Aqua/ipam"
	"github.com/Azure/Aqua/log"
	"github.com/Azure/Aqua/network"
)

// Binary version
const version = "v0.1"

// Libnetwork plugin names
const netPluginName = "aqua"
const ipamPluginName = "nullipam"

// Prints description and usage information.
func printHelp() {
	fmt.Println("Usage: aqua [net] [ipam]")
}

func main() {
	var netPlugin network.NetPlugin
	var ipamPlugin ipam.IpamPlugin
	var err error

	// Set defaults.
	logTarget := log.TargetStderr

	// Parse command line arguments.
	args := os.Args

	if len(args) == 1 {
		printHelp()
		return
	}

	handleDependencies()

	for i, arg := range args {
		if i == 0 {
			continue
		}

		switch arg {
		case "net":
			netPlugin, err = network.NewPlugin(netPluginName, version)
			if err != nil {
				fmt.Printf("Failed to create network plugin %v\n", err)
				return
			}

		case "ipam":
			ipamPlugin, err = ipam.NewPlugin(ipamPluginName, version)
			if err != nil {
				fmt.Printf("Failed to create IPAM plugin %v\n", err)
				return
			}

		case "--log-target=syslog":
			logTarget = log.TargetSyslog

		default:
			fmt.Printf("Unknown argument: %s\n", arg)
			printHelp()
			return
		}
	}

	// Create a channel to receive unhandled errors from the plugins.
	errorChan := make(chan error, 1)

	// Create logging provider.
	err = log.SetTarget(logTarget)
	if err != nil {
		fmt.Printf("Failed to configure logging: %v\n", err)
		return
	}

	// Start plugins.
	if netPlugin != nil {
		err = netPlugin.Start(errorChan)
		if err != nil {
			fmt.Printf("Failed to start network plugin %v\n", err)
			return
		}
	}

	if ipamPlugin != nil {
		err = ipamPlugin.Start(errorChan)
		if err != nil {
			fmt.Printf("Failed to start IPAM plugin %v\n", err)
			return
		}
	}

	// Shutdown on two conditions:
	//    a. Unhandled exceptions in plugins
	//    b. Explicit OS signal
	osSignalChannel := make(chan os.Signal, 1)

	// Relay these incoming signals to OS signal channel.
	signal.Notify(osSignalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait until receiving a signal.
	select {
	case sig := <-osSignalChannel:
		fmt.Printf("\nCaught signal <" + sig.String() + "> shutting down...\n")
	case err := <-errorChan:
		if err != nil {
			fmt.Printf("\nReceived unhandled error %v, shutting down...\n", err)
		}
	}

	// Cleanup.
	if netPlugin != nil {
		netPlugin.Stop()
	}

	if ipamPlugin != nil {
		ipamPlugin.Stop()
	}
}

func handleDependencies() {
	installEbtables()
}

func installEbtables() {
	contents, err := ioutil.ReadFile("/proc/version")
	if err == nil {
		value := string(contents)
		if strings.Contains(value, "ubuntu") || strings.Contains(value, "Ubuntu") {
			fmt.Print("Detected ubuntu " + value)
			core.ExecuteShellCommand("apt-get install ebtables")
		}
	} // else unsupported distro
}