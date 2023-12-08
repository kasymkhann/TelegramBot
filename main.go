package main

import (
	"flag"
	"log"

	"sqlit-lessonTEST/clients/telegram"
	enent_consumer "sqlit-lessonTEST/consumer/enent-consumer"
	"sqlit-lessonTEST/events/telegram2"
	"sqlit-lessonTEST/storage/files"
)

const storagePath = "storage"
const bathSize = 100

func main() {
	tgClients := telegram.New(botHost(), mustToken())

	eventsProcessor := telegram2.New(tgClients, files.New(storagePath))

	log.Println("Start service")

	consumer := enent_consumer.New(eventsProcessor, eventsProcessor, bathSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped ", err)
	}
}

func mustToken() string {
	token := flag.String("token-bot", "", "token mustToken access ")

	flag.Parse()

	if *token == "" {
		log.Fatal("token not working")
	}

	return *token
}

func botHost() string {
	tgBotHost := flag.String("telegram2-bot-host", "api.telegram2.org", "host")

	flag.Parse()

	return *tgBotHost
}
