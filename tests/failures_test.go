package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func TestDownstreamErrorReturns400(t *testing.T) {
	receiver := NewMockReceiver(WithStatus(http.StatusInternalServerError))
	defer receiver.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(receiver.URL())}))

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(1, time.Second), 1)
}

func TestPartialFanoutFailure(t *testing.T) {
	receiverOk := NewMockReceiver()
	defer receiverOk.Close()
	receiverFail := NewMockReceiver(WithStatus(http.StatusInternalServerError))
	defer receiverFail.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{
		postSender(receiverOk.URL()),
		postSender(receiverFail.URL()),
	}))

	res := postJSON(t, server.URL+"/alpha", `{}`)
	// one failing sender fails the whole request, but does not cancel the others
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Len(t, receiverOk.WaitForRequests(1, time.Second), 1)
	assert.Len(t, receiverFail.WaitForRequests(1, time.Second), 1)
}

func TestUnreachableDownstreamReturns400(t *testing.T) {
	unreachable := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := unreachable.URL
	unreachable.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(url)}))

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestMalformedJsonReturns400(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(receiver.URL())}))

	res := postJSON(t, server.URL+"/alpha", `{not json`)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}

func TestEmptyBodyIsAccepted(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(receiver.URL())}))

	res := postJSON(t, server.URL+"/alpha", "")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(1, time.Second), 1)
}

func TestClientTimeoutReturns400(t *testing.T) {
	receiver := NewMockReceiver(WithDelay(300 * time.Millisecond))
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Client.Timeout = "50ms"

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestRegisterRouteWithInvalidTemplateFails(t *testing.T) {
	crossover := core.NewCrossover(&core.Config{})

	badBody := postSender("http://localhost:1")
	badBody.Body = `{{ .body.unclosed`

	err := crossover.RegisterRoute(core.Route{
		Name:     "broken",
		Receiver: core.Receiver{Endpoint: "broken", Method: "POST"},
		Sender:   []core.Sender{badBody},
	}, nil)
	assert.Error(t, err)

	err = crossover.RegisterRoute(core.Route{
		Name:       "broken-condition",
		Receiver:   core.Receiver{Endpoint: "broken-condition", Method: "POST"},
		Sender:     []core.Sender{postSender("http://localhost:1")},
		Conditions: []string{`{{ eq .body.event`},
	}, nil)
	assert.Error(t, err)
}
