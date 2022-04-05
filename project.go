package bitbucket

import (
	"bytes"
	"encoding/json"
)

type Project struct {
	*Api
	Key string
}

func NewProject(host, key string) *Project {
	return &Project{
		NewApi(host, "projects"),
		key,
	}
}

const defaultBranch = "master"

func (p *Project) GetRepo(slug string) (resp *Response, err error) {
	url := p.Url([]string{p.Key, "repos", slug}, nil)
	return p.Get(url)
}

type postRepoReq struct {
	Name          string `json:"name"`
	ScmId         string `json:"scmId"`
	Forkable      bool   `json:"forkable"`
	DefaultBranch string `json:"defaultBranch"`
}

func (p *Project) PostRepo(name string) (resp *Response, err error) {
	req := postRepoReq{
		name,
		"git",
		true,
		defaultBranch,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	url := p.Url([]string{p.Key, "repos"}, nil)
	return p.Post(url, bytes.NewReader(body))
}
