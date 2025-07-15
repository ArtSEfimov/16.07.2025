package app

const sessionID = "session_id"

const (
	taskLimit        = 3
	taskLimitMessage = "task limit exceeded"
	filesLimit       = 3
)

const (
	createTaskPath    = "/create-task"
	getTaskStatusPath = "/get-status"
	addLinkPath       = "/add-link"
)

const (
	taskStatusCreated    = "created"
	taskStatusPending    = "pending"
	taskStatusProcessing = "processing"
	taskStatusCompleted  = "completed"
	taskStatusError      = "error"
)

const (
	errInvalidLinkFormat      = "invalid link format"
	errUnsupportedContentType = "unsupported content type"
	errInaccessibleLink       = "inaccessible link"
	errZipFileCreation        = "failed to create ZIP archive. Please try again later."
	errCannotAddLinkToTask    = "cannot add link to task"
)

const (
	errUserHasNoTaskByID = "user has no tasks by ID"
	errUserNotFound      = "user not found"
)

const baseOutputZipPath = "static"
