package main

import (
	"encoding/json"
	"log"
	"nats-demo/models"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	opt := nats.UserInfo("ruser", "T0pS3cr3t")
	nc, err := nats.Connect("nats://localhost:4222", opt)
	if err != nil {
		log.Fatal(err)
	}

	//JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Create durable consumer monitor
	sub, _ := js.QueueSubscribe("1.slack.eu", "queue-1", func(msg *nats.Msg) {

		var messageEvent models.MessageEvent
		err = json.Unmarshal(msg.Data, &messageEvent)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s -  Message from Org:%d, Provider:%s, Region: %s and Message:%s \n", msg.Sub.Queue, messageEvent.OrgId, messageEvent.Provider, messageEvent.Message, messageEvent.Region)
		err := msg.Ack()
		if err != nil {
			log.Fatal(err)
		}
	}, nats.Durable("durable-push"), nats.AckWait(time.Second))

	sub2, _ := js.QueueSubscribe("2.slack.eu", "queue-2", func(msg *nats.Msg) {

		err := msg.Ack()
		if err != nil {
			log.Fatal(err)
		}

		var messageEvent models.MessageEvent
		err = json.Unmarshal(msg.Data, &messageEvent)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s -  Message from Org:%d, Provider:%s, Region: %s and Message:%s \n", msg.Sub.Queue, messageEvent.OrgId, messageEvent.Provider, messageEvent.Message, messageEvent.Region)
	}, nats.Durable("durable-push"), nats.ManualAck(), nats.MaxDeliver(2), nats.AckWait(time.Second))

	subAll, _ := nc.Subscribe("*.slack.eu", func(msg *nats.Msg) {
		var messageEvent models.MessageEvent
		err = json.Unmarshal(msg.Data, &messageEvent)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s -  Message from Org:%d, Provider:%s, Region: %s and Message:%s \n", msg.Subject, messageEvent.OrgId, messageEvent.Provider, messageEvent.Message, messageEvent.Region)
	})

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	sub.Unsubscribe()
	sub2.Unsubscribe()
	subAll.Unsubscribe()
	nc.Drain()
}
