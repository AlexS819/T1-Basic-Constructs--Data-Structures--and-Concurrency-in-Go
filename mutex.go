package main

import (
	"sync"
	"time"
)

// runWithMutex — рівень 2: кожен URL перевіряється в окремій горутині.
// Спільна статистика захищена м'ютексом усередині Stats.
func runWithMutex(checker Checker, urls []string) (StatsValues, time.Duration) {
	start := time.Now()

	var wg sync.WaitGroup
	var stats Stats

	for _, url := range urls {
		// Add викликаємо до запуску горутини, інакше Wait може
		// завершитися раніше, ніж горутина почне роботу.
		wg.Add(1)

		// URL передаємо параметром, а не захоплюємо змінну циклу.
		go func(url string) {
			defer wg.Done()
			stats.Add(checker.Check(url))
		}(url)
	}

	wg.Wait()

	return stats.Snapshot(), time.Since(start)
}
