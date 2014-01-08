package statstruct

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