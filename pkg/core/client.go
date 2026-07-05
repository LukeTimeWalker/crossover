package core

import (
	"crypto/tls"
	"net/http"
	"time"
)

var clients map[*Sender]*http.Client

func GetClient(sender *Sender) (*http.Client, error) {
	if clients == nil {
		clients = make(map[*Sender]*http.Client)
	}

	if client, ok := clients[sender]; ok {
		return client, nil
	}

	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sender.Client.TlsSkipVerify},
	}

	timeout, err := time.ParseDuration(sender.Client.Timeout)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transportConfig,
		Timeout:   timeout,
	}

	clients[sender] = client

	return client, nil
}
