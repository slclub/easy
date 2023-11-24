package option

/*
Application scenario:

	Initialize
	Configurition.

Very poor performance.It does not suitable for business logical program or code of framework.
It is better to running only once.
*/
type Assignment interface {
	// set target object to Option
	Target(any)
	// run the assign movement.
	Apply()
	// Assign default values to Target object.
	Default(...ItemAssignment) Assignment

	Final(...ItemAssignment) Assignment
}

type ItemAssignment interface {
	Apply() (string, any)
}

/*
@return

	string ； it is a key
	any ; it is a value corresponding to the key
*/
type OptionFunc func() (string, any)

func (fn OptionFunc) Apply() (string, any) {
	// we can add somethings to here
	return fn()
}
