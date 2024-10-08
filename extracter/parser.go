/**
 * Copyright 2024 KusionStack Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extracter

import (
	"errors"

	"k8s.io/client-go/util/jsonpath"
)

// Parse is unlike the jsonpath.Parse, which supports multi-paths input.
// The input like `{.kind} {.apiVersion}` or
// `{range .spec.containers[*]}{.name}{end}` will result in an error.
func Parse(name, text string) (*parser, error) {
	p, err := jsonpath.Parse(name, text)
	if err != nil {
		return nil, err
	}

	if len(p.Root.Nodes) > 1 {
		return nil, errors.New("not support multi-paths input")
	}

	return &parser{p}, nil
}

type parser struct {
	*jsonpath.Parser
}
