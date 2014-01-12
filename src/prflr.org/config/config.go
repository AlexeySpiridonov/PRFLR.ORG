package config

/**
 * Global variables
 */
const (
    BaseDir                 = "src/prflr.org/"
    UserCookieName          = "prflr.User.ApiKey"
    DomainName              = "prflr.loc"

	DBName                  = "prflr"
	DBHosts                 = "127.0.0.1"
	DBTimers                = "timers"
	DBUsers                 = "users"

	UDPPort                 = ":4000"
	HTTPPort                = ":8080"

	CappedCollectionMaxByte = 100000000 // 100Mb
	CappedCollectionMaxDocs = 500000    // 500k
)