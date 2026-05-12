package main

import (
	"github.com/gucooing/spoon"
	"github.com/gucooing/spoon/external/tcp"
	"log"
)

func main() {
	gateWay := spoon.New(
		spoon.ID("1"),
		spoon.Name("Test App"),
		spoon.Version("1.0.0"),
		spoon.Servers(tcp.NewServer(
			[]tcp.ServerOption{
				tcp.SetAddress(":20001"),
			}...,
		)),
	)

	if err := gateWay.Run(); err != nil {
		log.Fatal(err)
	}
}
