# simple-ci-pipeline-runner

This project  is a exercise  to understand ``Pipeline`` pattern using ``Channels``. Pipeline has 3 stages
    - Clone
    - Test
    - Build

- **Clone**: In this stage the URL of the repository that is to be cloned as passed as argument, The repository is cloned with a random file name in the current working directory.

- **Test**: This stage is a optional , The current working directory is switched to the cloned project directory , The directory name is passed downstream through ``context``. And the configured command for the tests is run.

- **Build**: In this stage the command configured for build is run and this marks the end of the pipeline


- The  possible arguments that a  task can take is given in the below struct

```go
    type TaskArgs struct {
        RepositoryURL string `json:"repository_url,omitempty"`
        Command       string `json:"command,omitempty"`
    }
```

- Below is the Pipeline object

```go
type Pipeline struct {
	Tasks map[TaskStage]TaskArgs
	Err   *Error
}
```

## Note
This is not production ready , This code directly runs on the host machine which is serious security issue and there are multiple features missing and flaws are there , This is just a  way to put the pipeline pattern into action.
 
## Additional Links

[Demo](https://youtu.be/dPiFzJVhjZI)

[Sample Repo](https://github.com/VarthanV/sample-nodejs-app)

[Pipeline pattern reading](https://github.com/VarthanV/go-concurrency-exercises/tree/main/pipelines)

