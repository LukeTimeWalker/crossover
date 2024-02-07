package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

func GetClient(config ClientConfig) (*http.Client, error) {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.TlsSkipVerify},
	}

	client := &http.Client{Transport: transportConfig, Timeout: time.Duration(config.Timeout) * time.Second}

	return client, nil
}
