package utils

import (
	"math/rand"
	"sync"
	"time"
)

type JitterTicker struct {
	duration time.Duration
	jitter   time.Duration
	C        chan time.Time
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewJitterTicker(duration time.Duration, jitter time.Duration) *JitterTicker {
	if duration < jitter {
		panic(`JitterTicker Error: must be (duration >= jitter)`)
	}
	if duration <= 0 || jitter <= 0 {
		panic("JitterTicker Error: must be (duration/jitter > 0) ")
	}

	j := &JitterTicker{
		duration: duration,
		jitter:   jitter,
		C:        make(chan time.Time),
		stop:     make(chan struct{}),
	}
	j.Start()
	return j
}

func (j *JitterTicker) Start() {
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		for {
			random := rand.Int63n(int64(j.jitter * 2))
			// random := rand.Int63n(int64(j.jitter)) - int64(j.jitter)
			jitter := time.Duration(random) - j.jitter
			sleepTime := j.duration + jitter

			select {
			case <-time.After(sleepTime):
				select {
				case j.C <- time.Now():
				default:
				}
			case <-j.stop:
				return
			}
		}
	}()
}

func (j *JitterTicker) Stop() {
	select {
	case <-j.stop:
	default:
		close(j.stop)
		j.wg.Wait()
		close(j.C)
	}
}
