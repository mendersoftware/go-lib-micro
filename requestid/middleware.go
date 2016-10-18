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
package requestid

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/mendersoftware/go-lib-micro/requestlog"
	"github.com/satori/go.uuid"
)

const RequestIdHeader = "X-MEN-RequestID"

// RequestIdMiddleware sets the X-MEN-RequestID header if it's not present, and and adds the request id to the request's logger's context.
type RequestIdMiddleware struct {
}

// MiddlewareFunc makes RequestIdMiddleware implement the Middleware interface.
func (mw *RequestIdMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		reqId := r.Header.Get(RequestIdHeader)
		if reqId == "" {
			reqId = uuid.NewV4().String()
		}

		r.Env[RequestIdHeader] = reqId

		// enrich log context
		logger := r.Env[requestlog.ReqLog]
		if logger != nil {
			logger := logger.(*log.Logger)
			logger = logger.F(log.Ctx{"request_id": reqId})
			r.Env[requestlog.ReqLog] = logger
		}

		//return the reuqest ID in response too, the client can log it
		//for end-to-end req tracing
		w.Header().Add(RequestIdHeader, reqId)

		h(w, r)
	}
}
