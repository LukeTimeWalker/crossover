package core

import (
	"crypto/tls"
	"net/http"
	"time"
)

// func GetClient(config Client) (*http.Client, error) {
// 	transportConfig := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.TlsSkipVerify},
// 	}

// 	client := &http.Client{Transport: transportConfig, Timeout: time.Duration(config.Timeout) * time.Second}

// 	return client, nil
// }

func GetClient() (*http.Client, error) {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transportConfig, Timeout: time.Duration(10) * time.Second}

	return client, nil
}
