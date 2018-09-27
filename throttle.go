package amazonmws

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Throttling implements Limiter
type Throttling interface {
	Limiter(in []string, out chan []string)
}

// Throttle stores the throttling type for the Products API
type Throttle struct {
	Operation string
	Duration  time.Duration
	Throttled bool
	Denied    bool
	Rate      int /*per second*/
	Quota     int /*per hour*/
	Attempt   int
	SleepMap  map[int]int64
	Start     time.Time
	Synced    *sync.Mutex
	Timer     *time.Timer
	Ticker    *time.Ticker
}

// Request are the throttling rates for the Products API section operations that only throttle per request.
type Request struct {
	*Throttle
	Response string
	Parsed   struct{}
}

// Requests contain a list of request
type Requests struct {
	Requests []*Request
}

// // Item are the throttling rates for the Products API section operations that throttle per item.
// type Item struct {
// 	*Throttle
// }

// NewRequest return a *Request
func NewRequest() *Request {
	return &Request{NewThrottler(), "", struct{}{}}
}

// NewThrottler creates a new Throttler
func NewThrottler() *Throttle {
	return &Throttle{"", time.Second, false, false, 0, 0, 0, nil, time.Now(), &sync.Mutex{}, time.NewTimer(1 * time.Nanosecond), time.NewTicker(1 * time.Nanosecond)}
}

// Limiter aids in preventing excess throttling
func (t *Throttle) Limiter(in []string, out chan []string) {
	x := len(in)
	// pass := make(chan []string)
	y := 5
	// defer close(out)
	if x > 0 {
		s := splitList(x, y)
		// m := modList(float64(x), float64(y))
		if x > y {
			go func() {
				for i := 0; i < s; i++ {
					// for {
					l := len(in)
					if l > y {
						t.NewTicker()
						fmt.Printf("sending %v strings %v\n", y, in[0:y])
						out <- in[0:y]
						// fmt.Println("waiting")
					} else {
						// t.Sleeper()
						// t.NewTicker()
						fmt.Printf("sending %v strings\n%v\n", l, in[0:l])
						out <- in[0:l]
						// fmt.Println("waiting")
					}
					if l > y {
						in = append(in[:0], in[y:]...)
					} else {
						in = append(in[:0], in[l:]...)
					}
					// if l == 0 {
					// 	return
					// }
				}
				close(out)
			}()
		} else {
			// t.NewTicker()
			fmt.Printf("sending %v strings\n", x)
			out <- in
			close(out)
		}

	}
}

// Sleepy stores throttling durations
func (t *Throttle) Sleepy() {
	m := make(map[int]int64)
	m[0] = int64(2)
	m[1] = int64(5)
	m[2] = int64(10)
	m[3] = int64(30)
	// m[4] = int64(60)
	t.SleepMap = m
}

// Sleeper handles throttling durations
func (t *Throttle) Sleeper() {
	if len(t.SleepMap) == 0 {
		t.Sleepy()
	}
	v, ok := t.SleepMap[t.Attempt]
	if ok == false {
		t.Sleepy()
		t.Sleeper()
	}
	t.Duration = time.Duration(v) * time.Second
	t.Throttler()
	// t.NewTicker()
}

// NewTimer initializes a time.NewTimer
func (t *Throttle) NewTimer(d time.Duration) {
	if t.Timer.Stop() {
		t.Timer.Reset(d)
	}
	<-t.Timer.C
}

// NewTicker initializes a time.NewTicker
func (t *Throttle) NewTicker() {
	t.Ticker = time.NewTicker(time.Second)
	defer t.Ticker.Stop()
	done := make(chan bool, 1)
	go func() {
		t.Sleeper()
		// time.Sleep(time.Second)
		done <- true
		// for nt := range t.Ticker.C {
		// 	fmt.Printf("ticking for %v at %v\n", time.Second, nt)
		// }
	}()
	for {
		select {
		case <-done:
			fmt.Println("Finished ticking")
			return
		case t := <-t.Ticker.C:
			fmt.Println("sleeper is ticking", t)
		}
	}
	// <-t.Ticker.C

}

// NewSlowTicker initializes a time.NewTicker
func (t *Throttle) NewSlowTicker() {
	t.Ticker = time.NewTicker(time.Minute)
	defer t.Ticker.Stop()
	done := make(chan bool, 1)
	go func() {
		wait := t.Start.Add(time.Hour)
		// time.Sleep(wait.Sub(time.Now()))
		t.NewTimer(wait.Sub(time.Now()))
		done <- true
		// for nt := range t.Ticker.C {
		// 	fmt.Printf("ticking for %v at %v\n", time.Second, nt)
		// }
	}()
	for {
		select {
		case <-done:
			fmt.Println("slow sleeper finished ticking")
			return
		case t := <-t.Ticker.C:
			fmt.Println("slow sleeper is ticking", t)
		}
	}
	// <-t.Ticker.C

}

// Throttler sleeps
func (t *Throttle) Throttler() {
	if t.Attempt == 3 {
		fmt.Println("sleeping for :", t.Duration)
		time.Sleep(t.Duration)
		t.Attempt = 0
	} else {
		fmt.Println("sleeping for :", t.Duration)
		time.Sleep(t.Duration)
		t.Attempt++
	}
}
func splitList(x, y int) int {
	return x / y
}
func modList(x, y float64) float64 {
	return math.Mod(x, y)
}

// time your retries with the following time spacing: 1s, 4s, 10s, 30s.
/*Per-Request Throttling
These are the throttling rates for the Products API section operations that
only throttle per request.

Operation	Maximum request quota	Restore rate	Hourly request quota
ListMatchingProducts	20 requests	One request every five seconds	720 requests
per hour
GetProductCategoriesForSKU and GetProductCategoriesForASIN	20 requests	One request every five seconds	720 requests per hour

Per-Item Throttling
Operation	Maximum request quota	Restore rate	Hourly request quota
GetMatchingProduct	20 requests	Two items every second	7200 requests per hour
GetMatchingProductForId	20 requests	Five items every second	18000 requests per hour
GetCompetitivePricingForSKU and GetCompetitivePricingForASIN	20 requests	10 items every second	36000 requests per hour
GetLowestOfferListingsForSKU and GetLowestOfferListingsForASIN	20 requests	10 items every second	36000 requests per hour
GetLowestPricedOffersForSKU and GetLowestPricedOffersForASIN	10 requests	Five items every second	200 requests per hour
GetMyFeesEstimate	20 requests	10 items every second	36000 requests per hour
GetMyPriceForSKU and GetMyPriceForASIN	20 requests	10 items every second	36000 requests per hour */
