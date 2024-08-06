package zaplogger

import (
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginLast("zaplogger", "example", nil, UpdateExample)
}

// UpdateExample modifies the example generated files by replacing
// the log import reference when needed
// It also modify the initially generated main and service files
func UpdateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		updateExample(f)
	}
	return files, nil
}

func updateExample(file *codegen.File) {
	for _, section := range file.SectionTemplates {
		switch section.Name {
		case "server-main-services":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "go.uber.org/zap"})
			oldinit := "{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}()"
			section.Source = strings.Replace(section.Source, oldinit, initT, 1)
		case "basic-service-struct":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "go.uber.org/zap"})
			section.Source = basicServiceStructT
		case "basic-service-init":
			section.Source = basicServiceInitT
		case "basic-endpoint":
			section.Source = strings.Replace(
				section.Source,
				`log.Printf(ctx, "{{ .ServiceVarName }}.{{ .Name }}")`,
				`s.logger.Info("{{ .ServiceVarName}}.{{ .Name }}")`,
				1,
			)
		}
	}
}

const (
	initT = `
	var zlog *zap.SugaredLogger
	if *dbgF {
		l, _ := zap.NewDevelopment()
		zlog = l.Sugar().With(zap.String("service", {{ printf "%q" .Name }}))
	} else {
		l, _ := zap.NewProduction()
		zlog = l.Sugar().With(zap.String("service", {{ printf "%q" .Name }}))
	}
	{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}(zlog)`

	basicServiceInitT = `
{{ printf "New%s returns the %s service implementation." .StructName .Name | comment }}
func New{{ .StructName }}(logger *zap.SugaredLogger) {{ .PkgName }}.Service {
	return &{{ .VarName }}srvc{
		logger: logger,
	}
}
`

	basicServiceStructT = `
	{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Name | comment }}
	type {{ .VarName }}srvc struct {
		logger *zap.SugaredLogger
	}
`
)
