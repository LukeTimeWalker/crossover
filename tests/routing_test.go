package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func TestMethodNotAllowed(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(receiver.URL())}))

	res, err := http.Get(server.URL + "/alpha")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}

func TestUnknownEndpointNotFound(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	server := startCrossover(t, singleRouteConfig([]core.Sender{postSender(receiver.URL())}))

	res := postJSON(t, server.URL+"/unknown", `{}`)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}

func TestMultipleRoutesAreIndependent(t *testing.T) {
	receiverAlpha := NewMockReceiver()
	defer receiverAlpha.Close()
	receiverBeta := NewMockReceiver()
	defer receiverBeta.Close()

	config := core.Config{
		Routes: []core.Route{
			{
				Name:     "alpha",
				Receiver: core.Receiver{Endpoint: "alpha", Method: "POST"},
				Sender:   []core.Sender{postSender(receiverAlpha.URL())},
			},
			{
				Name:     "beta",
				Receiver: core.Receiver{Endpoint: "beta", Method: "POST"},
				Sender:   []core.Sender{postSender(receiverBeta.URL())},
			},
		},
	}

	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.Len(t, receiverAlpha.WaitForRequests(1, time.Second), 1)
	receiverBeta.AssertNoRequests(t, 100*time.Millisecond)
}

func TestDynamicRouteRegisterUnregister(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	crossover := core.NewCrossover(&core.Config{})
	server := httptest.NewServer(crossover.Router())
	defer server.Close()

	route := core.Route{
		Name:     "dynamic",
		Receiver: core.Receiver{Endpoint: "dynamic", Method: "POST"},
		Sender:   []core.Sender{postSender(receiver.URL())},
	}

	assert.NoError(t, crossover.RegisterRoute(route, nil))
	res := postJSON(t, server.URL+"/dynamic", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(1, time.Second), 1)

	crossover.UnregisterRoute(route)
	res = postJSON(t, server.URL+"/dynamic", `{}`)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	assert.NoError(t, crossover.RegisterRoute(route, nil))
	res = postJSON(t, server.URL+"/dynamic", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(2, time.Second), 2)
}
