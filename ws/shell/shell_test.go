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
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMenderShellRegisterProtocol(t *testing.T) {
	err := RegisterProtocol()
	assert.NoError(t, err)
}

func TestMenderShellProtocol_Encapsulate(t *testing.T) {
	testCases := []struct {
		name          string
		protoType     ws.ProtoType
		messageStatus MenderShellMessageStatus
		messageType   string
		sessionId     string
		data          []byte
		err           error
	}{
		{
			name:          "shell_command_successful",
			protoType:     ws.ProtoTypeShell,
			messageStatus: NormalMessage,
			messageType:   MessageTypeShellCommand,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("ls -al\n"),
		},
		{
			name:          "shell_command_error",
			protoType:     ws.ProtoTypeShell,
			messageStatus: ErrorMessage,
			messageType:   MessageTypeShellCommand,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("command not found"),
		},
		{
			name:          "spawn_command_successful",
			protoType:     ws.ProtoTypeShell,
			messageStatus: NormalMessage,
			messageType:   MessageTypeSpawnShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("shell spawned"),
		},
		{
			name:          "shell_command_error",
			protoType:     ws.ProtoTypeShell,
			messageStatus: ErrorMessage,
			messageType:   MessageTypeSpawnShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("cant start another shell"),
		},
		{
			name:          "stop_command_successful",
			protoType:     ws.ProtoTypeShell,
			messageStatus: NormalMessage,
			messageType:   MessageTypeStopShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("shell stopped"),
		},
		{
			name:          "shell_command_error",
			protoType:     ws.ProtoTypeShell,
			messageStatus: ErrorMessage,
			messageType:   MessageTypeStopShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("failed to stop shell"),
		},
		{
			name:          "unsupported_message_type",
			protoType:     ws.ProtoTypeShell,
			messageStatus: ErrorMessage,
			messageType:   MessageTypeStopShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("failed to stop shell"),
			err:           ws.ErrUnsupportedTypeError,
		},
		{
			name:          "unregistered_protocol",
			protoType:     ws.ProtoInvalid,
			messageStatus: ErrorMessage,
			messageType:   MessageTypeStopShell,
			sessionId:     uuid.NewV4().String(),
			data:          []byte("failed to stop shell"),
			err:           ws.ErrUnregisteredProtocol,
		},
	}

	err := RegisterProtocol()
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var messageIn interface{}
			messageIn = &MenderShellMessage{
				Type:      tc.messageType,
				SessionId: tc.sessionId,
				Status:    tc.messageStatus,
				Data:      tc.data,
			}

			if tc.err != nil {
				if tc.err == ws.ErrUnsupportedTypeError {
					messageIn = "this type is not registered"
				}
			}
			protoMsg, err := ws.Encapsulate(tc.protoType, messageIn)
			if tc.err == nil {
				assert.NoError(t, err)
				assert.NotNil(t, protoMsg)

				messageOut, err := ws.DeEncapsulate(protoMsg)
				assert.NoError(t, err)
				assert.Equal(t, messageIn, messageOut)
			} else {
				assert.Error(t, tc.err, err)
			}
		})
	}
}
