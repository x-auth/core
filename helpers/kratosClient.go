package helpers

import client "github.com/ory/kratos-client-go"

var KratosClient *client.APIClient

func InitKratosClient() {
	var cfg = client.Configuration{Host: Config.KratosURL}
	KratosClient = client.NewAPIClient(&cfg)
}
