// Copyright © 2021 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package select_options

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type Options struct {
	AskOpts survey.AskOpt
	Name    string
	Catalog string
	Version string
}

func (opts *Options) Ask(resourceInfo string, options []string) error {
	var ans string
	var qs = []*survey.Question{
		{
			Name: resourceInfo,
			Prompt: &survey.Select{
				Message: fmt.Sprintf("Select %s:", resourceInfo),
				Options: options,
			},
		},
	}

	if err := survey.Ask(qs, &ans, opts.AskOpts); err != nil {
		return err
	}

	switch resourceInfo {
	case "catalog":
		opts.Catalog = ans
	case "task":
		opts.Name = ans
	case "version":
		opts.Version = ans
	}

	return nil
}
