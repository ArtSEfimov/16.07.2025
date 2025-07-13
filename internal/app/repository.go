package app

type Repository struct {
	userObjects map[string][]Link
	taskCounter uint8
}

func NewRepository() *Repository {
	return &Repository{
		userObjects: make(map[string][]Link),
	}
}

func (repository *Repository) AddNewUser(userUUID string) {
	repository.userObjects[userUUID] = make([]Link, 0, 3)
}

func (repository *Repository) GetUserObjectCount(userUUID string) int {
	return len(repository.userObjects[userUUID])
}

func (repository *Repository) AddNewObject(userUUID string, link Link) {
	repository.userObjects[userUUID] = append(repository.userObjects[userUUID], link)
}

func (repository *Repository) AddTaskCounter() {
	repository.taskCounter++
}

func (repository *Repository) GetTaskCount() uint8 {
	return repository.taskCounter
}
