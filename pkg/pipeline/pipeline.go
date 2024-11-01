package pipeline

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

type TaskStage string

const (
	TaskStageClone TaskStage = "clone"
	TaskStageTest  TaskStage = "test"
	TaskStageBuild TaskStage = "build"
)

type TaskArgs struct {
	RepositoryURL string `json:"repository_url,omitempty"`
	Command       string `json:"command,omitempty"`
}

type Error struct {
	Stage TaskStage `json:"stage,omitempty"`
	Error error     `json:"message,omitempty"`
}

type Pipeline struct {
	Tasks map[TaskStage]TaskArgs
	Err   *Error
}

func generate(p ...Pipeline) <-chan Pipeline {
	ch := make(chan Pipeline)
	go func() {
		defer close(ch)
		for _, val := range p {
			ch <- val
		}
	}()
	return ch
}

// clone: Stage 1 of the pipeline
func clone(ctx context.Context, pipeline <-chan Pipeline) <-chan Pipeline {
	outStream := make(chan Pipeline)
	errorMsg := &Error{
		Stage: TaskStageClone,
	}

	go func() {
		defer close(outStream)
		for {
			select {
			case <-ctx.Done():
				return
			case p, ok := <-pipeline:
				if !ok {
					return
				}

				// Check if clone stage is configured
				cloneStage, ok := p.Tasks[TaskStageClone]
				if !ok {
					errorMsg.Error = ErrStageCloneRequired
					p.Err = errorMsg
					outStream <- p
					continue
				}

				color.Green("############### Stage1: Cloning Repo ######################")
				cmd := exec.Command(
					"git",
					"clone",
					cloneStage.RepositoryURL,
					"--progress")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					errorMsg.Error = err
					p.Err = errorMsg
					outStream <- p
					continue
				}

				outStream <- p
				color.Green("######### Stage1:sucessful #########################")
			}
		}
	}()

	return outStream
}

// test: Stage 2 of the pipeline
func test(ctx context.Context, pipeline <-chan Pipeline) <-chan Pipeline {
	outStream := make(chan Pipeline)
	go func() {
		defer close(outStream)
		for {
			select {
			case <-ctx.Done():
				return

			case p, ok := <-pipeline:
				if !ok {
					return
				}

				// Check if testing is configured it is optional can skip
				_, ok = p.Tasks[TaskStageTest]
				if !ok {
					outStream <- p
					continue
				}

			}
		}
	}()

	return pipeline
}

func Run(ctx context.Context) {
	ch := generate(Pipeline{
		Tasks: map[TaskStage]TaskArgs{
			"clone": {
				RepositoryURL: "https://github.com/VarthanV/simple-shell",
			},
		},
	})

	pipeline := clone(ctx, ch)

	for rslt := range pipeline {
		log.Println(rslt)
	}
}
