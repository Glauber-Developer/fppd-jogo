package main

import (
	"os"
	"sync"
	"time"
)

var interfaceLock sync.Mutex

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	go iniciarMoedasMoveis(&jogo)
	go iniciarFantasmas(&jogo)
	go iniciarTeleportes(&jogo)
	go iniciarSistemaBombas(&jogo)

	// Loop de renderização: 10 quadros/seg.
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			interfaceLock.Lock()
			interfaceDesenharJogo(&jogo)
			interfaceLock.Unlock()
		}
	}()

	// Entrada do jogador
	for {
		ev := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(ev, &jogo); !continuar {
			break
		}

		verificarTeleporte(&jogo)
		verificarBomba(&jogo)
		if verificarContatoFantasma(&jogo) {
			jogoMutex.Lock()
			jogo.StatusMsg = "Um fantasma te encontrou! Cuidado!"
			jogoMutex.Unlock()
		}
	}
}
