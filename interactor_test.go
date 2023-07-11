package interactor

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testLogger is an interface to pass the testing Type to a method that only needs a logger
type testLogger interface {
	Logf(string, ...interface{})
}

var (
	testUser     = envFallback("TEST_USER", "admin")
	testPass     = envFallback("TEST_PASS", "admin")
	testTeamUser = envFallback("TEST_TEAM_USER", "team")
	testTeamPass = envFallback("TEST_TEAM_PASS", "team")
	testBase     = envFallback("TEST_BASE", "https://www.domjudge.org/demoweb/api")
	testContest  = envFallback("TEST_CONTEST", "nwerc18")
	testProblem  = envFallback("TEST_PROBLEM", "accesspoints")

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
	api, err := ContestInteractor(testBase, testUser, testPass, testContest, false)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	return api
}

func teamInteractor(t assert.TestingT) ContestApi {
	api, err := ContestInteractor(testBase, testTeamUser, testTeamPass, testContest, false)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	return api
}

func TestContestInteractor(t *testing.T) {
	t.Run("invalid-contest-id", func(t *testing.T) {
		interactor, err := ContestInteractor(testBase, testUser, testPass, testContestWrong, false)
		assert.NotNil(t, err)
		assert.Nil(t, interactor)
	})

	t.Run("valid-contest-id", func(t *testing.T) {
		interactor, err := ContestInteractor(testBase, testUser, testPass, testContest, false)
		assert.Nil(t, err)
		assert.NotNil(t, interactor)
	})
}

func TestContestsInteractor(t *testing.T) {
	t.Run("invalid-base", func(t *testing.T) {
		interactor, err := ContestsInteractor("this-does-not-exists", "", "", false)
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
			t.Skip("no contests could be found, retrieving single contest cannot be tested")
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
		problems, err := List(api, Problem{})
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
			t.Skip("no problems could be found, retrieving single problem cannot be tested")
		}

		problem, err := GetById[Problem](api, problemId)
		assert.Nil(t, err)
		assert.EqualValues(t, problemId, problem.Id)
	})
}

func TestJudgementTypeRetrieval(t *testing.T) {
	api := interactor(t)

	var jtId string
	t.Run("all-judgement-types", func(t *testing.T) {
		judgementTypes, err := List(api, JudgementType{})
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
			t.Skip("no judgement types could be found, retrieving single judgement type cannot be tested")
		}

		jt, err := GetById[JudgementType](api, jtId)
		assert.Nil(t, err)
		assert.EqualValues(t, jtId, jt.Id)
	})
}

func TestSubmissionRetrieval(t *testing.T) {
	api := interactor(t)

	var submissionId string
	t.Run("all-submissions", func(t *testing.T) {
		submissions, err := List(api, Submission{})
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
			t.Skip("no submissions could be found, retrieving single submission cannot be tested")
		}

		t.Log(submissionId, testUser, testPass)

		submission, err := GetById[Submission](api, submissionId)
		assert.Nil(t, err)
		assert.EqualValues(t, submissionId, submission.Id)
	})
}

func TestJudgementRetrieval(t *testing.T) {
	api := interactor(t)

	var jId string
	t.Run("all-judgements", func(t *testing.T) {
		judgements, err := List(api, Judgement{})
		assert.Nil(t, err)
		assert.NotNil(t, judgements)

		for _, j := range judgements {
			if j.Id != "" {
				jId = j.Id
				return
			}
		}
	})

	t.Run("single-judgement", func(t *testing.T) {
		if jId == "" {
			t.Skip("no judgements could be found, retrieving single judgement cannot be tested")
		}

		j, err := GetById[Judgement](api, jId)
		assert.Nil(t, err)
		assert.EqualValues(t, jId, j.Id)
	})
}

func TestGroupRetrieval(t *testing.T) {
	api := interactor(t)

	var groupId string
	t.Run("all-groups", func(t *testing.T) {
		groups, err := List(api, Group{})
		assert.Nil(t, err)
		assert.NotNil(t, groups)

		for _, group := range groups {
			if group.Id != "" {
				groupId = group.Id
				return
			}
		}
	})

	t.Run("single-group", func(t *testing.T) {
		if groupId == "" {
			t.Skip("no groups could be found, retrieving single group cannot be tested")
		}

		group, err := GetById[Group](api, groupId)
		assert.Nil(t, err)
		assert.EqualValues(t, groupId, group.Id)
	})
}

func TestOrganizationRetrieval(t *testing.T) {
	api := interactor(t)

	var organizationId string
	t.Run("all-organizations", func(t *testing.T) {
		organizations, err := List(api, Organization{})
		assert.Nil(t, err)
		assert.NotNil(t, organizations)

		for _, organization := range organizations {
			if organization.Id != "" {
				organizationId = organization.Id
				return
			}
		}
	})

	t.Run("single-organization", func(t *testing.T) {
		if organizationId == "" {
			t.Skip("no organizations could be found, retrieving single organization cannot be tested")
		}

		organization, err := GetById[Organization](api, organizationId)
		assert.Nil(t, err)
		assert.EqualValues(t, organizationId, organization.Id)
	})
}

