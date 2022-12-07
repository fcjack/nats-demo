# nats-demo

Small demo application to demonstrate how create publishers/subscribers to NATS in Go.

## What is NATS

Nats is a single technology that enables applications to securely communicate across any combination of cloud vendors, on-premise, edge, web and mobile, and devices.

More details about NATS can be found [here](https://nats.io)

## Local setup

Docker compose contains a cluster of NATS that can be used to work locally with the demo.

All the configuration for the nats server can be done using the config file, in our demo the file is called `jetstream.config`

The publisher will send messages constantly for different topics and subjects, while we have multiple subscribers that will consume the messages and print the message consumed.

The subscribers are defined using jetstream context and in a specific case we have using NATS context that will create a subscriber using the wildcard in the subject name, that will allow us to consume all messages that matches the pattern.
