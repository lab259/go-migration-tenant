package mtnt

import (
	"github.com/fatih/color"
)

var (
	styleBold = color.New(color.Bold).SprintFunc()

	styleNormal = color.New(color.FgHiBlack).SprintFunc()

	styleAccountID = styleBold

	styleSuccess = color.New(color.Bold, color.FgGreen).SprintFunc()
	styleWarning = color.New(color.FgYellow).SprintFunc()
	styleError   = color.New(color.FgRed).SprintFunc()
)
