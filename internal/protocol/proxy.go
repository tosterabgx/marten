package protocol

import (
	"io"
	"net"
)

func Proxy(connA, connB net.Conn) {
	defer connA.Close()
	defer connB.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(connB, connA)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(connA, connB)
	}()

	<-done
}
