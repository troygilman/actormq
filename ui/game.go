package ui

import (
	"image"
	"io"
	"log"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/troygilman/actormq/client"
	"github.com/troygilman/actormq/cluster"
)

type Game struct {
	engine  *actor.Engine
	client  *actor.PID
	widgets []*WidgetContainer
}

func (g *Game) Update() error {
	for _, widget := range g.widgets {
		if err := widget.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for _, widget := range g.widgets {
		widget.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func Run() {
	game := setup()

	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func setup() *Game {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	discoveryPID := engine.Spawn(cluster.NewDiscovery(), "discovery")

	config := cluster.PodConfig{
		Topics:    []string{"test", "test1", "test2"},
		Discovery: discoveryPID,
		Logger:    slog.New(slog.NewJSONHandler(io.Discard, nil)),
		// Logger: slog.Default(),
	}

	pods := []*actor.PID{
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
		engine.Spawn(cluster.NewPod(config), "pod"),
	}

	clientPID := engine.Spawn(client.NewClient(client.ClientConfig{Nodes: pods}), "client")
	game := &Game{
		engine:  engine,
		client:  clientPID,
		widgets: []*WidgetContainer{},
	}

	game.createWidgetContainer(&topicList{}, image.Rect(100, 100, 200, 200))
	return game
}

func (g *Game) createWidgetContainer(widget Widget, rect image.Rectangle) {
	events := make(chan any, 100)
	adapter := g.engine.Spawn(NewAdapter(events), "adapter")
	g.widgets = append(g.widgets, &WidgetContainer{
		widget:  widget,
		rect:    rect,
		engine:  g.engine,
		client:  g.client,
		adapter: adapter,
		events:  events,
	})
}
