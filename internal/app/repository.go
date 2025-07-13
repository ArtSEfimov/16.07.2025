package app

type Repository struct {
	taskCounter map[string]uint8
}

func NewRepository() *Repository {
	return &Repository{
		taskCounter: make(map[string]uint8),
	}
}

func (repository *Repository) AddNewUser(userUUID string) {
	repository.taskCounter[userUUID] = 0
}

func (repository *Repository) GetUserTaskCount(userUUID string) uint8 {
	return repository.taskCounter[userUUID]
}
