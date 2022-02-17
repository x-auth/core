package authenticators

import (
	"github.com/x-auth/common/models"
	"github.com/x-auth/common/plugins"
	"plugin"
	"x-net.at/idp/helpers"
	"x-net.at/idp/logger"
)

var Authenticators = make(map[string]plugins.AuthPlugin)

func Init() {
	cfg := LoadAuthConfig()
	for _, authenticator := range cfg.Authenticators {
		p, err := plugin.Open(cfg.PluginDir + "/" + authenticator.Plugin + ".so")
		if err != nil {
			logger.Log.Fatal(err)
		}

		sym, err := p.Lookup("NewPlugin")
		if err != nil {
			logger.Log.Fatal(err)
		}

		constructor, ok := sym.(func(map[string]string) plugins.AuthPlugin)
		if !ok {
			logger.Log.Fatal("No plugin constructor found")
		}

		logger.Log.Debug("cfg:", authenticator.Config)
		Authenticators[authenticator.Name] = constructor(authenticator.Config)
	}
}

func Login(username string, password string, preflightRealm string) (models.Profile, bool) {
	var realmObj helpers.Realm
	for _, realm := range helpers.Config.Realms {
		if realm.Name == preflightRealm {
			realmObj = realm
		}
	}

	authenticator := Authenticators[realmObj.Authenticator]
	return authenticator.Login(username, password)
}
