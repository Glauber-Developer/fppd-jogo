package main

import (
    "math"
    "math/rand"
    "sync"
    "time"
)

var (
    Fantasma = Elemento{'ゴ', CorVermelho, CorPadrao, false, false}
    
    // Mutex para proteger acesso ao mapa
    mapaLock sync.Mutex
    
    // Canal para notificar os fantasmas da posição do jogador
    jogadorPosicao = make(chan struct{PosX, PosY int}, 1)
    
    // Canal para controlar os fantasmas (parar/iniciar)
    controleFantasmas = make(chan string)
)

// Inicia os fantasmas no mapa
func iniciarFantasmas(jogo *Jogo) {
    rand.Seed(time.Now().UnixNano())
    
    // Encontra posições iniciais para os fantasmas
    for y := 0; y < len(jogo.Mapa); y++ {
        for x := 0; x < len(jogo.Mapa[y]); x++ {
            // Coloca fantasmas em posições vazias com 1% de probabilidade
            if jogo.Mapa[y][x] == Vazio && rand.Float32() < 0.01 {
                mapaLock.Lock()
                jogo.Mapa[y][x] = Fantasma
                mapaLock.Unlock()
                
                // Inicia uma goroutine para cada fantasma
                go controlarFantasma(jogo, x, y)
            }
        }
    }
    
    // Goroutine para monitorar a posição do jogador e notificar os fantasmas
    go monitorarJogador(jogo)
}

// Monitora a posição do jogador e notifica os fantasmas
func monitorarJogador(jogo *Jogo) {
    ultimoX, ultimoY := jogo.PosX, jogo.PosY
    
    for {
        // Verifica se a posição do jogador mudou
        if jogo.PosX != ultimoX || jogo.PosY != ultimoY {
            ultimoX, ultimoY = jogo.PosX, jogo.PosY
            
            // Atualiza o canal com a nova posição
            select {
            case jogadorPosicao <- struct{PosX, PosY int}{jogo.PosX, jogo.PosY}:
                // Canal atualizado
            default:
                // Canal cheio, remove o valor antigo
                <-jogadorPosicao
                jogadorPosicao <- struct{PosX, PosY int}{jogo.PosX, jogo.PosY}
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// Controla o comportamento de um fantasma
func controlarFantasma(jogo *Jogo, x, y int) {
    posX, posY := x, y
    
    // Distância máxima para perseguir o jogador
    const raioDeteccao = 8
    modo := "patrulha" // modos: "patrulha" ou "perseguindo"
    
    for {
        select {
        case cmd := <-controleFantasmas:
            if cmd == "parar" {
                return
            }
        case pos := <-jogadorPosicao:
            // Calcula a distância até o jogador
            dist := math.Sqrt(float64((posX-pos.PosX)*(posX-pos.PosX) + (posY-pos.PosY)*(posY-pos.PosY)))
            if dist <= raioDeteccao {
                modo = "perseguindo"
            } else {
                modo = "patrulha"
            }
        case <-time.After(800 * time.Millisecond):
            // Movimento do fantasma após timeout
            var dx, dy int
            
            if modo == "perseguindo" {
                // Movimento em direção ao jogador
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
                
                // Se está bloqueado em uma direção, tenta a outra
                if dx != 0 && dy != 0 {
                    // Escolhe aleatoriamente qual direção tentar primeiro
                    if rand.Intn(2) == 0 {
                        dx = 0
                    } else {
                        dy = 0
                    }
                }
            } else {
                // Movimento aleatório de patrulha
                direcoes := []struct{ dx, dy int }{
                    {0, -1}, {0, 1}, {-1, 0}, {1, 0},
                }
                idx := rand.Intn(len(direcoes))
                dx, dy = direcoes[idx].dx, direcoes[idx].dy
            }
            
            novoX, novoY := posX+dx, posY+dy
            
            // Verifica se pode mover para a nova posição
            if novoX >= 0 && novoX < len(jogo.Mapa[0]) && 
               novoY >= 0 && novoY < len(jogo.Mapa) {
                
                mapaLock.Lock()
                if jogo.Mapa[novoY][novoX] == Vazio || jogo.Mapa[novoY][novoX] == Vegetacao {
                    // Remove da posição atual
                    jogo.Mapa[posY][posX] = Vazio
                    
                    // Coloca na nova posição
                    jogo.Mapa[novoY][novoX] = Fantasma
                    
                    // Atualiza a posição
                    posX, posY = novoX, novoY
                }
                mapaLock.Unlock()
            }
        }
    }
}

// Verifica se o jogador está em contato com um fantasma
func verificarContatoFantasma(jogo *Jogo) bool {
    // Verifica as posições adjacentes
    for dy := -1; dy <= 1; dy++ {
        for dx := -1; dx <= 1; dx++ {
            nx, ny := jogo.PosX+dx, jogo.PosY+dy
            
            // Verifica se a posição é válida
            if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
                if jogo.Mapa[ny][nx] == Fantasma {
                    return true
                }
            }
        }
    }
    return false
}