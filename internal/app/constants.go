package app

const sessionID = "session_id"

const (
	taskLimit        = 3
	taskLimitMessage = "task limit exceeded"
)
const objectsLimit = 3
const getLinksPath = "/links"

const (
	taskStatusCreated    = "created"
	taskStatusPending    = "pending"
	taskStatusProcessing = "processing"
	taskStatusDone       = "done"
	taskStatusError      = "error"
)

const (
	errInvalidLinkFormat   = "invalid link format"
	errUnsupportedFileType = "unsupported file type"
	errInaccessibleLink    = "inaccessible link"
)