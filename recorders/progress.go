/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package recorders

import (
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
	"github.com/pinterest/bender"
)

// NewLoadTestProgress returns progress and bar ready to be used in a recoreder.
func NewLoadTestProgress(count int) (*uiprogress.Progress, *uiprogress.Bar) {
	progress := uiprogress.New()

	// We want to print progress on stderr so results can be easily redirected
	progress.SetOut(os.Stderr)

	// Create new progress bar displaying ELAPSED, CURRENT/MAX and COMPLETED
	bar := progress.AddBar(count)

	bar.PrependElapsed()
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d / %d", b.Current(), count)
	})
	bar.AppendCompleted()

	return progress, bar
}

// NewProgressBarRecorder creates a new progress bar recorder.
func NewProgressBarRecorder(bar *uiprogress.Bar) bender.Recorder {
	return func(msg interface{}) {
		//nolint:gocritic
		switch msg.(type) {
		case *bender.EndRequestEvent:
			bar.Incr()
		}
	}
}
