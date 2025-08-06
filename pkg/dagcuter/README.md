# Dagcuter

Dagcuter is a Go library for executing Directed Acyclic Graphs (DAGs) of tasks. It allows you to define tasks with dependencies, execute them in the correct order, and handle pre-execution, execution, and post-execution phases.

## Features

- **Task Dependency Management**: Automatically resolves and executes tasks based on their dependencies.
- **Cycle Detection**: Validates the DAG to ensure there are no circular dependencies.
- **Concurrent Execution**: Executes independent tasks concurrently for better performance.
- **Customizable Task Lifecycle**: Supports `PreExecution`, `Execute`, and `PostExecution` phases for each task.
- **Retry Policy**: Supports retry policies for failed tasks.

## Installation

To use Dagcuter in your project, add it to your `go.mod` file:

```bash
go get github.com/busyster996/dagflow/pkg/dagcuter
```

## Usage

Here’s an example of how to use Dagcuter:

```go
package main

import (
    "context"
    "fmt"

    "github.com/busyster996/dagflow/pkg/dagcuter"
)

type ExampleTask struct {
    name         string
    dependencies []string
}

func (t *ExampleTask) Name() string {
    return t.name
}

func (t *ExampleTask) Dependencies() []string {
    return t.dependencies
}

func (t *ExampleTask) RetryPolicy() *dagcuter.RetryPolicy {
	return nil // No retry policy for mock tasks
}

func (t *ExampleTask) PreExecution(ctx context.Context, input map[string]any) error {
    fmt.Printf("PreExecution for task: %s\n", t.name)
    return nil
}

func (t *ExampleTask) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
    fmt.Printf("Executing task: %s\n", t.name)
    return map[string]any{"result": fmt.Sprintf("output of %s", t.name)}, nil
}

func (t *ExampleTask) PostExecution(ctx context.Context, output map[string]any) error {
    fmt.Printf("PostExecution for task: %s\n", t.name)
    return nil
}

func main() {
    tasks := map[string]dagcuter.Task{
        "task1": &ExampleTask{name: "task1", dependencies: []string{}},
        "task2": &ExampleTask{name: "task2", dependencies: []string{"task1"}},
        "task3": &ExampleTask{name: "task3", dependencies: []string{"task2"}},
    }

    dag, err := dagcuter.NewDagcuter(tasks)
    if err != nil {
        panic(err)
    }

    results, err := dag.Execute(context.Background())
    if err != nil {
        panic(err)
    }

    fmt.Println("Execution Results:", results)
    fmt.Println("Execution Order:", dag.ExecutionOrder())
}
```

## Task Interface

To define a task, implement the following interface:

```go
type Task interface {
    Name() string
    Dependencies() []string
    RetryPolicy() *dagcuter.RetryPolicy
    PreExecution(ctx context.Context, input map[string]any) error
    Execute(ctx context.Context, input map[string]any) (map[string]any, error)
    PostExecution(ctx context.Context, output map[string]any) error
}
```

## Example

You can find a complete example in the examples/simple directory.

## Testing

Dagcuter includes a comprehensive suite of unit tests. To run the tests, use:

```bash
go test ./...
```
