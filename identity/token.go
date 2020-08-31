// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package identity

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Identity struct {
	Subject  string `json:"sub" valid:"required"`
	Tenant   string `json:"mender.tenant,omitempty"`
	IsUser   bool   `json:"mender.user,omitempty"`
	IsDevice bool   `json:"mender.device,omitempty"`
	Plan     string `json:"mender.plan,omitempty"`
}

// Generate identity information from given JWT by extracting subject and tenant claims.
// Note that this function does not perform any form of token signature
// verification.
func ExtractIdentity(token string) (id Identity, err error) {
	var (
		b64Claims string
		claims    []byte
		jwt       []string
	)
	jwt = strings.Split(token, ".")
	if len(jwt) != 3 {
		return id, errors.New("identity: incorrect token format")
	}
	b64Claims = jwt[1]
	if pad := len(b64Claims) % 4; pad != 0 {
		b64Claims += strings.Repeat("=", 4-pad)
	}
	claims, err = base64.StdEncoding.DecodeString(b64Claims)
	if err != nil {
		return id, errors.Wrap(err,
			"identity: failed to decode base64 JWT claims")
	}
	err = json.Unmarshal(claims, &id)
	if err != nil {
		return id, errors.Wrap(err,
			"identity: failed to decode JSON JWT claims")
	}
	return id, id.Validate()
}

// Extract identity information from HTTP Authorization header. The header is
// assumed to contain data in format: `Bearer <token>`
func ExtractIdentityFromHeaders(headers http.Header) (Identity, error) {
	auth := strings.Split(headers.Get("Authorization"), " ")

	if len(auth) != 2 {
		return Identity{}, errors.Errorf("malformed authorization data")
	}

	if auth[0] != "Bearer" {
		return Identity{}, errors.Errorf("unknown authorization method %v", auth[0])
	}

	return ExtractIdentity(auth[1])
}

func (id Identity) Validate() error {
	if id.Subject == "" {
		return errors.New("identity: claim \"sub\" is required")
	}
	return nil
}
