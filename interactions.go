package interactor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (i inter) Contests() ([]Contest, error) {
	obj, err := i.GetObjects(new(Contest))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve contests; %w", err)
	}

	// obj should be a slice of *Contest, cast to it to slice of Contest
	ret := make([]Contest, len(obj))
	for k, v := range obj {
		vv, ok := v.(*Contest)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected contest, got: %T", v)
		}

		ret[k] = *vv
	}

	return ret, nil
}

func (i inter) ContestById(contestId string) (c Contest, err error) {
	// Retrieve all contests and check whether the contest exists, TODO decide on whether to optimize
	contests, err := i.Contests()
	if err != nil {
		return c, fmt.Errorf("could not retrieve contest")
	}

	for _, v := range contests {
		if v.Id == contestId {
			return v, nil
		}
	}

	return c, errNotFound
}

func (i inter) Problems() ([]Problem, error) {
	obj, err := i.GetObjects(new(Problem))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve problems; %w", err)
	}

	// obj should be a slice of *Contest, cast to it to slice of Contest
	ret := make([]Problem, len(obj))
	for k, v := range obj {
		vv, ok := v.(*Problem)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected problem, got: %T", v)
		}

		ret[k] = *vv
	}

	return ret, nil
}

func (i inter) ProblemById(problemId string) (p Problem, err error) {
	obj, err := i.GetObject(new(Problem), problemId)
	if err != nil {
		return p, fmt.Errorf("could not retrieve problem; %w", err)
	}

	vv, ok := obj.(*Problem)
	if !ok {
		return p, fmt.Errorf("unexpected type found, expected problem, got: %T", obj)
	}

	p = *vv
	return
}

func (i inter) Submissions() ([]Submission, error) {
	obj, err := i.GetObjects(new(Submission))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve submissions; %w", err)
	}

	// obj should be a slice of *Contest, cast to it to slice of Contest
	ret := make([]Submission, len(obj))
	for k, v := range obj {
		vv, ok := v.(*Submission)
		if !ok {
			return ret, fmt.Errorf("unexpected type found, expected submission, got: %T", v)
		}

		ret[k] = *vv
	}

	return ret, nil
}

func (i inter) SubmissionById(submissionId string) (s Submission, err error) {
	obj, err := i.GetObject(new(Submission), submissionId)
	if err != nil {
		return s, fmt.Errorf("could not retrieve submission; %w", err)
	}

	vv, ok := obj.(*Submission)
	if !ok {
		return s, fmt.Errorf("unexpected type found, expected submission, got: %T", obj)
	}

	s = *vv
	return
}

func (i inter) GetObject(interactor ApiType, id string) (ApiType, error) {
	objs, err := i.retrieve(interactor, id)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve problem; %w", err)
	}

	if len(objs) != 1 {
		return nil, fmt.Errorf("incorrect number of objects found, expected 1, got: %v", len(objs))
	}

	return objs[0], nil
}

func (i inter) GetObjects(interactor ApiType) ([]ApiType, error) {
	return i.retrieve(interactor, "")
}

func (i inter) retrieve(interactor ApiType, id string) ([]ApiType, error) {
	resp, err := i.Get(i.baseUrl + interactor.Path(i.contestId, id))
	if err != nil {
		return nil, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode); err != nil {
		return nil, err
	}

	// Some json should be returned, construct a decoder
	decoder := json.NewDecoder(resp.Body)

	// If id is not empty, only a single instance is expected to be returned
	if id != "" {
		in := interactor.Generate()
		return []ApiType{in}, decoder.Decode(in)
	}

	// We read everything into a slice of
	var temp []json.RawMessage
	if err := decoder.Decode(&temp); err != nil {
		return nil, err
	}

	// Create the actual slice to return
	ret := make([]ApiType, len(temp))
	for k, v := range temp {
		// Generate a new interactor
		in := interactor.Generate()
		if err := in.FromJSON(v); err != nil {
			return ret, err
		}

		ret[k] = in
	}

	return ret, nil
}

func statusToError(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errUnauthorized
	case http.StatusNotFound:
		return errNotFound
	default:
		return fmt.Errorf("invalid statuscode received: %d", status)
	}
}
