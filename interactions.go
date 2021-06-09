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
		return nil, fmt.Errorf("could not retrieve contests; %w", err)
	}

	// obj should be a slice of Contest, cast to it to slice of Contest
	ret := make([]Contest, len(obj))
	for k, v := range obj {
		vv, ok := v.(Contest)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected contest, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ContestById(contestId string) (c Contest, err error) {
	obj, err := i.GetObject(c, contestId)
	if err != nil {
		return c, fmt.Errorf("could not retrieve contest; %w", err)
	}

	vv, ok := obj.(Contest)
	if !ok {
		return c, fmt.Errorf("unexpected type found, expected contest, got: %T", obj)
	}

	c = vv
	return
}

func (i inter) Problems() ([]Problem, error) {
	obj, err := i.GetObjects(Problem{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve problems; %w", err)
	}

	// obj should be a slice of Problem, cast to it to slice of Problem
	ret := make([]Problem, len(obj))
	for k, v := range obj {
		vv, ok := v.(Problem)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected problem, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ProblemById(problemId string) (p Problem, err error) {
	obj, err := i.GetObject(p, problemId)
	if err != nil {
		return p, fmt.Errorf("could not retrieve problem; %w", err)
	}

	vv, ok := obj.(Problem)
	if !ok {
		return p, fmt.Errorf("unexpected type found, expected problem, got: %T", obj)
	}

	p = vv
	return
}

func (i inter) Submissions() ([]Submission, error) {
	obj, err := i.GetObjects(Submission{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve submissions; %w", err)
	}

	// obj should be a slice of Submission, cast to it to slice of Submission
	ret := make([]Submission, len(obj))
	for k, v := range obj {
		vv, ok := v.(Submission)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected submission, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) SubmissionById(submissionId string) (s Submission, err error) {
	obj, err := i.GetObject(s, submissionId)
	if err != nil {
		return s, fmt.Errorf("could not retrieve submission; %w", err)
	}

	vv, ok := obj.(Submission)
	if !ok {
		return s, fmt.Errorf("unexpected type found, expected submission, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) Languages() ([]Language, error) {
	obj, err := i.GetObjects(Language{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve languages; %w", err)
	}

	// obj should be a slice of Language, cast to it to slice of Language
	ret := make([]Language, len(obj))
	for k, v := range obj {
		vv, ok := v.(Language)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected language, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) LanguageById(languageId string) (l Language, err error) {
	obj, err := i.GetObject(l, languageId)
	if err != nil {
		return l, fmt.Errorf("could not retrieve language; %w", err)
	}

	vv, ok := obj.(Language)
	if !ok {
		return l, fmt.Errorf("unexpected type found, expected language, got: %T", obj)
	}

	l = vv
	return
}

func (i inter) JudgementTypes() ([]JudgementType, error) {
	obj, err := i.GetObjects(JudgementType{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve judgement types; %w", err)
	}

	// obj should be a slice of JudgementType, cast to it to slice of JudgementType
	ret := make([]JudgementType, len(obj))
	for k, v := range obj {
		vv, ok := v.(JudgementType)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected judgement type, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) JudgementTypeById(judgementTypeId string) (jt JudgementType, err error) {
	obj, err := i.GetObject(jt, judgementTypeId)
	if err != nil {
		return jt, fmt.Errorf("could not retrieve judgement type; %w", err)
	}

	vv, ok := obj.(JudgementType)
	if !ok {
		return jt, fmt.Errorf("unexpected type found, expected judgement type, got: %T", obj)
	}

	jt = vv
	return
}

func (i inter) Judgements() ([]Judgement, error) {
	obj, err := i.GetObjects(Judgement{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve judgements; %w", err)
	}

	// obj should be a slice of Judgement, cast to it to slice of Judgement
	ret := make([]Judgement, len(obj))
	for k, v := range obj {
		vv, ok := v.(Judgement)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected judgement, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) JudgementById(judgementId string) (j Judgement, err error) {
	obj, err := i.GetObject(j, judgementId)
	if err != nil {
		return j, fmt.Errorf("could not retrieve judgement; %w", err)
	}

	vv, ok := obj.(Judgement)
	if !ok {
		return j, fmt.Errorf("unexpected type found, expected judgement, got: %T", obj)
	}

	j = vv
	return
}

func (i inter) Clarifications() ([]Clarification, error) {
	obj, err := i.GetObjects(Clarification{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve clarifications; %w", err)
	}

	// obj should be a slice of Clarification, cast to it to slice of Clarification
	ret := make([]Clarification, len(obj))
	for k, v := range obj {
		vv, ok := v.(Clarification)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected clarification, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) ClarificationById(clarificationId string) (c Clarification, err error) {
	obj, err := i.GetObject(c, clarificationId)
	if err != nil {
		return c, fmt.Errorf("could not retrieve clarification; %w", err)
	}

	vv, ok := obj.(Clarification)
	if !ok {
		return c, fmt.Errorf("unexpected type found, expected clarification, got: %T", obj)
	}

	c = vv
	return
}

func (i inter) Groups() ([]Group, error) {
	obj, err := i.GetObjects(Group{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve groups; %w", err)
	}

	// obj should be a slice of Group, cast to it to slice of Group
	ret := make([]Group, len(obj))
	for k, v := range obj {
		vv, ok := v.(Group)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected group, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) GroupById(groupId string) (g Group, err error) {
	obj, err := i.GetObject(g, groupId)
	if err != nil {
		return g, fmt.Errorf("could not retrieve group; %w", err)
	}

	vv, ok := obj.(Group)
	if !ok {
		return g, fmt.Errorf("unexpected type found, expected group, got: %T", obj)
	}

	g = vv
	return
}

func (i inter) Organizations() ([]Organization, error) {
	obj, err := i.GetObjects(Organization{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve organizations; %w", err)
	}

	// obj should be a slice of Organization, cast to it to slice of Organization
	ret := make([]Organization, len(obj))
	for k, v := range obj {
		vv, ok := v.(Organization)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected organization, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) OrganizationById(organizationId string) (o Organization, err error) {
	obj, err := i.GetObject(o, organizationId)
	if err != nil {
		return o, fmt.Errorf("could not retrieve organization; %w", err)
	}

	vv, ok := obj.(Organization)
	if !ok {
		return o, fmt.Errorf("unexpected type found, expected organization, got: %T", obj)
	}

	o = vv
	return
}

func (i inter) Teams() ([]Team, error) {
	obj, err := i.GetObjects(Team{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve teams; %w", err)
	}

	// obj should be a slice of Team, cast to it to slice of Team
	ret := make([]Team, len(obj))
	for k, v := range obj {
		vv, ok := v.(Team)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected team, got: %T", v)
		}

		ret[k] = vv
	}

	return ret, nil
}

func (i inter) TeamById(teamId string) (t Team, err error) {
	obj, err := i.GetObject(t, teamId)
	if err != nil {
		return t, fmt.Errorf("could not retrieve team; %w", err)
	}

	vv, ok := obj.(Team)
	if !ok {
		return t, fmt.Errorf("unexpected type found, expected team, got: %T", obj)
	}

	t = vv
	return
}

func (i inter) Scoreboard() (s Scoreboard, err error) {
	obj, err := i.GetObject(s, "")
	if err != nil {
		return s, fmt.Errorf("could not retrieve scoreboard; %w", err)
	}

	vv, ok := obj.(Scoreboard)
	if !ok {
		return s, fmt.Errorf("unexpected type found, expected scoreboard, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) State() (s State, err error) {
	obj, err := i.GetObject(s, "")
	if err != nil {
		return s, fmt.Errorf("could not retrieve state; %w", err)
	}

	vv, ok := obj.(State)
	if !ok {
		return s, fmt.Errorf("unexpected type found, expected state, got: %T", obj)
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
		return c, fmt.Errorf("could not post clarification; %w", err)
	}

	vv, ok := obj.(Clarification)
	if !ok {
		return c, fmt.Errorf("unexpected type found, expected clarification, got: %T", obj)
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
		return s, fmt.Errorf("could not post submission; %w", err)
	}

	vv, ok := obj.(Submission)
	if !ok {
		return s, fmt.Errorf("unexpected type found, expected submission, got: %T", obj)
	}

	s = vv
	return
}

func (i inter) Submit(s Submittable) (ApiType, error) {
	return i.post(s, s)
}

func (i inter) GetObject(interactor ApiType, id string) (ApiType, error) {
	objs, err := i.retrieve(interactor, i.toPath(interactor)+id, true)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve; %w", err)
	}

	if len(objs) != 1 {
		return nil, fmt.Errorf("incorrect number of objects found, expected 1, got: %v", len(objs))
	}

	return objs[0], nil
}

func (i inter) toPath(interactor ApiType) string {
	var base string
	if interactor.InContest() {
		base = "contests/" + i.contestId + "/"
	}

	return base + interactor.Path() + "/"
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
			return nil, fmt.Errorf("could not read entire body; %w", err)
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
		return nil, fmt.Errorf("could not post request; %w", err)
	}

	defer resp.Body.Close()

	if err := responseToError(resp); err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read entire body; %w", err)
	}

	in, err := interactor.FromJSON(bts)
	return in, err
}

func responseToError(r *http.Response) error {
	var statusErr error
	switch r.StatusCode {
	case http.StatusOK:
		return statusErr
	case http.StatusUnauthorized:
		statusErr = errUnauthorized
	case http.StatusNotFound:
		statusErr = errNotFound
	default:
		statusErr = fmt.Errorf("invalid statuscode received: %d", r.StatusCode)
	}

	// Read the contents
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("API error encountered (%v), furthermore an error was encountered while reading the response: %w", statusErr, err)
	}

	return fmt.Errorf("API error '%s' (%w)", bts, statusErr)
}
