package gear

import (
	"fmt"
	"github.com/docker/docker/client"
	"launch/only"
	"launch/ux"
	"net/url"
	"os"
)

const (
	ProviderDocker = "docker"
)

type Provider struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Url     url.URL `json:"url"`
	Project string `json:"project"`
	Debug   bool `json:"debug"`
}


func (p *Provider) NewProvider() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		if p.Name == "" {
			p.Name = ProviderDocker
		}

		if p.Name == ProviderDocker {
			for range only.Once {
				if p.Host == "" {
					break
				}

				if p.Port == "" {
					break
				}

				var urlString *url.URL
				urlString, err = client.ParseHostURL(fmt.Sprintf("tcp://%s:%s", p.Host, p.Port))
				if err != nil {
					break
				}

				err = os.Setenv("DOCKER_HOST", urlString.String())
				if err != nil {
					break
				}
			}

			break
		}
	}

	return state
}


//func (me *Provider) NewGear() (*Gear, ux.State) {
//	var g Gear
//	var state ux.State
//
//	for range only.Once {
//		g.Docker, state = dockerClient.New()
//		if state.IsError() {
//			state.SetError("can not connect to Docker service provider")
//			break
//		}
//		g.Docker.Debug = me.Debug
//
//		g.Repo, state = githubClient.New()
//		state.ClearError()
//		//if state.IsError() {
//		//	break
//		//}
//	}
//
//	return &g, state
//}
