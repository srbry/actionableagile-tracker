package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Projects []Project

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c Client) Project(name string) (Project, error) {
	request, err := http.NewRequest("GET", url("/projects"), nil)
	if err != nil {
		return Project{}, err
	}

	request.Header.Add("X-TrackerToken", c.Config.Token)
	resp, err := c.conn.Do(request)
	if err != nil {
		return Project{}, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Project{}, err
	}
	var projects Projects
	json.Unmarshal(data, &projects)
	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}
	return Project{}, fmt.Errorf("Project %s not found", name)
}
