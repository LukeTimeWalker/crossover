package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func TestRouteConditionTrue(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := singleRouteConfig([]core.Sender{postSender(receiver.URL())}, `{{ eq .body.event "push" }}`)
	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{"event": "push"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(1, time.Second), 1)
}

func TestRouteConditionFalseSkipsSenders(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := singleRouteConfig([]core.Sender{postSender(receiver.URL())}, `{{ eq .body.event "push" }}`)
	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{"event": "pull"}`)
	// a failed condition is not an error: 200 OK, but nothing is forwarded
	assert.Equal(t, http.StatusOK, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}

func TestConditionsAreAnded(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := singleRouteConfig(
		[]core.Sender{postSender(receiver.URL())},
		`true`,
		`{{ eq .body.event "pull" }}`,
	)
	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{"event": "push"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}

func TestOrWithinSingleCondition(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := singleRouteConfig(
		[]core.Sender{postSender(receiver.URL())},
		`{{ or (eq .body.event "push") (eq .body.event "tag") }}`,
	)
	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{"event": "tag"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver.WaitForRequests(1, time.Second), 1)
}

func TestSenderLevelCondition(t *testing.T) {
	receiver1 := NewMockReceiver()
	defer receiver1.Close()
	receiver2 := NewMockReceiver()
	defer receiver2.Close()

	sender1 := postSender(receiver1.URL())
	sender2 := postSender(receiver2.URL())
	sender2.Conditions = []string{`{{ eq .body.event "deploy" }}`}

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender1, sender2}))

	res := postJSON(t, server.URL+"/alpha", `{"event": "push"}`)
	// a sender skipped by its condition does not fail the request
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Len(t, receiver1.WaitForRequests(1, time.Second), 1)
	receiver2.AssertNoRequests(t, 100*time.Millisecond)
}

func TestNonBooleanConditionIsError(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	config := singleRouteConfig([]core.Sender{postSender(receiver.URL())}, `maybe`)
	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	receiver.AssertNoRequests(t, 100*time.Millisecond)
}
