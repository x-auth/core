package controllers

import (
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"net/url"
	"nictec.net/auth/helpers"
)

var hydraAdmin admin.ClientService

func InitHydra(){
	adminURL, _ := url.Parse(helpers.Config.HydraURL)
	hydraAdmin = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}).Admin
}

func getAdmin() admin.ClientService {
	return hydraAdmin
}