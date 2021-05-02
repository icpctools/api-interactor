package interactor

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUser     = envFallback("TEST_USER", "admin")
	testPass     = envFallback("TEST_PASS", "admin")
	testTeamUser = envFallback("TEST_TEAM_USER", "team")
	testTeamPass = envFallback("TEST_TEAM_PASS", "team")
	testBase     = envFallback("TEST_BASE", "https://www.domjudge.org/demoweb/api")
	testContest  = envFallback("TEST_CONTEST", "nwerc18")

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

func interactor(t assert.TestingT) ContestApi {
	api, err := ContestInteractor(testBase, testUser, testPass, testContest, true)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	return api
}

func teamInteractor(t assert.TestingT) ContestApi {
	api, err := ContestInteractor(testBase, testTeamUser, testTeamPass, testContest, true)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	return api
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

		// Upgrading to a ContestInteractor should fail! The api should be nil such that it is not usable
		api, err := interactor.ToContest(testContest)
		assert.NotNil(t, err)
		assert.Nil(t, api)
	})

	t.Run("valid-base", func(t *testing.T) {
		// Since nothing is verified we expect the response to always be non-nil and nil
		interactor := interactor(t)

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
	api := interactor(t)

	var contestId string
	t.Run("all-contests", func(t *testing.T) {
		contests, err := api.Contests()
		assert.Nil(t, err)
		assert.NotNil(t, contests)

		for _, contest := range contests {
			if contest.Id != "" {
				contestId = contest.Id
				return
			}
		}
	})

	t.Run("single-contest", func(t *testing.T) {
		if contestId == "" {
			t.Skip("no contest could be found, retrieving single contest cannot be tested")
		}

		contest, err := api.ContestById(contestId)
		assert.Nil(t, err)
		assert.EqualValues(t, contestId, contest.Id)
	})
}

func TestProblemRetrieval(t *testing.T) {
	api := interactor(t)

	var problemId string
	t.Run("all-problems", func(t *testing.T) {
		problems, err := api.Problems()
		assert.Nil(t, err)
		assert.NotNil(t, problems)

		for _, problem := range problems {
			if problem.Id != "" {
				problemId = problem.Id
				return
			}
		}
	})

	t.Run("single-problem", func(t *testing.T) {
		if problemId == "" {
			t.Skip("no problem could be found, retrieving single problem cannot be tested")
		}

		problem, err := api.ProblemById(problemId)
		assert.Nil(t, err)
		assert.EqualValues(t, problemId, problem.Id)
	})
}

func TestJudgementTypeRetrieval(t *testing.T) {
	api := interactor(t)

	var jtId string
	t.Run("all-judgement-types", func(t *testing.T) {
		judgementTypes, err := api.JudgementTypes()
		assert.Nil(t, err)
		assert.NotNil(t, judgementTypes)

		for _, jt := range judgementTypes {
			if jt.Id != "" {
				jtId = jt.Id
				return
			}
		}
	})

	t.Run("single-judgement-type", func(t *testing.T) {
		if jtId == "" {
			t.Skip("no judgement type could be found, retrieving single judgement type cannot be tested")
		}

		jt, err := api.JudgementTypeById(jtId)
		assert.Nil(t, err)
		assert.EqualValues(t, jtId, jt.Id)
	})
}

func TestSubmissionRetrieval(t *testing.T) {
	api := interactor(t)

	var submissionId string
	t.Run("all-submissions", func(t *testing.T) {
		submissions, err := api.Submissions()
		assert.Nil(t, err)
		assert.NotNil(t, submissions)

		for _, submission := range submissions {
			if submission.Id != "" {
				submissionId = submission.Id
				t.Logf("Found submission %+v\n", submission)
				return
			}
		}
	})

	t.Run("single-submissions", func(t *testing.T) {
		if submissionId == "" {
			t.Skip("no submission could be found, retrieving single submission cannot be tested")
		}

		t.Log(submissionId, testUser, testPass)

		submission, err := api.SubmissionById(submissionId)
		assert.Nil(t, err)
		assert.EqualValues(t, submissionId, submission.Id)
	})
}

func TestJudgementRetrieval(t *testing.T) {
	api := interactor(t)

	var jtId string
	t.Run("all-judgements", func(t *testing.T) {
		judgements, err := api.Judgements()
		assert.Nil(t, err)
		assert.NotNil(t, judgements)

		for _, jt := range judgements {
			if jt.Id != "" {
				jtId = jt.Id
				return
			}
		}
	})

	t.Run("single-judgement", func(t *testing.T) {
		if jtId == "" {
			t.Skip("no judgement could be found, retrieving single judgement cannot be tested")
		}

		jt, err := api.JudgementById(jtId)
		assert.Nil(t, err)
		assert.EqualValues(t, jtId, jt.Id)
	})
}

func TestLanguageRetrieval(t *testing.T) {
	api := interactor(t)

	var languageId string
	t.Run("all-languages", func(t *testing.T) {
		languages, err := api.Languages()
		assert.Nil(t, err)
		assert.NotNil(t, languages)

		t.Log(languages)

		for _, language := range languages {
			if language.Id != "" {
				languageId = language.Id
				t.Logf("Found language %+v\n", language)
				return
			}
		}
	})

	t.Run("single-language", func(t *testing.T) {
		if languageId == "" {
			t.Skip("no language could be found, retrieving single language cannot be tested")
		}

		t.Log(languageId, testUser, testPass)

		language, err := api.LanguageById(languageId)
		assert.Nil(t, err)
		assert.EqualValues(t, languageId, language.Id)
	})
}

func TestClarificationRetrieval(t *testing.T) {
	api := interactor(t)

	var clarificationId string
	t.Run("all-clarifications", func(t *testing.T) {
		clarifications, err := api.Clarifications()
		assert.Nil(t, err)
		assert.NotNil(t, clarifications)

		for _, clarification := range clarifications {
			if clarification.Id != "" {
				clarificationId = clarification.Id
				return
			}
		}
	})

	t.Run("single-clarification", func(t *testing.T) {
		if clarificationId == "" {
			t.Skip("no clarification could be found, retrieving single clarification cannot be tested")
		}

		clarification, err := api.ClarificationById(clarificationId)
		assert.Nil(t, err)
		assert.EqualValues(t, clarificationId, clarification.Id)
	})
}

func TestSendClarification(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		api := interactor(t)

		id, err := api.PostClarification("A", "testing clarification")
		assert.Empty(t, id)
		assert.NotNil(t, err)

		t.Logf("Sent clarification, got error: '%v'", err)
	})

	t.Run("authorized", func(t *testing.T) {
		api := interactor(t)

		id, err := api.PostClarification("accesspoints", "testing clarification")
		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent clarification, got id: '%v'", id)
	})

	t.Run("authorized-struct", func(t *testing.T) {
		api := interactor(t)
		clar := Clarification{
			ProblemId: "accesspoints",
			Text:      "This is only a test",
		}

		bts, err := json.Marshal(clar)
		assert.Nil(t, err)
		assert.NotNil(t, bts)

		id, err := api.Submit(&clar)

		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent clarification, got id: '%v'", id)
	})
}

