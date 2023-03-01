package troubleshoot

import (
	"testing"

	tserrors "github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	checkError = errors.New("check function failed")

	tsContext = &TroubleshootContext{}

	passingBasicCheck = &CheckListEntry{
		Name: "passingBasicCheck",
		Do: func(*TroubleshootContext) error {
			return nil
		},
	}

	failingBasicCheck = &CheckListEntry{
		Name: "failingBasicCheck",
		Do: func(*TroubleshootContext) error {
			return checkError
		},
	}

	passingCheckDependendOnPassingCheck = &CheckListEntry{
		Name: "passingCheckDependendOnPassingCheck",
		Do: func(*TroubleshootContext) error {
			return nil
		},
		Prerequisites: []*CheckListEntry{passingBasicCheck},
	}

	failingCheckDependendOnPassingCheck = &CheckListEntry{
		Name: "failingCheckDependendOnPassingCheck",
		Do: func(*TroubleshootContext) error {
			return checkError
		},
		Prerequisites: []*CheckListEntry{passingBasicCheck},
	}

	failingCheckDependendOnFailingCheck = &CheckListEntry{
		Name: "failingCheckDependendOnFailingCheck",
		Do: func(*TroubleshootContext) error {
			return checkError
		},
		Prerequisites: []*CheckListEntry{failingBasicCheck},
	}
)

func Test_runChecks(t *testing.T) {
	t.Run("no checks", func(t *testing.T) {
		checks := []*CheckListEntry{}
		results := NewChecksResults()
		err := runChecks(results, tsContext, checks)
		require.NoError(t, err)
	})
	t.Run("a few passing checks", func(t *testing.T) {
		checks := []*CheckListEntry{
			passingBasicCheck,
			passingCheckDependendOnPassingCheck,
		}
		results := NewChecksResults()
		err := runChecks(results, tsContext, checks)
		require.NoError(t, err)
	})
	t.Run("passing and failing checks", func(t *testing.T) {
		checks := []*CheckListEntry{
			passingBasicCheck,
			passingCheckDependendOnPassingCheck,
			failingCheckDependendOnPassingCheck,
			failingBasicCheck,
			failingCheckDependendOnFailingCheck, // should be skipped and error should not be reported
		}
		results := NewChecksResults()
		resetLogger()
		err := runChecks(results, tsContext, checks)
		require.Error(t, err)

		aggregatedError := tserrors.AggregatedError{}
		isAggredatedError := errors.As(err, &aggregatedError)

		require.True(t, isAggredatedError)
		require.Len(t, aggregatedError.Errs, 2)
		require.ErrorIs(t, aggregatedError.Errs[0], checkError)
		require.ErrorIs(t, aggregatedError.Errs[1], checkError)
	})
	t.Run("check should not be run if prerequisite check failed", func(t *testing.T) {
		checks := []*CheckListEntry{
			failingBasicCheck,
			failingCheckDependendOnFailingCheck, // should be skipped and error should not be reported
		}
		results := NewChecksResults()
		err := runChecks(results, tsContext, checks)
		require.Error(t, err)
	})
}
