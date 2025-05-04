package balancer

import "sync"

type RoundRobin struct {
	backends []*Backend
	current  uint64
	mux      sync.Mutex
}

func NewRoundRobin(backends []*Backend) *RoundRobin {
	return &RoundRobin{
		backends: backends,
		current:  0,
	}
}

// Вычисляет следующий доступный бэкенд
func (r *RoundRobin) GetNext() *Backend {
	r.mux.Lock()
	defer r.mux.Unlock()

	if len(r.backends) == 0 {
		return nil
	}

	//Вычисление следующего бэкенда
	next := int(r.current % uint64(len(r.backends)))
	r.current++

	/*
		Проверка доступности выбранного бекэнда
		Возвращаем его если он доступен
	*/
	if r.backends[next].IsAlive() {
		return r.backends[next]
	}

	//Поиск следующего доступного бекэнда
	for i := next + 1; i != next; i++ {
		if i >= len(r.backends) {
			i = 0
		}
		if r.backends[i].IsAlive() {
			r.current = uint64(i + 1)
			return r.backends[i]
		}
	}

	return nil
}
