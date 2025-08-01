package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

// LokiClient is a lightweight client for sending log entries to Grafana Loki
// over HTTP using the /loki/api/v1/push endpoint.
type LokiClient struct {
	URL    string
	Labels map[string]string
	Client *http.Client
}

// LokiStream represents a single stream entry for a Loki log payload.
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

// LokiPayload is the top-level structure sent to Loki’s /push API.
type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

// NewLokiClient creates and returns a new LokiClient configured with a
// target Loki endpoint and set of constant stream labels.
func NewLokiClient(url string, labels map[string]string) *LokiClient {
	return &LokiClient{
		URL:    url,
		Labels: labels,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

// SendLog pushes a single log message to Loki, wrapped in a stream
// using the client’s default labels and the current timestamp.
func (lc *LokiClient) SendLog(message string) error {
	timestamp := time.Now().UnixNano()
	payload := LokiPayload{
		Streams: []LokiStream{
			{
				Stream: lc.Labels,
				Values: [][2]string{
					{fmt.Sprintf("%d", timestamp), message},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", lc.URL+"/loki/api/v1/push", bytes.NewBuffer(body))
	if err != nil {
		log.Error(err)
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := lc.Client.Do(req)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Error(err)
		return fmt.Errorf("loki returned status: %s", resp.Status)
	}

	return nil
}

// LokiSyncer is a custom zapcore.WriteSyncer implementation that writes logs
// to a Loki backend using a LokiClient.
type LokiSyncer struct {
	Client *LokiClient
}

// NewLokiSyncer wraps a LokiClient in a zapcore.WriteSyncer interface,
// allowing integration with Uber Zap’s logging system.
func NewLokiSyncer(client *LokiClient) zapcore.WriteSyncer {
	return &LokiSyncer{Client: client}
}

// Write sends the given log bytes to Loki using the client's SendLog method.
func (l *LokiSyncer) Write(p []byte) (n int, err error) {
	err = l.Client.SendLog(string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Sync is a no-op for LokiSyncer and always returns nil to satisfy the zapcore.WriteSyncer interface.
func (l *LokiSyncer) Sync() error {
	return nil
}
