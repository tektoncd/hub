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

package search

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/tektoncd/hub/api/pkg/cli/app"
	"github.com/tektoncd/hub/api/pkg/cli/flag"
	"github.com/tektoncd/hub/api/pkg/cli/formatter"
	"github.com/tektoncd/hub/api/pkg/cli/hub"
	"github.com/tektoncd/hub/api/pkg/cli/printer"
	"github.com/tektoncd/hub/api/pkg/parser"
)

const resTemplate = `{{- $rl := len .Resources }}{{ if eq $rl 0 -}}
No Resources found
{{ else -}}
NAME	KIND	DESCRIPTION	TAGS
{{ range $_, $r := .Resources -}}
{{ formatName $r.Name $r.LatestVersion.Version }}	{{ $r.Kind }}	{{ formatDesc $r.LatestVersion.Description }}	{{ formatTags $r.Tags }}	
{{ end }}
{{- end -}}
`

var (
	funcMap = template.FuncMap{
		"formatName": formatter.FormatName,
		"formatDesc": formatter.FormatDesc,
		"formatTags": formatter.FormatTags,
	}
	tmpl = template.Must(template.New("List Resources").Funcs(funcMap).Parse(resTemplate))
)

type options struct {
	cli    app.CLI
	Limit  uint
	Match  string
	Output string
	Tags   []string
	Kinds  []string
	Args   []string
}

func Command(cli app.CLI) *cobra.Command {

	opts := &options{cli: cli}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search resource by combination of its name, kind, and tags",
		Long:  ``,
		Annotations: map[string]string{
			"commandType": "main",
		},
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return opts.run()
		},
	}

	cmd.Flags().UintVarP(&opts.Limit, "limit", "l", 0, "Max number of resources to fetch")
	cmd.Flags().StringVar(&opts.Match, "match", "contains", "Accept type of search. 'exact' or 'contains'.")
	cmd.Flags().StringArrayVar(&opts.Kinds, "kinds", nil, "Accepts a comma separated list of kinds")
	cmd.Flags().StringArrayVar(&opts.Tags, "tags", nil, "Accepts a comma separated list of tags")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "table", "Accepts output format: [table, json]")

	return cmd
}

func (opts *options) run() error {

	if err := opts.validate(); err != nil {
		return err
	}

	hubClient := opts.cli.Hub()

	result := hubClient.Search(hub.SearchOption{
		Name:  opts.name(),
		Kinds: opts.Kinds,
		Tags:  opts.Tags,
		Match: opts.Match,
		Limit: opts.Limit,
	})

	out := opts.cli.Stream().Out

	if opts.Output == "json" {
		return printer.New(out).JSON(result.Raw())
	}

	typed, err := result.Typed()
	if err != nil {
		return err
	}

	var templateData = struct {
		Resources hub.SearchResponse
	}{
		Resources: typed,
	}

	return printer.New(out).Tabbed(tmpl, templateData)
}

func (opts *options) validate() error {

	if flag.AllEmpty(opts.Args, opts.Kinds, opts.Tags) {
		return fmt.Errorf("please specify a name, tag or a kind to search")
	}

	if err := flag.InList("match", opts.Match, []string{"contains", "exact"}); err != nil {
		return err
	}

	if err := flag.InList("output", opts.Output, []string{"table", "json"}); err != nil {
		return err
	}

	opts.Kinds = flag.TrimArray(opts.Kinds)
	opts.Tags = flag.TrimArray(opts.Tags)

	for _, k := range opts.Kinds {
		if !parser.IsSupportedKind(k) {
			return fmt.Errorf("invalid value %q set for option kinds. supported kinds: [%s]",
				k, strings.ToLower(strings.Join(parser.SupportedKinds(), ", ")))
		}
	}
	return nil
}

func (opts *options) name() string {
	if len(opts.Args) == 0 {
		return ""
	}
	return strings.TrimSpace(opts.Args[0])
}
