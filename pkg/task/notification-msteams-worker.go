package task

import (
	"context"
	"templrjs/pkg/config"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/msteams"
)

func SendMsTeamsNotification() {

	ctx := context.Background()
	notifier := notify.New()

	// Set webhook url.
	webhookUrl := config.Config.MSTeams.WebhookUrl
	msteamsService := msteams.New()

	msteamsService.AddReceivers(webhookUrl)

	notifier.UseServices(msteamsService)

	//subject := gjson.Get(msg.Payload, "message.subject")
	//body := gjson.Get(msg.Payload, "message.body")

	// Send a message
	_ = notifier.Send(
		ctx,
		"SDET - Automation Test Hello world",
		"Here are some examples of formatted stuff like "+
			"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***",
	)
}
