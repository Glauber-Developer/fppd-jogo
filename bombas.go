package main

import (
	"math/rand"
	"time"
)

var (
	Bomba = Elemento{'💣', CorVermelho, CorPadrao, false, false}
)

func iniciarBombas(jogo *Jogo) {
	rand.Seed(time.Now().UnixNano())

	// Gera as bombas aleatoriamente no mapa
	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			if jogo.Mapa[y][x] == Vazio && rand.Float32() < 0.05 { // 5% de chance
				jogoMutex.Lock()
				jogo.Mapa[y][x] = Bomba
				jogoMutex.Unlock()

				go controlarBomba(jogo, x, y)
			}
		}
	}
}

// Controla a explosão da bomba com um tempo limite
func controlarBomba(jogo *Jogo, x, y int) {
	// Espera 10 segundos para explodir
	time.Sleep(10 * time.Second)

	// Quando a bomba explodir, substitui por vazio e afeta a área ao redor
	jogoMutex.Lock()
	if jogo.Mapa[y][x] == Bomba {
		jogo.Mapa[y][x] = Vazio
		// Afeta as células ao redor (se estiverem no mapa)
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				nx, ny := x+dx, y+dy
				if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
					if jogo.Mapa[ny][nx] == Personagem {
						// A bomba afeta o personagem
						jogo.StatusMsg = "Você foi atingido pela bomba!"
					}
					// Pode remover outros elementos, dependendo da lógica desejada
					jogo.Mapa[ny][nx] = Vazio
				}
			}
		}
	}
	jogoMutex.Unlock()
}
