/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package core

import (
	"sync"
)

var (
	postinitLock      = make(chan struct{}, 1)
	postinitWaitGroup = &sync.WaitGroup{}
)

// DeferPostInit defers an execution of postinit function until a StartPostInit
// is called.
func DeferPostInit(postinit func()) {
	postinitWaitGroup.Add(1)
	go func() {
		<-postinitLock
		postinit()
		postinitWaitGroup.Done()
		postinitLock <- struct{}{}
	}()
}

// StartPostInit starts all postinit functions and blocks execution until all of
// them are finished.
func StartPostInit() {
	postinitLock <- struct{}{}
	postinitWaitGroup.Wait()
}
