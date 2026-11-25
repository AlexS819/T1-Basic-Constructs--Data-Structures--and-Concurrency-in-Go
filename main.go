package main

import (
	"fmt"
	"sort"
	"time"
)

// Таймаут на один запит. Якщо сервер не відповідає за цей час,
// Check поверне помилку, а програма не «зависне» назавжди.
const requestTimeout = 5 * time.Second

func main() {
	// 2 реальні сайти, один повільний endpoint (він добре показує
	// перевагу паралельності), 404 для перевірки невдалого статусу
	// і неіснуючий домен для перевірки помилки з'єднання.
	urls := []string{
		"https://go.dev",
		"https://github.com",
		"https://httpbin.org/status/404",
		"https://httpbin.org/delay/1",
		"https://invalid-fake-host-test.org",
	}

	fmt.Println("=== HTTP URL Health Checker ===")
	fmt.Printf("Адрес у списку: %d | таймаут запиту: %v\n", len(urls), requestTimeout)

	// ---------- Рівень 1 ----------
	fmt.Println("\n--- Рівень 1: послідовна перевірка ---")
	seqResults, seqStats, seqTime := runSequential(NewHTTPChecker(requestTimeout), urls)
	printResults(seqResults)
	printStats(seqStats)
	fmt.Printf("  Час рівня 1: %v\n", seqTime.Round(time.Millisecond))

	// ---------- Рівень 2 ----------
	fmt.Println("\n--- Рівень 2: горутини + WaitGroup + Mutex ---")
	mutexStats, mutexTime := runWithMutex(NewHTTPChecker(requestTimeout), urls)
	fmt.Println("  Усі горутини завершилися, статистику захищено м'ютексом.")
	printStats(mutexStats)
	fmt.Printf("  Час рівня 2: %v\n", mutexTime.Round(time.Millisecond))

	// ---------- Рівень 3 ----------
	fmt.Println("\n--- Рівень 3: збір результатів через канал ---")
	results, stats, chanTime := runWithChannels(NewHTTPChecker(requestTimeout), urls)

	// Порядок надходження з каналу непередбачуваний,
	// тому сортуємо лише для стабільного вигляду звіту.
	sort.Slice(results, func(i, j int) bool {
		return results[i].URL < results[j].URL
	})

	printResults(results)
	printFinalReport(stats, seqTime, chanTime)
}
