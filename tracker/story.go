package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StoryTransitions []StoryTransition

type StoryTransition struct {
	Kind           string `json:"kind"`
	State          string `json:"state"`
	StoryID        int    `json:"story_id"`
	ProjectID      int    `json:"project_id"`
	ProjectVersion int    `json:"project_version"`
	OccurredAt     string `json:"occurred_at"`
	PerformedBy    int    `json:"performed_by_id"`
}

func (c Client) StoryTransitions(projectID, storyID int) (StoryTransitions, error) {
	request, err := http.NewRequest("GET", url(fmt.Sprintf(
		"/projects/%v/stories/%v/transitions",
		projectID, storyID,
	)), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-TrackerToken", c.Config.Token)
	resp, err := c.conn.Do(request)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var storyTransitions StoryTransitions
	json.Unmarshal(data, &storyTransitions)
	return storyTransitions, nil
}
