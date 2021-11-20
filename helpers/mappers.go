package helpers

var ScopeClaimMapper = map[string][]string{
	"profile": []string{"Name", "FamilyName", "GivenName", "NickName", "Picture"},
	"email":   []string{"Email"},
	//"address": []string{"address"}, To be implemented later
	"phone": []string{"PhoneNumber"},
}

var ScopeClaimMapperGerman = map[string][]string{
	"profile": []string{"Name", "Nachname", "Vorname", "Nickname", "Profilbild"},
	"email":   []string{"E-mail Adresse"},
	//"address": []string{"address"}, To be implemented later
	"phone": []string{"Telefonnummer"},
}

var ScopeClaimMapperEnglish = map[string][]string{
	"profile": []string{"name", "last name", "first name", "nickname", "profile image"},
	"email":   []string{"e-mail address"},
	//"address": []string{"address"}, To be implemented later
	"phone": []string{"telephone number"},
}

func GetClaims(grantScopes []string) []string {
	var claims []string
	for _, grantScope := range grantScopes {
		claims = append(claims, ScopeClaimMapper[grantScope]...)
	}

	return claims
}
