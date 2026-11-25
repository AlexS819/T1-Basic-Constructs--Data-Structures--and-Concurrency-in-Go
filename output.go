package main

import (
	"fmt"
	"time"
)

// printResults друкує результати перевірок у тому порядку, в якому вони прийшли.
func printResults(results []CheckResult) {
	for _, res := range results {
		if res.Err != nil {
			fmt.Printf("  -> %s: ПОМИЛКА (%v)\n", res.URL, res.Err)
			continue
		}

		fmt.Printf("  -> %s: %s (час: %v)\n",
			res.URL,
			res.StatusText(),
			res.Duration.Round(time.Millisecond),
		)
	}
}

// printStats друкує компактний підсумок.
func printStats(stats StatsValues) {
	fmt.Printf("  Разом: %d | 200 OK: %d | інші статуси: %d | помилки: %d\n",
		stats.Total, stats.OK, stats.BadStatus, stats.Errors)
}

// printFinalReport формує підсумковий звіт і порівнює два підходи.
func printFinalReport(stats StatsValues, sequential, concurrent time.Duration) {
	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("ПІДСУМКОВИЙ ЗВІТ МОНІТОРИНГУ")
	fmt.Println("------------------------------------------")
	fmt.Printf("Усього перевірено адрес: %d\n", stats.Total)
	fmt.Printf("Успішних відповідей (200 OK): %d\n", stats.OK)
	fmt.Printf("Неуспішних HTTP-статусів: %d\n", stats.BadStatus)
	fmt.Printf("Мережевих помилок: %d\n", stats.Errors)
	fmt.Println("------------------------------------------")
	fmt.Printf("Час послідовно (рівень 1): %v\n", sequential.Round(time.Millisecond))
	fmt.Printf("Час конкурентно (рівень 3): %v\n", concurrent.Round(time.Millisecond))

	if concurrent > 0 {
		speedup := float64(sequential) / float64(concurrent)
		fmt.Printf("Прискорення: %.2fx\n", speedup)
	}

	fmt.Println("==========================================")
}