func TestTeamRetrieval(t *testing.T) {
	api := interactor(t)

	var teamId string
	t.Run("all-teams", func(t *testing.T) {
		teams, err := List(api, Team{})
		assert.Nil(t, err)
		assert.NotNil(t, teams)

		for _, team := range teams {
			if team.Id != "" {
				teamId = team.Id
				return
			}
		}
	})

	t.Run("single-team", func(t *testing.T) {
		if teamId == "" {
			t.Skip("no teams could be found, retrieving single team cannot be tested")
		}

		team, err := GetById[Team](api, teamId)
		assert.Nil(t, err)
		assert.EqualValues(t, teamId, team.Id)
	})
}

func TestScoreboardRetrieval(t *testing.T) {
	api := interactor(t)

	sb, err := api.Scoreboard()
	assert.Nil(t, err)

	t.Log(sb)
}

func TestStateRetrieval(t *testing.T) {
	api := interactor(t)

	sb, err := api.State()
	assert.Nil(t, err)

	t.Log(sb)
}

func TestLanguageRetrieval(t *testing.T) {
	api := interactor(t)

	var languageId string
	t.Run("all-languages", func(t *testing.T) {
		languages, err := List(api, Language{})
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
			t.Skip("no languages could be found, retrieving single language cannot be tested")
		}

		t.Log(languageId, testUser, testPass)

		language, err := GetById[Language](api, languageId)
		assert.Nil(t, err)
		assert.EqualValues(t, languageId, language.Id)
	})
}

func TestClarificationRetrieval(t *testing.T) {
	api := interactor(t)

	var clarificationId string
	t.Run("all-clarifications", func(t *testing.T) {
		clarifications, err := List(api, Clarification{})
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
			t.Skip("no clarifications could be found, retrieving single clarification cannot be tested")
		}

		clarification, err := GetById[Clarification](api, clarificationId)
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

		id, err := api.PostClarification(testProblem, "testing clarification")
		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent clarification, got id: '%v'", id)
	})

	t.Run("authorized-struct", func(t *testing.T) {
		api := interactor(t)
		clar := Clarification{
			ProblemId: testProblem,
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
		id, err := api.PostSubmission(testProblem, "cpp", "", sampleSubmission)
		assert.Nil(t, err)
		assert.NotEmpty(t, id)

		t.Logf("Sent submission, got id: '%v'", id)
	})

	t.Run("authorized-struct", func(t *testing.T) {
		api := teamInteractor(t)
		var sampleSubmission LocalFileReference
		_ = sampleSubmission.FromString("sample.cpp", "int main() { return 0; }")
		submission := Submission{
			ProblemId:  testProblem,
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

func TestInvalidCert(t *testing.T) {
	// This test forces x509 key errors by using a proxy with an invalid certificate

	// Create handling with a reverse proxy
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rr, err := http.NewRequest(request.Method, testBase+request.URL.String(), request.Body)
		assert.Nil(t, err)

		resp, err := http.DefaultClient.Do(rr)
		assert.Nil(t, err)

		writer.WriteHeader(resp.StatusCode)
		_, err = io.Copy(writer, resp.Body)
		assert.Nil(t, err)
	})

	// Start an internal tls server that is only valid for "example.com"
	var err error
	ser := httptest.NewUnstartedServer(http.DefaultServeMux)
	ser.Listener, err = net.Listen("tcp", "127.0.0.1:8888")
	assert.Nil(t, err)
	go ser.StartTLS()
	defer ser.Close()

	t.Run("contest", func(t *testing.T) {
		api, err := ContestInteractor("https://localhost:8888", testUser, testPass, testContest, false)
		assert.NotNil(t, err)
		assert.Nil(t, api)

		isX509, err := unwrapsToX509Error(t, err)
		assert.True(t, isX509)
		assert.NotNil(t, err)
	})

	t.Run("contests", func(t *testing.T) {
		api, err := ContestsInteractor("https://localhost:8888", testUser, testPass, false)
		assert.Nil(t, err)
		assert.NotNil(t, api)

		contests, err := api.Contests()
		assert.Nil(t, contests)
		assert.NotNil(t, err)

		isX509, err := unwrapsToX509Error(t, err)
		assert.True(t, isX509)
		assert.NotNil(t, err)
	})

	t.Run("problems", func(t *testing.T) {
		api := interactor(t)
		api.(*inter).baseUrl = "https://localhost:8888/"

		problems, err := List(api, Problem{})
		assert.Nil(t, problems)
		assert.NotNil(t, err)

		isX509, err := unwrapsToX509Error(t, err)
		assert.True(t, isX509)
		assert.NotNil(t, err)
	})

}

func unwrapsToX509Error(t testLogger, err error) (bool, error) {
	for err != nil {
		// The following is an exhaustive list of all x509 errors
		switch err.(type) {
		case x509.ConstraintViolationError:
			return true, err
		case x509.HostnameError:
			return true, err
		case x509.SystemRootsError:
			return true, err
		case x509.CertificateInvalidError:
			return true, err
		case x509.InsecureAlgorithmError:
			return true, err
		case x509.UnknownAuthorityError:
			return true, err
		}

		t.Logf("Non x509 error found: %v (%T)", err, err)
		err = errors.Unwrap(err)
	}

	return false, nil
}
