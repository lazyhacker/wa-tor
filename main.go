// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main // package lazyhacker.dev/wa-tor

import (
	"flag"
	"image"
	"image/color"
	"log"
	"strconv"

	"lazyhacker.dev/wa-tor/internal/wator"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TileSize = 32 // pixels width/height per tile
)

var (
	startFish   = flag.Int("fish", 50, "Initial # of fish.")
	startSharks = flag.Int("sharks", 20, "Initial # of sharks.")
	fsr         = flag.Int("fish-spawn-rate", 30, "fish spawn rate")
	ssr         = flag.Int("shark-spawn-rate", 50, "shark spawn rate")
	health      = flag.Int("health", 30, "# of cycles shark can go with feeding before dying.")
	width       = flag.Int("width", 32, "number of tiles horizontally (cols)")
	height      = flag.Int("height", 24, "number of tiles verticals (rows)")
)

// Tile is a position on the screen corresponding to the position of the Wa-tor
// world.
type Tile struct {
	sprite   int
	tileType int
	x, y     float64
}

// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	world        wator.Wator
	tileMap      []Tile
	sharkSprite  []*ebiten.Image
	fishSprite   []*ebiten.Image
	pause        bool
	frames       [][]Tile
	frameCounter int
	speedTPS     int
	pixelsMove   int
}

func (g *Game) AnimationSteps() int {

	return TileSize / g.pixelsMove
}

// Set up the initial tileMap and randomly seed it with sharks and fish.
// If called again, it will reset the map and re-seed.
func (g *Game) Init(numfish, numshark, width, height int) {

	g.speedTPS = 10
	g.pixelsMove = 4

	// Set up the sprites.
	g.fishSprite = make([]*ebiten.Image, g.AnimationSteps())
	g.sharkSprite = make([]*ebiten.Image, g.AnimationSteps())

	ss, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/Shark.png")
	if err != nil {
		log.Fatalln("Unable to load shark image.")
	}
	for i := 0; i < 8; i++ {
		g.sharkSprite[i] = ss.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 32)).(*ebiten.Image)
	}

	fs, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/Orange.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}
	for i := 0; i < 8; i++ {
		g.fishSprite[i] = fs.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 16)).(*ebiten.Image)
	}

	// Initialize the world.
	g.world = wator.Wator{}
	if err := g.world.Init(width, height, numfish, numshark, *fsr, *ssr, *health); err != nil {
		log.Fatalf(err.Error())
	}

	ws := g.world.Update()
	g.tileMap = g.StateToTiles(ws.Current)
	g.frames = append(g.frames, g.StateToTiles(ws.Current))
}

// StateToTiles converts the positions of the Wa-tor to the set of tiles
// that can buse used to give a visual repesentation on the screen.
func (g *Game) StateToTiles(w wator.WorldState) []Tile {

	tiles := make([]Tile, len(w))
	for i := 0; i < len(w); i++ {
		x, y := g.TileCoordinate(i)
		tiles = append(tiles, Tile{
			sprite:   0,
			tileType: w[i],
			x:        x,
			y:        y,
		})
	}
	return tiles
}

// DeltaToTiles generates intermediate tile maps to animate the movement of
// fishes and sharks so that it doesn't look like they teleported between
// tiles.
func (g *Game) DeltaToTiles(delta []wator.Delta) {

	steps := g.AnimationSteps() + 1
	for i := 1; i < steps; i++ {
		offset := float64(i * g.pixelsMove)
		tile := make([]Tile, g.world.Height*g.world.Width)
		for _, d := range delta {
			x, y := g.TileCoordinate(d.From)

			switch d.Action {
			case wator.MOVE_EAST:
				x += offset
			case wator.MOVE_WEST:
				x -= offset
			case wator.MOVE_NORTH:
				y -= offset
			case wator.MOVE_SOUTH:
				y += offset
			default:
				continue
			}

			tile[d.From] = Tile{
				sprite:   i - 1,
				tileType: d.Object,
				x:        x,
				y:        y,
			}
		}
		g.frames = append(g.frames, tile)
	}

	return
}

// TileCoordinate converts the map tile index to the logical location (row, col)
// and return the pixel location (x,y).
func (g *Game) TileCoordinate(idx int) (float64, float64) {

	row := (idx / g.world.Width) * TileSize
	col := (idx % g.world.Width) * TileSize

	return float64(col), float64(row)
}

// RenderMap will paint the world and the creatures to the screen.
func (g *Game) RenderMap(screen *ebiten.Image, m []Tile) {
	opts := &ebiten.DrawImageOptions{}
	for _, t := range m {
		opts.GeoM.Reset()
		opts.GeoM.Translate(t.x, t.y)
		switch t.tileType {
		case wator.FISH:
			screen.DrawImage(g.fishSprite[t.sprite], opts)
		case wator.SHARK:
			screen.DrawImage(g.sharkSprite[t.sprite], opts)
		}
	}
}

// Update is called by Ebiten every 'tick' based on Ticks Per Seconds (TPS).
// By default, Ebiten tries to run at 60 TPS so Update will be called every
// 1/60th of a second.  TPS can be changed with the SetTPS method.
func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.pause = !g.pause
	}

	// Don't advance the world every update because that moves too fast
	// for users to see the changes each Chronon.
	g.frameCounter++
	if g.frameCounter < g.speedTPS {
		return nil
	}
	g.frameCounter = 0

	if !g.pause {
		worldStates := g.world.Update()
		delta := worldStates.ChangeLog
		g.DeltaToTiles(delta)
		g.frames = append(g.frames, g.StateToTiles(worldStates.Current))
		g.tileMap = g.StateToTiles(worldStates.Current)
	}

	return nil
}

// Draw is called by Ebiten at the refresh rate of the display to render
// the images on the screen.  For example, when the display rate is 60Hz,
// Ebiten will call Draw 60 times per second.  When a display has a 120Hz
// refresh rate, Draw will be called twice as often as Update.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})
	ebitenutil.DebugPrint(screen, strconv.FormatUint(uint64(g.world.Chronon), 10))

	if len(g.frames) == 0 {
		g.RenderMap(screen, g.tileMap)
		return
	}
	g.RenderMap(screen, g.frames[0])
	g.frames = g.frames[1:]

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
	return TileSize * g.world.Width, TileSize * g.world.Height
}

func main() {
	flag.Parse()
	ebiten.SetWindowSize(TileSize**width, TileSize**height)
	ebiten.SetWindowTitle("Wa-Tor")
	ebiten.SetWindowResizable(true)

	wator := &Game{}
	wator.Init(*startFish, *startSharks, *width, *height)

	if err := ebiten.RunGame(wator); err != nil {
		log.Fatal(err)
	}
}
