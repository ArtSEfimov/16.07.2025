package app

type Repository struct {
	allTasks map[string][]*Task
}

func NewRepository() *Repository {
	return &Repository{
		allTasks: make(map[string][]*Task),
	}
}

func (repository *Repository) AddTask(userUUID string, task *Task) {
	repository.allTasks[userUUID] = append(repository.allTasks[userUUID], task)
}

func (repository *Repository) GetTaskByID(userUUID string, id uint64) *Task {
	allUserTasks := repository.GetUserTasks(userUUID)
	for _, task := range allUserTasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (repository *Repository) GetActiveTaskCount() int {
	var allActiveTasks = 0
	for _, tasks := range repository.allTasks {
		for _, task := range tasks {
			if task.Status != taskStatusCompleted {
				allActiveTasks++
			}
		}
	}
	return allActiveTasks
}

func (repository *Repository) isUserHasTask(userUUID string) bool {
	if len(repository.allTasks[userUUID]) == 0 {
		return false
	}
	return true
}

func (repository *Repository) GetUserTasks(userUUID string) []*Task {
	return repository.allTasks[userUUID]
}
