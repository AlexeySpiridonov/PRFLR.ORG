package config

/**
 * Global variables
 */
const (
    BaseDir                 = ""
    UserCookieName          = "prflr.User.ApiKey"
    DomainName              = "prflr.loc"

    DebugLogFilePath        = "/var/log/prflr/debug.log"
    ErrorLogFilePath        = "/var/log/prflr/error.log"

    DBName                  = "prflr"
    DBHosts                 = "127.0.0.1"
    DBTimers                = "timers"
    DBUsers                 = "users"

    UDPPort                 = ":4000"
    HTTPPort                = ":8080"

    CappedCollectionMaxByte = 200000000 // 200Mb
    CappedCollectionMaxDocs = 1000000   // 1M timers   

    RegisterEmailFrom       = "robot@prflr.org"
    RegisterEmailTo         = "info@prflr.org"
    RegisterEmailSubject    = "New Account Registered!"

    RecoveryEmailFrom       = "robot@prflr.org"
    RecoveryEmailSubject    = "Your PRFLR Account Password"

    FlurryAppKey            = "VG4NW39X4ZYWBYSNNT5P"
)