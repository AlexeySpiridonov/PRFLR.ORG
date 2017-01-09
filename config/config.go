package config

/**
 * Global variables
 */
const (
	BaseDir        = ""
	UserCookieName = "prflr.auth.cookie"
	DomainName     = "prflr.org"

	DBName   = "prflr"
	DBHosts  = "127.0.0.1"
	DBTimers = "timers"
	DBUsers  = "users"

	UDPPort  = ":4000"
	HTTPPort = ":8080"

	CappedCollectionMaxByte = 50000000
	CappedCollectionMaxDocs = 100000    

	RegisterEmailFrom    = "robot@prflr.org"
	RegisterEmailTo      = "info@prflr.org"
	RegisterEmailSubject = "New Account Registered!"

	RecoveryEmailFrom    = "robot@prflr.org"
	RecoveryEmailSubject = "Your PRFLR Account Password"

	FlurryAppKey = "VG4NW39X4ZYWBYSNNT5P"
)
