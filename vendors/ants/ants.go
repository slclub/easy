package ants

import (
	"errors"
	ants "github.com/panjf2000/ants/v2"
	"github.com/slclub/easy/log"
	"runtime"
)

var ants_pool *ants.Pool

// please remenber use release method to free your memory
func Pool() *ants.Pool {

	if ants_pool == nil {
		var err error = errors.New("")
		ants_pool, err = ants.NewPool(runtime.NumCPU()*8, ants.WithLogger(log.Log()))
		if err != nil {
			panic(any(err))
		}
	}
	return ants_pool
}
