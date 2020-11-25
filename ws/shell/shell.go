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

package shell

import (
	"github.com/mendersoftware/go-lib-micro/ws"
)

const version = "1.0"

type MenderShellMessageStatus int

const (
	NormalMessage MenderShellMessageStatus = iota
	ErrorMessage
)

var (
	MessageTypeShellCommand = "shell"
	MessageTypeSpawnShell   = "new"
	MessageTypeStopShell    = "stop"
)

// MenderShellMessage represents a message between the device and the backend
type MenderShellMessage struct {
	//type of message, used to determine the meaning of data
	Type string `json:"type" msgpack:"type"`
	//session id, as returned to the caller in a response to the MessageTypeSpawnShell
	//message.
	SessionId string `json:"session_id" msgpack:"session_id"`
	//message status, currently normal and error message types are supported
	Status MenderShellMessageStatus `json:"status_code" msgpack:"status_code"`
	//the message payload, if
	// * .Type===MessageTypeShellCommand interpreted as keystrokes and passed
	//   to the stdin of the terminal running the shell.
	// * .Type===MessageTypeSpawnShell interpreted as user_id and passed
	//   to the session.NewMenderShellSession.
	Data []byte `json:"data" msgpack:"data"`
}

type MenderShellProtocol struct {
	version string
}

func (w *MenderShellProtocol) DeEncapsulate(m *ws.ProtoMsg) (interface{}, error) {
	message := &MenderShellMessage{}
	message.SessionId = m.Header.SessionID
	message.Type = m.Header.MsgType
	message.Status = MenderShellMessageStatus(m.Header.Status)
	message.Data = m.Body
	return message, nil
}

func (w *MenderShellProtocol) Encapsulate(m interface{}) (*ws.ProtoMsg, error) {
	switch shellMessage := m.(type) {
	case *MenderShellMessage:
		message := &ws.ProtoMsg{
			Header: ws.ProtoHdr{
				Proto:     ws.ProtoTypeShell,
				MsgType:   shellMessage.Type,
				Status:    int(shellMessage.Status),
				SessionID: shellMessage.SessionId,
				Properties: nil,
			},
			Body: shellMessage.Data,
		}
		return message, nil
	default:
		return nil, ws.ErrUnsupportedTypeError
	}
}

func RegisterProtocol() error {
	return ws.RegisterProtocol(ws.ProtoTypeShell, &MenderShellProtocol{version: version})
}
