package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var DefaultURL = "https://www.pivotaltracker.com"

type Client struct {
	conn   *http.Client
	Config Config
}

type Config struct {
	Username string
	Password string
	Token    string
}

type me struct {
	Token string `json:"api_token"`
}

func NewClient(config Config) *Client {
	client := &Client{
		conn:   &http.Client{},
		Config: config,
	}
	myself, _ := client.Me()
	client.Config.Token = myself.Token
	return client
}

func (c Client) Me() (me, error) {
	request, err := http.NewRequest("GET", url("/me"), nil)
	if err != nil {
		return me{}, err
	}

	request.SetBasicAuth(c.Config.Username, c.Config.Password)
	resp, err := c.conn.Do(request)
	if err != nil {
		return me{}, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return me{}, err
	}
	var myself me
	json.Unmarshal(data, &myself)
	return myself, err
}

func url(path string) string {
	return fmt.Sprintf("%s/services/v5%s", DefaultURL, path)
}
