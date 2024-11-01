package main

import (
	"context"
	"fmt"

	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/pipeline"
)

func main() {
	fmt.Println("hello world")
	pipeline.Run(context.Background())

}
