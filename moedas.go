package main

import (
    "math/rand"
    "time"
)

// Inicia a movimentação das moedas no mapa
func iniciarMoedasMoveis(jogo *Jogo) {
    rand.Seed(time.Now().UnixNano())
    
    // Primeiro, encontrar todas as moedas no mapa e iniciar uma goroutine para cada uma
    for y := 0; y < len(jogo.Mapa); y++ {
        for x := 0; x < len(jogo.Mapa[y]); x++ {
            if jogo.Mapa[y][x] == Moeda {
                // Para cada moeda, inicia uma goroutine para movê-la
                go moverMoeda(jogo, x, y)
            }
        }
    }
    
    // Também cria algumas moedas adicionais que aparecem periodicamente
    go gerarMoedasPeriodicas(jogo)
}

// Move uma moeda pelo mapa
func moverMoeda(jogo *Jogo, posX, posY int) {
    // Continua movendo até a moeda ser coletada
    for {
        // Espera um tempo antes de mover
        time.Sleep(1 * time.Second)
        
        // Direções possíveis: cima, baixo, esquerda, direita
        direcoes := []struct{ dx, dy int }{
            {0, -1}, {0, 1}, {-1, 0}, {1, 0},
        }
        
        // Escolhe uma direção aleatória
        idx := rand.Intn(len(direcoes))
        dx, dy := direcoes[idx].dx, direcoes[idx].dy
        
        // Calcula nova posição
        novoX, novoY := posX+dx, posY+dy
        
        // Verifica se pode mover para a nova posição (dentro do mapa e vazia)
        if novoY >= 0 && novoY < len(jogo.Mapa) && 
           novoX >= 0 && novoX < len(jogo.Mapa[novoY]) {
            
            // Só move se a nova posição estiver vazia ou tiver vegetação
            if jogo.Mapa[novoY][novoX] == Vazio || jogo.Mapa[novoY][novoX] == Vegetacao {
                // Remove da posição atual
                jogo.Mapa[posY][posX] = Vazio
                
                // Coloca na nova posição
                jogo.Mapa[novoY][novoX] = Moeda
                
                // Atualiza a posição
                posX, posY = novoX, novoY
            }
        }
        
        // Verificação para sair do loop se a moeda foi coletada
        // (se na posição atual não tem mais moeda)
        if jogo.Mapa[posY][posX] != Moeda {
            break
        }
    }
}

// Gera novas moedas aleatoriamente pelo mapa
func gerarMoedasPeriodicas(jogo *Jogo) {
    for {
        // Espera um tempo entre gerar novas moedas
        time.Sleep(10 * time.Second)
        
        // Tenta encontrar uma posição vazia
        maxTentativas := 20
        for i := 0; i < maxTentativas; i++ {
            // Escolhe uma posição aleatória
            y := rand.Intn(len(jogo.Mapa))
            if len(jogo.Mapa[y]) == 0 {
                continue
            }
            
            x := rand.Intn(len(jogo.Mapa[y]))
            
            // Verifica se a posição está vazia
            if jogo.Mapa[y][x] == Vazio {
                // Coloca uma moeda nova
                jogo.Mapa[y][x] = Moeda
                
                // Inicia uma goroutine para movimentar essa moeda
                go moverMoeda(jogo, x, y)
                break
            }
        }
    }
}