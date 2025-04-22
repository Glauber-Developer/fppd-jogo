// main.go - Loop principal do jogo
package main

import (
    "os"
    "sync"
)

var interfaceLock sync.Mutex

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}


    // Inicia os elementos autônomos concorrentes
    go iniciarMoedasMoveis(&jogo)
    go iniciarFantasmas(&jogo)
    go iniciarTeleportes(&jogo)
    go iniciarSistemaBombas(&jogo)
	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	 // Loop principal de entrada
	 for {
        evento := interfaceLerEventoTeclado()
        if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
            break
        }
        
        // Verifica interações com elementos concorrentes
        verificarTeleporte(&jogo)
        verificarBomba(&jogo)
        
        if verificarContatoFantasma(&jogo) {
            jogo.StatusMsg = "Um fantasma te encontrou! Cuidado!"
        }
        
        interfaceLock.Lock()
        interfaceDesenharJogo(&jogo)
        interfaceLock.Unlock()
    }
}