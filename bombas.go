package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

var (
    BombaInativa   = Elemento{'◎', CorVermelho, CorPadrao, false, false}
    BombaAtiva     = Elemento{'⊕', CorVermelho, CorPadrao, false, false}
    Explosao       = Elemento{'✸', CorVermelho, CorPadrao, false, false}
    
    // Canal para ativar/desativar a geração de bombas
    controleBombas = make(chan bool)
    
    // Mutex para proteger acesso às bombas
    bombaLock sync.Mutex
    
    // Mapeamento de bombas ativas
    bombasAtivas = make(map[struct{ x, y int }]bool)
)

// Inicia o sistema de bombas
func iniciarSistemaBombas(jogo *Jogo) {
    rand.Seed(time.Now().UnixNano())
    
    // Goroutine para gerar bombas aleatoriamente
    go gerarBombas(jogo)
}

// Gera bombas aleatoriamente pelo mapa
func gerarBombas(jogo *Jogo) {
    bombas := make(chan struct{ x, y int }, 10) // Canal para novas bombas
    ativo := true
    
    // Goroutine para processar explosões de bombas
    go processarExplosoes(jogo, bombas)
    
    for {
        select {
        case estado := <-controleBombas:
            ativo = estado
            
        case <-time.After(7 * time.Second):
            if !ativo {
                continue
            }
            
            // Tenta encontrar uma posição vazia para a bomba
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
                    bombaLock.Lock()
                    // Coloca uma bomba inativa
                    jogo.Mapa[y][x] = BombaInativa
                    bombaLock.Unlock()
                    
                    // Inicia o timer da bomba
                    go func(x, y int) {
                        time.Sleep(3 * time.Second)
                        
                        bombaLock.Lock()
                        // Ativa a bomba se ela ainda existir
                        if jogo.Mapa[y][x] == BombaInativa {
                            jogo.Mapa[y][x] = BombaAtiva
                            bombasAtivas[struct{ x, y int }{x, y}] = true
                            
                            // Envia para o canal de bombas
                            select {
                            case bombas <- struct{ x, y int }{x, y}:
                                // Bomba enviada para explosão
                            default:
                                // Canal cheio, tenta novamente após um tempo
                                go func(x, y int) {
                                    time.Sleep(100 * time.Millisecond)
                                    bombas <- struct{ x, y int }{x, y}
                                }(x, y)
                            }
                        }
                        bombaLock.Unlock()
                    }(x, y)
                    
                    break
                }
            }
        }
    }
}

// Processa explosões de bombas
func processarExplosoes(jogo *Jogo, bombas <-chan struct{ x, y int }) {
    for bomba := range bombas {
        // Espera um tempo antes da explosão
        time.Sleep(2 * time.Second)
        
        // Explode a bomba
        explodir(jogo, bomba.x, bomba.y)
    }
}

// Explode uma bomba e afeta a área próxima
func explodir(jogo *Jogo, x, y int) {
    bombaLock.Lock()
    defer bombaLock.Unlock()
    
    // Verifica se a bomba ainda existe na posição
    if jogo.Mapa[y][x] != BombaAtiva {
        delete(bombasAtivas, struct{ x, y int }{x, y})
        return
    }
    
    // Cria a explosão
    jogo.Mapa[y][x] = Explosao
    delete(bombasAtivas, struct{ x, y int }{x, y})
    
    // Afeta áreas próximas
    raio := 2 // raio da explosão
    for dy := -raio; dy <= raio; dy++ {
        for dx := -raio; dx <= raio; dx++ {
            // Pula o centro da explosão (já foi processado)
            if dx == 0 && dy == 0 {
                continue
            }
            
            nx, ny := x+dx, y+dy
            
            // Verifica se a posição é válida
            if ny >= 0 && ny < len(jogo.Mapa) && 
               nx >= 0 && nx < len(jogo.Mapa[ny]) {
                // Se for uma parede ou inimigo, não afeta
                if jogo.Mapa[ny][nx] == Parede || jogo.Mapa[ny][nx] == Inimigo {
                    continue
                }
                
                // Verifica se havia uma bomba e a remove
                if jogo.Mapa[ny][nx] == BombaInativa || jogo.Mapa[ny][nx] == BombaAtiva {
                    delete(bombasAtivas, struct{ x, y int }{nx, ny})
                }
                
                // Coloca explosão no local
                jogo.Mapa[ny][nx] = Explosao
                
                // Se o jogador está nessa posição, ele é atingido
                if jogo.PosX == nx && jogo.PosY == ny {
                    jogo.StatusMsg = "Você foi atingido por uma explosão!"
                }
            }
        }
    }
    
    // Limpa a explosão após um tempo
    go func() {
        time.Sleep(500 * time.Millisecond)
        bombaLock.Lock()
        for dy := -raio; dy <= raio; dy++ {
            for dx := -raio; dx <= raio; dx++ {
                nx, ny := x+dx, y+dy
                
                // Verifica se a posição é válida
                if ny >= 0 && ny < len(jogo.Mapa) && 
                   nx >= 0 && nx < len(jogo.Mapa[ny]) {
                    // Se for uma explosão, limpa
                    if jogo.Mapa[ny][nx] == Explosao {
                        jogo.Mapa[ny][nx] = Vazio
                    }
                }
            }
        }
        bombaLock.Unlock()
    }()
}

// Verifica se o jogador está sobre uma bomba e permite desarmá-la
func verificarBomba(jogo *Jogo) bool {
    if jogo.Mapa[jogo.PosY][jogo.PosX] == BombaInativa {
        bombaLock.Lock()
        jogo.Mapa[jogo.PosY][jogo.PosX] = Vazio
        jogo.StatusMsg = "Você desarmou uma bomba!"
        bombaLock.Unlock()
        return true
    }
    
    if jogo.Mapa[jogo.PosY][jogo.PosX] == BombaAtiva {
        bombaLock.Lock()
        jogo.StatusMsg = fmt.Sprintf("Cuidado! Esta bomba vai explodir!")
        bombaLock.Unlock()
        return true
    }
    
    return false
}