package pipeline

import (
	"context"
	"fmt"
	"os/exec"
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

				cmd := exec.Command("git", "clone", cloneStage.RepositoryURL)
				stdout, err := cmd.Output()
				if err != nil {
					errorMsg.Error = err
					p.Err = errorMsg
					outStream <- p
					continue
				}
				fmt.Println(string(stdout))
			}
		}
	}()

	return outStream
}
