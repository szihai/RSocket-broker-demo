package main

import (
	"context"
	"log"
	"sync"
	"time"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

func main() {
	const metadataRegister = "register"
	wg := &sync.WaitGroup{}
	clientB, err := rsocket.Connect().
		KeepAlive(3000*time.Second, 2000*time.Second, 3).
		SetupPayload(payload.NewString("client b", metadataRegister)).
		Transport("tcp://127.0.0.1:8888").
		Start()

	defer func() {
		_ = clientB.Close()
	}()

	if err != nil {
		panic(err)
	}
	clientB.RequestResponse(payload.NewString("Asking service A", "Service A")).
		DoOnSuccess(func(ctx context.Context, s rx.Subscription, elem payload.Payload) {
			m, _ := elem.MetadataUTF8()
			log.Printf("got Response from A: data=%s, metadata=%s\n", elem.DataUTF8(), m)
		}).
		Subscribe(context.Background())
}
