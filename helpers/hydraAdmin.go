package helpers

import (
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"net/url"
)

var hydraAdmin admin.ClientService

func InitHydra() error {
	adminURL, _ := url.Parse(Config.HydraURL)
	hydraAdmin = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}).Admin
	_, err := hydraAdmin.IsInstanceAlive(admin.NewIsInstanceAliveParams())
	return err
}

func GetAdmin() admin.ClientService {
	return hydraAdmin
}
