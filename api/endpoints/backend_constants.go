package endpoints

const (
	usernameRegex          = "^[a-zA-Z0-9_.]+$"
	emailRegex             = "[a-zA-Z0-9.-_]{1,}@[a-zA-Z.-]{2,}[.]{1}[a-zA-Z]{2,}"
	lettersRegex           = "^[a-zA-Z]+$"
	alphanumericRegex      = "^[a-zA-Z0-9]+$"
	alphanumericSpaceRegex = "^[A-Za-z0-9 ]+$"
	passRegex              = "^[a-zA-Z0-9!\\-@#$%^&*_]*$"
	urlRegex               = "^([a-zA-Z0-9][a-zA-Z0-9-_]*\\.)*[a-zA-Z0-9]*[a-zA-Z0-9-_]*[[a-zA-Z0-9]+$"
	ipRegex                = "^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"

	upperBoundNameLen = 100
	lowerBoundNameLen = 2

	oauthBoundLen = 3000

	lowerBoundEmailLen = 5

	ticketParam  = "ticket"
	scannerParam = "scanner"
	idParam      = "id"
	userParam    = "user"
	orgParam     = "org"
	sourceParam  = "source"

	baseAdmin  = 16
	admin      = 8
	manager    = 4
	reader     = 2
	reporter   = 1
	allAllowed = admin | manager | reader | reporter
)
