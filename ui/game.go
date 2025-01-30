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
	widgets []*Widget
}

func (g *Game) Update() error {
	for _, widget := range g.widgets {
		if err := widget.game.Update(); err != nil {
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
	// for _, widget := range g.widgets {
	// 	widget.width, widget.height = widget.game.Layout(widget.width, widget.height)
	// }
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

	return &Game{
		engine: engine,
		client: engine.Spawn(client.NewClient(client.ClientConfig{Nodes: pods}), "client"),
		widgets: []*Widget{
			NewWidget(&topicList{}, image.Rect(0, 0, 100, 100)),
		},
	}
}
