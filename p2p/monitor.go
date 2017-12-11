package p2p

import (
	"sync"
	"time"
)

type Monitor struct {
	mtx       sync.Mutex
	count     uint
	done      chan bool
	isrunning bool
	lastok    time.Duration
	lastop    time.Duration
}

func (m *Monitor) Update(op bool) {

	m.mtx.Lock()
	defer m.mtx.Unlock()
	if op {
		m.lastok = time.Duration(time.Now().Unix())
		if m.count >= 1 {
			m.count--
		} else {
			m.count = 0
		}
	}

	m.lastop = time.Duration(time.Now().Unix())
	if !op {
		m.count++
	}

}

func NewMonitor() *Monitor {

	var m = &Monitor{
		done:      make(chan bool),
		lastok:    time.Duration(time.Now().Unix()),
		lastop:    time.Duration(time.Now().Unix()),
		count:     0,
		isrunning: true,
	}
	m.Start()
	return m
}

func (m *Monitor) Start() {
	go func(m *Monitor) {
		for {
			tick := time.NewTicker(time.Second * 5)
			select {
			case <-tick.C:
				if m.lastop-m.lastok > 600 || m.count > 10 {

					m.Stop()

				}
			case <-m.done:
				break
			}
		}

	}(m)

}

func (m *Monitor) IsRunning() bool {
	return m.isrunning
}
func (m *Monitor) Stop() {
	m.isrunning = false
	m.done <- false
}