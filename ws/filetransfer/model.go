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

package filetransfer

const (
	MessageTypeGet          = "get_file"
	MessageTypePut          = "put_file"
	MessageTypeStat         = "file_info_req"
	MessageTypeStatResponse = "file_info_res"
	MessageTypeChunk        = "file_chunk"
	MessageTypeError        = "error"
)

// The Error struct is passed in the Body of MsgProto in case the message type is ErrorMessage
type Error struct {
	// The error description, as in "Permission denied while opening a file"
	Error string `msgpack:"err" json:"error"`
	// Type of message that raised the error
	MessageType string `msgpack:"msgtype,omitempty" json:"message_type,omitempty"`
	// Message id is passed in the MsgProto Properties, and in case it is available and
	// error occurs it is passed for reference in the Body of the error message
	MessageID string `msgpack:"msgid,omitempty" json:"message_id,omitempty"`
}

// Get file requests the flow of MessageTypeFileChunk messages to be started from the remote end
type GetFile struct {
	// The file path to the file we are requesting
	Path string `msgpack:"path,omitempty" json:"path,omitempty"`
}

// Put file marks the flow of MessageTypeFileChunk messages to be started to the remote end
type PutFile struct {
	// The file path to the file we are requesting
	Path string `msgpack:"path,omitempty" json:"path,omitempty"`
}

// Stat file requests the file stat structure from the remote end
type StatFile struct {
	// The file path to the file we are requesting
	Path string `msgpack:"path,omitempty" json:"path,omitempty"`
}

// Stat file response is the reply to the StatFile from the remote end
type StatFileResponse struct {
	// The file path to the file we are sending status for
	Path string `msgpack:"path,omitempty" json:"path,omitempty"`
	// The file size
	Size int64 `msgpack:"size,omitempty" json:"size,omitempty"`
	// The file type
	Type string `msgpack:"type,omitempty" json:"type,omitempty"`
	// The file owner
	Owner string `msgpack:"owner,omitempty" json:"owner,omitempty"`
	// The file group
	Group string `msgpack:"group,omitempty" json:"group,omitempty"`
}

// File chunk carry the contents of the file
type FileChunk struct {
	// The current offset in the file
	Offset string `msgpack:"offset,omitempty" json:"offset,omitempty"`
	// The chunk size
	Size string `msgpack:"chunksize,omitempty" json:"chunk_size,omitempty"`
	// Array of bytes of file contents
	Data []byte `msgpack:"data,omitempty" json:"data,omitempty"`
}
