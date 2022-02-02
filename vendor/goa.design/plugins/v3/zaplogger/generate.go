package zaplogger

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type fileToModify struct {
	file        *codegen.File
	path        string
	serviceName string
	isMain      bool
}

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginFirst("zaplogger", "gen", nil, Generate)
	codegen.RegisterPluginLast("zaplogger-updater", "example", nil, UpdateExample)
}

// Generate generates zap logger specific files.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, GenerateFiles(genpkg, r)...)
		}
	}
	return files, nil
}

// UpdateExample modifies the example generated files by replacing
// the log import reference when needed
// It also modify the initially generated main and service files
func UpdateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {

	filesToModify := []*fileToModify{}

	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {

			// Add the generated main files
			for _, svr := range r.API.Servers {
				pkg := codegen.SnakeCase(codegen.Goify(svr.Name, true))
				filesToModify = append(filesToModify,
					&fileToModify{path: filepath.Join("cmd", pkg, "main.go"), serviceName: svr.Name, isMain: true})
				filesToModify = append(filesToModify,
					&fileToModify{path: filepath.Join("cmd", pkg, "http.go"), serviceName: svr.Name, isMain: true})
				filesToModify = append(filesToModify,
					&fileToModify{path: filepath.Join("cmd", pkg, "grpc.go"), serviceName: svr.Name, isMain: true})
			}

			// Add the generated service files
			for _, svc := range r.API.HTTP.Services {
				servicePath := codegen.SnakeCase(svc.Name()) + ".go"
				filesToModify = append(filesToModify, &fileToModify{path: servicePath, serviceName: svc.Name(), isMain: false})
			}

			// Update the added files
			for _, fileToModify := range filesToModify {
				for _, file := range files {
					if file.Path == fileToModify.path {
						fileToModify.file = file
						updateExampleFile(genpkg, r, fileToModify)
						break
					}
				}
			}
		}
	}

	return files, nil
}

// GenerateFiles create log specific files
func GenerateFiles(genpkg string, root *expr.RootExpr) []*codegen.File {
	fw := make([]*codegen.File, 1)
	fw[0] = GenerateLoggerFile(genpkg)
	return fw
}

// GenerateLoggerFile returns the generated zap logger file.
func GenerateLoggerFile(genpkg string) *codegen.File {
	path := filepath.Join(codegen.Gendir, "log", "logger.go")
	title := fmt.Sprint("Zap logger implementation")
	sections := []*codegen.SectionTemplate{
		codegen.Header(title, "log", []*codegen.ImportSpec{
			{Path: "go.uber.org/zap"},
			{Path: "goa.design/goa/v3/http/middleware"},
		}),
	}

	sections = append(sections, &codegen.SectionTemplate{
		Name:   "zaplogger",
		Source: loggerT,
	})

	return &codegen.File{Path: path, SectionTemplates: sections}
}

func updateExampleFile(genpkg string, root *expr.RootExpr, f *fileToModify) {

	header := f.file.SectionTemplates[0]
	logPath := path.Join(genpkg, "log")

	data := header.Data.(map[string]interface{})
	specs := data["Imports"].([]*codegen.ImportSpec)

	for _, spec := range specs {
		if spec.Path == "log" {
			spec.Name = "log"
			spec.Path = logPath
		}
	}

	if f.isMain {

		codegen.AddImport(header, &codegen.ImportSpec{Path: "go.uber.org/zap"})

		for _, s := range f.file.SectionTemplates {
			s.Source = strings.Replace(s.Source, `logger = log.New(os.Stderr, "[{{ .APIPkg }}] ", log.Ltime)`, `logger = log.New("{{ .APIPkg }}", false)`, 1)
			s.Source = strings.Replace(s.Source, "adapter = middleware.NewLogger(logger)", "adapter = logger", 1)
			s.Source = strings.Replace(s.Source, "handler = middleware.RequestID()(handler)",
				`handler = middleware.PopulateRequestContext()(handler)
				handler = middleware.RequestID(middleware.UseXRequestIDHeaderOption(true))(handler)`, 1)
			s.Source = strings.Replace(s.Source, `logger.Printf("[%s] ERROR: %s", id, err.Error())`,
				`logger.With(zap.String("id",id)).Error(err.Error())`, 1)
			s.Source = strings.Replace(s.Source, "logger.Print(", "logger.Info(", -1)
			s.Source = strings.Replace(s.Source, "logger.Printf(", "logger.Infof(", -1)
			s.Source = strings.Replace(s.Source, "logger.Println(", "logger.Info(", -1)
		}
	} else {
		for _, s := range f.file.SectionTemplates {
			s.Source = strings.Replace(s.Source, "logger.Print(", "logger.Info(", -1)
			s.Source = strings.Replace(s.Source, "logger.Printf(", "logger.Infof(", -1)
			s.Source = strings.Replace(s.Source, "logger.Println(", "logger.Info(", -1)
		}
	}
}

const loggerT = `
// Logger is an adapted zap logger
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new zap logger
func New(serviceName string, production bool) *Logger {

	if production {
		l, _ := zap.NewProduction()
		return &Logger{l.Sugar().With(zap.String("service", serviceName))}
	} else {
		l, _ := zap.NewDevelopment()
		return &Logger{l.Sugar().With(zap.String("service", serviceName))}
	}
}

// Log is called by the log middleware to log HTTP requests key values
func (logger *Logger) Log(keyvals ...interface{}) error {
	logger.Infow("HTTP Request", keyvals...)
	return nil
}
`