func TestPostSubmission(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		api := teamInteractor(t)

		id, err := api.PostSubmission("A", "cpp", "", LocalFileReference{})
		assert.Empty(t, id)
		assert.NotNil(t, err)

		t.Logf("Sent submission, got error: '%v'", err)
	})

	t.Run("authorized", func(t *testing.T) {
		api := teamInteractor(t)

		var sampleSubmission LocalFileReference
		_ = sampleSubmission.FromString("sample.cpp", "int main() { return 0; }")
		goModFile, _ := os.Open("go.mod")
		_ = sampleSubmission.FromFile(goModFile)
		id, err := api.PostSubmission("accesspoints", "cpp", "", sampleSubmission)
		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent submission, got id: '%v'", id)
	})

	t.Run("authorized-struct", func(t *testing.T) {
		api := teamInteractor(t)
		var sampleSubmission LocalFileReference
		_ = sampleSubmission.FromString("sample.cpp", "int main() { return 0; }")
		submission := Submission{
			ProblemId:  "accesspoints",
			LanguageId: "cpp",
			Time:       new(ApiTime),
			Files: []FileReference{
				{
					Mime: "application/zip",
					Data: sampleSubmission,
				},
			},
		}

		bts, err := json.Marshal(submission)
		assert.Nil(t, err)
		assert.NotNil(t, bts)

		id, err := api.Submit(&submission)
		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent submission, got id: '%v'", id)
	})
}
