/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package helpers

var ScopeClaimMapper = map[string][]string{
	"profile": []string{"Name", "FamilyName", "GivenName", "NickName", "Picture"},
	"email":   []string{"Email"},
	//"address": []string{"address"}, To be implemented later
	"phone":  []string{"PhoneNumber"},
	"groups": []string{"Groups"},
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
