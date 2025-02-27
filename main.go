// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math/rand"
	"strconv"
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
	startFish   = flag.Int("fish", 80, "Initial # of fish.")
	startSharks = flag.Int("sharks", 40, "Initial # of sharks.")
	fsr         = flag.Int("fish-spawn-rate", 35, "fish spawn rate")
	ssr         = flag.Int("shark-spawn-rate", 45, "shark spawn rate")
	health      = flag.Int("health", 45, "# of cycles shark can go with feeding before dying.")
	shark       *ebiten.Image
	fish        *ebiten.Image
)

const (
	TILE_NONE = iota
	TILE_SHARK
	TILE_FISH
)

// TIle represents a place on the map.  It has an image and descriptor of
// what is on the tile.
type Tile struct {
	image    *ebiten.Image
	tileType int
}

type Creature struct {
	image         *ebiten.Image
	height, width uint
	age           int
}

func (c *Creature) SetAge(a int) {
	c.age = a
}

func (c *Creature) Age() int {
	return c.age
}

// ----------- Sharks -------------------
type Shark struct {
	health int
	Creature
}

func NewShark() *Shark {
	return &Shark{
		*health,
		Creature{
			image:  shark,
			height: 32,
			width:  32,
		},
	}
}

func (s *Shark) Spawn() bool {

	if s.age%*ssr == 0 && s.age > 0 {
		return true
	}
	return false
}

func (s *Shark) New() MapItem {

	return NewShark()
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

func (f *Fish) Spawn() bool {

	if f.age%*fsr == 0 && f.age > 0 {
		return true
	}
	return false
}

func (f *Fish) New() MapItem {

	return NewFish()
}

// MapItem holds an tile item: FISH, SHARK or NONE
type MapItem interface {
	Age() int
	SetAge(int)
	Spawn() bool
	New() MapItem
}

// ----------- Game  -------------------
// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	tileMap []MapItem // Game map is a NxM but represented linearly.
	Chronon int
}

// Set up the initial tileMap and randomly seed it with sharks and fish.
// If called again, it will reset the map and re-seed.
func (g *Game) Init(numfish, numshark, width, height int) {

	if numfish+numshark > width*height {
		log.Fatalf("Too many creatures to fit on map!")
	}

	if *health > *ssr {
		log.Fatalf("shark spawn rate is faster then health rate so shark will always spawn befor hunger.")
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

	g.tileMap = make([]MapItem, mapSize)

	// seed fishes on the tile map.
	for i := 0; i < *startFish; i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place FISH.")
			break
		}

		p := sequence.Next()
		g.tileMap[p] = NewFish()
	}

	// seed the sharks on the tile map.
	for i := 0; i < *startSharks; i++ {

		if sequence.Length() == 0 {
			log.Println("No more tiles left on map to place SHARK.")
			break
		}

		p := sequence.Next()
		g.tileMap[p] = NewShark()
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

// Adjacent returns up, down, left, right tile locations from the position.
func Adjacent(pos int) []int {

	totalTiles := MapWidth * MapHeight
	up := pos - MapWidth
	down := pos + MapWidth
	left := pos - 1
	right := pos + 1

	// Check if needs to loop around to the bottom of the map.
	if up < 0 {
		up += totalTiles
	}

	// Check to see if needs to loop around to the top of the map.
	if down >= totalTiles {
		down -= totalTiles
	}

	// Check if it needs to go to wrap around to the end of the row.
	if (right % MapWidth) == 0 {
		right -= MapWidth
	}

	// Check if it needs to wrap around to the start of the row.
	if (left % MapWidth) < 0 {
		left += MapWidth
	}

	return []int{up, down, left, right}
}

// PickPosition randomly picks the element from the given slice.
func PickPosition(curr int, numbers []int) int {

	if len(numbers) == 0 {
		return curr
	}
	rand.Seed(time.Now().UnixNano())
	return numbers[rand.Intn(len(numbers))]
}

/* ------------------- Ebiten ------------------- */

// Update is called by Ebiten every 'tick' based on Ticks Per Seconds (TPS).
// By default, Ebiten tries to run at 60 TPS so Update will be called every
// 1/60th of a second.  TPS can be changed with the SetTPS method.
func (g *Game) Update() error {

	//if g.Chronon%20 != 0 {
	//		g.Chronon++
	//		return nil
	//}

	for i, tile := range g.tileMap {
		if tile == nil {
			continue
		}
		// Handle movement and creature-specific behaviors.
		adjacent := Adjacent(i)
		var openTiles []int
		var newPos int
		switch c := tile.(type) {
		case *Fish:
			// Fish can only move to non-occupied squares.
			for j := 0; j < len(adjacent); j++ {
				if g.tileMap[adjacent[j]] == nil {
					openTiles = append(openTiles, adjacent[j])
				}
			}

			newPos = PickPosition(i, openTiles)
		case *Shark:

			// If shark doesn't eat, it dies.
			(*c).health--
			if (*c).health == 0 {
				g.tileMap[i] = nil
				continue
			}

			// Shark cannot move to tiles that have other sharks
			for j := 0; j < len(adjacent); j++ {
				if v, ok := g.tileMap[adjacent[j]].(*Shark); !ok || v == nil {
					openTiles = append(openTiles, adjacent[j])
				}
			}

			newPos = PickPosition(i, openTiles)
			if _, ok := g.tileMap[newPos].(*Fish); ok {
				(*c).health = *health + 1
				g.tileMap[newPos] = nil
			}
		}

		if newPos != i {

			g.tileMap[newPos], g.tileMap[i] = g.tileMap[i], g.tileMap[newPos]
			if tile.Spawn() {
				g.tileMap[i] = tile.New()
			}

		}

		tile.SetAge(tile.Age() + 1)

	}
	g.Chronon++
	return nil
}

// Draw is called by Ebiten at the refresh rate of the display to render
// the images on the screen.  For example, when the display rate is 60Hz,
// Ebiten will call Draw 60 times per second.  When a display has a 120Hz
// refresh rate, Draw will be called twice as often as Update.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := &ebiten.DrawImageOptions{}

	// Draw each of the map tiles with the sprite of the creature (fish/shark).
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
	ebitenutil.DebugPrint(screen, strconv.Itoa(g.Chronon))

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
