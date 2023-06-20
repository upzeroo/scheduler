package scheduler

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type (
	// Task -- general task struct
	Task struct {
		QueuePriority string
		TaskType      string
		Data          map[string]interface{}
		MaxRetries    int
		ProcessAt     time.Time
	}
)

// Add -- creates asynq task and adds it into the queue
func (s *Scheduler) Add(task Task) error {
	logger := logrus.WithFields(logrus.Fields{
		"module":    "scheduler",
		"component": "client",
	})

	if task.QueuePriority == "" {
		task.QueuePriority = QueueDefault
	}

	payload, err := json.Marshal(task.Data)
	if err != nil {
		logger.WithError(err).Errorf("json.Marshal failed: %v", err)

		return err
	}

	redisOpts, err := s.GetRedisClientOpts()
	if err != nil {
		logger.WithError(err).Fatal("error parsing REDIS_URL")
	}

	client := asynq.NewClient(
		redisOpts,
	)

	defer client.Close()

	info, err := client.Enqueue(
		asynq.NewTask(
			task.TaskType,
			payload,
			asynq.MaxRetry(task.MaxRetries),
			asynq.Timeout(20*time.Minute),
		),
		asynq.Queue(task.QueuePriority),
		asynq.ProcessAt(task.ProcessAt),
	)
	if err != nil {
		logger.WithError(err).Fatalf("could not add task: %v", err)

		return err
	}

	logger.Infof("task added: id=%s queue=%s", info.ID, info.Queue)

	return nil
}
