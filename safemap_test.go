package main

import "testing"

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()
	fill(m, 4)

	// Каждый ключ должен иметь значение 3.
	for k := 1; k <= year; k++ {
		v, _ := m.Get(k)
		if v != 3 {
			t.Errorf("ключ %d: ожидали 3, получили %d", k, v)
		}
	}

	if want := year * 3; m.AccessCount() != want {
		t.Errorf("AccessCount: ожидали %d, получили %d", want, m.AccessCount())
	}
	if m.AddCount() != year {
		t.Errorf("AddCount: ожидали %d, получили %d", year, m.AddCount())
	}
}
