package main

import (
	"fmt"
	"os"

	"spark-cluster/cmd/sparkctl/commands"
)

func main() {
	if err := commands.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
