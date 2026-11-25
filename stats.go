package main

import "sync"

// StatsValues — «знімок» статистики без м'ютекса.
// Його зручно повертати з функцій і передавати в інші горутини,
// бо копіювання sync.Mutex — це помилка.
type StatsValues struct {
	Total     int
	OK        int
	BadStatus int
	Errors    int
}

// Add враховує один результат.
// Метод не потокобезпечний: у рівні 2 він викликається під захистом
// Stats.mu, а в рівні 3 — лише з головної горутини.
func (s *StatsValues) Add(r CheckResult) {
	s.Total++

	switch {
	case r.Err != nil:
		s.Errors++
	case r.OK():
		s.OK++
	default:
		s.BadStatus++
	}
}

// Stats — спільний лічильник для рівня 2.
// М'ютекс захищає поля values від одночасного запису з різних горутин.
type Stats struct {
	mu     sync.Mutex
	values StatsValues
}

// Add безпечно оновлює статистику з будь-якої горутини.
func (s *Stats) Add(r CheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values.Add(r)
}

// Snapshot повертає копію даних, тому після виходу з метода
// блокування вже не потрібне.
func (s *Stats) Snapshot() StatsValues {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.values
}
