package main

import "time"

// runSequential — рівень 1: звичайний цикл, запити виконуються один за одним.
func runSequential(checker Checker, urls []string) ([]CheckResult, StatsValues, time.Duration) {
	start := time.Now()

	results := make([]CheckResult, 0, len(urls))
	var stats StatsValues

	for _, url := range urls {
		res := checker.Check(url)
		results = append(results, res)
		stats.Add(res)
	}

	return results, stats, time.Since(start)
}
