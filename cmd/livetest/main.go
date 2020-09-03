package main

import (
	"github.com/axelspringer/swerve/cmd/livetest/testclient"

	"github.com/axelspringer/swerve/cmd/livetest/config"
)

func main() {
	cfg := config.Get()
	client := testclient.New(cfg)
	client.Version()
	client.InsertRedirects()
	client.TestRedirects()
}
