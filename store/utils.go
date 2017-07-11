// Copyright 2016 Mender Software AS
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
package store

import (
	"context"

	"github.com/mendersoftware/go-lib-micro/identity"
)

// DbFromContext generates database name using tenant field from identity extracted
// from context and original database name
func DbFromContext(ctx context.Context, origDbName string) string {
	identity := identity.FromContext(ctx)
	if identity != nil && identity.Tenant != "" {
		return origDbName + "-" + identity.Tenant
	} else {
		return origDbName
	}
}

// DbForTenant generates database name using simply the tenant name and
//the original db name
func DbForTenant(tenantId, origDbName string) string {
	return origDbName + "-" + tenantId
}
