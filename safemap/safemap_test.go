package safemap_test

import (
	"math/rand/v2"
	"sync"
	"testing"

	"testik/safemap"
)

// year = 1969 — год первой высадки людей на Луну (миссия Apollo 11, 20 июля 1969 г.)
// Источник: https://ru.wikipedia.org/wiki/Аполлон-11
const year = 1969

// fill запускает n горутин.
//
// Горутина id пропускает ключи где k%n == id —
// каждый ключ обходят ровно n-1 горутин.
// Порядок обхода внутри каждой горутины случаен.
//
// Инварианты после завершения (n=4):
//
//	Get(k)       == 3         для всех k в [1, year]
//	AccessCount  == year*3
//	AddCount     == year
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
			// Каждая горутина перемешивает независимо — горутины конкурируют
			// за одни и те же ключи в произвольном порядке.
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

func TestSafeMap(t *testing.T) {
	t.Parallel()
	m := safemap.New()
	fill(m, 4)

	for k := 1; k <= year; k++ {
		v, ok := m.Get(k)
		if !ok {
			t.Errorf("ключ %d отсутствует", k)
			continue
		}
		if v != 3 {
			t.Errorf("ключ %d: ожидали 3, получили %d", k, v)
		}
	}

	if want := int64(year * 3); m.AccessCount() != want {
		t.Errorf("AccessCount: ожидали %d, получили %d", want, m.AccessCount())
	}
	if want := int64(year); m.AddCount() != want {
		t.Errorf("AddCount: ожидали %d, получили %d", want, m.AddCount())
	}
}
