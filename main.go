  /** Example
 *  window size = 1 hour
 *  0:00:00 - Put("foo", 42)
 *  0:05:00 - Put("bar", 76)
 *  0:50:00 - Get("foo") => 42, Get("bar") => 76, GetAverage() => 59
 *  1:02:00 - Get("foo") => nothing, Get("bar") => 76, GetAverage() => 76
 *  1:06:00 - Get("foo") => nothing, Get("bar") => nothing, GetAverage() => 0
 */
package main

import (
	"fmt"
	"time"
)

type WindowedMap interface {
	Put(key string, value int64)
	Get(key string) int64
	GetAverage() float64
}

type windowedMap struct {
	window_size_ms int64
	data           map[string]entry
	lastInserted   []entry
}

type entry struct {
	insertedTime int64
	value        int64
}

func NewWindowedMap(window_size_ms int64) WindowedMap {
	return &windowedMap{
		window_size_ms: window_size_ms,
		lastInserted:   make([]entry, 0),
		data:           make(map[string]entry),
	}
}

func (m *windowedMap) Put(key string, value int64) {
	// TODO: implement me
	e := entry{
		insertedTime: time.Now().UnixMilli(),
		value:        value,
	}

	m.data[key] = e
	m.lastInserted = append(m.lastInserted, e)
}

func (m *windowedMap) Get(key string) int64 {
	// TODO: implement me
	entry := m.data[key]
	if m.window_size_ms+entry.insertedTime < time.Now().UnixMilli() {
		return 0
	} else {
		return entry.value
	}
}

func (m *windowedMap) GetAverage() float64 {
	// TODO: implement me
	now := time.Now().UnixMilli()
	counter := float64(0)
	total := float64(0)
	end := len(m.lastInserted) - 1
	for i := end; i >= 0; i-- {
		if m.window_size_ms+m.lastInserted[i].insertedTime > now {
			total = float64(m.lastInserted[i].value) + total
			counter++
		} else {
			m.lastInserted = m.lastInserted[i:end]
			break
		}
	}

	if counter == 0 {
		return 0
	} else {
		return total / counter
	}
}

func main() {
	winMap := NewWindowedMap(1000)
	winMap.Put("foo", 42)
	time.Sleep(300 * time.Millisecond)
	winMap.Put("bar", 76)

	time.Sleep(200 * time.Millisecond)

	fmt.Println(winMap.Get("foo"))   //42
	fmt.Println(winMap.Get("bar"))   // 76,
	fmt.Println(winMap.GetAverage()) // 59

	time.Sleep(500 * time.Millisecond)
	fmt.Println(winMap.Get("foo"))   // nothing,
	fmt.Println(winMap.Get("bar"))   // 76,
	fmt.Println(winMap.GetAverage()) // 76

	time.Sleep(301 * time.Millisecond)
	fmt.Println(winMap.Get("foo"))   // nothing,
	fmt.Println(winMap.Get("bar"))   // 76,
	fmt.Println(winMap.GetAverage()) // 76
}
