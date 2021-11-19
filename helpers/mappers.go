package helpers

var ScopeClaimMapper = map[string][]string{
	"profile": []string{"Name", "FamilyName", "GivenName", "NickName", "Picture"},
	"email":   []string{"Email"},
	//"address": []string{"address"}, To be implemented later
	"phone": []string{"PhoneNumber"},
}

func GetClaims(grantScopes []string) []string {
	var claims []string
	for _, grantScope := range grantScopes {
		claims = append(claims, ScopeClaimMapper[grantScope]...)
	}

	return claims
}
