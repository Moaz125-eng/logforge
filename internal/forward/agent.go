package forward

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Agent struct {
	peers  []string
	client *http.Client
	mu     sync.Mutex
	fail   map[string]int
}

func NewAgent(peers []string) *Agent {
	return &Agent{
		peers:  peers,
		client: &http.Client{Timeout: 5 * time.Second},
		fail:   make(map[string]int),
	}
}

func (a *Agent) Forward(ctx context.Context, entry logentry.Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	for _, peer := range a.peers {
		if a.isDown(peer) {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, peer+"/ingest", bytes.NewReader(data))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.client.Do(req)
		if err != nil {
			a.recordFail(peer)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 500 {
			a.recordFail(peer)
		} else {
			a.clearFail(peer)
		}
	}
	return nil
}

func (a *Agent) recordFail(peer string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.fail[peer]++
}

func (a *Agent) clearFail(peer string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.fail, peer)
}

func (a *Agent) isDown(peer string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.fail[peer] >= 3
}
