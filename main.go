package main

import (
	"context"
	"log"
	"os"

	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/objects"
	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/pipeline"
	"github.com/fatih/color"
	"github.com/google/uuid"
)

func main() {
	dirName := uuid.NewString()

	defer func(dirName string) {
		color.Blue("Cleaning up...")
		err := os.RemoveAll(dirName)
		if err != nil {
			log.Println("error in removing dir ", err)
		}
	}(dirName)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// FIXME: can be refactored and improve readablity with custom contexts
	// will do in the upcoming iteration
	ctx = context.WithValue(ctx, objects.PipelineValueDirectoryName, dirName)
	pipeline.Run(ctx)
	panic("stacktrace")
}
