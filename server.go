package scheduler

import (
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// StartServer -- method, used to start the scheduler server
func (s *Scheduler) StartServer(queueHandlers map[string]asynq.Handler) error {
	logger := logrus.WithFields(logrus.Fields{
		"module":    "scheduler",
		"component": "server",
	})

	redisOpts, err := s.GetRedisClientOpts()
	if err != nil {
		logger.WithError(err).Fatal("error parsing REDIS_URL")
	}

	srv := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Concurrency: s.Config.WorkersNum,
			Queues: map[string]int{
				QueueRealtime: 9,
				QueueCritical: 6,
				QueueDefault:  3,
				QueueLow:      1,
			},
			Logger: logger,
		},
	)

	mux := asynq.NewServeMux()

	for queueName, taskEntity := range queueHandlers {
		mux.Handle(queueName, taskEntity)
	}

	if err := srv.Run(mux); err != nil {
		logger.WithError(err).Fatalf("could not run server: %v", err)

		return err
	}

	return nil
}
