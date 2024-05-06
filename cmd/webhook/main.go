package main

import (
	"fmt"

	"github.com/probstenhias/external-dns-anexia-webhook/cmd/webhook/init/configuration"
	"github.com/probstenhias/external-dns-anexia-webhook/cmd/webhook/init/dnsprovider"
	"github.com/probstenhias/external-dns-anexia-webhook/cmd/webhook/init/logging"
	"github.com/probstenhias/external-dns-anexia-webhook/cmd/webhook/init/server"
	"github.com/probstenhias/external-dns-anexia-webhook/pkg/webhook"
	log "github.com/sirupsen/logrus"
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

	srv := server.Init(config, webhook.New(provider))
	server.ShutdownGracefully(srv)
}
