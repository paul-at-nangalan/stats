package stats

import (
	"testing"
	"time"
)

func TestCounter_Inc(t *testing.T) {
	cntr := NewCounter("test")

	for i := 0 ; i < 100; i++{
		cntr.Inc()
	}
	t10 := time.NewTimer(10 * time.Second)
	testdone := false
	for {
		select {
		case <-t10.C:
			rate := cntr.getrate()
			if rate > 10.1 || rate < 9.9{
				t.Error("Incorrect rat expected ~10 but got ", rate)
			}
			testdone = true
			break
		}
		if testdone{
			break
		}
	}
}

func TestBucketCounter_Inc(t *testing.T) {
	bcntr := NewBucketCounter(0, 1200, 200, "test")

	for i := 0; i < 100; i++{
		bcntr.Inc(-100) ///bucket 0
		bcntr.Inc(24.3) ///bucket 1
		bcntr.Inc(102.5) /// bucket 1
		bcntr.Inc(506) /// bucket 3
		bcntr.Inc(890) /// bucket 5
	}
	t10 := time.NewTimer(10 * time.Second)
	testdone := false
	for{
		select {
		case <-t10.C:
			for i := 0; i < len(bcntr.buckets); i++{
				rate := bcntr.buckets[i].getrate()
				if i == 0{
					if rate > 10.1 || rate < 9.9{
						t.Error("Invalid rate for -100, expected 10 at bucket 0 but got ", rate)
					}
				}else if i == 1{
					if rate > 20.2 || rate < 19.8{
						t.Error("Invalid rate for 24.3 + 102.5, expected 20 at bucket 1 but got ", rate)
					}
				}else if i == 3 {
					if rate > 10.1 || rate < 9.9{
						t.Error("Invalid rate for 506, expected 10 at bucket 4 but got ", rate)
					}
				}else if i == 5{
					if rate > 10.1 || rate < 9.9{
						t.Error("Invalid rate for 890, expected 10 at bucket 5 but got ", rate)
					}
				}else{
					if rate != 0{
						t.Error("Expected 0 for ", i, " but got ", rate)
					}
				}

			}
			testdone = true
			break
		}
		if testdone{
			break
		}
	}
}