package spoon

import (
	"testing"

	"github.com/gucooing/spoon/external/tcp"
)

func Test_GateWay_TCP(t *testing.T) {
	gateWay := New(
		ID("1"),
		Name("Test App"),
		Version("1.0.0"),
		Servers(tcp.NewServer(
			[]tcp.ServerOption{
				tcp.SetAddress(":20001"),
			}...,
		)),
	)

	if err := gateWay.Run(); err != nil {
		t.Fatal(err)
	}
}
