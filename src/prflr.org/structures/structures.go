package structures

/**
 * Web panel Struct
 */
type Stat struct {
	Src   string
	Timer string
	Count int
	Total float32
	Min   float32
	Avg   float32
	Max   float32
}

/**
 * UDP Package struct
 */
type Timer struct {
	Thrd   string
	Src    string
	Timer  string
	Time   float32
	Info   string
	Apikey string
}

/**
 * User struct
 */
type User struct {
	Email   string
	Password    string
	Apikey  string
	Token   float32
	Info   string
}