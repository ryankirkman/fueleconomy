package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/handlers"
	"github.com/teasherm/fueleconomy/workers"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "0.0.0.0:8000", "Address to listen for HTTP requests on")
)

func main() {
	global.InitLogger(os.Stdout)

	dbConfig, err := global.GetDbConfig()
	if err != nil {
		global.Logger.Fatalln(err)
	}

	err = global.InitDb("postgres", dbConfig)
	if err != nil {
		global.Logger.Fatalln(err)
	}

	flag.Parse()
	workers.StartDispatcher(*NWorkers)

	global.Logger.Println("server listening at: ", *HTTPAddr)

	if err := http.ListenAndServe(*HTTPAddr, handlers.NewRouter()); err != nil {
		global.Logger.Fatalln(err.Error())
	}
}
