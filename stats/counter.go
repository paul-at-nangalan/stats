package stats

import (
	"fmt"
	"log"
	"math"
	"sync/atomic"
	"time"
)

type Counter struct {
	count     int64
	lastcount int64
	name      string
	t         time.Time
}

func NewCounter(name string) *Counter {
	t := time.Now()
	return &Counter{
		count: 0,
		t:     t,
		name:  name,
	}
}

func (p *Counter) Inc() {
	atomic.AddInt64(&p.count, 1)
}

func (p *Counter) print() {
	rate := p.getrate()
	fmt.Println(p.name, ": ", rate)
}

func (p *Counter) getname() string {
	return p.name
}

func (p *Counter) getrate() float64 {
	lclcount := p.count
	t2 := time.Now()
	tdiff := t2.Sub(p.t)
	rate := float64(lclcount-p.lastcount) / tdiff.Seconds()
	p.t = t2
	p.lastcount = lclcount
	return rate
}

func (p *Counter) hasrate() bool {
	if p.count > 0 {
		return true
	}
	return false
}

type Stat interface {
	print()
	getname() string
	getrate() float64
	hasrate() bool
}

var stop chan bool

func Run(printintvl time.Duration, counters []Stat) {
	stop = make(chan bool, 5)
	quit := false
	ticker := time.NewTicker(printintvl)
	for !quit {
		select {
		case <-ticker.C:
			for _, stat := range counters {
				stat.print()
			}
		case <-stop:
			quit = true
		}
	}
}
func Stop() {
	if stop != nil {
		stop <- true
	}
}

type BucketCounter struct {
	buckets []*Counter
	ranges  []float64
	name    string
	count   int64

	min, max, step float64
}

func NewBucketCounter(min, max, step float64, name string) *BucketCounter {
	buckets := make([]*Counter, 1)
	ranges := make([]float64, 1)
	buckets[0] = NewCounter(fmt.Sprint(name, ": < ", min))
	ranges[0] = min
	for i := min + step; i < max; i += step {
		name := fmt.Sprint(">=", i, "<", i+step)
		ranges = append(ranges, i)
		buckets = append(buckets, NewCounter(name))
	}
	bucketname := fmt.Sprint(">", max)
	ranges = append(ranges, math.MaxFloat64)
	buckets = append(buckets, NewCounter(bucketname))
	return &BucketCounter{
		buckets: buckets,
		ranges:  ranges,
		name:    name,
		min:     min,
		max:     max,
		step:    step,
	}
}

func (p *BucketCounter) Inc(val float64) {
	atomic.AddInt64(&p.count, 1)
	if val <= p.min {
		p.buckets[0].Inc()
		return
	}
	if val >= p.max {
		p.buckets[len(p.buckets)-1].Inc()
		return
	}

	relpos := (val - p.min)
	if relpos < 0 {
		log.Panic("This should be impossible ... the index is negative despite checking for val <= min")
	}
	indx := int(relpos / p.step)
	fmt.Println("ffs index is ", indx, " relpos is ", relpos, " step is ", p.step, " val is ", val,
		" min is ", p.min, " max is ", p.max)
	for i := indx; i < len(p.ranges); i++ {
		if val <= p.ranges[i] {
			break
		}
		indx++
	}
	p.buckets[indx].Inc()
}

func (p *BucketCounter) print() {
	fmt.Println("---------------------------", p.name, "--------------------------------")
	sep := ""
	for i := 0; i < len(p.buckets); i++ {
		counter := p.buckets[i]
		if counter.hasrate() {
			fmt.Print(sep, counter.getname())
			sep = ","
		}
	}
	fmt.Println()
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println()
	sep = ""
	for i := 0; i < len(p.buckets); i++ {
		counter := p.buckets[i]
		if counter.hasrate() {
			fmt.Print(sep, counter.getrate())
			sep = ","
		}
	}
	fmt.Println()
	fmt.Println("--------------------- END ", p.name, "--------------------------------")
}

func (p *BucketCounter) getname() string {
	return p.name
}

func (p *BucketCounter) getrate() float64 {
	log.Panic("BucketCounter get rate not implemented")
	return 0
}

func (p *BucketCounter) hasrate() bool {
	if p.count > 0 {
		return true
	}
	return false
}
