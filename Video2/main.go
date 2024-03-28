package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Close Window"
)

type game struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func initializeSDL() error {
	var err error
	var sdlFlags uint32 = sdl.INIT_EVERYTHING

	if err = sdl.Init(sdlFlags); err != nil {
		return fmt.Errorf("Error initializing SDL2: %v", err)
	}

	return err
}

func closeSDL() {
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

func (g *game) close() {
	if g != nil {
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

	g.run()
}
