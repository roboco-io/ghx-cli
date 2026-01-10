package discussion

const (
	// Format constants
	formatJSON    = "json"
	formatTable   = "table"
	formatDetails = "details"

	// State constants
	stateOpen   = "open"
	stateClosed = "closed"
	stateAll    = "all"

	// Display constants
	defaultListLimit          = 20
	defaultCommentLimit       = 50
	tableSeparatorWidth       = 100
	titleMaxLength            = 40
	titleTruncateLength       = 37
	categoryMaxLength         = 15
	categoryTruncateLength    = 12
	authorMaxLength           = 15
	authorTruncateLength      = 12
	bodyPreviewLength         = 100
	viewSeparatorWidth        = 80
	categorySeparatorWidth    = 90
	descriptionTruncateLength = 40

	// Close reasons
	closeReasonResolved  = "resolved"
	closeReasonOutdated  = "outdated"
	closeReasonDuplicate = "duplicate"
)
