package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Moaz125-eng/logforge/internal/bench"
)

func main() {
	target := flag.String("target", "http://localhost:8080", "logforge http base url")
	workers := flag.Int("workers", 8, "concurrent workers")
	total := flag.Int("total", 1000, "total requests")
	profile := flag.Bool("profile", false, "run default load profiles")
	out := flag.String("out", "", "optional report output path")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if *profile {
		for _, p := range bench.DefaultProfiles {
			runProfile(ctx, *target, p, *out)
		}
		return
	}

	runner := bench.NewRunner(*target, *workers, *total)
	result := runner.Run(ctx)
	bench.PrintReport(result)
	if *out != "" {
		_ = bench.WriteReport(*out, result)
	}
}

func runProfile(ctx context.Context, target string, profile bench.LoadProfile, out string) {
	fmt.Printf("profile %s workers=%d total=%d\n", profile.Name, profile.Workers, profile.Total)
	runner := bench.NewRunner(target, profile.Workers, profile.Total)
	result := runner.Run(ctx)
	bench.PrintReport(result)
	if out != "" {
		path := fmt.Sprintf("%s-%s.json", out, profile.Name)
		_ = bench.WriteReport(path, result)
	}
}
