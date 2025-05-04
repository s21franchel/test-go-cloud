package balancer

import (
	"log"
	"net/http"
	"time"
)

type HealthChecker struct {
	backends []*Backend
	interval time.Duration
	stopChan chan struct{}
	logger   *log.Logger
}

func NewHealthChecker(backends []*Backend, interval time.Duration, logger *log.Logger) *HealthChecker {
	return &HealthChecker{
		backends: backends,
		interval: interval,
		stopChan: make(chan struct{}),
		logger:   logger,
	}
}

// Запускает health checker
func (hc *HealthChecker) Start() {
	hc.logger.Printf("[%s] Starting health checker with interval %v", LogLevelInfo, hc.interval)

	//Интервал между проверками
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	/*
		Обрабатываются два канала:
		первый срабатывает по таймеру, запуская проверку бекэндов
		второй принимает сигнал остановки
	*/
	for {
		select {
		case <-ticker.C:
			hc.check()
		case <-hc.stopChan:
			hc.logger.Printf("[%s] Health checker stopped", LogLevelInfo)
			return
		}
	}
}

// Закрывает канал health checker
func (hc *HealthChecker) Stop() {
	close(hc.stopChan)
}

// Проверяет доступность конкретного бэкенда
func (hc *HealthChecker) isBackendAlive(b *Backend) bool {
	//Устанавка таймаута на 2 секунды для HTTP запроса
	timeout := 2 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	//Отпрака Get запроса на эндпоинт /health
	resp, err := client.Get(b.URL.String() + "/health")
	if err != nil {
		hc.logger.Printf("[%s] Backend %s health check failed: %v",
			LogLevelWarn, b.URL.Host, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		hc.logger.Printf("[%s] Backend %s returned non-OK status: %d",
			LogLevelWarn, b.URL.Host, resp.StatusCode)
		return false
	}

	//Бэкенд доступен, если статус 200 ОК
	return true
}

// Проверяет все бэкенды и обновляет их статусы
func (hc *HealthChecker) check() {
	for _, b := range hc.backends {
		alive := hc.isBackendAlive(b)
		oldStatus := b.IsAlive()
		b.SetAlive(alive)

		//Логирование если статус поменялся
		if oldStatus != alive {
			status := "UP"
			if !alive {
				status = "DOWN"
			}
			hc.logger.Printf("[%s] Backend %s status changed to %s",
				LogLevelInfo, b.URL.Host, status)
		}
	}
}
