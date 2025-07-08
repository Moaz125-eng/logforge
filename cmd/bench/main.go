package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "logforge bench: run after benchmark package is wired")
	os.Exit(1)
}
