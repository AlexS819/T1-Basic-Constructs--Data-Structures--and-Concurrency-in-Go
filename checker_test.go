package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPCheckerStatusCodes(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{name: "200 OK", status: http.StatusOK},
		{name: "404 Not Found", status: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte("test"))
			}))
			defer server.Close()

			checker := NewHTTPChecker(2 * time.Second)
			res := checker.Check(server.URL)

			if res.Err != nil {
				t.Fatalf("не очікували помилку: %v", res.Err)
			}
			if res.StatusCode != tt.status {
				t.Errorf("отримали статус %d, хотіли %d", res.StatusCode, tt.status)
			}
			if res.URL != server.URL {
				t.Errorf("URL у результаті %q, хотіли %q", res.URL, server.URL)
			}
			if res.Duration <= 0 {
				t.Errorf("тривалість має бути додатною, отримали %v", res.Duration)
			}
		})
	}
}

func TestHTTPCheckerConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close() // тепер порт закритий, запит має впасти

	checker := NewHTTPChecker(time.Second)
	res := checker.Check(url)

	if res.Err == nil {
		t.Fatal("очікували помилку з'єднання, але її немає")
	}
	if res.StatusCode != 0 {
		t.Errorf("при помилці статус має бути 0, отримали %d", res.StatusCode)
	}
}

func TestHTTPCheckerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	checker := NewHTTPChecker(40 * time.Millisecond)
	res := checker.Check(server.URL)

	if res.Err == nil {
		t.Fatal("очікували помилку таймауту, але її немає")
	}
}

func TestCheckResultOK(t *testing.T) {
	tests := []struct {
		name string
		res  CheckResult
		want bool
	}{
		{name: "200 без помилки", res: CheckResult{StatusCode: 200}, want: true},
		{name: "404", res: CheckResult{StatusCode: 404}, want: false},
		{name: "помилка з'єднання", res: CheckResult{Err: http.ErrHandlerTimeout}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.OK(); got != tt.want {
				t.Errorf("OK() = %v, хотіли %v", got, tt.want)
			}
		})
	}
}
