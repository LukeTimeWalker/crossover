package crossover

import (
	"crypto/tls"
	"net/http"
	"time"
)

func GetClient(config Client) (*http.Client, error) {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.TlsSkipVerify},
	}

	client := &http.Client{Transport: transportConfig, Timeout: time.Duration(config.Timeout) * time.Second}

	return client, nil
}
