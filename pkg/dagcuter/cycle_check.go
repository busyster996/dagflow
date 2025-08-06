package dagcuter

func HasCycle(tasks map[string]Task) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(taskName string) bool
	dfs = func(taskName string) bool {
		if recStack[taskName] {
			// Cycle detected
			return true
		}
		if visited[taskName] {
			// Already processed
			return false
		}

		visited[taskName] = true
		recStack[taskName] = true

		for _, dep := range tasks[taskName].Dependencies() {
			if dfs(dep) {
				return true
			}
		}

		recStack[taskName] = false
		return false
	}

	for taskName := range tasks {
		if !visited[taskName] {
			if dfs(taskName) {
				return true
			}
		}
	}

	return false
}
