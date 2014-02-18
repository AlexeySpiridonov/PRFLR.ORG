package config

/**
 * Global variables
 */
const (
    BaseDir                 = "src/prflr.org/"
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

    CappedCollectionMaxByte = 450000000 // 450Mb
    CappedCollectionMaxDocs = 2000000   // 2M timers   

    RegisterEmailFrom       = "robot@prflr.org"
    RegisterEmailTo         = "info@prflr.org"
    RegisterEmailSubject    = "New Account Registered!"

    RecoveryEmailFrom       = "robot@prflr.org"
    RecoveryEmailSubject    = "Your PRFLR Account Password"

    FlurryAppKey            = "VG4NW39X4ZYWBYSNNT5P"
)