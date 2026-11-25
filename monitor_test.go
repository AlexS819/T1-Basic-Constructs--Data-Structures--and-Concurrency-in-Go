package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Перевіряємо, що обидва конкурентні підходи дають однакову статистику
// на однаковому наборі адрес.
func TestRunWithMutexAndChannels(t *testing.T) {
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okServer.Close()

	notFoundServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer notFoundServer.Close()

	// Слухаємо вільний порт і одразу закриваємо його,
	// щоб отримати стабільну помилку «connection refused».
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("не вдалося відкрити тестовий порт: %v", err)
	}
	badURL := "http://" + listener.Addr().String()
	listener.Close()

	urls := []string{okServer.URL, notFoundServer.URL, badURL}

	mutexStats, mutexTime := runWithMutex(NewHTTPChecker(2*time.Second), urls)
	checkStats(t, mutexStats)

	if mutexTime <= 0 {
		t.Errorf("час роботи рівня 2 має бути додатним, отримали %v", mutexTime)
	}

	results, channelStats, channelTime := runWithChannels(NewHTTPChecker(2*time.Second), urls)
	checkStats(t, channelStats)

	if channelTime <= 0 {
		t.Errorf("час роботи рівня 3 має бути додатним, отримали %v", channelTime)
	}

	if len(results) != len(urls) {
		t.Fatalf("отримали %d результатів, хотіли %d", len(results), len(urls))
	}

	seen := make(map[string]int)
	for _, res := range results {
		seen[res.URL]++
	}

	for _, url := range urls {
		if seen[url] != 1 {
			t.Errorf("URL %s зустрічається %d разів, хотіли 1", url, seen[url])
		}
	}
}

func checkStats(t *testing.T, got StatsValues) {
	t.Helper()

	if got.Total != 3 {
		t.Errorf("Total = %d, хотіли 3", got.Total)
	}
	if got.OK != 1 {
		t.Errorf("OK = %d, хотіли 1", got.OK)
	}
	if got.BadStatus != 1 {
		t.Errorf("BadStatus = %d, хотіли 1", got.BadStatus)
	}
	if got.Errors != 1 {
		t.Errorf("Errors = %d, хотіли 1", got.Errors)
	}
}
