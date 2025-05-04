package balancer

import (
	"context"
	"loabalancer/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

type LoadBalancer struct {
	rr            *RoundRobin
	server        *http.Server
	healthChecker *HealthChecker
	wg            sync.WaitGroup
	shutdownChan  chan struct{}
	logger        *log.Logger
}

// Создает новый блансировщик нагрузки
func NewLoadBalancer(config *config.Config) (*LoadBalancer, error) {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	backends := make([]*Backend, 0, len(config.Backends))

	//Инициализация каждого бэкенда и парсинг URL
	for _, rawUrl := range config.Backends {
		serverUrl, err := url.Parse(rawUrl)
		if err != nil {
			logger.Println("[%s] Failed to parse backend URL %s: %v", LogLevelError, rawUrl, err)
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		backends = append(backends, &Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		logger.Printf("[%s] Registered backend: %s", LogLevelInfo, serverUrl.String())
	}

	//Создание балансировщика
	rr := NewRoundRobin(backends)
	lb := &LoadBalancer{
		shutdownChan: make(chan struct{}),
		logger:       logger,
	}

	//Настройка HTTP маршрутизатора
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.balanceRequest)

	//Настройка HTTP сервера
	lb.server = &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	lb.rr = rr
	lb.healthChecker = NewHealthChecker(backends, 10*time.Second, logger)

	lb.logger.Printf("[%s] Load balancer initialized with %d backends", LogLevelInfo, len(backends))
	return lb, nil
}

// Запускает балансировщик нагрузки
func (lb *LoadBalancer) Start() error {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	//Запускает фоновую проверку бэкендов
	lb.wg.Add(1)
	go func() {
		defer lb.wg.Done()
		lb.healthChecker.Start()
	}()

	//Запускает HTTP-сервер, который принимает запросы и распределяет их между бэкендами
	lb.wg.Add(1)
	go func() {
		defer lb.wg.Done()
		lb.logger.Printf("[%s] Starting load balancer on %s", LogLevelInfo, lb.server.Addr)
		if err := lb.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lb.logger.Printf("[%s] Server error: %v", LogLevelError, err)
		}
	}()

	//Принимает сигналы ОС или внутренний сигнал остановки из другого места в коде
	select {
	case <-interruptChan:
		lb.logger.Printf("[%s] Received interrupt signal, shutting down...", LogLevelInfo)
	case <-lb.shutdownChan:
		lb.logger.Printf("[%s] Received shutdown request, shutting down...", LogLevelInfo)
	}

	return lb.gracefulShutdown()
}

// Завершает работу балансировщика
func (lb *LoadBalancer) gracefulShutdown() error {
	lb.logger.Printf("[%s] Starting graceful shutdown...", LogLevelInfo)

	//Завершается горутина health checker
	lb.healthChecker.Stop()
	lb.logger.Printf("[%s] Health checker stopped", LogLevelInfo)
	/*
		Создается контекст завершения рабооты сервера
		Если сервер не завершится за 30 сек - работа завершится принудительно
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Остановка HTTP сервера
	if err := lb.server.Shutdown(ctx); err != nil {
		lb.logger.Printf("[%s] Shutdown error: %v", LogLevelError, err)
		return err
	}

	//Ожидание завершения всех горутин
	lb.wg.Wait()
	lb.logger.Printf("[%s] Server gracefully stopped", LogLevelInfo)
	return nil
}

// Закрывает канал балансировщика
func (lb *LoadBalancer) Shutdown() {
	lb.logger.Printf("[%s] Shutdown requested", LogLevelInfo)
	close(lb.shutdownChan)
}

// Распределяет HTTP-запросы между доступными бэкендами
func (lb *LoadBalancer) balanceRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	//Возвращает доступный бэкенд и перенаправляет на него запрос
	peer := lb.rr.GetNext()

	if peer != nil {
		lb.logger.Printf("[%s] Routing request to %s %s -> %s", LogLevelInfo, r.Method, r.URL.Path, peer.URL.Host)
		peer.ReverseProxy.ServeHTTP(w, r)

		lb.logger.Printf("[%s] Request %s %s completed in %v", LogLevelInfo, r.Method, r.URL.Path, time.Since(start))
		return
	}

	lb.logger.Printf("[%s] No available backends for request %s %s", LogLevelInfo, r.Method, r.URL.Path)
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
