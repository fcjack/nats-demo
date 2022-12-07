package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nats-demo/models"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var streamsMap = map[string][]string{
	"1-slack": {"1.slack.eu"},
	"2-slack": {"2.slack.eu", "2.slack.us"},
	"3-slack": {"3.slack.eu"},
}

func main() {
	opt := nats.UserInfo("ruser", "T0pS3cr3t")
	nc, err := nats.Connect("nats://localhost:4222", opt)

	if err != nil {
		log.Fatal("check error: ", err)
	}

	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("check error: ", err)
	}

	for streamName, streamSubjects := range streamsMap {
		// Creates stream
		err = createStream(js, streamName, streamSubjects)
		if err != nil {
			log.Fatal("check error: ", err)
		}
	}

	for i := 0; ; i++ {
		for orgId := 1; orgId <= 3; orgId++ {

			messageEvent := models.MessageEvent{
				Provider: "slack",
				OrgId:    orgId,
				Message:  fmt.Sprintf("this is my message number: %d", i),
				Region:   "eu",
			}

			createMessageEvent(js, messageEvent)

		}

		time.Sleep(time.Second)
	}
}

func createStream(js nats.JetStreamContext, streamName string, streamSubjects []string) error {
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Printf("failed to read stream: %s", err.Error())
	}

	if stream == nil {
		log.Printf("creating stream %q and subjects %s", streamName, strings.Join(streamSubjects, " "))

		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  streamSubjects,
			Retention: nats.LimitsPolicy,
			MaxAge:    time.Minute * 10000,
		})

		if err != nil {
			log.Printf("failed to create stream: %s", err.Error())
			return err
		}
	}

	return nil
}

func createMessageEvent(js nats.JetStreamContext, messageEvent models.MessageEvent) error {
	meessageEventJSON, _ := json.Marshal(messageEvent)
	subj := fmt.Sprintf("%d.%s.%s", messageEvent.OrgId, messageEvent.Provider, messageEvent.Region)
	_, err := js.Publish(subj, meessageEventJSON)
	if err != nil {
		return err
	}

	return nil
}
