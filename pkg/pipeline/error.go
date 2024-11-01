package pipeline

import "errors"

var (
	ErrInvalidRepoURL     = errors.New("invalid repo_url")
	ErrStageCloneRequired = errors.New("stage clone required")
)
