package counter

const (
	kindIncrement = iota
	kindGet = 1
)


type request struct {
	resp chan int
	requestKind int
	key string
	delta int
}

type Counter struct {
	req chan request
}

func New() *Counter {
	// create request channel
	c:=&Counter {
		req: make (chan request),
	}
	go counter(c)
	return c
}

func (c *Counter) Get (key string) int {
	resp:= make (chan int)

	c.req <- request {
		resp: resp,
		requestKind: kindGet,
		key: key,
		delta: 0,
	}
	return <- resp
}


func (c *Counter) Incr(key string, delta int) int {
	resp:= make (chan int)

	c.req <- request {
		resp: nil,
		requestKind: kindIncrement,
		key: key,
		delta: delta,
	}
	return <- resp
}


func counter(c *Counter) {
	data :=make(map[string]int)
	for req := range c.req {
		switch (req.requestKind) {
		case kindIncrement:
			data[req.key] =data[req.key] + req.delta	
		case kindGet:
			req.resp <-data[req.key]
		}
	}
}
