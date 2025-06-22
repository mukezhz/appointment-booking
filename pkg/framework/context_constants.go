package framework

const (
	// Authentication and Authorization
	Claims      = "Claims"
	UID         = "UID"
	UIDHeader   = "X-UID"
	Token       = "Token"
	CognitoPass = "CognitoPass"
	Role        = "Role"

	// Context Keys
	RequestIDKey   = "request_id"
	LoggerKey      = "logger"
	CorrelationKey = "correlation_id"

	// File upload
	File = "@uploaded_file"

	// Pagination
	Limit = "Limit"
	Page  = "Page"

	// Rate Limiting
	RateLimit = "RateLimit"
)
