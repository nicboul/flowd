package main

import (
	"net/http"
	"os"
	"time"

	"github.com/nicboul/flowdata/internal/flowdata"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "flowd"
	app.Usage = "flowd collect, aggregate, store and serve flow data"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "timeout",
			Value:  600,
			Usage:  "timeout in seconds",
			EnvVar: "FLOWD-TIMEOUT",
		},
		cli.StringFlag{
			Name:   "port",
			Value:  "8080",
			Usage:  "port number the service is listening to",
			EnvVar: "FLOWD-PORT",
		},
		cli.StringFlag{
			Name:   "listen",
			Value:  "127.0.0.1",
			Usage:  "IP address the service is listening to",
			EnvVar: "FLOWD-LISTEN",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("starting service flowd")

		params := flowdata.FlowDataParams{
			Timeout: time.Second * time.Duration(int64(c.Int("timeout"))),
			Store:   store.NewFlowDataStore(),
			Queue:   queue.NewFlowDataQueue(100),
		}

		server := flowdata.NewFlowDataServer(params)
		listen := c.String("listen") + ":" + c.String("port")

		log.Info("listening on: ", listen)

		return http.ListenAndServe(listen, server.MuxRouter)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
