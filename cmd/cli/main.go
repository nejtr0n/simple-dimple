package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nejtr0n/simple-dimple/aggregator"
	"github.com/nejtr0n/simple-dimple/app"
	"github.com/nejtr0n/simple-dimple/generator"
	"github.com/nejtr0n/simple-dimple/pubsub"
	"github.com/nejtr0n/simple-dimple/queue"
	"github.com/nejtr0n/simple-dimple/storage"
	"github.com/urfave/cli"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := &cli.App{
		Name:  "simple-dimple",
		Usage: "Simple pub sub broker with aggregation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "config",
				Value:  "config.json",
				Usage:  "path to config file",
				EnvVar: "APP_CONFIG",
			},
		},
		Action: func(c *cli.Context) error {
			// read config
			file := c.String("config")
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			var config app.Config
			err = json.Unmarshal(data, &config)
			if err != nil {
				return err
			}

			// start process
			pipe := make(queue.Queue, config.Queue.Size)
			var workers []<-chan generator.Message
			for _, generatorItem := range config.Generators {
				gen := make(generator.Generator, len(generatorItem.DataSources))
				for _, s := range generatorItem.DataSources {
					gen.Add(ctx, generatorItem.TimeoutS, generator.NewWorker(s.Id, s.InitValue, s.MaxChangeStep))
				}
				close(gen)
				workers = append(workers, gen.Start(ctx, generatorItem.SendPeriodS))
			}

			pipe.Merge(ctx, workers...)

			ps := pubsub.NewPubSub()
			for _, agr := range config.Aggregators {
				ps.Subscribe(agr.SubIds, aggregator.NewAggregator(ctx, agr.AggregatePeriodS))
			}
			result := ps.Process(ctx, pipe)

			// async send to storage
			store, err := storage.NewStorage(config.StorageType)
			if err != nil {
				return err
			}
			defer store.Close()
			var wg sync.WaitGroup
			for summary := range result {
				wg.Add(1)
				go func(s storage.Summary) {
					defer wg.Done()
					store.Store(s)

				}(summary)
			}
			wg.Wait()

			return nil
		},
	}

	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		osCall := <-c
		log.Printf("system call:%+v\r\n", osCall)
		cancel()
	}()

	err := application.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
