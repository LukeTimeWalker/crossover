package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func TestBodyTranslation(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Body = `{"repository": "{{ .body.repo }}", "action": "{{ .body.event }}", "author": "{{ .body.commit.author }}"}`
	sender.Header["Authorization"] = "Bearer askjfiwjfiwfjnaief"

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	payload := `{"event": "push", "repo": "crossover", "commit": {"author": "lukas"}}`
	res := postJSON(t, server.URL+"/alpha", payload)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, http.MethodPost, reqs[0].Method)
	assert.Equal(t, "application/json", reqs[0].Headers.Get("Content-Type"))
	assert.Equal(t, "Bearer askjfiwjfiwfjnaief", reqs[0].Headers.Get("Authorization"))
	assert.JSONEq(t, `{"repository": "crossover", "action": "push", "author": "lukas"}`, string(reqs[0].Body))
}

func TestUrlParamTemplating(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Body = `{"token": "{{ .url.token }}"}`

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha?token=xyz", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.JSONEq(t, `{"token": "xyz"}`, string(reqs[0].Body))
}

func TestMultiValueUrlParams(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Body = `{"tags": "{{ index .url.tag 0 }},{{ index .url.tag 1 }}"}`

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha?tag=a&tag=b", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.JSONEq(t, `{"tags": "a,b"}`, string(reqs[0].Body))
}

func TestEnvTemplating(t *testing.T) {
	t.Setenv("CROSSOVER_TEST_TOKEN", "s3cret")

	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Header["Authorization"] = "Bearer {{ .env.CROSSOVER_TEST_TOKEN }}"

	config := singleRouteConfig([]core.Sender{sender})
	config.EnvironementVariables = []string{"CROSSOVER_TEST_TOKEN"}

	server := startCrossover(t, config)

	res := postJSON(t, server.URL+"/alpha", `{}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, "Bearer s3cret", reqs[0].Headers.Get("Authorization"))
}

func TestTemplatedSenderEndpoint(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL() + "/hooks/{{ .body.target }}")

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha", `{"target": "deploy"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, "/hooks/deploy", reqs[0].Path)
}

func TestTemplatedHeaders(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Header["X-Repo"] = "{{ .body.repo }}"

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha", `{"repo": "crossover"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, "crossover", reqs[0].Headers.Get("X-Repo"))
}

func TestNewlineInBodyIsPreserved(t *testing.T) {
	receiver := NewMockReceiver()
	defer receiver.Close()

	sender := postSender(receiver.URL())
	sender.Body = `{{ .body.msg }}`

	server := startCrossover(t, singleRouteConfig([]core.Sender{sender}))

	res := postJSON(t, server.URL+"/alpha", `{"msg": "line1\nline2"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	reqs := receiver.WaitForRequests(1, time.Second)
	require.Len(t, reqs, 1)
	assert.Equal(t, `line1\nline2`, string(reqs[0].Body))
}
