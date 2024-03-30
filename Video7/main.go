package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Sound Effects and Music"
)

type game struct {
	window          *sdl.Window
	renderer        *sdl.Renderer
	backgroundImage *sdl.Texture
	textImage       *sdl.Texture
	fontSize        int
	fontColor       sdl.Color
	textRect        sdl.Rect
	textVel         int32
	textXVel        int32
	textYVel        int32
	spriteImage     *sdl.Texture
	spriteRect      sdl.Rect
	spriteVel       int32
	keystate        []uint8
	goSound         *mix.Chunk
	sdlSound        *mix.Chunk
	music           *mix.Music
	rng             *rand.Rand
}

func initializeSDL() error {
	var err error
	var sdlFlags uint32 = sdl.INIT_EVERYTHING
	imgFlags := img.INIT_PNG
	mixFlags := mix.INIT_OGG

	if err = sdl.Init(sdlFlags); err != nil {
		return fmt.Errorf("Error initializing SDL2: %v", err)
	}

	if err = img.Init(imgFlags); err != nil {
		return fmt.Errorf("Error initializing SDL_image: %v", err)
	}

	if err = ttf.Init(); err != nil {
		return fmt.Errorf("Error initializing SDL_ttf: %v", err)
	}

	if err = mix.Init(mixFlags); err != nil {
		return fmt.Errorf("Error initializing SDL_mixer: %v", err)
	}

	if err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT,
		mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		return fmt.Errorf("Error opening Audio: %v", err)
	}

	return err
}

func closeSDL() {
	mix.CloseAudio()
	mix.Quit()
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
	g.spriteVel = 5

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

	g.keystate = sdl.GetKeyboardState()

	return err
}

func (g *game) loadMedia() error {
	var err error

	if g.backgroundImage, err = img.LoadTexture(g.renderer, "images/background.png"); err != nil {
		return fmt.Errorf("Error loading Texture: %v", err)
	}

	font, err := ttf.OpenFont("fonts/freesansbold.ttf", g.fontSize)
	if err != nil {
		return fmt.Errorf("Error opening Font: %v", err)
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

	if g.spriteImage, err = img.LoadTexture(g.renderer, "images/Go-logo.png"); err != nil {
		return fmt.Errorf("Error loading Texture: %v", err)
	}

	if _, _, g.spriteRect.W, g.spriteRect.H, err = g.spriteImage.Query(); err != nil {
		return fmt.Errorf("Error querying Texture: %v", err)
	}

	if g.goSound, err = mix.LoadWAV("sounds/Go.ogg"); err != nil {
		return fmt.Errorf("Error loading Chunk: %v", err)
	}

	if g.sdlSound, err = mix.LoadWAV("sounds/SDL.ogg"); err != nil {
		return fmt.Errorf("Error loading Chunk: %v", err)
	}

	if g.music, err = mix.LoadMUS("music/freesoftwaresong-8bit.ogg"); err != nil {
		return fmt.Errorf("Error loading Music: %v", err)
	}

	if err = g.music.Play(-1); err != nil {
		return fmt.Errorf("Error playing Music: %v", err)
	}

	return err
}

func (g *game) randColor() {
	g.renderer.SetDrawColor(uint8(g.rng.Intn(256)),
		uint8(g.rng.Intn(256)), uint8(g.rng.Intn(256)), 255)
	g.goSound.Play(-1, 0)
}

func (g *game) updateText() {
	g.textRect.X += g.textXVel
	g.textRect.Y += g.textYVel

	if g.textRect.X < 0 {
		g.textXVel = g.textVel
		g.sdlSound.Play(-1, 0)
	} else if (g.textRect.X + g.textRect.W) > windowWidth {
		g.textXVel = -g.textVel
		g.sdlSound.Play(-1, 0)
	}
	if g.textRect.Y < 0 {
		g.textYVel = g.textVel
		g.sdlSound.Play(-1, 0)
	} else if (g.textRect.Y + g.textRect.H) > windowHeight {
		g.textYVel = -g.textVel
		g.sdlSound.Play(-1, 0)
	}
}

func (g *game) updateSprite() {
	if g.keystate[sdl.SCANCODE_LEFT] == 1 || g.keystate[sdl.SCANCODE_A] == 1 {
		g.spriteRect.X -= g.spriteVel
	}
	if g.keystate[sdl.SCANCODE_RIGHT] == 1 || g.keystate[sdl.SCANCODE_D] == 1 {
		g.spriteRect.X += g.spriteVel
	}
	if g.keystate[sdl.SCANCODE_UP] == 1 || g.keystate[sdl.SCANCODE_W] == 1 {
		g.spriteRect.Y -= g.spriteVel
	}
	if g.keystate[sdl.SCANCODE_DOWN] == 1 || g.keystate[sdl.SCANCODE_S] == 1 {
		g.spriteRect.Y += g.spriteVel
	}
}

func (g *game) pauseMusic() {
	if mix.PausedMusic() {
		mix.ResumeMusic()
	} else {
		mix.PauseMusic()
	}
}

func (g *game) close() {
	if g != nil {
		mix.HaltMusic()
		mix.HaltChannel(-1)
		g.music.Free()
		g.music = nil
		g.sdlSound.Free()
		g.sdlSound = nil
		g.goSound.Free()
		g.goSound = nil
		g.spriteImage.Destroy()
		g.spriteImage = nil
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
					switch e.Keysym.Scancode {
					case sdl.SCANCODE_ESCAPE:
						return
					case sdl.SCANCODE_SPACE:
						g.randColor()
					case sdl.SCANCODE_M:
						g.pauseMusic()
					}
				}
			}
		}

		g.updateText()
		g.updateSprite()

		g.renderer.Clear()

		g.renderer.Copy(g.backgroundImage, nil, nil)
		g.renderer.Copy(g.textImage, nil, &g.textRect)
		g.renderer.Copy(g.spriteImage, nil, &g.spriteRect)

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
