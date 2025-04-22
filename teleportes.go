package main

import (
    "math/rand"
    "sync"
    "time"
)

var (
    TeleporteDesativado = Elemento{'○', CorCinzaEscuro, CorPadrao, false, false}
    TeleporteAtivado    = Elemento{'●', CorVerde, CorPadrao, false, false}
    
    // Canal para ativação dos teleportes
    canalAtivacao = make(chan int, 1)
    
    // Mutex para proteger acesso aos teleportes
    teleporteLock sync.Mutex
    
    // Lista de teleportes no mapa
    teleportes []struct{ x, y int }
)

// Inicia os teleportes no mapa
func iniciarTeleportes(jogo *Jogo) {
    // Limpa a lista de teleportes
    teleportes = nil
    
    // Coloca teleportes em posições aleatórias do mapa
    numTeleportes := 5 // Número de teleportes no mapa
    
    for i := 0; i < numTeleportes; i++ {
        for {
            // Escolhe uma posição aleatória
            y := rand.Intn(len(jogo.Mapa))
            if len(jogo.Mapa[y]) == 0 {
                continue
            }
            
            x := rand.Intn(len(jogo.Mapa[y]))
            
            // Verifica se a posição está vazia
            if jogo.Mapa[y][x] == Vazio {
                teleporteLock.Lock()
                // Coloca um teleporte desativado
                jogo.Mapa[y][x] = TeleporteDesativado
                // Adiciona à lista de teleportes
                teleportes = append(teleportes, struct{ x, y int }{x, y})
                teleporteLock.Unlock()
                break
            }
        }
    }
    
    // Inicia a goroutine para controlar os teleportes
    go controlarTeleportes(jogo)
}

// Controla a ativação e desativação dos teleportes
func controlarTeleportes(jogo *Jogo) {
    for {
        select {
        // Ativa um teleporte específico
        case id := <-canalAtivacao:
            if id >= 0 && id < len(teleportes) {
                teleporteLock.Lock()
                pos := teleportes[id]
                // Ativa o teleporte
                jogo.Mapa[pos.y][pos.x] = TeleporteAtivado
                teleporteLock.Unlock()
                
                // Espera 3 segundos com timeout
                select {
                case <-time.After(3 * time.Second):
                    // Desativa o teleporte após o timeout
                    teleporteLock.Lock()
                    jogo.Mapa[pos.y][pos.x] = TeleporteDesativado
                    teleporteLock.Unlock()
                }
            }
        
        // A cada 5 segundos, ativa um teleporte aleatório
        case <-time.After(5 * time.Second):
            if len(teleportes) > 0 {
                id := rand.Intn(len(teleportes))
                teleporteLock.Lock()
                pos := teleportes[id]
                jogo.Mapa[pos.y][pos.x] = TeleporteAtivado
                teleporteLock.Unlock()
                
                // Agenda a desativação após 3 segundos
                go func(id int) {
                    time.Sleep(3 * time.Second)
                    teleporteLock.Lock()
                    if id < len(teleportes) {
                        pos := teleportes[id]
                        jogo.Mapa[pos.y][pos.x] = TeleporteDesativado
                    }
                    teleporteLock.Unlock()
                }(id)
            }
        }
    }
}

// Verifica se o jogador está em um teleporte e o teleporta
func verificarTeleporte(jogo *Jogo) bool {
    for i, pos := range teleportes {
        if pos.x == jogo.PosX && pos.y == jogo.PosY && jogo.Mapa[pos.y][pos.x] == TeleporteAtivado {
            // Teleporta o jogador para outro teleporte aleatório
            if len(teleportes) > 1 {
                // Escolhe um destino diferente do atual
                destinos := make([]int, 0, len(teleportes)-1)
                for j := range teleportes {
                    if j != i {
                        destinos = append(destinos, j)
                    }
                }
                
                destino := destinos[rand.Intn(len(destinos))]
                destPos := teleportes[destino]
                
                // Teleporta o jogador
                teleporteLock.Lock()
                jogo.PosX, jogo.PosY = destPos.x, destPos.y
                jogo.StatusMsg = "Você foi teleportado!"
                teleporteLock.Unlock()
                
                // Desativa os teleportes usados
                go func() {
                    teleporteLock.Lock()
                    jogo.Mapa[pos.y][pos.x] = TeleporteDesativado
                    jogo.Mapa[destPos.y][destPos.x] = TeleporteDesativado
                    teleporteLock.Unlock()
                }()
                
                return true
            }
        }
    }
    return false
}