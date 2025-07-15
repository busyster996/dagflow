package dagcuter

import "context"

type Task interface {
	Name() string
	Dependencies() []string
	// RetryPolicy returns the retry policy for the task
	RetryPolicy() *RetryPolicy
	// PreExecution is called before Execute
	PreExecution(ctx context.Context, input map[string]any) error
	// Execute is called after PreExecution
	Execute(ctx context.Context, input map[string]any) (map[string]any, error)
	// PostExecution is called after Execute
	PostExecution(ctx context.Context, output map[string]any) error
}
