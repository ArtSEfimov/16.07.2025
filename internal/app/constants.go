package app

const sessionID = "session_id"

const (
	taskLimit        = 3
	taskLimitMessage = "task limit exceeded"
)
const filesLimit = 3
const (
	getLinksPath = "/create-task"
	getTaskStatusPath = "/get-status"
	addLinkPath = "/add-link"
)

const (
	taskStatusCreated    = "created"
	taskStatusPending    = "pending"
	taskStatusProcessing = "processing"
	taskStatusCompleted       = "completed"
	taskStatusError      = "error"
)

const (
	errInvalidLinkFormat      = "invalid link format"
	errUnsupportedContentType = "unsupported content type"
	errInaccessibleLink       = "inaccessible link"
	errLinkLimitExceeded      = "link limit exceeded"
)

const (
	const errUserHasNoTasks = "user has no any tasks"
	const ErrTaskAlreadyExists = fmt.Sprintf("Task already exists. Use the URL %s to add a new link.", addLinkPath)
)
