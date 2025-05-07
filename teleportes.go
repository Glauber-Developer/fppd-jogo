package main

import (
	"math/rand"
	"time"
)

var (
	TeleporteAtivo = Elemento{'⧫', CorAzul, CorPadrao, false, false}
)

func iniciarTeleportes(jogo *Jogo) {
	rand.Seed(time.Now().UnixNano())

	// Gera os teletransportes aleatoriamente no mapa
	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			if jogo.Mapa[y][x] == Vazio && rand.Float32() < 0.05 { // 5% de chance
				jogoMutex.Lock()
				jogo.Mapa[y][x] = TeleporteAtivo
				jogoMutex.Unlock()

				go controlarTeleporte(jogo, x, y)
			}
		}
	}
}

// Controla o funcionamento do teleporte
func controlarTeleporte(jogo *Jogo, x, y int) {
	// O teleporte fica ativo por 15 segundos
	time.Sleep(15 * time.Second)

	// Quando o tempo passar, o teleporte desaparece
	jogoMutex.Lock()
	if jogo.Mapa[y][x] == TeleporteAtivo {
		jogo.Mapa[y][x] = Vazio
	}
	jogoMutex.Unlock()
}

// Verifica se o jogador está próximo de um teleporte
func verificarTeleporte(jogo *Jogo) {
	jogoMutex.RLock()
	defer jogoMutex.RUnlock()

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := jogo.PosX+dx, jogo.PosY+dy
			if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
				if jogo.Mapa[ny][nx] == TeleporteAtivo {
					jogoMutex.Lock()
					jogo.StatusMsg = "Você foi teletransportado!"
					jogo.Mapa[ny][nx] = Vazio
					jogo.Mapa[jogo.PosY][jogo.PosX] = Vazio
					jogo.PosX, jogo.PosY = nx, ny
					jogo.Mapa[ny][nx] = Personagem
					jogoMutex.Unlock()
					return
				}
			}
		}
	}
}
