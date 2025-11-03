package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Notifier struct {
	target string
	client *http.Client
	mu     sync.Mutex
	sent   uint64
}

func NewNotifier(target string) *Notifier {
	return &Notifier{
		target: target,
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (n *Notifier) Notify(events []Event) {
	if n.target == "" || len(events) == 0 {
		return
	}
	for _, ev := range events {
		data, err := json.Marshal(ev)
		if err != nil {
			continue
		}
		req, err := http.NewRequest(http.MethodPost, n.target, bytes.NewReader(data))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := n.client.Do(req)
		if err == nil {
			resp.Body.Close()
			n.mu.Lock()
			n.sent++
			n.mu.Unlock()
		}
	}
}

func (n *Notifier) Sent() uint64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.sent
}
