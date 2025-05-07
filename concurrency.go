package main

import "sync"

// Protege TODA a estrutura de jogo (mapa, posição,
// status, coleta de itens etc.).  Usamos RWMutex
// para permitir múltiplas leituras concorrentes
// enquanto apenas uma escrita é garantida.
var jogoMutex sync.RWMutex
