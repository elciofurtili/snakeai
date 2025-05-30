package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 400
	screenHeight = 400
	gridSize     = 20
)

type Point struct {
	X, Y int
}

type Game struct {
	snake     []Point
	direction Point
	food      Point
	score     int
	tick      int
	gameOver  bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	// Movimentação com teclado
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction.Y != 1 {
		g.direction = Point{0, -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction.Y != -1 {
		g.direction = Point{0, 1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction.X != 1 {
		g.direction = Point{-1, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction.X != -1 {
		g.direction = Point{1, 0}
	}

	g.tick++
	if g.tick%5 != 0 {
		return nil // Controla velocidade
	}

	// Nova posição da cabeça
	head := g.snake[0]
	newHead := Point{head.X + g.direction.X, head.Y + g.direction.Y}

	// Colisão com parede
	if newHead.X < 0 || newHead.Y < 0 || newHead.X >= screenWidth/gridSize || newHead.Y >= screenHeight/gridSize {
		g.gameOver = true
		return nil
	}

	// Colisão com o próprio corpo
	for _, s := range g.snake {
		if s == newHead {
			g.gameOver = true
			return nil
		}
	}

	// Move a cobra
	g.snake = append([]Point{newHead}, g.snake...)

	// Comer comida
	if newHead == g.food {
		g.score++
		g.spawnFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1] // remove cauda
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fundo
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Comida
	ebitenutil.DrawRect(screen, float64(g.food.X*gridSize), float64(g.food.Y*gridSize), gridSize, gridSize, color.RGBA{255, 0, 0, 255})

	// Cobra
	for i, p := range g.snake {
		col := color.RGBA{0, 255, 0, 255}
		if i == 0 {
			col = color.RGBA{0, 180, 0, 255} // cabeça
		}
		ebitenutil.DrawRect(screen, float64(p.X*gridSize), float64(p.Y*gridSize), gridSize, gridSize, col)
	}

	// Game Over
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER!", 150, 180)
	}

	// Score
	ebitenutil.DebugPrintAt(screen, "Score: "+itoa(g.score), 10, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenWidth, screenHeight
}

func (g *Game) spawnFood() {
	for {
		x := rand.Intn(screenWidth / gridSize)
		y := rand.Intn(screenHeight / gridSize)
		food := Point{x, y}
		collision := false
		for _, s := range g.snake {
			if s == food {
				collision = true
				break
			}
		}
		if !collision {
			g.food = food
			break
		}
	}
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		snake:     []Point{{5, 5}},
		direction: Point{1, 0},
		score:     0,
	}
	game.spawnFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game in Go (Ebiten)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
