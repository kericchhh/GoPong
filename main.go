package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 360
	winningScore = 5
)

type Game struct {
	paddleL Paddle
	paddleR Paddle
	ball    Ball
	scoreL  int
	scoreR  int
	winner  int
}

type Paddle struct {
	x      float32
	y      float32
	width  float32
	height float32
	color  color.Color
}

type Ball struct {
	x       float32
	y       float32
	velX    float32
	velY    float32
	width   float32
	height  float32
	color   color.Color
	counter int
}

func (paddle *Paddle) Move(direction int, speed float32) {
	switch direction {
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
	maxY := float32(screenHeight) - paddle.height
	if paddle.y > maxY {
		paddle.y = maxY
	}
}

func (ball *Ball) Move() {
	ball.x += ball.velX
	ball.y += ball.velY

	if ball.y < 0 {
		ball.y = 0
		ball.velY = -ball.velY
	}
	maxY := float32(screenHeight) - ball.height
	if ball.y > maxY {
		ball.y = maxY
		ball.velY = -ball.velY
	}
}

func (ball *Ball) Serve() {
	ball.x = (float32(screenWidth) - ball.width) / 2
	ball.y = (float32(screenHeight) - ball.height) / 2
	ball.velX = 4
	ball.velY = 2
	ball.counter = 0
	if rand.IntN(2) == 0 {
		ball.velX = -ball.velX
	}
	if rand.IntN(2) == 0 {
		ball.velY = -ball.velY
	}
}

func (ball *Ball) overlaps(paddle Paddle, previousX float32) bool {
	return min(previousX, ball.x) < paddle.x+paddle.width &&
		max(previousX+ball.width, ball.x+ball.width) > paddle.x &&
		ball.y < paddle.y+paddle.height &&
		ball.y+ball.height > paddle.y
}

func (g *Game) Update() error {
	if g.winner != 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.resetMatch()
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.paddleL.Move(-1, 4.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.paddleL.Move(1, 4.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.paddleR.Move(-1, 4.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.paddleR.Move(1, 4.0)
	}
	previousX := g.ball.x
	g.ball.Move()
	hitPaddle := false
	if g.ball.velX < 0 && g.ball.overlaps(g.paddleL, previousX) {
		g.ball.x = g.paddleL.x + g.paddleL.width
		g.ball.velX = -g.ball.velX
		hitPaddle = true
	} else if g.ball.velX > 0 && g.ball.overlaps(g.paddleR, previousX) {
		g.ball.x = g.paddleR.x - g.ball.width
		g.ball.velX = -g.ball.velX
		hitPaddle = true
	}
	if hitPaddle {
		g.ball.counter++
		if g.ball.counter == 4 {
			speed := math.Hypot(float64(g.ball.velX), float64(g.ball.velY))
			factor := float32((speed + 0.5) / speed)
			g.ball.velX *= factor
			g.ball.velY *= factor
			g.ball.counter = 0
		}
	}
	if g.ball.x+g.ball.width <= 0 {
		g.scoreR++
		if g.scoreR == winningScore {
			g.winner = 2
		}
		g.ball.Serve()
	} else if g.ball.x >= screenWidth {
		g.scoreL++
		if g.scoreL == winningScore {
			g.winner = 1
		}
		g.ball.Serve()
	}
	return nil
}

func (g *Game) resetMatch() {
	g.scoreL = 0
	g.scoreR = 0
	g.winner = 0
	g.paddleL.y = (float32(screenHeight) - g.paddleL.height) / 2
	g.paddleR.y = (float32(screenHeight) - g.paddleR.height) / 2
	g.ball.Serve()
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, g.paddleL.x, g.paddleL.y, g.paddleL.width, g.paddleL.height, g.paddleL.color, false)
	vector.FillRect(screen, g.paddleR.x, g.paddleR.y, g.paddleR.width, g.paddleR.height, g.paddleR.color, false)
	vector.FillRect(screen, (screenWidth-10)/2, 0, 10, screenHeight, color.Gray{Y: 120}, false)
	vector.FillRect(screen, g.ball.x, g.ball.y, g.ball.width, g.ball.height, g.ball.color, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.scoreL), screenWidth/2-50, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.scoreR), screenWidth/2+40, 10)
	if g.winner != 0 {
		message := "Left player wins!"
		if g.winner == 2 {
			message = "Right player wins!"
		}
		ebitenutil.DebugPrintAt(screen, message, screenWidth/2-60, screenHeight/2-20)
		ebitenutil.DebugPrintAt(screen, "Press Space to restart", screenWidth/2-80, screenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		paddleL: Paddle{
			x:      0.0,
			y:      155.0,
			width:  10.0,
			height: 50.0,
			color:  color.White,
		},
		paddleR: Paddle{
			x:      630.0,
			y:      155.0,
			width:  10.0,
			height: 50.0,
			color:  color.White,
		},
		ball: Ball{
			width:   10.0,
			height:  10.0,
			color:   color.White,
			counter: 0,
		},
	}
	game.ball.Serve()
	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	ebiten.SetWindowTitle("Go Pong")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
