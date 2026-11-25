package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// CheckResult — результат однієї перевірки URL.
// Якщо запит не вдався, StatusCode залишається 0, а причина лежить у Err.
type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

// OK повертає true, якщо запит завершився без помилки і сервер віддав саме 200.
func (r CheckResult) OK() bool {
	return r.Err == nil && r.StatusCode == http.StatusOK
}

// StatusText — людське представлення статусу для виводу в консоль.
func (r CheckResult) StatusText() string {
	if r.Err != nil {
		return "помилка"
	}

	text := http.StatusText(r.StatusCode)
	if text == "" {
		return fmt.Sprintf("%d", r.StatusCode)
	}
	return fmt.Sprintf("%d %s", r.StatusCode, text)
}

// Checker — контракт перевірки адреси.
// Будь-який тип із методом Check(string) CheckResult автоматично
// задовольняє цей інтерфейс: ключове слово implements у Go не потрібне.
type Checker interface {
	Check(url string) CheckResult
}

// HTTPChecker — реалізація Checker через стандартний net/http.
type HTTPChecker struct {
	client *http.Client
}

// Перевірка на етапі компіляції: *HTTPChecker справді задовольняє Checker.
// Якщо сигнатуру Check випадково змінити, цей рядок не скомпілюється.
var _ Checker = (*HTTPChecker)(nil)

// NewHTTPChecker створює клієнт з окремим таймаутом і власним пулом з'єднань.
// Окремий Transport потрібен, щоб різні рівні лабораторної не ділили
// вже відкриті з'єднання і заміри часу були чеснішими.
func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	return &HTTPChecker{
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

// Check виконує GET-запит і повертає структуру з результатом.
// Час вимірюється навколо всього запиту, включаючи очікування відповіді.
func (c *HTTPChecker) Check(url string) CheckResult {
	start := time.Now()

	res := CheckResult{URL: url}

	resp, err := c.client.Get(url)
	res.Duration = time.Since(start)

	if err != nil {
		res.Err = err
		return res
	}
	defer resp.Body.Close()

	res.StatusCode = resp.StatusCode

	// Тіло потрібно дочитати до кінця, інакше з'єднання не повернеться
	// в пул повторного використання. Сам вміст сторінки нам не потрібен.
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		res.Err = fmt.Errorf("не вдалося прочитати відповідь: %w", err)
	}

	return res
}
