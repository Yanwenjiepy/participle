package main

import (
	"flag"
	"log"
	"participle/internal"
	"participle/logger"
	"participle/pkg/rpc_server"
	"time"
)

func main() {

	csvFilePath := flag.String("csv", "", "csv file path")
	logFilePath := flag.String("log", "", "log file path")
	logLevel := flag.String("level", "info", "log level: debug, info, warning, error, fatal")
	serverPort := flag.String("port", ":9001", "rpc server port")

	flag.Parse()

	if *csvFilePath == "" || *logFilePath == "" {
		log.Println("must be input csv file path and log file path!")
		return
	}

	if (*serverPort)[0] != ':' {
		log.Println("Wrong format of server port!")
		return
	}

	logger.Log = logger.InitLog(*logFilePath, *logLevel)

	var accountList []internal.Account

	err := internal.LoadCSV(*csvFilePath, &accountList)
	if err != nil {
		logger.Log.Error("wrong format or wrong path of csv file!")
		return
	}

	tokenListOne := internal.GetAccessToken(&accountList)
	tokenList := make([]string, 0, 2*len(tokenListOne))
	tokenList = append(tokenList, tokenListOne...)
	tokenList = append(tokenList, tokenListOne...)

	internal.New(2 * len(tokenListOne))

	go func() {
		for {

			for _, token := range tokenList {
				err = internal.AddToken(token)
				if err != nil {

					logger.Log.Info(err.Error())
					break
				}
			}

			time.Sleep(1 * time.Second)
		}

	}()

	time.Sleep(1 * time.Second)

	err = rpc_server.Server(*serverPort)
	if err != nil {
		log.Fatal(err)
	}
}
