package app

const sessionID = "session_id"

const (
	taskLimit        = 3
	taskLimitMessage = "task limit exceeded"
)
const filesLimit = 3
const getLinksPath = "/links"

const (
	taskStatusCreated    = "created"
	taskStatusPending    = "pending"
	taskStatusProcessing = "processing"
	taskStatusDone       = "done"
	taskStatusError      = "error"
)

const (
	errInvalidLinkFormat      = "invalid link format"
	errUnsupportedContentType = "unsupported content type"
	errInaccessibleLink       = "inaccessible link"
	errLinkLimitExceeded      = "link limit exceeded"
)
