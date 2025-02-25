// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main

import (
	"flag"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	TileSize  = 32 // pixels width/height per tile
	MapWidth  = 32 // number of tiles horizontally
	MapHeight = 24 // number of tiles vertically
)

var (
	startFish   = flag.Int("fish", 700, "Initial # of fish.")
	startSharks = flag.Int("sharks", 68, "Initial # of sharks.")
	birthFish   = flag.Int("fbreed", 25, "# of cycles for fish to reproduce.")
	birthShark  = flag.Int("sbreed", 35, "# of cycles for shark to reproduce.")
	starve      = flag.Int("starve", 50, "# of cycles shark can go with feeding before dying.")
)

type Creature struct {
	image         *ebiten.Image
	height, width uint
}

type Shark struct {
	velocity uint
	health   uint
	Creature
}

type Fish struct {
	Creature
}

type Tile interface {
	Image() *ebiten.Image
}

func (f Fish) Image() *ebiten.Image {

	return f.image
}

func (s Shark) Image() *ebiten.Image {

	return s.image
}

// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	fishes  []Fish
	sharks  []Shark
	tileMap []Tile // Game map is a NxM but represented linearly.
}

func (g *Game) Init(numfish, numshark, width, height int) {

	if numfish+numshark > width*height {
		log.Fatalf("Too many creatures to fit on map!")
	}

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

	mapSize := width * height

	// Have a sequence of numbers from 0 to mapSize correspond to
	// locations on the tileMap that isn't occupied.
	sequence := make([]uint, mapSize)
	for i := 0; i < mapSize; i++ {
		sequence[i] = uint(i)
	}
	// Shuffle the sequence
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})

	g.tileMap = make([]Tile, mapSize)
	g.fishes = make([]Fish, numfish)
	g.sharks = make([]Shark, numshark)

	rand.Seed(time.Now().UnixNano())

	// seed the fishes
	for i := 0; i < len(g.fishes); i++ {

		if len(sequence) == 0 {
			log.Println("No more tiles left on map to place FISH.")
			break
		}

		g.fishes[i] = Fish{
			Creature{
				image:  fish,
				height: 16,
				width:  32,
			},
		}

		t := sequence[0]        // get the tile number
		sequence = sequence[1:] // remove the tile number since it's been taken

		g.tileMap[t] = g.fishes[i]
	}

	// seed the sharks
	for i := 0; i < len(g.sharks); i++ {

		if len(sequence) == 0 {
			log.Println("No more tiles left on map to place SHARK.")
			break
		}

		g.sharks[i] = Shark{
			2,
			uint(*starve),
			Creature{
				image:  shark,
				height: 32,
				width:  32,
			},
		}

		t := sequence[0]        // get the tile number
		sequence = sequence[1:] // remove the tile number since it's been taken

		g.tileMap[t] = g.sharks[i]
	}
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

	for i, c := range g.tileMap {
		if c != nil {
			opts.GeoM.Reset()
			opts.GeoM.Translate(TileCoordinate(i))
			screen.DrawImage(c.Image(), opts)
		}
	}
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
	ebiten.SetWindowSize(TileSize*MapWidth, TileSize*MapHeight)
	ebiten.SetWindowTitle("Wa-Tor")

	wator := &Game{}
	wator.Init(*startFish, *startSharks, MapWidth, MapHeight)

	if err := ebiten.RunGame(wator); err != nil {
		log.Fatal(err)
	}
}

func TileCoordinate(idx int) (float64, float64) {

	row := (idx / MapWidth) * TileSize
	col := (idx % MapWidth) * TileSize

	return float64(col), float64(row)
}
