package main

import (
	"fmt"
	"os"

	"apigateway/internal/cli"
	"apigateway/internal/config"
)

func main() {
	if err := config.Init("config.yaml"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
	}

	c, err := cli.NewRouteCLI(config.GlobalConfig.Etcd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	if err := c.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
