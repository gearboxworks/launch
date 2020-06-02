package gear

import (
	"github.com/newclarity/scribeHelpers/helperDocker"
	"github.com/newclarity/scribeHelpers/ux"
	"net/url"
	"os"
)

const (
	ProviderDocker = "docker"
)

type Provider struct {
	Name    string  `json:"name"`
	Host    string  `json:"host"`
	Port    string  `json:"port"`
	Url     url.URL `json:"url"`
	Project string  `json:"project"`

	Debug   bool    //`json:"debug"`
	State   *ux.State
}


func (p *Provider) NewProvider(debugMode bool) *ux.State {
	for range OnlyOnce {
		var err error
		p.State = p.State.EnsureNotNil()
		p.State.DebugSet(debugMode)
		p.Debug = debugMode

		if p.Name == "" {
			p.Name = ProviderDocker
		}

		if p.Name == ProviderDocker {
			for range OnlyOnce {
				if p.Host == "" {
					break
				}

				if p.Port == "" {
					break
				}

				var urlString *url.URL
				//urlString, err = client.ParseHostURL(fmt.Sprintf("tcp://%s:%s", p.Host, p.Port))
				urlString, err = helperDocker.ParseHostURL("tcp://%s:%s", p.Host, p.Port)
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

	return p.State
}
