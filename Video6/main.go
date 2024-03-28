package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Moving Text and Icon"
)

type game struct {
	window          *sdl.Window
	renderer        *sdl.Renderer
	backgroundImage *sdl.Texture
	fontSize        int
	fontColor       sdl.Color
	textImage       *sdl.Texture
	textRect        sdl.Rect
	textVel         int32
	textXVel        int32
	textYVel        int32
	rng             *rand.Rand
}

func initializeSDL() error {
	var err error
	var sdlFlags uint32 = sdl.INIT_EVERYTHING
	imgFlags := img.INIT_PNG

	if err = sdl.Init(sdlFlags); err != nil {
		return fmt.Errorf("Error initializing SDL2: %v", err)
	}

	if err = img.Init(imgFlags); err != nil {
		return fmt.Errorf("Error initializing SDL_image: %v", err)
	}

	if err = ttf.Init(); err != nil {
		return fmt.Errorf("Error initializing SDL_ttf: %v", err)
	}

	return err
}

func closeSDL() {
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func newGame() *game {
	g := &game{}

	g.fontSize = 80
	g.fontColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	g.textVel = 3
	g.textXVel = g.textVel
	g.textYVel = g.textVel

	return g
}

func (g *game) init() error {
	var err error

	if g.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED, windowWidth, windowHeight, sdl.WINDOW_SHOWN); err != nil {
		return fmt.Errorf("Error creating Window: %v", err)
	}

	if g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fmt.Errorf("Error creating Renderer: %v", err)
	}

	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	iconSurf, err := img.Load("images/Go-logo.png")
	if err != nil {
		return fmt.Errorf("Error loading Surface: %v", err)
	}
	defer iconSurf.Free()

	g.window.SetIcon(iconSurf)

	return err
}

func (g *game) loadMedia() error {
	var err error

	if g.backgroundImage, err = img.LoadTexture(g.renderer, "images/background.png"); err != nil {
		return fmt.Errorf("Error loading Texture: %v", err)
	}

	font, err := ttf.OpenFont("fonts/freesansbold.ttf", g.fontSize)
	if err != nil {
		return fmt.Errorf("Error creating Font: %v", err)
	}
	defer font.Close()

	fontSurf, err := font.RenderUTF8Blended("SDL", g.fontColor)
	if err != nil {
		return fmt.Errorf("Error creating text Surface: %v", err)
	}
	defer fontSurf.Free()

	g.textRect.W = fontSurf.W
	g.textRect.H = fontSurf.H

	if g.textImage, err = g.renderer.CreateTextureFromSurface(fontSurf); err != nil {
		return fmt.Errorf("Error creating Texture from Surface: %v", err)
	}

	return err
}

func (g *game) randColor() {
	g.renderer.SetDrawColor(uint8(g.rng.Intn(256)),
		uint8(g.rng.Intn(256)), uint8(g.rng.Intn(256)), 255)
}

func (g *game) updateText() {
	g.textRect.X += g.textXVel
	g.textRect.Y += g.textYVel
	if g.textRect.X < 0 {
		g.textXVel = g.textVel
	} else if (g.textRect.X + g.textRect.W) > windowWidth {
		g.textXVel = -g.textVel
	}
	if g.textRect.Y < 0 {
		g.textYVel = g.textVel
	} else if (g.textRect.Y + g.textRect.H) > windowHeight {
		g.textYVel = -g.textVel
	}
}

func (g *game) close() {
	if g != nil {
		g.textImage.Destroy()
		g.textImage = nil
		g.backgroundImage.Destroy()
		g.backgroundImage = nil
		g.renderer.Destroy()
		g.renderer = nil
		g.window.Destroy()
		g.window = nil
	}
}

func (g *game) run() {
	for true {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_ESCAPE:
						return
					case sdl.K_SPACE:
						g.randColor()
					}
				}
			}
		}

		g.updateText()

		g.renderer.Clear()

		g.renderer.Copy(g.backgroundImage, nil, nil)
		g.renderer.Copy(g.textImage, nil, &g.textRect)

		g.renderer.Present()

		sdl.Delay(16)
	}
}

func main() {
	var err error

	defer closeSDL()
	if err = initializeSDL(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	g := newGame()
	defer g.close()
	if err = g.init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err = g.loadMedia(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	g.run()
}
