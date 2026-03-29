package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const year = 1969

type SafeMap struct {
	mu       sync.Mutex //not use sync.Map is optimized for two cases: when a key is written once but read many times, or when goroutines work with disjoint sets of keys.
	data     map[int]int
	accesses int
	adds     int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{data: make(map[int]int)}
}

// Update атомарно применяет fn к текущему значению ключа.
func (m *SafeMap) Update(key int, fn func(cur int) int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.data[key]
	if !ok {
		m.adds++
	}
	m.accesses++
	m.data[key] = fn(cur)
}

// Get возвращает значение и факт существования ключа.
func (m *SafeMap) Get(key int) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok
}

// AccessCount возвращает общее число вызовов Update.
func (m *SafeMap) AccessCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.accesses
}

// AddCount возвращает число созданных ключей.
func (m *SafeMap) AddCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.adds
}

// fill запускает n горутин.
func fill(m *SafeMap, n int) {
	var wg sync.WaitGroup
	for id := range n {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			keys := make([]int, 0, year*(n-1)/n+1)
			for k := 1; k <= year; k++ {
				if k%n != id {
					keys = append(keys, k)
				}
			}
			// Перемешиваем — доступ не должен быть последовательным.
			rand.Shuffle(len(keys), func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})
			for _, k := range keys {
				m.Update(k, func(cur int) int { return cur + 1 })
			}
		}(id)
	}
	wg.Wait()
}

func main() {
	m := NewSafeMap()
	fill(m, 4)

	val1, _ := m.Get(1)
	valYear, _ := m.Get(year)
	fmt.Printf("data[1]=%d data[%d]=%d | AccessCount=%d (ожид. %d) | AddCount=%d (ожид. %d)\n",
		val1, year, valYear,
		m.AccessCount(), year*3,
		m.AddCount(), year,
	)
}
