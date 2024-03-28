package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Background"
)

type game struct {
	window          *sdl.Window
	renderer        *sdl.Renderer
	backgroundImage *sdl.Texture
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

	return err
}

func closeSDL() {
	img.Quit()
	sdl.Quit()
}

func newGame() *game {
	g := &game{}

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

	return err
}

func (g *game) loadMedia() error {
	var err error

	if g.backgroundImage, err = img.LoadTexture(g.renderer, "images/background.png"); err != nil {
		return fmt.Errorf("Error loading Texture: %v", err)
	}

	return err
}

func (g *game) close() {
	if g != nil {
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
					}
				}
			}
		}

		g.renderer.Clear()

		g.renderer.Copy(g.backgroundImage, nil, nil)

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
