package interactor

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type (
	ContestsApi interface {
		Contests() ([]Contest, error)
		ContestById(contestId string) (Contest, error)
		ToContest(cid string) (ContestApi, error)
	}

	ContestApi interface {
		ContestsApi

		Problems() ([]Problem, error)
		ProblemById(problemId string) (Problem, error)

		JudgementTypes() ([]JudgementType, error)
		JudgementTypeById(judgementTypeId string) (JudgementType, error)

		Submissions() ([]Submission, error)
		SubmissionById(submissionId string) (Submission, error)

		Judgements() ([]Judgement, error)
		JudgementById(judgementId string) (Judgement, error)

		Clarifications() ([]Clarification, error)
		ClarificationById(clarificationId string) (Clarification, error)

		Languages() ([]Language, error)
		LanguageById(languageId string) (Language, error)

		Groups() ([]Group, error)
		GroupById(groupId string) (Group, error)

		Organizations() ([]Organization, error)
		OrganizationById(organizationId string) (Organization, error)

		Teams() ([]Team, error)
		TeamById(teamId string) (Team, error)

		GetObject(interactor ApiType, id string) (ApiType, error)
		GetObjects(interactor ApiType) ([]ApiType, error)

		Submit(submittable Submittable) (ApiType, error)
		PostClarification(problemId, text string) (Clarification, error)
		PostSubmission(problemId, languageId, entrypoint string, files LocalFileReference) (Submission, error)

		Scoreboard() (Scoreboard, error)
		State() (State, error)
	}

	inter struct {
		http.Client
		contestId string
		username  string
		password  string
		baseUrl   string
	}

	// Implementation of the http.RoundTripper interface, used for always adding basic-auth
	basicAuthTransport struct {
		username, password string

		T http.RoundTripper
	}
)

var (
	errUnauthorized = errors.New("request not authorized")
	errNotFound     = errors.New("object not found")
)

// RoundTrip adds the basic auth headers
func (b basicAuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if b.username != "" && b.password != "" {
		request.SetBasicAuth(b.username, b.password)
	}

	return b.T.RoundTrip(request)
}

func ContestInteractor(baseUrl, username, password, contestId string, insecure bool) (ContestApi, error) {
	i := &inter{
		baseUrl:   strings.TrimRight(baseUrl, "/") + "/",
		username:  username,
		password:  password,
		contestId: contestId,
		Client:    buildClient(username, password, insecure),
	}

	if _, err := i.ContestById(contestId); err != nil {
		// If the contest cannot be found, ensure the interactor cannot be used
		return nil, fmt.Errorf("could not find contest; %w", err)
	}

	return i, nil
}

func ContestsInteractor(baseUrl, username, password string, insecure bool) (ContestsApi, error) {
	return &inter{
		baseUrl:  strings.TrimRight(baseUrl, "/") + "/",
		username: username,
		password: password,
		Client:   buildClient(username, password, insecure),
	}, nil
}

// ToContest "upgrades" a ContestsApi to a ContestApi for a specific contest. When called from a ContestApi it can be
// used to change the current contest associated with that ContestApi.
func (i *inter) ToContest(cid string) (ContestApi, error) {
	i.contestId = cid

	if _, err := i.ContestById(cid); err != nil {
		return nil, fmt.Errorf("contest could not be found; %w", err)
	}

	return i, nil
}

func buildClient(username, password string, insecure bool) http.Client {
	// Create a transport for (possibly) insecure communication and adding of basic-auth headers
	transport := http.DefaultTransport.(*http.Transport)
	// if insecure {
	// Check and reuse existing tls config
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = new(tls.Config)
	}

	transport.TLSClientConfig.InsecureSkipVerify = insecure
	// }

	return http.Client{
		Transport: basicAuthTransport{username, password, transport},
	}
}
