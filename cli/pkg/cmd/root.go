// Copyright © 2020 The Tekton Authors.
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

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tektoncd/hub/cli/pkg/cmd/search"
)

// Root represents the base command when called without any subcommands
func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "tkn-hub",
		Short:        "CLI for tekton hub",
		Long:         ``,
		SilenceUsage: true,
	}

	cmd.AddCommand(search.Command())

	return cmd
}
