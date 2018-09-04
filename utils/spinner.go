/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils

import (
	"context"
	"time"

	spin "github.com/tj/go-spin"

	"github.com/facebookincubator/fbender/log"
)

// Default refresh rate
const spinnerRefresh = 100 * time.Millisecond

// NewBackgroundSpinner creates a new spinner which runs in background refreshing
// its output on a constant rate. The spinner is prefixed with the provided
// description. Returns a function which cancels the spinner. The default value
// for refresh is used if provided refresh is less than or equal to zero.
func NewBackgroundSpinner(description string, refresh time.Duration) context.CancelFunc {
	if refresh <= 0 {
		refresh = spinnerRefresh
	}

	ctx, cancel := context.WithCancel(context.Background())
	spinner := spin.New()
	spinner.Set(spin.Spin1)
	sync := make(chan bool)

	go func() {
		for {
			log.Errorf("\r%s... ", description)
			select {
			case <-ctx.Done():
				log.Errorf("Done.\n")
				sync <- false
				return
			default:
				log.Errorf("%s", spinner.Next())
			}
			time.Sleep(refresh)
		}
	}()

	return func() {
		// cancel context and wait for goroutine to clean the output
		cancel()
		<-sync
	}
}
