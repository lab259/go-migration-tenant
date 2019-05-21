package rlog

import (
	"github.com/fatih/color"
	rlog2 "github.com/lab259/rlog"

	mtnt "github.com/lab259/go-migration-tenant"
)

var styleAccountID = color.New(color.Bold).SprintFunc()

type rLogReporter struct {
	logger rlog2.Logger
}

func NewRLogReporter(logger rlog2.Logger) *rLogReporter {
	return &rLogReporter{
		logger: logger,
	}
}

func (reporter *rLogReporter) BeforeAccount(account mtnt.Account) {
	reporter.logger.Infof("Migrating %s", styleAccountID(account.Identification()))
}

func (reporter *rLogReporter) AfterAccount(account mtnt.Account) {}
