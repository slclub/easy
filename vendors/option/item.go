package option

var (
	// If the corresponding value is 0, use default
	DEFAULT_IGNORE_ZERO OptionFunc
)

func init() {
	DEFAULT_IGNORE_ZERO = defaultForce
}
func defaultForce() (string, any) {
	return "OPTION_DEFAULT_IGNORE_ZERO", true
}
