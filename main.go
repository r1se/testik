package main

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"testik/safemap"
)

// year = 1969 — год первой высадки людей на Луну (миссия Apollo 11, 20 июля 1969 г.)
// Источник: https://ru.wikipedia.org/wiki/Аполлон-11
const year = 1969

func fill(m *safemap.SafeMap, n int) {
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
	m := safemap.New()

	// Демо: 4 горутины, каждый ключ получает значение 3.
	// Логика распределения — в тесте (safemap/safemap_test.go).
	fill(m, 4)

	val1, _ := m.Get(1)
	valYear, _ := m.Get(year)
	fmt.Printf("data[1]=%d  data[%d]=%d  AccessCount=%d  AddCount=%d\n",
		val1, year, valYear, m.AccessCount(), m.AddCount())
}
