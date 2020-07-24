package authtest

import (
	"fmt"
	"sync/atomic"
)

type factory struct {
	num int32
}

func (f *factory) getPhoneNumber() string {
	atomic.AddInt32(&f.num, 1)
	return fmt.Sprintf("8%.10d", f.num)
}

func (f *factory) getEmail() string {
	atomic.AddInt32(&f.num, 1)
	return fmt.Sprintf("test%.3d@example.com", f.num)
}

func (f *factory) getLogin() string {
	atomic.AddInt32(&f.num, 1)
	return fmt.Sprintf("username%.3d", f.num)
}
