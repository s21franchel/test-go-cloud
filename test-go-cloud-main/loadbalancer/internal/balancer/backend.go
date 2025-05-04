package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// Устанавливает стаус доступности бэкенда
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// Возвращает параметр alive
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}
