# scheduler

# Examples

Here we show some basic code snippets - creating and using scheduler sch, starting the server, composing the client, adding a job into queue, and finally actuall task worker implementation

#### Scheduler, later used into http handler/server - 
```
...

	workersNum, err := strconv.Atoi(os.Getenv("SCHEDULER_WORKERS_NUM"))
	if err != nil {
		log.Fatal("cant process default workers number in pool")
	}

	// scheduler
	sch := scheduler.New(
		scheduler.Config{
			RedisURL:   os.Getenv("REDIS_URL"),
			WorkersNum: workersNum,
		},
	)

...
```

#### Server
Here we are creating the sch and strarting server workers, defining per each needed variables, worker pool nums etc.
```
...

	workersNum, err := strconv.Atoi(os.Getenv("SCHEDULER_WORKERS_NUM"))
	if err != nil {
		log.Fatal("cant process default workers number in pool")
	}

	// scheduler
	sch := scheduler.New(
		scheduler.Config{
			RedisURL:   os.Getenv("REDIS_URL"),
			WorkersNum: workersNum,
		},
	)

	err = sch.StartServer(map[string]asynq.Handler{
		tasks.TypeRSSParse: tasks.RSSParseTask{},
	})
	if err != nil {
    ...
...


```

#### Client - Handler Add Job into queue -
Here we reuse injected scheduler sch into server instance, to add a job
```
...

func (serv *QueueServer) Add(resp http.ResponseWriter, req *http.Request) {

    ...

	err = serv.scheduler.Add(scheduler.Task{
		//QueuePriority: scheduler.QueueCritical,
		TaskType: tasks.TypeRSSParse,
		Data: map[string]interface{}{
			"rss_uris": reqData.RssURIs,
		},
		MaxRetries: 3,
		// if not defined, will be now
		// if we want after 5 secs, or in fully defined date/time in the future
		//ProcessAt:  time.Now().Add(5 * time.Second),
	})
	if err != nil {
		serv.logger.WithError(err).Error("error scheduling task: loadsController.TypeLoadsCsvExport")

		httputil.RenderErr(resp, "cant add task to queue", http.StatusInternalServerError)
		return
	}

    ...
}
```

#### Sample task implentation
Just a snippet, showing task type definition, task struct with needed fields/if any/ and the actuall method, which MUST satisfy the asynq.Handler interface.

```
...


const (
	// TypeRSSParse -- rss parse task type
	TypeRSSParse = "rss:parse"
)

type (
	// RSSParseTask -- implements asynq.Handler interface
	RSSParseTask struct {
		// whatever fields needed here ...
	}
)

// ProcessTask -- exec method, implementing asynq.Handler interface
func (task RSSParseTask) ProcessTask(ctx context.Context, t *asynq.Task) error {

    ... 

}

```
