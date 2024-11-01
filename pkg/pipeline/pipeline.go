package pipeline

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/VarthanV/simple-ci-pipeline-runner/pkg/objects"
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

				dirNameFromCtx := ctx.Value(objects.PipelineValueDirectoryName)
				if dirNameFromCtx == nil {
					log.Println("filename not set")
					errorMsg.Error = ErrFileNameRequired
					p.Err = errorMsg
					outStream <- p
					continue
				}

				dirName, _ := dirNameFromCtx.(string)

				color.Green("############### Stage1: Cloning Repo ######################")
				cmd := exec.CommandContext(
					ctx,
					"git",
					"clone",
					cloneStage.RepositoryURL,
					dirName,
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

				// Check if testing is configured it is optional can skip
				testStage, ok := p.Tasks[TaskStageTest]
				if !ok {
					outStream <- p
					color.Yellow("Stage2: Testing phase not found , Skipping it...")
					continue
				}

				color.Green("############### Stage2: Running Tests ######################")
				// Change to dir

				dirNameFromCtx := ctx.Value(objects.PipelineValueDirectoryName)
				if dirNameFromCtx == nil {
					log.Println("filename not set")
					errorMsg.Error = ErrFileNameRequired
					p.Err = errorMsg
					outStream <- p
					continue
				}

				dirName, _ := dirNameFromCtx.(string)
				err := os.Chdir(dirName)
				if err != nil {
					log.Println("erorr in changing dir ", err)
					errorMsg.Error = err
					p.Err = errorMsg
					outStream <- p
					continue
				}

				//FIXME: base setup command is hardcoded need to get from user when
				// extending
				cmd := exec.CommandContext(ctx,
					"sh", "-c",
					fmt.Sprintf("npm i && %s", testStage.Command))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					errorMsg.Error = err
					p.Err = errorMsg
					outStream <- p
					continue
				}

				outStream <- p
				color.Green("######### Stage2:successful #########################")

			}
		}
	}()

	return outStream
}

// build: Stage 3 of the pipeline
func build(ctx context.Context, pipeline <-chan Pipeline) <-chan Pipeline {
	outStream := make(chan Pipeline)
	errorMsg := &Error{}

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
				if p.Err != nil {
					outStream <- p
					continue
				}

				// Get build phase task
				taskStageBuild, ok := p.Tasks[TaskStageBuild]
				if !ok {
					errorMsg.Error = ErrStageBuildRequired
					p.Err = errorMsg
					outStream <- p
					continue
				}

				color.Green("############### Stage 3: Build  ######################")

				cmd := exec.CommandContext(ctx,
					"sh", "-c",
					taskStageBuild.Command)
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
				color.Green("######### Stage3:successful #########################")
			}
		}
	}()
	return outStream
}

func Run(ctx context.Context) {
	ch := generate(Pipeline{
		Tasks: map[TaskStage]TaskArgs{
			"clone": {
				RepositoryURL: "https://github.com/VarthanV/sample-nodejs-app",
			},
			"test": {
				Command: "npm run test",
			},
			"build": {
				Command: "npm run build",
			},
		},
	})

	pipeline := build(ctx, test(ctx, clone(ctx, ch)))

	for rslt := range pipeline {
		if rslt.Err != nil {
			color.Red("Pipeline failed in stage %s with error %s", rslt.Err.Stage, rslt.Err.Error)
		}
	}
}
