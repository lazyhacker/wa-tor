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
	startFish   = flag.Int("fish", 50, "Initial # of fish.")
	startSharks = flag.Int("sharks", 15, "Initial # of sharks.")
	fsr         = flag.Int("fish-spawn-rate", 25, "fish spawn rate")
	ssr         = flag.Int("shark-spawn-rate", 35, "shark spawn rate")
	health      = flag.Int("health", 50, "# of cycles shark can go with feeding before dying.")
	shark       *ebiten.Image
	fish        *ebiten.Image
)

type CreatureMap interface {
	Image() *ebiten.Image
}

type Creature struct {
	image         *ebiten.Image
	height, width uint
	spawn         int
}

func (c Creature) Image() *ebiten.Image {

	return c.image
}

// ----------- Sharks -------------------
type Shark struct {
	velocity uint
	health   uint
	Creature
}

func NewShark() *Shark {
	return &Shark{
		2,
		uint(*health),
		Creature{
			image:  shark,
			height: 32,
			width:  32,
			spawn:  *ssr,
		},
	}
}

// ----------- Fish -------------------
type Fish struct {
	Creature
}

func NewFish() *Fish {
	return &Fish{
		Creature{
			image:  fish,
			height: 16,
			width:  32,
		},
	}
}

// ----------- Game  -------------------
type Tile interface {
	Image() *ebiten.Image
}

// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	fishes      []*Fish
	sharks      []*Shark
	tileMap     []Tile // Game map is a NxM but represented linearly.
	creatureMap []CreatureMap
}

// Set up the initial tileMap and randomly seed it with sharks and fish.
// If called again, it will reset the map and re-seed.
func (g *Game) Init(numfish, numshark, width, height int) {

	if numfish+numshark > width*height {
		log.Fatalf("Too many creatures to fit on map!")
	}

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

	mapSize := width * height

	// Have a sequence of numbers from 0 to mapSize correspond to
	// locations on the tileMap that isn't occupied.
	sequence := Sequence{}
	sequence.Init(mapSize)

	//g.tileMap = make([]Tile, mapSize)
	g.creatureMap = make([]CreatureMap, mapSize)
	g.fishes = make([]*Fish, numfish)
	g.sharks = make([]*Shark, numshark)

	// seed fishes on the tile map.
	for i := 0; i < len(g.fishes); i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place FISH.")
			break
		}

		g.fishes[i] = NewFish()
		t := sequence.Next()
		g.creatureMap[t] = g.fishes[i]
	}

	// seed the sharks on the tile map.
	for i := 0; i < len(g.sharks); i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place SHARK.")
			break
		}

		g.sharks[i] = NewShark()
		t := sequence.Next()
		g.creatureMap[t] = g.sharks[i]
	}
}

type Sequence struct {
	sequence []int
}

// Init creates a slice of sequential integers and then shuffle them.
func (s *Sequence) Init(size int) {
	s.sequence = make([]int, size)
	for i := 0; i < size; i++ {
		s.sequence[i] = int(i)
	}
	rand.Seed(time.Now().UnixNano())

	// Shuffle the sequence
	rand.Shuffle(len(s.sequence), func(i, j int) {
		s.sequence[i], s.sequence[j] = s.sequence[j], s.sequence[i]
	})
}

// Next return the next value in the sequence.
func (s *Sequence) Next() int {
	n := s.sequence[0]          // get the tile number
	s.sequence = s.sequence[1:] // remove the tile number since it's been taken

	return n
}

func (s *Sequence) Length() int {

	return len(s.sequence)
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
	opts := &ebiten.DrawImageOptions{}

	// Draw each of the map tiles with the sprite of the creature (fish/shark).
	for i, c := range g.creatureMap {
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
