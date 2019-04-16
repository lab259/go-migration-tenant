package mtnt

import (
	"fmt"
	"io"
)

type Reporter interface {
	BeforeAccount(account Account)
	AfterAccount(account Account)
}

type defaultReporter struct {
	Writer io.Writer
}

func NewDefaultReporter(writer io.Writer) *defaultReporter {
	return &defaultReporter{
		Writer: writer,
	}
}

func (reporter *defaultReporter) BeforeAccount(account Account) {
	fmt.Fprintf(reporter.Writer, "Migrating %s", styleAccountID(account.Identification()))
}

func (reporter *defaultReporter) AfterAccount(account Account) {
	fmt.Fprintln(reporter.Writer)
}
