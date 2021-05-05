package interactor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Scoreboard struct {
		EventId     Identifier `json:"event_id"`
		Time        ApiTime    `json:"time"`
		ContestTime ApiRelTime `json:"contest_time"`
		State       State      `json:"state"`
		Rows        []Row      `json:"rows"`
	}

	State struct {
		Started      *ApiTime `json:"started"`
		Ended        *ApiTime `json:"ended"`
		Frozen       *ApiTime `json:"frozen"`
		Thawed       *ApiTime `json:"thawed,omitempty"`
		Finalized    *ApiTime `json:"finalized"`
		EndOfUpdates *ApiTime `json:"end_of_updates"`
	}

	Row struct {
		Rank     int            `json:"rank"`
		TeamId   Identifier     `json:"team_id"`
		Score    Score          `json:"score,omitempty"`
		Problems []ScoreProblem `json:"problems"`
	}

	Score struct {
		NumSolved int     `json:"num_solved,omitempty"`
		TotalTime int     `json:"total_time,omitempty"`
		Score     float64 `json:"score,omitempty"`
	}

	ScoreProblem struct {
		ProblemId  Identifier `json:"problem_id"`
		NumJudged  int        `json:"num_judged"`
		NumPending int        `json:"num_pending"`
		Solved     bool       `json:"solved"`
		Score      float64    `json:"score,omitempty"`
		Time       int        `json:"time"`
	}
)

func (s Scoreboard) FromJSON(data []byte) (ApiType, error) {
	err := json.Unmarshal(data, &s)
	return s, err
}

func (s Scoreboard) String() string {
	rows := make([]string, len(s.Rows))
	for k, row := range s.Rows {
		rows[k] = row.String()
	}

	return fmt.Sprintf(`
    event id: %v
        time: %v
contest time: %v
       state: %v
        rows: %v
`, s.EventId, s.Time, s.ContestTime, s.State, strings.Join(rows, ""))
}

func (s Scoreboard) InContest() bool {
	return true
}

func (s Scoreboard) Path() string {
	return "scoreboard"
}

func (s Scoreboard) Generate() ApiType {
	return Scoreboard{}
}

func (r Row) String() string {
	problems := make([]string, len(r.Problems))
	for k, problem := range r.Problems {
		problems[k] = problem.String()
	}

	return fmt.Sprintf(`
        rank: %v
     team id: %v
       score: %v
    problems: %v
`, r.Rank, r.TeamId, r.Score, strings.Join(problems, ""))
}

func (s Score) String() string {
	return fmt.Sprintf(`
    num solved: %v
    total time: %v
         score: %v
`, s.NumSolved, s.TotalTime, s.Score)
}

func (s ScoreProblem) String() string {
	return fmt.Sprintf(`
         problem id: %v
         num judged: %v
        num pending: %v
             solved: %v
              score: %v
               time: %v
`, s.ProblemId, s.NumJudged, s.NumPending, s.Solved, s.Score, s.Time)
}

func (s State) String() string {
	return fmt.Sprintf(`
           started: %v
             ended: %v
            frozen: %v
            thawed: %v
         finalized: %v
    end of updates: %v
`, s.Started, s.Ended, s.Frozen, s.Thawed, s.Finalized, s.EndOfUpdates)
}

func (s State) FromJSON(data []byte) (ApiType, error) {
	err := json.Unmarshal(data, &s)
	return s, err
}

func (s State) Path() string {
	return "state"
}

func (s State) Generate() ApiType {
	return State{}
}

func (s State) InContest() bool {
	return true
}
