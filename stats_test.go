package main

import (
	"errors"
	"sync"
	"testing"
)

func TestStatsConcurrentAdd(t *testing.T) {
	const n = 300

	var stats Stats
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			switch i % 3 {
			case 0:
				stats.Add(CheckResult{StatusCode: 200})
			case 1:
				stats.Add(CheckResult{StatusCode: 404})
			default:
				stats.Add(CheckResult{Err: errors.New("тестова помилка")})
			}
		}(i)
	}

	wg.Wait()
	got := stats.Snapshot()

	if got.Total != n {
		t.Errorf("Total = %d, хотіли %d", got.Total, n)
	}
	if got.OK != n/3 {
		t.Errorf("OK = %d, хотіли %d", got.OK, n/3)
	}
	if got.BadStatus != n/3 {
		t.Errorf("BadStatus = %d, хотіли %d", got.BadStatus, n/3)
	}
	if got.Errors != n/3 {
		t.Errorf("Errors = %d, хотіли %d", got.Errors, n/3)
	}
}
