// personagem.go
package main

import "fmt"

func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMutex.Lock()
		if jogo.Mapa[ny][nx].coletavel {
			jogo.StatusMsg = "Você coletou uma moeda! +1 ponto"
		}
		jogoMutex.Unlock()

		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)

		jogoMutex.Lock()
		jogo.PosX, jogo.PosY = nx, ny
		jogoMutex.Unlock()
	} else {
		jogoMutex.Lock()
		jogo.StatusMsg = "Movimento inválido!"
		jogoMutex.Unlock()
	}
}

func personagemInteragir(jogo *Jogo) {
	jogoMutex.Lock()
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
	jogoMutex.Unlock()
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
