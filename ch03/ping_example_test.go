package ch03

import (
	"context"
	"io"
	"time"
	"log"
)

func ExamplePinger() {
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe() // in lieu of net.Conn
	done := make(chan struct{})
	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second // initial ping interval

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			log.Printf("resetting timer (%s)\n", d)
			resetTimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf)

		if err != nil {
			log.Println(err)
		}

		log.Printf("received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))
	}

	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		log.Printf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	cancel()
	<-done // ensure the pinger exits after canceling the context
}