// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Creature struct {
	image *ebiten.Image
	X, Y  float64
}

type Shark struct {
	velocity uint
	hunger   uint
	Creature
}

type Fish struct {
	Creature
}

// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	fishes []Fish
	sharks []Shark
}

// Update is called by Ebiten every 'tick' based on Ticks Per Seconds (TPS).
// By default, Ebiten tries to run at 60 TPS so Update will be called every
// 1/60th of a second.  TPS can be changed with the SetTPS method.
func (g *Game) Update() error {
	return nil
}

// Draw is called by Ebiten at the refresh rate of the display to render
// the images on the screen.  For example, when the display rate is 60Hz,
// Ebiten will call Draw 60 times per second.  When a display has a 120Hz
// refresh rate, Draw will be called twice as often as Update.
func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	for _, s := range g.sharks {
		opts.GeoM.Reset()
		opts.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(s.image, opts)
	}

	for _, f := range g.fishes {
		opts.GeoM.Reset()
		opts.GeoM.Translate(f.X, f.Y)
		screen.DrawImage(f.image, opts)
	}

}

// Layout is the logical screen size which can be different from the actual
// screen size.  Ebiten will handle the scaling automatically.  For example,
// if the actual window size is 640x480, the layout can be 320x240 and Ebiten
// will scale the images so that it fits into the window.  This is also useful
// when the window can be resized but the game's logical screen size stays
// constant.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Wa-Tor")

	ss, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/Shark.png")
	if err != nil {
		log.Fatalln("Unable to load shark image.")
	}
	shark := ss.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image)

	fs, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/Orange.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}
	fish := fs.SubImage(image.Rect(0, 0, 32, 16)).(*ebiten.Image)

	wator := &Game{
		fishes: []Fish{
			Fish{
				Creature{
					image: fish,
					X:     10,
					Y:     20,
				},
			},
			Fish{
				Creature{
					image: fish,
					X:     50,
					Y:     70,
				},
			},
			Fish{
				Creature{
					image: fish,
					X:     150,
					Y:     70,
				},
			},
		},
		sharks: []Shark{
			Shark{
				1,
				10,
				Creature{
					image: shark,
					X:     100,
					Y:     200,
				},
			},
		},
	}

	if err := ebiten.RunGame(wator); err != nil {
		log.Fatal(err)
	}
}
