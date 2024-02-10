package stats

import (
	"fmt"
	"sync/atomic"
)

type CumulativeCounter struct {
	count int64
	name  string
}

func (c *CumulativeCounter) print() {
	fmt.Println(c.name, ": ", c.count)

}

func (c *CumulativeCounter) getname() string {
	return c.name
}

func (c *CumulativeCounter) getrate() float64 {
	return float64(c.count) //// not a rate counter
}

func (c *CumulativeCounter) hasrate() bool {
	return false
}

func (c *CumulativeCounter) Inc(amount int64) {
	atomic.AddInt64(&c.count, amount)
}

func NewCumulativeCounter(name string) *CumulativeCounter {
	return &CumulativeCounter{
		count: 0,
		name:  name,
	}
}
