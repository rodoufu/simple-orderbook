package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/rodoufu/simple-orderbook/pkg/engine"
	"github.com/rodoufu/simple-orderbook/pkg/event"
	"github.com/rodoufu/simple-orderbook/pkg/io"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	formatter := &logrus.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logger := logrus.New()
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stderr)
	log := logger.WithFields(logrus.Fields{})

	log.Info("staring service")

	// The io.ReadTransactions creates a goroutine to read the file
	transactions, err := io.ReadTransactions(ctx, "input_file.csv")
	if err != nil {
		log.WithError(err).Fatal("problem loading transactions parser")
	}

	// Writing to stdout in a specific goroutine.
	toOutput := make(chan event.Output)
	go func() {
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case output, ok := <-toOutput:
				if !ok {
					return
				}
				if msg := output.Output(); len(msg) > 0 {
					fmt.Println(msg)
				}
			}
		}
	}()

	mktEngine, events := engine.NewListEngine()
	defer mktEngine.Close()
	go func() {
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case evt, ok := <-events:
				if !ok {
					return
				}
				if output, ok := evt.(event.Output); ok {
					toOutput <- output
				}
			}
		}
	}()

	for transaction := range transactions {
		if err = mktEngine.ProcessTransaction(ctx, transaction); err != nil {
			log.WithError(err).Error("problem processing transaction")
		}
	}
}
