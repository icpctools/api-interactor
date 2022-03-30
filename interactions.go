package interactor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (i inter) Contests() ([]Contest, error) {
	obj, err := i.GetObjects(Contest{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Contest, cast to it to slice of Contest
	ret := make([]Contest, len(obj))
	for k, v := range obj {
		vv, ok := v.(Contest)
		if !ok {
			return ret, fmt.Errorf("expected contest, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ContestById(contestId string) (c Contest, err error) {
	obj, err := i.GetObject(c, contestId)
	if err != nil {
		return c, err
	}

	vv, ok := obj.(Contest)
	if !ok {
		return c, fmt.Errorf("expected contest, got: %T", obj)
	}

	c = vv
	return
}

func (i inter) Contest() (c Contest, err error) {
	return i.ContestById(i.contestId)
}

func (i inter) Accounts() ([]Account, error) {
	obj, err := i.GetObjects(Account{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Account, cast to it to slice of Account
	ret := make([]Account, len(obj))
	for k, v := range obj {
		vv, ok := v.(Account)
		if !ok {
			return ret, fmt.Errorf("expected account, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) AccountById(accountId string) (a Account, err error) {
	obj, err := i.GetObject(a, accountId)
	if err != nil {
		return a, err
	}

	vv, ok := obj.(Account)
	if !ok {
		return a, fmt.Errorf("expected account, got: %T", obj)
	}

	a = vv
	return
}

func (i inter) Account() (a Account, err error) {
	objs, err := i.retrieve(Account{}, "contests/"+i.contestId+"/account", true)

	if err != nil {
		return a, err
	}

	if len(objs) != 1 {
		return a, fmt.Errorf("expected 1 object, got: %v", len(objs))
	}

	vv, ok := objs[0].(Account)
	if !ok {
		return a, fmt.Errorf("expected account, got: %T", objs[0])
	}

	a = vv
	return
}

func (i inter) Problems() ([]Problem, error) {
	obj, err := i.GetObjects(Problem{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Problem, cast to it to slice of Problem
	ret := make([]Problem, len(obj))
	for k, v := range obj {
		vv, ok := v.(Problem)
		if !ok {
			return ret, fmt.Errorf("expected problem, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ProblemById(problemId string) (p Problem, err error) {
	obj, err := i.GetObject(p, problemId)
	if err != nil {
		return p, err
	}

	vv, ok := obj.(Problem)
	if !ok {
		return p, fmt.Errorf("expected problem, got: %T", obj)
	}

	p = vv
	return
}

func (i inter) Submissions() ([]Submission, error) {
	obj, err := i.GetObjects(Submission{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Submission, cast to it to slice of Submission
	ret := make([]Submission, len(obj))
	for k, v := range obj {
		vv, ok := v.(Submission)
		if !ok {
			return ret, fmt.Errorf("expected submission, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) SubmissionById(submissionId string) (s Submission, err error) {
	obj, err := i.GetObject(s, submissionId)
	if err != nil {
		return s, err
	}

	vv, ok := obj.(Submission)
	if !ok {
		return s, fmt.Errorf("expected submission, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) Languages() ([]Language, error) {
	obj, err := i.GetObjects(Language{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Language, cast to it to slice of Language
	ret := make([]Language, len(obj))
	for k, v := range obj {
		vv, ok := v.(Language)
		if !ok {
			return ret, fmt.Errorf("expected language, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) LanguageById(languageId string) (l Language, err error) {
	obj, err := i.GetObject(l, languageId)
	if err != nil {
		return l, err
	}

	vv, ok := obj.(Language)
	if !ok {
		return l, fmt.Errorf("expected language, got: %T", obj)
	}

	l = vv
	return
}

func (i inter) JudgementTypes() ([]JudgementType, error) {
	obj, err := i.GetObjects(JudgementType{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of JudgementType, cast to it to slice of JudgementType
	ret := make([]JudgementType, len(obj))
	for k, v := range obj {
		vv, ok := v.(JudgementType)
		if !ok {
			return ret, fmt.Errorf("expected judgement type, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) JudgementTypeById(judgementTypeId string) (jt JudgementType, err error) {
	obj, err := i.GetObject(jt, judgementTypeId)
	if err != nil {
		return jt, err
	}

	vv, ok := obj.(JudgementType)
	if !ok {
		return jt, fmt.Errorf("judgement type, got: %T", obj)
	}

	jt = vv
	return
}

func (i inter) Judgements() ([]Judgement, error) {
	obj, err := i.GetObjects(Judgement{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Judgement, cast to it to slice of Judgement
	ret := make([]Judgement, len(obj))
	for k, v := range obj {
		vv, ok := v.(Judgement)
		if !ok {
			return ret, fmt.Errorf("expected judgement, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) JudgementById(judgementId string) (j Judgement, err error) {
	obj, err := i.GetObject(j, judgementId)
	if err != nil {
		return j, err
	}

	vv, ok := obj.(Judgement)
	if !ok {
		return j, fmt.Errorf("expected judgement, got: %T", obj)
	}

	j = vv
	return
}

func (i inter) Clarifications() ([]Clarification, error) {
	obj, err := i.GetObjects(Clarification{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Clarification, cast to it to slice of Clarification
	ret := make([]Clarification, len(obj))
	for k, v := range obj {
		vv, ok := v.(Clarification)
		if !ok {
			return ret, fmt.Errorf("expected clarification, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ClarificationById(clarificationId string) (c Clarification, err error) {
	obj, err := i.GetObject(c, clarificationId)
	if err != nil {
		return c, err
	}

	vv, ok := obj.(Clarification)
	if !ok {
		return c, fmt.Errorf("expected clarification, got: %T", obj)
	}

	c = vv
	return
}

func (i inter) Groups() ([]Group, error) {
	obj, err := i.GetObjects(Group{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Group, cast to it to slice of Group
	ret := make([]Group, len(obj))
	for k, v := range obj {
		vv, ok := v.(Group)
		if !ok {
			return ret, fmt.Errorf("expected group, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) GroupById(groupId string) (g Group, err error) {
	obj, err := i.GetObject(g, groupId)
	if err != nil {
		return g, err
	}

	vv, ok := obj.(Group)
	if !ok {
		return g, fmt.Errorf("expected group, got: %T", obj)
	}

	g = vv
	return
}

func (i inter) Organizations() ([]Organization, error) {
	obj, err := i.GetObjects(Organization{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Organization, cast to it to slice of Organization
	ret := make([]Organization, len(obj))
	for k, v := range obj {
		vv, ok := v.(Organization)
		if !ok {
			return ret, fmt.Errorf("expected organization, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) OrganizationById(organizationId string) (o Organization, err error) {
	obj, err := i.GetObject(o, organizationId)
	if err != nil {
		return o, err
	}

	vv, ok := obj.(Organization)
	if !ok {
		return o, fmt.Errorf("expected organization, got: %T", obj)
	}

	o = vv
	return
}

func (i inter) Teams() ([]Team, error) {
	obj, err := i.GetObjects(Team{})
	if err != nil {
		return nil, err
	}

	// obj should be a slice of Team, cast to it to slice of Team
	ret := make([]Team, len(obj))
	for k, v := range obj {
		vv, ok := v.(Team)
		if !ok {
			return ret, fmt.Errorf("expected team, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) TeamById(teamId string) (t Team, err error) {
	obj, err := i.GetObject(t, teamId)
	if err != nil {
		return t, err
	}

	vv, ok := obj.(Team)
	if !ok {
		return t, fmt.Errorf("expected team, got: %T", obj)
	}

	t = vv
	return
}

func (i inter) Scoreboard() (s Scoreboard, err error) {
	obj, err := i.GetObject(s, "")
	if err != nil {
		return s, err
	}

	vv, ok := obj.(Scoreboard)
	if !ok {
		return s, fmt.Errorf("expected scoreboard, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) State() (s State, err error) {
	obj, err := i.GetObject(s, "")
	if err != nil {
		return s, err
	}

	vv, ok := obj.(State)
	if !ok {
		return s, fmt.Errorf("expected state, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) PostClarification(problemId, text string) (c Clarification, err error) {
	obj, err := i.post(c, Clarification{
		ProblemId: problemId,
		Text:      text,
	})
	if err != nil {
		return c, err
	}

	vv, ok := obj.(Clarification)
	if !ok {
		return c, fmt.Errorf("expected clarification, got: %T", obj)
	}

	c = vv
	return
}

func (i inter) PostSubmission(problemId, languageId, entrypoint string, files LocalFileReference) (s Submission, err error) {
	obj, err := i.post(s, Submission{
		ProblemId:  problemId,
		LanguageId: languageId,
		EntryPoint: entrypoint,
		Files: []FileReference{
			{
				Mime: "application/zip",
				Data: files,
			},
		},
	})
	if err != nil {
		return s, err
	}

	vv, ok := obj.(Submission)
	if !ok {
		return s, fmt.Errorf("expected submission, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) Submit(s Submittable) (ApiType, error) {
	return i.post(s, s)
}

func (i inter) GetObject(interactor ApiType, id string) (ApiType, error) {
	objs, err := i.retrieve(interactor, i.toPath(interactor)+"/"+id, true)

	if err != nil {
		return nil, err
	}

	if len(objs) != 1 {
		return nil, fmt.Errorf("expected 1 object, got: %v", len(objs))
	}

	return objs[0], nil
}

func (i inter) toPath(interactor ApiType) string {
	var base string
	if interactor.InContest() {
		base = "contests/" + i.contestId + "/"
	}

	return base + interactor.Path()
}

func (i inter) GetObjects(interactor ApiType) ([]ApiType, error) {
	return i.retrieve(interactor, i.toPath(interactor), false)
}

func (i inter) retrieve(interactor ApiType, path string, single bool) ([]ApiType, error) {
	resp, err := i.Get(i.baseUrl + path)
	if err != nil {
		return nil, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := responseToError(resp); err != nil {
		return nil, err
	}

	// If single is true, only a single instance is expected to be returned
	if single {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read response body; %w", err)
		}

		in, err := interactor.FromJSON(bts)
		return []ApiType{in}, err
	}

	// Some json should be returned, construct a decoder
	decoder := json.NewDecoder(resp.Body)

	// We read everything into a slice of
	var temp []json.RawMessage
	if err := decoder.Decode(&temp); err != nil {
		return nil, err
	}

	// Create the actual slice to return
	ret := make([]ApiType, len(temp))
	for k, v := range temp {
		// Generate a new interactor
		vv, err := interactor.FromJSON(v)
		if err != nil {
			return ret, err
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) post(interactor ApiType, encodableBody Submittable) (ApiType, error) {
	var buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(encodableBody)
	if err != nil {
		return nil, fmt.Errorf("could not marshal body; %w", err)
	}

	// Post the body
	resp, err := i.Post(i.baseUrl+i.toPath(interactor), "application/json", buf)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := responseToError(resp); err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body; %w", err)
	}

	return interactor.FromJSON(bts)
}

func responseToError(r *http.Response) error {
	var statusErr error
	switch r.StatusCode {
	case http.StatusOK:
		return statusErr
	case http.StatusBadRequest:
		statusErr = errBadRequest
	case http.StatusUnauthorized:
		statusErr = errUnauthorized
	case http.StatusForbidden:
		statusErr = errForbidden
	case http.StatusNotFound:
		statusErr = errNotFound
	case http.StatusConflict:
		statusErr = errConflict
	default:
		statusErr = fmt.Errorf("invalid statuscode received: %d", r.StatusCode)
	}

	// Read the contents
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("API error (%v), could not read response body: %w", statusErr, err)
	}

	var e Error
	err = json.Unmarshal(bts, &e)
	if err != nil {
		return fmt.Errorf("API error (%v), couldn't parse details: %w", statusErr, err)
	}

	return fmt.Errorf("%s (error code %d)", e.Message, e.Code)
}
