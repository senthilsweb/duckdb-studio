package task

import (
	"templrjs/pkg/broker"
	"templrjs/pkg/utils"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

var server *machinery.Server

func SendNewsletter(payload string) (string, error) {
	var strToConvert string
	strToConvert = utils.DecodeAsStringFromBase64EncodedString(payload)
	subject := gjson.Get(strToConvert, "message.subject")
	log.Info("subject = " + subject.String())
	return subject.String(), nil
}

func LaunchNewsletterWorkers() {

	server = broker.GetBroker()
	log.Info("Broker received")
	server.RegisterTask("newsletter", SendNewsletter)
	log.Info("Task registered")
	worker := server.NewWorker("worker-1", 5)
	log.Info("Launching worker")
	err := worker.Launch()
	if err != nil {
		log.Fatal(err, "Could not launch worker!")
	}
}
