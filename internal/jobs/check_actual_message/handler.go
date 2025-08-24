package worker

import (
	"context"
	"encoding/json"
	"time"
)

func (s Job) CheckActualMessageJob(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			events, err := s.Storage.GetEventsToday(ctx)
			if err != nil {
				//TODO log
				continue
			}

			for _, event := range events {
				now := time.Now().Truncate(time.Minute)

				if now.After(event.TimeToNotify.Truncate(time.Minute)) && !event.IsSent {
					b, err := json.Marshal(event)
					if err != nil {
						//TODO log
						continue
					}

					err = s.PubSub.Publish(ctx, "check_actual_message", b)
					if err != nil {
						//TODO log
						continue
					}
				}
			}
		}

	}
}
