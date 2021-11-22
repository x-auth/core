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

package mock

import (
	"x-net.at/idp/logger"
	"x-net.at/idp/models"
)

func Login(username string, password string, config map[string]string) (models.Profile, bool) {
	if username == config["username"] && password == config["password"] {
		return models.Profile{
			Name:        "Foo Bar",
			FamilyName:  "Bar",
			GivenName:   "Foo",
			NickName:    "foobar",
			Email:       "foobar@example.com",
			PhoneNumber: "000000000",
		}, true
	} else {
		logger.Warning.Println("Login failed, username or password false")
		return models.Profile{}, false
	}
}

func getMockProfile(username string) models.Profile {
	return models.Profile{
		Name:        "Foo Bar",
		FamilyName:  "Bar",
		GivenName:   "Foo",
		NickName:    "foobar",
		Email:       "foobar@example.com",
		PhoneNumber: "000000000",
	}
}
