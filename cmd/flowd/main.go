package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicboul/flowdata/internal/flowdata"
	"github.com/nicboul/flowdata/internal/queue"
	"github.com/nicboul/flowdata/internal/store"
	log "github.com/sirupsen/logrus"
)

func main() {

	timeoutFlag := flag.Int64("timeout", 600, "timeout in second")
	portFlag := flag.String("port", "8080", "port number the service is listening to")
	listenFlag := flag.String("listen", "127.0.0.1", "IP address the service is listenign to")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.SetFormatter(&log.JSONFormatter{})
	log.Info("starting service flowd")

	params := flowdata.FlowDataParams{
		Timeout: time.Second * time.Duration(*timeoutFlag),
		Store:   store.NewFlowDataStore(),
		Queue:   queue.NewFlowDataQueue(5000),
		Listen:  *listenFlag + ":" + *portFlag,
	}

	server := flowdata.NewFlowDataServer(params)

	log.Info("listening on: ", params.Listen)
	server.Serve(ctx)

	fmt.Printf("bye bye!\n")
}
