package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/objects"
	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/pipeline"
	"github.com/fatih/color"
	"github.com/google/uuid"
)

func main() {
	dirName := uuid.NewString()
	cleanupDone := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	once := sync.Once{}

	cleanup := func(dirName string) {
		defer close(cleanupDone)
		color.Blue("Cleaning up...")
		// Change dir to parent
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		parentDir := filepath.Dir(currentDir)

		err = os.Chdir(parentDir)
		if err != nil {
			log.Println("unablr to change to parent dir")
		}

		err = os.RemoveAll(dirName)
		if err != nil {
			log.Println("error in removing dir ", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	defer once.Do(func() {
		cleanup(dirName)
	})

	go func() {
		<-sig
		once.Do(func() {
			cancel()
			cleanup(dirName)

		})
	}()

	// FIXME: can be refactored and improve readablity with custom contexts
	// will do in the upcoming iteration
	ctx = context.WithValue(ctx, objects.PipelineValueDirectoryName, dirName)

	pipeline.Run(ctx)

	<-cleanupDone
}
