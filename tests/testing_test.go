package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

type CapturedRequest struct {
	Method  string
	Path    string
	Headers http.Header
	Body    []byte
}

type MockReceiver struct {
	server   *httptest.Server
	mu       sync.Mutex
	received []CapturedRequest
}

func NewMockReceiver() *MockReceiver {
	m := &MockReceiver{}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		m.mu.Lock()
		m.received = append(m.received, CapturedRequest{Method: r.Method, Path: r.URL.Path, Headers: r.Header.Clone(), Body: body})
		m.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	return m
}

func (m *MockReceiver) URL() string {
	return m.server.URL
}

func (m *MockReceiver) Close() {
	m.server.Close()
}

func (m *MockReceiver) WaitForRequests(count int, timeout time.Duration) []CapturedRequest {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		m.mu.Lock()
		n := len(m.received)
		m.mu.Unlock()
		if n >= count {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]CapturedRequest{}, m.received...)
}

func defaultClient() core.Client {
	return core.Client{
		Timeout:       "10s",
		TlsSkipVerify: true,
	}
}

func TestCrossover(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := core.Config{
		Routes: []core.Route{{
			Name: "alpha",
			Receiver: core.Receiver{
				Endpoint: "alpha",
				Method:   "POST",
			},
			Sender: []core.Sender{{
				Endpoint: receiver.URL(),
				Method:   "POST",
				Body:     `{"repository": "{{ .body.repo }}", "action": "{{ .body.event }}"}`,
				Header: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer askjfiwjfiwfjnaief",
				},
				Client: defaultClient(),
			}},
		}},
	}

	crossover := core.NewCrossover(&config)

	crossoverServer := httptest.NewServer(crossover.Router())
	defer crossoverServer.Close()

	payload := `{"event": "push", "repo": "crossover"}`
	res, err := http.Post(crossoverServer.URL+"/alpha", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, http.MethodPost, reqs[0].Method)
	assert.Equal(t, "application/json", reqs[0].Headers.Get("Content-Type"))
	assert.Equal(t, "Bearer askjfiwjfiwfjnaief", reqs[0].Headers.Get("Authorization"))
	assert.JSONEq(t, `{"repository": "crossover", "action": "push"}`, string(reqs[0].Body))
}

func TestCrossoverMultipleSenders(t *testing.T) {
	receiver1 := NewMockReceiver()
	defer receiver1.Close()
	receiver2 := NewMockReceiver()
	defer receiver2.Close()

	config := core.Config{
		Routes: []core.Route{{
			Name: "alpha",
			Receiver: core.Receiver{
				Endpoint: "alpha",
				Method:   "POST",
			},
			Sender: []core.Sender{
				{
					Endpoint: receiver1.URL(),
					Method:   "POST",
					Header: map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer askjfiwjfiwfjnaief",
					},
					Client: defaultClient(),
				},
				{
					Endpoint: receiver2.URL(),
					Method:   "GET",
					Header: map[string]string{
						"Content-Type": "application/json",
					},
					Client: defaultClient(),
				},
			},
		}},
	}

	crossover := core.NewCrossover(&config)

	crossoverServer := httptest.NewServer(crossover.Router())
	defer crossoverServer.Close()

	payload := `{"event": "push", "repo": "crossover"}`
	res, err := http.Post(crossoverServer.URL+"/alpha", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs1 := receiver1.WaitForRequests(1, time.Second)
	require.Len(t, reqs1, 1)
	assert.Equal(t, http.MethodPost, reqs1[0].Method)
	assert.Equal(t, "Bearer askjfiwjfiwfjnaief", reqs1[0].Headers.Get("Authorization"))

	reqs2 := receiver2.WaitForRequests(1, time.Second)
	require.Len(t, reqs2, 1)
	assert.Equal(t, http.MethodGet, reqs2[0].Method)
	assert.Equal(t, "application/json", reqs2[0].Headers.Get("Content-Type"))
}
