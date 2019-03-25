package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/time/rate"

	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/store"
	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/transaction"
	"github.com/tetrateio/training/samples/modernbank/tools/trafficGen/pkg/user"
)

var (
	userStore store.Interface
	wg        sync.WaitGroup
	host      = flag.String("host", "localhost", "The Modernbank External IP")
)

func main() {
	flag.Parse()
	userStore = store.NewInMemory()
	ctx, cancel := context.WithCancel(context.Background())

	userCreator := user.NewCreator(*host, userStore, rate.Limit(1.0))                // 1qps
	transactionCreator := transaction.NewCreator(*host, userStore, rate.Limit(10.0)) // 10qps

	go waitForCompletion(ctx, userCreator.Run)
	go waitForCompletion(ctx, transactionCreator.Run)
	go waitForCompletion(ctx, transactionCreator.RunListTransactions)

	// Wait for SIGTERM or SIGINT then clean up
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	wg.Wait()
}

func waitForCompletion(ctx context.Context, fn func(context.Context)) {
	wg.Add(1)
	fn(ctx)
	wg.Done()
}

// Creates transactions between created users at a rate of 20 p/s
