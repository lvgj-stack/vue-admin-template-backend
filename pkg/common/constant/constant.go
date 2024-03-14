package constant

const (
	// XUsernameKey defines the key in gin context which represents the owner of the token.
	XUsernameKey = "X-Username"

	// XRequestIDKey defines the key in gin context which represents the uuid of the request.
	XRequestIDKey = "X-Request-ID"
)

type TableStatus string

const (
	Published TableStatus = "published"
	Draft     TableStatus = "draft"
	Deleted   TableStatus = "deleted"
)
