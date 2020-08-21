package database

// dynamodb field names
const (
	keyNameRedirectsTable = "redirect_from"
	attrNameDescription   = "description"
	attrNameRedirect      = "redirect_to"
	attrNamePromotable    = "promotable"
	attrNameCode          = "code"
	attrNameCreated       = "created"
	attrNameModified      = "modified"
	attrNamePathMap       = "cpath_map"
	keyNameCertCacheTable = "domain"
	attrNameCacheValue    = "cert"
	attrNameCreatedAt     = "created_at"
	keyNameUsersTable     = "username"
	attrNamePwd           = "pwd"
	defaultDynamoUser     = "admin"
)
