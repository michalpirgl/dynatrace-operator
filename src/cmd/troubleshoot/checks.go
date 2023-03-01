package troubleshoot

import (
	"strings"

	tserrors "github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot/errors"
	"github.com/Dynatrace/dynatrace-operator/src/functional"
	"github.com/go-logr/logr"
)

type Result int

const (
	PASSED Result = iota + 1
	FAILED
	SKIPPED
)

type troubleshootFunc func(troubleshootCtx *TroubleshootContext) error

type CheckListEntry struct {
	Check         Check
	Prerequisites []*CheckListEntry

	Do   troubleshootFunc
	Name string
}

type Check interface {
	Name() string
	Do(log logr.Logger) error
}

func runChecks(results ChecksResults, troubleshootCtx *TroubleshootContext, checks []*CheckListEntry) error {
	errs := tserrors.NewAggregatedError()

	for _, checkListEntry := range checks {
		if shouldSkip(results, checkListEntry) {
			results.set(checkListEntry, SKIPPED)
			continue
		}

		err := checkListEntry.Check.Do().Do(troubleshootCtx)
		if err != nil {
			logErrorf(err.Error())
			errs.Add(err)
			results.set(checkListEntry, FAILED)
		} else {
			results.set(checkListEntry, PASSED)
		}
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func shouldSkip(results ChecksResults, check *CheckListEntry) bool {
	failedOrSkippedPrerequisites := results.failedOrSkippedPrerequisites(check)

	if len(failedOrSkippedPrerequisites) == 0 {
		return false
	}

	getCheckName := func(check *CheckListEntry) string {
		return check.Name
	}
	prerequisitesNames := strings.Join(functional.Map(failedOrSkippedPrerequisites, getCheckName), ",")
	logWarningf("Skipped '%s' check because prerequisites aren't met: [%s]", check.Name, prerequisitesNames)

	return true
}
