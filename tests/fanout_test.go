package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func TestMultiSenderFanout(t *testing.T) {
	receiver1 := NewMockReceiver()
	defer receiver1.Close()
	receiver2 := NewMockReceiver()
	defer receiver2.Close()

	sender1 := postSender(receiver1.URL())
	sender1.Header["Authorization"] = "Bearer askjfiwjfiwfjnaief"

	sender2 := postSender(receiver2.URL())
	sender2.Method = "GET"

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender1, sender2}))

	res := postJSON(t, server.URL+"/alpha", `{"event": "push", "repo": "crossover"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs1 := receiver1.WaitForRequests(1, time.Second)
	require.Len(t, reqs1, 1)
	assert.Equal(t, http.MethodPost, reqs1[0].Method)
	assert.Equal(t, "Bearer askjfiwjfiwfjnaief", reqs1[0].Headers.Get("Authorization"))

	reqs2 := receiver2.WaitForRequests(1, time.Second)
	require.Len(t, reqs2, 1)
	assert.Equal(t, http.MethodGet, reqs2[0].Method)
	assert.Equal(t, "application/json", reqs2[0].Headers.Get("Content-Type"))
	assert.Empty(t, reqs2[0].Headers.Get("Authorization"))
}
