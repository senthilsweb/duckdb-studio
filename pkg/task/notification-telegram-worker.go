package task

import (
	"context"
	cfg "templrjs/pkg/config"

	"github.com/go-redis/redis/v8"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/sinhashubham95/jsonic"
	log "github.com/sirupsen/logrus"
)

func SendTelegramNotification() {
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
		// Create a telegram service. Ignoring error for demo simplicity.
		telegramService, _ := telegram.New("1832587595:AAHzuM6zPsFqKxvIjPjTNQqOzfvCyUojtzY")

		// Passing a telegram chat id as receiver for our messages.
		// Basically where should our message be sent?
		telegramService.AddReceivers(-502211658)

		// Create our notifications distributor.
		notifier := notify.New()

		// Tell our notifier to use the telegram service. You can repeat the above process
		// for as many services as you like and just tell the notifier to use them.
		// Inspired by http middlewares used in higher level libraries.
		notifier.UseServices(telegramService)

		j, err := jsonic.New([]byte(msg.Payload))
		if err != nil {
			return
		}
		data, err := j.Child("Statement")
		s, err := data.GetString(".")

		// Send a test message.
		_ = notifier.Send(
			context.Background(),
			"Subject/Title",
			s,
		)

	}
}
