/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/facebookincubator/fbender/cmd/core/errors"
	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	tester "github.com/facebookincubator/fbender/tester/http"
	"github.com/spf13/cobra"
)

const formats = "'GET RelativeURL' or 'POST RelativeURL FormData'"

func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	ssl, err := cmd.Flags().GetBool("ssl")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	r, err := input.NewRequestGenerator(o.Input, inputTransformer(ssl, o.Target), requestCreator)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	t := &tester.Tester{
		Timeout: o.Timeout,
	}

	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(ssl bool, target string) input.Transformer {
	protocol := "http"
	if ssl {
		protocol = "https"
	}

	return func(input string) (interface{}, error) {
		i := strings.Index(input, " ")
		if i < 0 {
			return nil, fmt.Errorf("%w, want: %s, got: %q", errors.ErrInvalidFormat, formats, input)
		}

		method, data := input[:i], input[i+1:]
		switch method {
		case "GET":
			return parseGetRequest(protocol, target, data)
		case "POST":
			return parsePostRequest(protocol, target, data)
		}

		return nil, fmt.Errorf("%w, want: (GET|POST), got: %q", errors.ErrInvalidFormat, method)
	}
}

type request interface {
	Create() (*http.Request, error)
}

type getRequest struct {
	url string
}

func (r *getRequest) Create() (*http.Request, error) {
	//nolint:noctx
	return http.NewRequest("GET", r.url, nil)
}

func parseGetRequest(protocol, target, data string) (interface{}, error) {
	rawurl, err := joinURL(protocol, target, data)
	if err != nil {
		return nil, err
	}

	return &getRequest{url: rawurl}, nil
}

type postRequest struct {
	url  string
	body string
}

func (r *postRequest) Create() (*http.Request, error) {
	//nolint:noctx
	req, err := http.NewRequest("POST", r.url, strings.NewReader(r.body))
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(r.body)))

	return req, nil
}

func parsePostRequest(protocol, target, data string) (interface{}, error) {
	i := strings.Index(data, " ")
	if i < 0 {
		return nil, fmt.Errorf("%w, want: %s, got: \"POST %s\"", errors.ErrInvalidFormat, formats, data)
	}

	form, err := url.ParseQuery(data[i+1:])
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	rawurl, err := joinURL(protocol, target, data[:i])
	if err != nil {
		return nil, err
	}

	return &postRequest{url: rawurl, body: form.Encode()}, nil
}

// NewRequest uses reader interface for message body, which is being used up.
// Therefore we cannot reuse a once created request and need to invoke
// NewRequest everytime before sending it.
func requestCreator(r interface{}) (interface{}, error) {
	if r, ok := r.(request); ok {
		return r.Create()
	}

	return nil, fmt.Errorf("%w, want: request, got: %T", errors.ErrInvalidType, r)
}

func joinURL(protocol, target, path string) (string, error) {
	path = strings.TrimPrefix(path, "/")
	rawurl := fmt.Sprintf("%s://%s/%s", protocol, target, path)
	_, err := url.Parse(rawurl)

	//nolint:wrapcheck
	return rawurl, err
}
