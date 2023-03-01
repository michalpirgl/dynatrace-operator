package troubleshoot

import "github.com/Dynatrace/dynatrace-operator/src/functional"

type ChecksResults struct {
	checkResultMap map[*CheckListEntry]Result
}

func NewChecksResults() ChecksResults {
	return ChecksResults{checkResultMap: map[*CheckListEntry]Result{}}
}

func (checkResults ChecksResults) set(check *CheckListEntry, result Result) {
	checkResults.checkResultMap[check] = result
}

func (checkResults ChecksResults) failedOrSkippedPrerequisites(check *CheckListEntry) []*CheckListEntry {
	isFailedOrSkipped := func(check *CheckListEntry) bool {
		return checkResults.checkResultMap[check] == FAILED || checkResults.checkResultMap[check] == SKIPPED
	}
	return functional.Filter(check.Prerequisites, isFailedOrSkipped)
}

func (checkResults ChecksResults) hasErrors() bool {
	for _, result := range checkResults.checkResultMap {
		if result == FAILED {
			return true
		}
	}
	return false
}
