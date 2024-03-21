package task

import (
	"context"
	cfg "templrjs/pkg/config"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func DockerUpgrade() {
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

	}
}
