package app

type userTaskData struct {
	userLinks map[string][]Link
	userTask  map[string]uint64 // Created or Pending
}

func (repository *Repository) AddNewUser(userUUID string) {
	repository.userLinks[userUUID] = make([]Link, 0, 3)
}

func (repository *Repository) AddUserTask(userUUID string, task Task) {
	repository.userTask[userUUID] = task.ID
	repository.AddTask(task)
}

func (repository *Repository) GetUserLinksCount(userUUID string) int {
	return len(repository.userLinks[userUUID])
}

func (repository *Repository) AddNewUserLink(userUUID string, link Link) {
	repository.userLinks[userUUID] = append(repository.userLinks[userUUID], link)
}

func (repository *Repository) isUserHasTask(userUUID string) bool {
	if repository.userTask[userUUID] == 0 {
		return false
	}
	return true
}

func (repository *Repository) GetUserTaskID(userUUID string) uint64 {
	return repository.userTask[userUUID]
}
