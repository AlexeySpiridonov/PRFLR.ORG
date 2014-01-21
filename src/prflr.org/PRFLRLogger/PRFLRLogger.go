package PRFLRLogger

import(
	"prflr.org/config"
	"log"
	"os"
)


func Debug(v ...interface{}) {
	f, err := getFile(config.DebugLogFilePath)

	if err != nil {
		return
	}
	defer f.Close()

	// @TODO: Include runtime.Stack() trace

	log.SetOutput(f)
	log.Print(v)
}

func Error(v ...interface{}) {
	f, err := getFile(config.ErrorLogFilePath)

	if err != nil {
		return
	}
	defer f.Close()

	// @TODO: Include runtime.Stack() trace

	log.SetOutput(f)
	log.Print(v)
}

// Just a wrapper
func Fatal(v ...interface{}) {
	log.Fatal(v)
}

/* NOT EXPORTED */
func getFile(filename string) (file *os.File, err error) {
	f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)

	if err != nil {
	    return nil, err
	}

	return f, nil
}