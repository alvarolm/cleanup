package cleanup

type CleanTask func()

type CleanErrTask func(err error)

// Cleaner holds whatever needs to be done
type Cleaner struct {
	onerror []CleanErrTask
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

func (c *Cleaner) SetReturnError(err error) {
	*c.Errptr = err
}

// NewCleaner creates a cleanup instance,
// use it at the beginning of your function
// hooking up your error pointer.
func NewCleaner(errPtr *error) (curecorder *Cleaner) {
	curecorder = new(Cleaner)
	curecorder.Errptr = errPtr
	return
}

// Clean performs the cleaning,
// call it as a defer at the beginning of your function.
func (c *Cleaner) Clean() {
	for i := len(c.always) - 1; i >= 0; i-- {
		(c.always[i])()
	}
	if (*c.Errptr) == nil {
		for i := len(c.onnil) - 1; i >= 0; i-- {
			(c.onnil[i])()
		}
	} else {
		for i := len(c.onerror) - 1; i >= 0; i-- {
			(c.onerror[i])(*c.Errptr)
		}
	}
}

// OnError executes fx when err is not nil,
// just call it when you need it.
func (c *Cleaner) OnError(fx CleanErrTask) {
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
