package cleanup

type CleanTask func()

// Cleaner holds whatever needs to be done
type Cleaner struct {
	cleaned bool
	onerror []CleanTask
	onnil   []CleanTask
	always  []CleanTask
	Errptr  *error
}

func (c *Cleaner) GetError() error {
	if c.Errptr != nil {
		return *c.Errptr
	}
	return nil
}

func (c *Cleaner) ComplementOrSetErr(complementor func(mainErr error, complementErr ...error) error, complementErr error) {

	ierr := c.GetError()
	if ierr != nil {

		c.SetReturnError(complementor(ierr, complementErr))

	} else {

		c.SetReturnError(complementErr)

	}

}

func (c *Cleaner) SetReturnError(err error) {
	*c.Errptr = err
}

// NewCleaner creates a cleanup instance,
// use it at the beginning of your function
// hooking up your error pointer.
func NewCleaner(errPtr *error) (c *Cleaner) {
	c = new(Cleaner)
	c.Errptr = errPtr
	return
}

// Clean performs the cleaning,
// call it as a defer at the beginning of your function
// or whitin the function context.
// its executed once, async unsafe.
func (c *Cleaner) Clean() {
	if c.cleaned {
		return
	}
	for i := len(c.always) - 1; i >= 0; i-- {
		(c.always[i])()
	}
	if c.Errptr != nil {
		if (*c.Errptr) == nil {
			for i := len(c.onnil) - 1; i >= 0; i-- {
				(c.onnil[i])()
			}
		} else {
			for i := len(c.onerror) - 1; i >= 0; i-- {
				(c.onerror[i])()
			}
		}
	}
}

// OnError executes fx when err is not nil,
// just call it when you need it.
func (c *Cleaner) OnError(fx CleanTask) {
	c.onerror = append(c.onerror, fx)
}

// OnNil executes fx when err is nil,
// just call it when you need it.
func (c *Cleaner) OnNil(fx CleanTask) {
	c.onnil = append(c.onnil, fx)
}

// Always executes always,
// just call it when you need it.
func (c *Cleaner) Always(fx CleanTask) {
	c.always = append(c.always, fx)
}

func Always(errPtr *error, do func(error)) {
	do(*errPtr)
}

func OnError(errPtr *error, do func(error)) {
	if (*errPtr) != nil {
		do(*errPtr)
	}
}

func OnNil(errPtr *error, do func()) {
	if (*errPtr) == nil {
		do()
	}
}

func OnTrue(boolPtr *bool, do func()) {
	if *boolPtr {
		do()
	}
}

func OnFalse(boolPtr *bool, do func()) {
	if !(*boolPtr) {
		do()
	}
}

func ExecIfSet(funcPtr *func()) {
	if funcPtr != nil {
		(*funcPtr)()
	}
}
