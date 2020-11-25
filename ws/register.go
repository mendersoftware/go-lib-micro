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

package ws

import (
	"errors"
)

var (
	ErrUnsupportedTypeError = errors.New("unsupported type")
	ErrUnregisteredProtocol = errors.New("protocol not registered")
)

type Protocol interface {
	DeEncapsulate(m *ProtoMsg) (interface{}, error)
	Encapsulate(m interface{}) (*ProtoMsg, error)
}

var protoMap = map[ProtoType]Protocol{}

func RegisterProtocol(t ProtoType, protocol Protocol) error {
	protoMap[t] = protocol
	return nil
}

func DeEncapsulate(m *ProtoMsg) (interface{}, error) {
	if v, ok := protoMap[m.Header.Proto]; ok {
		return v.DeEncapsulate(m)
	} else {
		return nil, ErrUnregisteredProtocol
	}
}

func Encapsulate(t ProtoType, m interface{}) (*ProtoMsg, error) {
	if v, ok := protoMap[t]; ok {
		return v.Encapsulate(m)
	} else {
		return nil, ErrUnregisteredProtocol
	}
}
