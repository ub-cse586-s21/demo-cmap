package cmap

import "log"

// Concurrent map data structure map[string]string

type ConcurrentMap interface {
	Get(key string) (string, bool)
	Put(key string, value string)
	Delete(key string) bool
	Free() error
}

func New() ConcurrentMap {
	c := make(chan *cReq)
	go cMapRunner(c)
	return &cMap{c}
}

type cMap struct {
	c chan<- *cReq
}

type cReq struct {
	key    string  // key to get, or key to put
	value  string  // value to put on put
	action cAction // put or get
	c      chan<- *cResult
}

type cAction int

const (
	cmap_PUT = cAction(iota + 1)
	cmap_GET
	cmap_DELETE
)

type cResult struct {
	value string
	found bool
}

// Flow for interacting with ConcurrentMap:
// 1. Create request object
// 2. Send on channel to CMap goroutine
// 3. Block waiting for response
//
// In parallel, CMap goroutine will:
// 1. Read request
// 2. Perform action
// 3. Send response
// 4. Goto 1

func (cm *cMap) Get(key string) (string, bool) {
	c := make(chan *cResult)
	cm.c <- &cReq{action: cmap_GET, key: key, c: c}
	res := <-c
	return res.value, res.found
}

func (cm *cMap) Put(key string, value string) {
	c := make(chan *cResult)
	cm.c <- &cReq{action: cmap_PUT, key: key, value: value, c: c}
	<-c
}

func (cm *cMap) Delete(key string) bool {
	c := make(chan *cResult)
	cm.c <- &cReq{action: cmap_DELETE, key: key, c: c}
	res := <-c
	return res.found
}

func (cm *cMap) Free() error {
	close(cm.c)
	return nil
}

func cMapRunner(c <-chan *cReq) {
	m := make(map[string]string)
	for req := range c {
		switch req.action {
		case cmap_GET:
			v, ok := m[req.key]
			req.c <- &cResult{value: v, found: ok}
		case cmap_PUT:
			m[req.key] = req.value
			req.c <- nil
		case cmap_DELETE:
			_, ok := m[req.key]
			delete(m, req.key)
			req.c <- &cResult{found: ok}
		default:
			log.Fatalf("Unknown request: %v", req.action)
		}
	}
}
