package main

import (
	"context"
	"fmt"
	"log"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

func main() {
	const metadataRegister = "register"
	bus := rsocket.NewBus()

	err := rsocket.Receive().
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.EnhancedRSocket) rsocket.RSocket {
			// register socket to bus.
			metadata, _ := setup.MetadataUTF8()
			if metadata == metadataRegister {
				merge := struct {
					id string
					sk rsocket.RSocket
				}{setup.DataUTF8(), sendingSocket}
				log.Println(setup.DataUTF8())
				sendingSocket.OnClose(func() {
					log.Println("removing service " + merge.id)
					bus.Remove(merge.id, merge.sk)
				})
				bus.Put(merge.id, merge.sk)
			}
			// bind responder: redirect to target socket.
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) rx.Mono {
					id, _ := msg.MetadataUTF8()
					log.Println("asking for " + id)
					sk, ok := bus.Get(id)
					if !ok {
						return rx.NewMono(func(ctx context.Context, sink rx.MonoProducer) {
							sink.Error(fmt.Errorf("MISSING_SERVICE_SOCKET_%s", id))
						})
					}
					return sk.RequestResponse(msg)
				}),
				rsocket.RequestStream(func(msg payload.Payload) rx.Flux {
					id, _ := msg.MetadataUTF8()
					sk, ok := bus.Get(id)
					if !ok {
						return rx.NewFlux(func(ctx context.Context, producer rx.Producer) {
							producer.Error(fmt.Errorf("MISSING_SERVICE_SOCKET_%s", id))
						})
					}
					return sk.RequestStream(msg)
				}),
			)
		}).
		Transport("127.0.0.1:8888").
		Serve()

	if err != nil {
		panic(err)
	}
}
