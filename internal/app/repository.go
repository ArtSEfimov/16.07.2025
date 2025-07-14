package app

type Repository struct {
	allTasks map[uint64]Task
	userTaskData
}

func NewRepository() *Repository {
	return &Repository{
		allTasks: make(map[uint64]Task),
		userTaskData: userTaskData{
			userLinks: make(map[string][]Link),
			userTask:  make(map[string]uint64),
		},
	}
}

func (repository *Repository) AddTask(task Task) {
	repository.allTasks[task.ID] = task
}

func (repository *Repository) GetTaskByID(id uint64) Task {
	return repository.allTasks[id]
}

func (repository *Repository) DeleteTask(task Task) {
	repository.userTask = make(map[string]uint64)
	delete(repository.allTasks, task.ID)
}

func (repository *Repository) GetTaskCount() int {
	return len(repository.allTasks)
}
