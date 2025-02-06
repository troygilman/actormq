package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/troygilman/actormq/ui/types"
)

func Run() {
	// game := setup()

	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(types.NewContainer(&topicList{})); err != nil {
		log.Fatal(err)
	}
}

// func setup() *Game {
// 	engine, err := actor.NewEngine(actor.NewEngineConfig())
// 	if err != nil {
// 		panic(err)
// 	}

// 	discoveryPID := engine.Spawn(cluster.NewDiscovery(), "discovery")

// 	config := cluster.PodConfig{
// 		Topics:    []string{"test", "test1", "test2"},
// 		Discovery: discoveryPID,
// 		Logger:    slog.New(slog.NewJSONHandler(io.Discard, nil)),
// 		// Logger: slog.Default(),
// 	}

// 	pods := []*actor.PID{
// 		engine.Spawn(cluster.NewPod(config), "pod"),
// 		engine.Spawn(cluster.NewPod(config), "pod"),
// 		engine.Spawn(cluster.NewPod(config), "pod"),
// 	}

// 	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{Nodes: pods}), "client")
// 	game := &Game{
// 		engine:  engine,
// 		client:  clientPID,
// 		widgets: []*WidgetContainer{},
// 	}

// 	game.createWidgetContainer(&topicList{}, image.Rect(100, 100, 200, 200))
// 	return game
// }

// func (g *Game) createWidgetContainer(widget Widget, rect image.Rectangle) {
// 	events := make(chan any, 100)
// 	adapter := g.engine.Spawn(NewAdapter(events), "adapter")
// 	g.widgets = append(g.widgets, &WidgetContainer{
// 		widget:  widget,
// 		rect:    rect,
// 		engine:  g.engine,
// 		client:  g.client,
// 		adapter: adapter,
// 		events:  events,
// 	})
// }
