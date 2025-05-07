package main

import (
	"math/rand"
	"time"
)

func iniciarMoedasMoveis(jogo *Jogo) {
	rand.Seed(time.Now().UnixNano())

	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			if jogo.Mapa[y][x] == Moeda {
				go moverMoeda(jogo, x, y)
			}
		}
	}
	go gerarMoedasPeriodicas(jogo)
}

func moverMoeda(jogo *Jogo, posX, posY int) {
	for {
		time.Sleep(1 * time.Second)
		dirs := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		v := dirs[rand.Intn(len(dirs))]
		nx, ny := posX+v.dx, posY+v.dy

		jogoMutex.Lock()
		if ny >= 0 && ny < len(jogo.Mapa) &&
			nx >= 0 && nx < len(jogo.Mapa[ny]) &&
			(jogo.Mapa[ny][nx] == Vazio || jogo.Mapa[ny][nx] == Vegetacao) {
			jogo.Mapa[posY][posX] = Vazio
			jogo.Mapa[ny][nx] = Moeda
			posX, posY = nx, ny
		}
		collected := jogo.Mapa[posY][posX] != Moeda
		jogoMutex.Unlock()

		if collected {
			break
		}
	}
}

func gerarMoedasPeriodicas(jogo *Jogo) {
	for {
		time.Sleep(10 * time.Second)
		for i := 0; i < 20; i++ {
			y := rand.Intn(len(jogo.Mapa))
			if len(jogo.Mapa[y]) == 0 {
				continue
			}
			x := rand.Intn(len(jogo.Mapa[y]))

			jogoMutex.Lock()
			if jogo.Mapa[y][x] == Vazio {
				jogo.Mapa[y][x] = Moeda
				jogoMutex.Unlock()
				go moverMoeda(jogo, x, y)
				break
			}
			jogoMutex.Unlock()
		}
	}
}
