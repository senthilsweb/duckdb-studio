package task

import (
	"context"
	cfg "templrjs/pkg/config"

	"github.com/go-redis/redis/v8"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/slack"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func SendSlackNotification() {
	//opt, _ := redis.ParseURL("rediss://:b3cd263ba8f14f32b7dee9907577d089@us1-deciding-kodiak-34245.upstash.io:34245")
	opt, _ := redis.ParseURL(cfg.Config.Queue.Uri)
	client := redis.NewClient(opt)

	ctx := context.Background()
	// Subscribe to the Topic given
	topic := client.Subscribe(ctx, "templrjs_github")
	// Get the Channel to use
	channel := topic.Channel()

	// Itterate any messages sent on the channel
	for msg := range channel {
		log.Info("Message Received" + msg.Payload)

		notifier := notify.New()

		// Provide your Slack OAuth Access Token
		slackService := slack.New("xoxb-2294734803764-2294265292647-zeiX0cFHIXo04iyJq7reMyNi")

		// Passing a Slack channel id as receiver for our messages.
		// Where to send our messages.
		slackService.AddReceivers("C05S91WK57D")

		// Tell our notifier to use the Slack service. You can repeat the above process
		// for as many services as you like and just tell the notifier to use them.
		notifier.UseServices(slackService)

		subject := gjson.Get(msg.Payload, "message.subject")
		body := gjson.Get(msg.Payload, "message.body")

		// Send a message
		_ = notifier.Send(
			context.Background(),
			subject.String(),
			body.String(),
		)

	}
}
