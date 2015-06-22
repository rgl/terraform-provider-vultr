package main

import (
	"fmt"
	"github.com/JamesClonk/vultr/lib"
	"log"
)

type Config struct {
	ApiKey string
}

// Client() returns a new client for accessing Vultr.
func (c *Config) Client() (*lib.Client, error) {
	client := lib.NewClient(
		c.ApiKey,
		&lib.Options{
			UserAgent: fmt.Sprintf("vultr-go/%s terraform", lib.Version)})

	log.Printf("[INFO] Vultr Client configured for URL: %s", client.Endpoint)

	return client, nil
}
