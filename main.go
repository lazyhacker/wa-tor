// Wa-tor is an implementation of the Wa-Tor simulation A.K. Dewdney presented
// in Scientific America in 1984.  This project is an exercise to learn Ebiten,
// a 2D game engine for Go.
package main // package lazyhacker.dev/wator

import (
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"golang.org/x/image/font/basicfont"
	"lazyhacker.dev/wa-tor/internal/wator"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	TileSize = 32 // pixels width/height per tile
	EAST     = iota
	WEST
	NORTH
	SOUTH
)

var (
	startFish   = flag.Int("fish", 50, "Initial # of fish.")
	startSharks = flag.Int("sharks", 20, "Initial # of sharks.")
	fsr         = flag.Int("fish-spawn-rate", 30, "fish spawn rate")
	ssr         = flag.Int("shark-spawn-rate", 50, "shark spawn rate")
	health      = flag.Int("health", 30, "# of cycles shark can go with feeding before dying.")
	width       = flag.Int("width", 16, "number of tiles horizontally (cols)")
	height      = flag.Int("height", 12, "number of tiles verticals (rows)")
)

// Frame is a position on the screen corresponding to the position of the Wa-tor
// world.
type Frame struct {
	sprite    int
	tileType  int
	x, y      float64
	direction int
}

// Game holds the game state.  For Ebiten, this needs to be an ebiten.Game
// interface.
type Game struct {
	world            wator.Wator
	currentScreen    []Frame
	sharkSprite      []*ebiten.Image
	fishSprite       []*ebiten.Image
	pause            bool
	frames           [][]Frame
	ctickCounter     int
	drawFrameCounter int
	tpsPerChronon    int
	tpsPerFrame      int
	pixelsMove       int
	startFish        int
	startShark       int
	width            int
	height           int
}

func (g *Game) AnimationSteps() int {

	return TileSize / g.pixelsMove
}

// Set up the initial tileMap and randomly seed it with sharks and fish.
// If called again, it will reset the map and re-seed.
func (g *Game) Init(numfish, numshark, width, height int) {

	g.frames = nil
	g.startFish = numfish
	g.startShark = numshark
	g.width = width
	g.height = height
	g.tpsPerChronon = 60
	g.pixelsMove = 4
	g.tpsPerFrame = 8

	g.ctickCounter = 0
	g.drawFrameCounter = 0
	g.pause = true
	// Set up the sprites.
	g.fishSprite = make([]*ebiten.Image, g.AnimationSteps()*2)
	g.sharkSprite = make([]*ebiten.Image, g.AnimationSteps()*2)

	ss, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/Shark.png")
	if err != nil {
		log.Fatalln("Unable to load shark image.")
	}
	for i := 0; i < g.AnimationSteps(); i++ {
		g.sharkSprite[i] = ss.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 32)).(*ebiten.Image)
	}

	ss_r, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Shark - 32x32/SharkReverse.png")
	if err != nil {
		log.Fatalln("Unable to load reverse shark image.")
	}
	for i := 0; i < g.AnimationSteps(); i++ {
		g.sharkSprite[i+g.AnimationSteps()] = ss_r.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 32)).(*ebiten.Image)
	}

	fs, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/Orange.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}
	for i := 0; i < g.AnimationSteps(); i++ {
		g.fishSprite[i] = fs.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 16)).(*ebiten.Image)
	}
	fs_r, _, err := ebitenutil.NewImageFromFile("assets/spearfishing/Sprites/Fish3 - 32x16/OrangeReverse.png")
	if err != nil {
		log.Fatalln("Unable to load fish image.")
	}
	for i := 0; i < g.AnimationSteps(); i++ {
		g.fishSprite[i+g.AnimationSteps()] = fs_r.SubImage(image.Rect(i*TileSize, 0, i*TileSize+TileSize, 16)).(*ebiten.Image)
	}

	// Initialize the world.
	g.world = wator.Wator{}
	if err := g.world.Init(width, height, numfish, numshark, *fsr, *ssr, *health); err != nil {
		log.Fatal(err.Error())
	}
	//ws := g.world.Update()
	//g.currentScreen = g.StateToFrame(ws.Current)
}

// StateToFrame converts the positions of the Wa-tor to the set of tiles
// that can buse used to give a visual repesentation on the screen.
func (g *Game) StateToFrame(w wator.WorldState) []Frame {

	tiles := make([]Frame, len(w))
	for i := 0; i < len(w); i++ {
		x, y := g.TileCoordinate(i)
		tiles = append(tiles, Frame{
			sprite:   0,
			tileType: w[i],
			x:        x,
			y:        y,
		})
	}
	return tiles
}

