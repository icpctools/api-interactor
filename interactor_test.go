package interactor

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	testUser    = envFallback("TEST_USER", "")
	testPass    = envFallback("TEST_PASS", "")
	testBase    = envFallback("TEST_BASE", "https://www.domjudge.org/demoweb/api")
	testContest = envFallback("TEST_CONTEST", "nwerc18")

	testContestWrong = envFallback("TEST_CONTEST_WRONG", "NON_EXISTENT_CONTEST_ID_I_HOPE")

	// Ensure the interfaces are adhered to
	_ ContestApi  = new(inter)
	_ ContestsApi = new(inter)
)

func envFallback(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	return fb
}

func TestContestInteractor(t *testing.T) {
	t.Run("invalid-contest-id", func(t *testing.T) {
		interactor, err := ContestInteractor(testBase, testUser, testPass, testContestWrong, true)
		assert.NotNil(t, err)
		assert.Nil(t, interactor)
	})

	t.Run("valid-contest-id", func(t *testing.T) {
		interactor, err := ContestInteractor(testBase, testUser, testPass, testContest, true)
		assert.Nil(t, err)
		assert.NotNil(t, interactor)
	})
}

func TestContestsInteractor(t *testing.T) {
	t.Run("invalid-base", func(t *testing.T) {
		interactor, err := ContestsInteractor("this-does-not-exists", "", "", true)
		assert.Nil(t, err)
		assert.NotNil(t, interactor)

		// Upgrading to a ContesInteractor should fail! The api should be nil such that it is not usable
		api, err := interactor.ToContest(testContest)
		assert.NotNil(t, err)
		assert.Nil(t, api)
	})

	t.Run("valid-base", func(t *testing.T) {
		// Since nothing is verified we expect the response to always be non-nil and nil
		interactor, err := ContestsInteractor(testBase, testUser, testPass, true)
		assert.Nil(t, err)
		assert.NotNil(t, interactor)

		// ToContest should not result in nil as long as it exists
		api, err := interactor.ToContest(testContest)
		assert.Nil(t, err)
		assert.NotNil(t, api)

		// ToContest should fail when a non-existent contest is given
		api, err = interactor.ToContest(testContestWrong)
		assert.NotNil(t, err)
		assert.Nil(t, api)
	})
}

func TestContestRetrieval(t *testing.T) {
	api, err := ContestsInteractor(testBase, testUser, testPass, true)
	assert.Nil(t, err)

	var contestId string
	t.Run("all-contests", func(t *testing.T) {
		contests, err := api.Contests()
		assert.Nil(t, err)
		assert.NotNil(t, contests)

		if len(contests) > 0 {
			contestId = contests[0].Id
		}
	})

	t.Run("single-contest", func(t *testing.T) {
		if contestId == "" {
			t.Skip("no single contest could be found, retrieving single contest cannot be tested")
		}

		contest, err := api.ContestById(contestId)
		assert.Nil(t, err)
		assert.EqualValues(t, contestId, contest.Id)
	})
}

func TestProblemRetrieval(t *testing.T) {
	api, err := ContestInteractor(testBase, testUser, testPass, testContest, true)
	assert.Nil(t, err)

	var problemId string
	t.Run("all-problems", func(t *testing.T) {
		problems, err := api.Problems()
		assert.Nil(t, err)
		assert.NotNil(t, problems)

		if len(problems) > 0 {
			problemId = problems[0].Id
		}
	})

	t.Run("single-problem", func(t *testing.T) {
		if problemId == "" {
			t.Skip("no single problem could be found, retrieving single problem cannot be tested")
		}

		problem, err := api.ProblemById(problemId)
		assert.Nil(t, err)
		assert.EqualValues(t, problemId, problem.Id)
	})
}

func TestSubmissionRetrieval(t *testing.T) {
	api, err := ContestInteractor(testBase, testUser, testPass, testContest, true)
	assert.Nil(t, err)

	var submissionId string
	t.Run("all-problems", func(t *testing.T) {
		submissions, err := api.Submissions()
		assert.Nil(t, err)
		assert.NotNil(t, submissions)

		if len(submissions) > 0 {
			submissionId = submissions[0].Id
		}
	})

	t.Run("single-problem", func(t *testing.T) {
		if submissionId == "" {
			t.Skip("no single submission could be found, retrieving single submission cannot be tested")
		}

		submission, err := api.SubmissionById(submissionId)
		assert.Nil(t, err)
		assert.EqualValues(t, submissionId, submission.Id)
	})
}
