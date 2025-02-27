// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main // package lazyhacker.dev/wa-tor

import (
	"flag"
	"image"
	"image/color"
	"log"

	"lazyhacker.dev/wa-tor/internal/wator"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TileSize  = 32 // pixels width/height per tile
	MapWidth  = 32 // number of tiles horizontally
	MapHeight = 24 // number of tiles vertically
)

var (
	startFish   = flag.Int("fish", 80, "Initial # of fish.")
	startSharks = flag.Int("sharks", 40, "Initial # of sharks.")
	fsr         = flag.Int("fish-spawn-rate", 35, "fish spawn rate")
	ssr         = flag.Int("shark-spawn-rate", 45, "shark spawn rate")
	health      = flag.Int("health", 45, "# of cycles shark can go with feeding before dying.")
	width       = flag.Int("width", 32, "number of tiles horizontally (cols)")
	height      = flag.Int("height", 24, "number of tiles verticals (rows)")
	shark       *ebiten.Image
	fish        *ebiten.Image
)

// ----------- Game  -------------------
// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	world wator.Wator
}

// Set up the initial tileMap and randomly seed it with sharks and fish.
// If called again, it will reset the map and re-seed.
func (g *Game) Init(numfish, numshark, width, height int) {

	// Set up the sprites.
	ss, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/Shark.png")
	if err != nil {
		log.Fatalln("Unable to load shark image.")
	}
	shark = ss.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image)

	fs, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/Orange.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}
	fish = fs.SubImage(image.Rect(0, 0, 32, 16)).(*ebiten.Image)

	g.world = wator.Wator{}
	g.world.Init(width, height, numfish, numshark, *fsr, *ssr, *health)

}

// TileCoordinate converts the map tile index to the logical location (row, col)
// and return the pixel location (x,y).
func TileCoordinate(idx int) (float64, float64) {

	row := (idx / MapWidth) * TileSize
	col := (idx % MapWidth) * TileSize

	return float64(col), float64(row)
}

/* ------------------- Ebiten ------------------- */

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
	screen.Fill(color.RGBA{120, 180, 255, 255})
	//opts := &ebiten.DrawImageOptions{}

	// Draw each of the map tiles with the sprite of the creature (fish/shark).
	/*
		for i, t := range g.tileMap {
			opts.GeoM.Reset()
			opts.GeoM.Translate(TileCoordinate(i))
			switch v := t.(type) {
			case *Fish:
				screen.DrawImage(v.image, opts)
			case *Shark:
				screen.DrawImage(v.image, opts)
			}
		}
	*/
	ebitenutil.DebugPrint(screen, "hello!")

}

// Layout is the logical screen size which can be different from the actual
// screen size.  Ebiten will handle the scaling automatically.  For example,
// if the actual window size is 640x480, the layout can be 320x240 and Ebiten
// will scale the images so that it fits into the window.  This is also useful
// when the window can be resized but the game's logical screen size stays
// constant.
// Ebiten will fit the logical screen to fit in the window.  If the logical
// screen is small then the window, the images are scaled up.  If the logical
// screen is larger, the images are scaled down.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return TileSize * MapWidth, TileSize * MapHeight
}

func main() {
	//ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowSize(TileSize*MapWidth, TileSize*MapHeight)
	ebiten.SetWindowTitle("Wa-Tor")

	wator := &Game{}
	wator.Init(*startFish, *startSharks, MapWidth, MapHeight)

	if err := ebiten.RunGame(wator); err != nil {
		log.Fatal(err)
	}
}
