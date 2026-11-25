package main

import (
	"sync"
	"time"
)

// runWithChannels — рівень 3: горутини не змінюють спільну пам'ять,
// а надсилають результати в канал. Підрахунок виконує лише main.
func runWithChannels(checker Checker, urls []string) ([]CheckResult, StatsValues, time.Duration) {
	start := time.Now()

	resultsCh := make(chan CheckResult)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			resultsCh <- checker.Check(url)
		}(url)
	}

	// Окрема горутина-координатор чекає на всі перевірки
	// і тільки після цього закриває канал.
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	results := make([]CheckResult, 0, len(urls))
	var stats StatsValues

	// range завершиться автоматично, коли канал закриють.
	for res := range resultsCh {
		results = append(results, res)
		stats.Add(res)
	}

	return results, stats, time.Since(start)
}
