// Package safemap предоставляет map[int]int, защищённую мьютексом.
//
// Выбор sync.Mutex вместо sync.Map обоснован паттерном использования:
// множество горутин многократно пишут по одним и тем же ключам —
// это худший сценарий для sync.Map, которая оптимизирована под
// "write once, read many" или непересекающиеся наборы ключей.
package safemap

import (
	"sync"
	"sync/atomic"
)

// SafeMap — конкурентно-безопасная map[int]int со счётчиками обращений и добавлений.
type SafeMap struct {
	mu       sync.Mutex
	data     map[int]int
	accesses atomic.Int64 // счётчик обращений; atomic позволяет читать без захвата mu
	adds     atomic.Int64 // счётчик добавлений; аналогично
}

func New() *SafeMap {
	return &SafeMap{data: make(map[int]int)}
}

// Update атомарно применяет fn к текущему значению ключа (0 если ключа нет).
// При первом обращении к ключу инкрементирует AddCount.
// Всегда инкрементирует AccessCount.
// Инкремент передаётся через fn — отдельного метода инкремента нет.
func (m *SafeMap) Update(key int, fn func(cur int) int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.data[key]
	if !ok {
		m.adds.Add(1)
	}
	m.accesses.Add(1)
	m.data[key] = fn(cur)
}

// Get возвращает значение и факт существования ключа.
func (m *SafeMap) Get(key int) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok
}

// AccessCount возвращает число вызовов Update. Читается без захвата mu.
func (m *SafeMap) AccessCount() int64 {
	return m.accesses.Load()
}

// AddCount возвращает число созданных ключей. Читается без захвата mu.
func (m *SafeMap) AddCount() int64 {
	return m.adds.Load()
}
