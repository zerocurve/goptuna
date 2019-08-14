package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func objective(trial goptuna.Trial) (float64, error) {
	ctx := trial.GetContext()

	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)

	cmd := exec.CommandContext(ctx, "sleep", "1")
	err := cmd.Run()
	if err != nil {
		return -1, err
	}
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		logger.Fatal("failed to open db", zap.Error(err))
	}
	defer db.Close()
	rdb.RunAutoMigrate(db)

	// create a study
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionStorage(rdb.NewStorage(db)),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSetLogger(logger),
	)
	if err != nil {
		logger.Fatal("failed to create study", zap.Error(err))
	}

	// create a context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	study.WithContext(ctx)

	// set signal handler
	signalch := make(chan os.Signal, 1)
	defer close(signalch)
	signal.Notify(signalch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// run optimize with context
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		sig, ok := <-signalch
		if !ok {
			return
		}
		cancel()
		logger.Error("Catch a kill signal", zap.String("signal", sig.String()))
	}()
	go func() {
		defer wg.Done()
		err = study.Optimize(objective, 10)
	}()
	wg.Wait()
	if err != nil {
		logger.Fatal("got error while run optimize", zap.Error(err))
	}

	v, _ := study.GetBestValue()
	fmt.Println("Best evaluation value:", v)
}
