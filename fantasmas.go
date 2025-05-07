package main

import (
	"math"
	"math/rand"
	"time"
)

var (
	Fantasma          = Elemento{'ゴ', CorVermelho, CorPadrao, false, false}
	jogadorPosicao    = make(chan struct{ PosX, PosY int }, 1)
	controleFantasmas = make(chan string)
)

func iniciarFantasmas(jogo *Jogo) {
	rand.Seed(time.Now().UnixNano())

	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			if jogo.Mapa[y][x] == Vazio && rand.Float32() < 0.01 {
				jogoMutex.Lock()
				jogo.Mapa[y][x] = Fantasma
				jogoMutex.Unlock()
				go controlarFantasma(jogo, x, y)
			}
		}
	}
	go monitorarJogador(jogo)
}

func monitorarJogador(jogo *Jogo) {
	ux, uy := jogo.PosX, jogo.PosY
	for {
		jogoMutex.RLock()
		px, py := jogo.PosX, jogo.PosY
		jogoMutex.RUnlock()

		if px != ux || py != uy {
			ux, uy = px, py
			select {
			case jogadorPosicao <- struct{ PosX, PosY int }{px, py}:
			default:
				<-jogadorPosicao
				jogadorPosicao <- struct{ PosX, PosY int }{px, py}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func controlarFantasma(jogo *Jogo, x, y int) {
	posX, posY := x, y
	const raioDeteccao = 8
	modo := "patrulha"

	for {
		select {
		case cmd := <-controleFantasmas:
			if cmd == "parar" {
				return
			}
		case pos := <-jogadorPosicao:
			dist := math.Hypot(float64(posX-pos.PosX), float64(posY-pos.PosY))
			if dist <= raioDeteccao {
				modo = "perseguindo"
			} else {
				modo = "patrulha"
			}
		case <-time.After(800 * time.Millisecond):
			var dx, dy int
			if modo == "perseguindo" {
				if posX < jogo.PosX {
					dx = 1
				} else if posX > jogo.PosX {
					dx = -1
				}
				if posY < jogo.PosY {
					dy = 1
				} else if posY > jogo.PosY {
					dy = -1
				}
				if dx != 0 && dy != 0 && rand.Intn(2) == 0 {
					dx = 0
				}
			} else {
				dirs := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
				v := dirs[rand.Intn(len(dirs))]
				dx, dy = v.dx, v.dy
			}

			nx, ny := posX+dx, posY+dy
			jogoMutex.Lock()
			if nx >= 0 && nx < len(jogo.Mapa[0]) &&
				ny >= 0 && ny < len(jogo.Mapa) &&
				(jogo.Mapa[ny][nx] == Vazio || jogo.Mapa[ny][nx] == Vegetacao) {

				jogo.Mapa[posY][posX] = Vazio
				jogo.Mapa[ny][nx] = Fantasma
				posX, posY = nx, ny
			}
			jogoMutex.Unlock()
		}
	}
}

func verificarContatoFantasma(jogo *Jogo) bool {
	jogoMutex.RLock()
	defer jogoMutex.RUnlock()

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := jogo.PosX+dx, jogo.PosY+dy
			if ny >= 0 && ny < len(jogo.Mapa) &&
				nx >= 0 && nx < len(jogo.Mapa[ny]) &&
				jogo.Mapa[ny][nx] == Fantasma {
				return true
			}
		}
	}
	return false
}
