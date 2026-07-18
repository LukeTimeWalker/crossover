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
	status   int
	delay    time.Duration
}

type MockOption func(*MockReceiver)

// WithStatus makes the receiver answer every request with the given status code.
func WithStatus(status int) MockOption {
	return func(m *MockReceiver) {
		m.status = status
	}
}

// WithDelay makes the receiver wait before answering, to trigger client timeouts.
func WithDelay(delay time.Duration) MockOption {
	return func(m *MockReceiver) {
		m.delay = delay
	}
}

func NewMockReceiver(opts ...MockOption) *MockReceiver {
	m := &MockReceiver{status: http.StatusOK}
	for _, opt := range opts {
		opt(m)
	}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.delay > 0 {
			time.Sleep(m.delay)
		}
		body, _ := io.ReadAll(r.Body)
		m.mu.Lock()
		m.received = append(m.received, CapturedRequest{Method: r.Method, Path: r.URL.Path, Headers: r.Header.Clone(), Body: body})
		m.mu.Unlock()
		w.WriteHeader(m.status)
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

// AssertNoRequests waits a grace period and then asserts nothing arrived.
func (m *MockReceiver) AssertNoRequests(t *testing.T, wait time.Duration) {
	t.Helper()
	time.Sleep(wait)
	m.mu.Lock()
	defer m.mu.Unlock()
	assert.Empty(t, m.received, "receiver should not have received any requests")
}

func defaultClient() core.Client {
	return core.Client{
		Timeout:       "10s",
		TlsSkipVerify: true,
	}
}

// postSender is a minimal valid sender posting JSON to the given endpoint.
func postSender(endpoint string) core.Sender {
	return core.Sender{
		Endpoint: endpoint,
		Method:   "POST",
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Client: defaultClient(),
	}
}

// singleRouteConfig wraps senders in one "POST /alpha" route.
func singleRouteConfig(senders []core.Sender, conditions ...string) core.Config {
	return core.Config{
		Routes: []core.Route{{
			Name: "alpha",
			Receiver: core.Receiver{
				Endpoint: "alpha",
				Method:   "POST",
			},
			Sender:     senders,
			Conditions: conditions,
		}},
	}
}

func startCrossover(t *testing.T, config core.Config) *httptest.Server {
	t.Helper()
	crossover := core.NewCrossover(&config)
	server := httptest.NewServer(crossover.Router())
	t.Cleanup(server.Close)
	return server
}

func postJSON(t *testing.T, url string, payload string) *http.Response {
	t.Helper()
	var body io.Reader
	if payload != "" {
		body = strings.NewReader(payload)
	}
	res, err := http.Post(url, "application/json", body)
	require.NoError(t, err)
	res.Body.Close()
	return res
}