// DeltaToFrames generates intermediate tile maps to animate the movement of
// fishes and sharks so that it doesn't look like they teleported between
// tiles.
func (g *Game) DeltaToFrames(delta []wator.Delta) [][]Frame {

	var dir int
	steps := g.AnimationSteps()
	var frames [][]Frame
	for i := 0; i < steps; i++ {
		offset := float64(i * g.pixelsMove)
		frame := make([]Frame, g.world.Height*g.world.Width)
		for _, d := range delta {
			x, y := g.TileCoordinate(d.From)

			switch d.Action {
			case wator.MOVE_EAST:
				x += offset
				dir = EAST
			case wator.MOVE_WEST:
				x -= offset
				dir = WEST
			case wator.MOVE_NORTH:
				y -= offset
				dir = NORTH
			case wator.MOVE_SOUTH:
				y += offset
				dir = SOUTH
			default:
				continue
			}

			frame[d.From] = Frame{
				sprite:    i,
				tileType:  d.Object,
				x:         x,
				y:         y,
				direction: dir,
			}
		}
		frames = append(frames, frame)
	}

	return frames
}

// TileCoordinate converts the map tile index to the logical location (row, col)
// and return the pixel location (x,y).
func (g *Game) TileCoordinate(idx int) (float64, float64) {

	row := (idx / g.world.Width) * TileSize
	col := (idx % g.world.Width) * TileSize

	return float64(col), float64(row)
}

// DrawFrame will paint the world and the creatures to the screen.
func (g *Game) DrawFrame(screen *ebiten.Image, m []Frame) {
	opts := &ebiten.DrawImageOptions{}

	for _, t := range m {
		spriteIdx := t.sprite
		opts.GeoM.Reset()
		opts.GeoM.Translate(t.x, t.y)
		if t.direction == WEST {
			spriteIdx += g.AnimationSteps()
		}
		switch t.tileType {
		case wator.FISH:
			screen.DrawImage(g.fishSprite[spriteIdx], opts)
		case wator.SHARK:
			screen.DrawImage(g.sharkSprite[spriteIdx], opts)
		}
	}
}

func (g *Game) ShowOptionsScreen(screen *ebiten.Image) {

	msg := "<SPACE> to begin/resume.\nR to restart.\nQ to quit."
	text.Draw(screen, msg, basicfont.Face7x13, 20, 50, color.Black)
}

// Update is called by Ebiten every 'tick' based on Ticks Per Seconds (TPS).
// By default, Ebiten tries to run at 60 TPS so Update will be called every
// 1/60th of a second.  TPS can be changed with the SetTPS method.
func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.pause = !g.pause
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Init(g.startFish, g.startShark, g.width, g.height)
		g.currentScreen = nil
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if !g.pause {
		g.ctickCounter++
		g.drawFrameCounter++

		if g.drawFrameCounter%g.tpsPerFrame == 0 {
			g.drawFrameCounter = 0
			if len(g.frames) > 0 {
				g.currentScreen = g.frames[0]
				g.frames = g.frames[1:]
			}
		}

		// Don't advance the world every update because that moves too fast
		// for users to see the changes each Chronon.
		if g.ctickCounter%g.tpsPerChronon == 0 {
			g.ctickCounter = 0
			// Advance the world 1 chronon and get the delta
			worldStates := g.world.Update()
			delta := worldStates.ChangeLog
			for _, f := range g.DeltaToFrames(delta) {
				g.frames = append(g.frames, f)
			}
		}
	}

	return nil
}

// Draw is called by Ebiten at the refresh rate of the display to render
// the images on the screen.  For example, when the display rate is 60Hz,
// Ebiten will call Draw 60 times per second.  When a display has a 120Hz
// refresh rate, Draw will be called twice as often as Update.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})
	if g.pause {
		g.ShowOptionsScreen(screen)
	}
	for x := 0; x < g.world.Width*TileSize; x += TileSize {
		ebitenutil.DrawLine(screen, float64(x), 0, float64(x), float64(g.world.Height*TileSize), color.White)
	}
	for y := 0; y < g.world.Height*TileSize; y += TileSize {
		ebitenutil.DrawLine(screen, 0, float64(y), float64(g.world.Width*TileSize), float64(y), color.White)
	}
	ebitenutil.DebugPrint(screen, strconv.FormatUint(uint64(g.world.Chronon), 10))

	g.DrawFrame(screen, g.currentScreen)

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
	//ebiten.SetWindowSize(TileSize**width, TileSize**height)
	ebiten.SetWindowSize(1024, 768)

	ebiten.SetWindowTitle("Wa-Tor")
	ebiten.SetWindowResizable(true)

	game := &Game{}
	game.Init(*startFish, *startSharks, *width, *height)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
