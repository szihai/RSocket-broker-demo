package main

import (
	"sync"
	"time"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

func main() {
	const metadataRegister = "register"
	wg := &sync.WaitGroup{}
	wg.Add(1)
	providerA, err := rsocket.Connect().
		KeepAlive(3000*time.Second, 2000*time.Second, 3).
		SetupPayload(payload.NewString("Service A", metadataRegister)).
		Acceptor(func(socket rsocket.RSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(rsocket.RequestResponse(func(msg payload.Payload) rx.Mono {
				return rx.JustMono(payload.NewString(msg.DataUTF8(), "Hello from A"))
			}))
		}).
		Transport("tcp://127.0.0.1:8888").
		Start()

	if err != nil {
		panic(err)
	}
	wg.Wait()

	defer func() {
		_ = providerA.Close()
		wg.Done()
	}()

}
