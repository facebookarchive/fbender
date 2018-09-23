/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package input

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/facebookincubator/fbender/cmd/core/runner"
	"github.com/facebookincubator/fbender/log"
)

// Transformer converts input line into a request
type Transformer func(string) (interface{}, error)

// Modifier changes request right before sending
type Modifier func(interface{}) (interface{}, error)

// NewRequestGenerator reads data from the specified input and converts it into
// requests using given transformer. The lines which aren't formatted correctly
// are skipped. The requests are then reused in a round-robin manner inside the
// generator. If modifiers are provided they are applied to the request every
// time just before being returned.
func NewRequestGenerator(filename string, transformer Transformer, mods ...Modifier) (runner.RequestGenerator, error) {
	file, err := open(filename)
	if err != nil {
		return nil, err
	}
	defer close(file)
	data := parse(file, transformer)
	if len(data) < 1 {
		return nil, errors.New("at least one valid input line is required")
	}
	return func(i int) interface{} {
		var err error
		request := data[i%len(data)]
		for _, mod := range mods {
			request, err = mod(request)
			if err != nil {
				panic(err)
			}
		}
		return request
	}, nil
}

// parse reads data from the specified input and converts it into requests
// using given transformer. The lines which are not formatted correctly are
// skipped and a warning message is printed to stderr.
func parse(input io.Reader, transformer Transformer) []interface{} {
	requests := make([]interface{}, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		request, err := transformer(line)
		if err != nil {
			log.Errorf("Warning: Error parsing input line %q: %v\n", line, err)
		} else {
			requests = append(requests, request)
		}
	}
	return requests
}

func open(filename string) (*os.File, error) {
	if len(filename) == 0 {
		log.Errorf("Reading input lines until EOF:\n")
		return os.Stdin, nil
	}
	return os.Open(filename)
}

func close(file io.Closer) {
	if file == os.Stdin {
		return
	}
	if err := file.Close(); err != nil {
		log.Errorf("Warning: Error closing input file: %v\n", err)
	}
}
