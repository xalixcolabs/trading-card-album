package events

import (
	"sync"
	"time"
)

// Subscription representa una conexión SSE suscrita a eventos de un usuario.
// Data recibe los eventos; Closed se cierra cuando el hub decide eliminarla
// (p. ej. limpieza por inactividad), indicando al consumidor que cierre.
type Subscription struct {
	Data     chan []byte
	Closed   chan struct{}
	lastSeen int64
	once     sync.Once
}

// Hub es un pub/sub en memoria para notificaciones por SSE, dirigido por
// usuario. Permite publicar eventos (p. ej. "tu QR fue escaneado") y limpiar
// suscripciones inactivas de forma periódica.
type Hub struct {
	mu   sync.RWMutex
	subs map[string]map[*Subscription]struct{}
}

var globalHub = &Hub{subs: make(map[string]map[*Subscription]struct{})}

// Subscribe registra una suscripción para el usuario y la devuelve.
func Subscribe(userID string) *Subscription {
	sub := &Subscription{
		Data:     make(chan []byte, 1),
		Closed:   make(chan struct{}),
		lastSeen: time.Now().UnixMilli(),
	}
	globalHub.mu.Lock()
	defer globalHub.mu.Unlock()
	if globalHub.subs[userID] == nil {
		globalHub.subs[userID] = make(map[*Subscription]struct{})
	}
	globalHub.subs[userID][sub] = struct{}{}
	return sub
}

// Unsubscribe elimina la suscripción del usuario sin señalizar cierre.
func Unsubscribe(userID string, sub *Subscription) {
	globalHub.mu.Lock()
	defer globalHub.mu.Unlock()
	deleteSubscriptionLocked(userID, sub)
}

// Publish envía data a las suscripciones del usuario sin bloquear y marca su
// última actividad.
func Publish(userID string, data []byte) {
	globalHub.mu.Lock()
	defer globalHub.mu.Unlock()
	now := time.Now().UnixMilli()
	for sub := range globalHub.subs[userID] {
		select {
		case sub.Data <- data:
			sub.lastSeen = now
		default:
		}
	}
}

// Cleanup elimina las suscripciones inactivas (sin actividad en maxIdle) y
// cierra su canal Closed para que el consumidor SSE cierre; el cliente se
// reconecta vía EventSource.
func Cleanup(maxIdle time.Duration) {
	globalHub.cleanup(maxIdle)
}

func (h *Hub) cleanup(maxIdle time.Duration) {
	cutoff := time.Now().Add(-maxIdle).UnixMilli()

	h.mu.RLock()
	type target struct {
		userID string
		sub    *Subscription
	}
	var targets []target
	for userID, subs := range h.subs {
		for sub := range subs {
			if sub.lastSeen < cutoff {
				targets = append(targets, target{userID: userID, sub: sub})
			}
		}
	}
	h.mu.RUnlock()

	for _, t := range targets {
		t.sub.once.Do(func() { close(t.sub.Closed) })
		h.mu.Lock()
		deleteSubscriptionLocked(t.userID, t.sub)
		h.mu.Unlock()
	}
}

// StartCleanup ejecuta Cleanup cada interval en segundo plano.
func StartCleanup(interval, maxIdle time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			globalHub.cleanup(maxIdle)
		}
	}()
}

func deleteSubscriptionLocked(userID string, sub *Subscription) {
	if subs, ok := globalHub.subs[userID]; ok {
		delete(subs, sub)
		if len(subs) == 0 {
			delete(globalHub.subs, userID)
		}
	}
}