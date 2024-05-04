package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cretz/bine/tor"
)

func main() {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting Tor, please wait...")
	t, err := tor.Start(nil, nil)
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer t.Close()

	// Wait for Tor to start up (this is a basic example, in practice, you might need more robust waiting logic)
	time.Sleep(5 * time.Second)

	// Set up HTTP client to use Tor SOCKS proxy
	dialer, err := t.Dialer(nil, nil)
	if err != nil {
		log.Panicf("Unable to create Tor dialer: %v", err)
	}

	httpTransport := &http.Transport{
		DialContext: dialer.DialContext,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   30 * time.Second,
	}

	// Set HTTP proxy
	proxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		log.Panicf("Error parsing proxy URL: %v", err)
	}
	httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	fmt.Println("HTTP client set up to use Tor")

	// Request https://search.brave4u7jddbv7cyviptqjc7jusxh72uik7zt6adtckl5f4nwy2v72qd.onion/
	resp, err := httpClient.Get("https://search.brave4u7jddbv7cyviptqjc7jusxh72uik7zt6adtckl5f4nwy2v72qd.onion/")
	fmt.Println("Requesting https://search.brave4u7jddbv7cyviptqjc7jusxh72uik7zt6adtckl5f4nwy2v72qd.onion/...")
	if err != nil {
		log.Panicf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response body:", resp.Body)
	fmt.Println("Response status:", resp.Status)
}
