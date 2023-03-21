package main

import (
	"log"
	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/outscale"
	"github.com/blushft/go-diagrams/nodes/apps"
)

func main() {
	d, err := diagram.New(diagram.Filename("app"), diagram.Label("App"), diagram.Direction("LR"))
	if err != nil {
		log.Fatal(err)
	}

	dns := outscale.Network.InternetService(diagram.NodeLabel("DNS"))
	lb := outscale.Network.LoadBalancer(diagram.NodeLabel("NLB"))
	cache := apps.Inmemory.Redis(diagram.NodeLabel("Cache"))
	db := apps.Database.Postgresql(diagram.NodeLabel("Database"))

	dc := diagram.NewGroup("OUTSCALE")
	dc.NewGroup("services").
		Label("Service Layer").
		Add(
			outscale.Compute.Compute(diagram.NodeLabel("Server 1")),
			outscale.Compute.Compute(diagram.NodeLabel("Server 2")),
			outscale.Compute.Compute(diagram.NodeLabel("Server 3")),
		).
		ConnectAllFrom(lb.ID(), diagram.Forward()).
		ConnectAllTo(cache.ID(), diagram.Forward())

	dc.NewGroup("data").Label("Data Layer").Add(cache, db).Connect(cache, db)

	d.Connect(dns, lb, diagram.Forward()).Group(dc)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
