package project

// ProgressFromTasks calculates progress percentage from tasks (completed / total * 100).
// Returns 0 if there are no tasks.
func ProgressFromTasks(tasks []Task) int {
	if len(tasks) == 0 {
		return 0
	}
	completed := 0
	for _, t := range tasks {
		if t.Status == TaskStatusCompleted {
			completed++
		}
	}
	return completed * 100 / len(tasks)
}
