package parser

import (
	"fmt"
	"strings"
)

type IssueType int

const (
	Critical IssueType = iota
	Warning
	Info
)

func (t IssueType) String() string {
	return [...]string{"critical", "warning", "info"}[t]
}

type Issue struct {
	Type    IssueType
	Message string
}

type Result struct {
	Errors []error
	Issues []Issue
}

func (r *Result) add(t IssueType, format string, args ...interface{}) {
	if r.Issues == nil {
		r.Issues = []Issue{}
	}
	r.Issues = append(r.Issues, Issue{Type: t, Message: fmt.Sprintf(format, args...)})

}

func (r *Result) Critical(format string, args ...interface{}) {
	r.add(Critical, format, args...)
}

func (r *Result) Warn(format string, args ...interface{}) {
	r.add(Warning, format, args...)
}

func (r *Result) Info(format string, args ...interface{}) {
	r.add(Info, format, args...)
}

func (r *Result) Error() string {

	if r.Errors == nil {
		return ""
	}

	buf := strings.Builder{}

	for _, err := range r.Errors {
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func (r *Result) AddError(err error) {
	if r.Errors == nil {
		r.Errors = []error{}
	}
	r.Errors = append(r.Errors, err)
}

func (r *Result) Combine(other Result) {
	r.Issues = append(r.Issues, other.Issues...)
	r.Errors = append(r.Errors, other.Errors...)
}
