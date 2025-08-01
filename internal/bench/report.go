package bench

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintReport(result Result) {
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func WriteReport(path string, result Result) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

type LoadProfile struct {
	Name    string `json:"name"`
	Workers int    `json:"workers"`
	Total   int    `json:"total"`
}

var DefaultProfiles = []LoadProfile{
	{Name: "warmup", Workers: 2, Total: 200},
	{Name: "steady", Workers: 8, Total: 2000},
	{Name: "burst", Workers: 32, Total: 8000},
}
