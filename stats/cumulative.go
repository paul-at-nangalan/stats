package stats

import (
	"fmt"
	"sync"
)

type CumulativeCounter struct {
	count float64
	lock  sync.RWMutex
	name  string
}

func (c *CumulativeCounter) print() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	fmt.Println(c.name, ": ", c.count)
}

func (c *CumulativeCounter) getname() string {
	return c.name
}

func (c *CumulativeCounter) getrate() float64 {
	return c.count //// not a rate counter
}

func (c *CumulativeCounter) hasrate() bool {
	return false
}

func (c *CumulativeCounter) Inc(amount float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.count += amount
}

func NewCumulativeCounter(name string) *CumulativeCounter {
	return &CumulativeCounter{
		count: 0,
		name:  name,
	}
}
