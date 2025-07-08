package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "logforge agent: configure LOGFORGE_FORWARD_PEERS and run server build first")
	os.Exit(1)
}
