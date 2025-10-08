package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.anx.io/external-dns-webhook/cmd/webhook/init/configuration"
	"go.anx.io/external-dns-webhook/cmd/webhook/init/dnsprovider"
	"go.anx.io/external-dns-webhook/cmd/webhook/init/logging"
	"go.anx.io/external-dns-webhook/cmd/webhook/init/server"
	"go.anx.io/external-dns-webhook/pkg/webhook"
)

const banner = `
    _    _   _ _______  _____    _
   / \  | \ | | ____\ \/ /_ _|  / \
  / _ \ |  \| |  _|  \  / | |  / _ \
 / ___ \| |\  | |___ /  \ | | / ___ \
/_/   \_\_| \_|_____/_/\_\___/_/   \_\
external-dns-anexia-webhook
version: %s (%s)

`

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	fmt.Printf(banner, Version, Gitsha)

	logging.Init()

	config := configuration.Init()
	provider, err := dnsprovider.Init(config)
	if err != nil {
		log.Fatalf("failed to initialize provider: %v", err)
	}

	srv := server.InitAndServe(config, webhook.New(provider))
	server.ShutdownGracefully(srv)
}
