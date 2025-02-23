package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Creature struct {
	img  *ebiten.Image
	x, y float64
}

type Shark struct {
	velocity uint
	hunger   uint
	Creature
}

type Fish struct {
	Creature
}

type Game struct {
	fish  Fish
	shark Shark
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.shark.x, g.shark.y)
	screen.DrawImage(g.shark.img.SubImage(
		image.Rect(0, 0, 32, 32)).(*ebiten.Image),
		opts,
	)

	opts.GeoM.Reset()
	opts.GeoM.Translate(g.fish.x, g.fish.y)
	screen.DrawImage(g.fish.img.SubImage(
		image.Rect(0, 0, 32, 16)).(*ebiten.Image),
		opts,
	)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Wa-Tor")

	shark, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/Shark.png")
	if err != nil {
		log.Fatalln("Unable to load shark image.")
	}

	fish, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/Orange.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}

	if err := ebiten.RunGame(&Game{
		fish: Fish{
			Creature{
				img: fish,
				x:   0,
				y:   0,
			},
		},
		shark: Shark{
			2,
			10,
			Creature{
				img: shark,
				x:   150,
				y:   50,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
