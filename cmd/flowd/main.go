package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicboul/flowdata/internal/flowdataread"
	"github.com/nicboul/flowdata/internal/flowdatawrite"
	"github.com/nicboul/flowdata/internal/store"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type FlowDataParams struct {
	Timeout time.Duration
	Store   *store.FlowDataStore
}

func NewFlowDataServer(p FlowDataParams) *mux.Router {

	flowDataWriteHandler := &flowdatawrite.FlowDataWrite{
		Store: p.Store,
	}

	flowDataReadHandler := &flowdataread.FlowDataRead{
		Store: p.Store,
	}

	muxRouter := mux.NewRouter()
	muxRouter.Methods("POST").PathPrefix("/flows").Handler(flowDataWriteHandler)
	muxRouter.Methods("GET").PathPrefix("/flows").Handler(flowDataReadHandler)

	return muxRouter
}

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

		params := FlowDataParams{
			Timeout: time.Second * time.Duration(int64(c.Int("timeout"))),
			Store:   store.NewFlowDataStore(),
		}

		server := NewFlowDataServer(params)
		listen := c.String("listen") + ":" + c.String("port")

		log.Info("listening on: ", listen)

		return http.ListenAndServe(listen, server)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
