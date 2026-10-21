package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	paddleL Paddle
	paddleR Paddle
}

type Paddle struct {
	x float32
	y float32
	width float32
	height float32
	color color.Color
}

func (paddle *Paddle) Move(direction int, speed float32) {
	switch direction{
	case -1:
		paddle.y -= speed
	case 1:
		paddle.y += speed
	default:
		return
	}
	if paddle.y < 0 {
		paddle.y = 0
	}
	maxY := 360.0 - paddle.height
	if paddle.y > maxY {
		paddle.y = maxY
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW){
		g.paddleL.Move(-1, 4.0)
	} 
	if ebiten.IsKeyPressed(ebiten.KeyS){
		g.paddleL.Move(1, 4.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp){
		g.paddleR.Move(-1, 4.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown){
		g.paddleR.Move(1, 4.0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, g.paddleL.x, g.paddleL.y, g.paddleL.width, g.paddleL.height, g.paddleL.color, false)
	vector.FillRect(screen, g.paddleR.x, g.paddleR.y, g.paddleR.width, g.paddleR.height, g.paddleR.color, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 360
}

func main() {
	game := &Game{
		paddleL: Paddle{
			x: 0.0,
			y: 155.0,
			width: 10.0,
			height: 50.0,
			color: color.White,
		},
		paddleR: Paddle{
			x: 630.0,
			y: 155.0,
			width: 10.0,
			height: 50.0,
			color: color.White,
		},
	}
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, world!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
