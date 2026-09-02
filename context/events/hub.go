package events

import "sync"

// Hub es un pub/sub en memoria para notificaciones por SSE. Permite
// publicar eventos dirigidos a un usuario (p. ej. "tu QR fue escaneado")
// y que sus conexiones SSE los reciban.
type Hub struct {
	mu   sync.RWMutex
	subs map[string]map[chan []byte]struct{}
}

var globalHub = &Hub{subs: make(map[string]map[chan []byte]struct{})}

// Subscribe registra un canal para recibir eventos del usuario y lo devuelve.
func Subscribe(userID string) chan []byte {
	ch := make(chan []byte, 1)
	globalHub.mu.Lock()
	defer globalHub.mu.Unlock()
	if globalHub.subs[userID] == nil {
		globalHub.subs[userID] = make(map[chan []byte]struct{})
	}
	globalHub.subs[userID][ch] = struct{}{}
	return ch
}

// Unsubscribe elimina el canal del usuario.
func Unsubscribe(userID string, ch chan []byte) {
	globalHub.mu.Lock()
	defer globalHub.mu.Unlock()
	if subs, ok := globalHub.subs[userID]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(globalHub.subs, userID)
		}
	}
}

// Publish envía data a todos los suscriptores del usuario sin bloquear.
func Publish(userID string, data []byte) {
	globalHub.mu.RLock()
	defer globalHub.mu.RUnlock()
	for ch := range globalHub.subs[userID] {
		select {
		case ch <- data:
		default:
		}
	}
}