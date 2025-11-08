package export

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

func WriteEntries(w io.Writer, format Format, entries []logentry.Entry) error {
	switch format {
	case FormatCSV:
		return writeCSV(w, entries)
	default:
		return writeJSON(w, entries)
	}
}

func writeJSON(w io.Writer, entries []logentry.Entry) error {
	enc := json.NewEncoder(w)
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			return err
		}
	}
	return nil
}

func writeCSV(w io.Writer, entries []logentry.Entry) error {
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"timestamp", "level", "service", "message"})
	for _, entry := range entries {
		row := []string{
			entry.Timestamp.Format(time.RFC3339),
			string(entry.Level),
			entry.Service,
			entry.Message,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
