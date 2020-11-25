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

// ProtoType defines how the ProtoMsg should be interpreted.
type ProtoType uint16

const (
	// ProtoInvalid signifies an invalid (uninitialized) ProtoMsg.
	ProtoInvalid ProtoType = iota
	// ProtoTypeShell is used for communicating remote terminal session data.
	ProtoTypeShell
)

// ProtoHdr provides the info about what the ProtoMsg contains and
// to which protocol the message should be routed.
type ProtoHdr struct {
	// Proto defines which protocol this message belongs
	// to (required). E.g.:
	// * MenderShell for remote terminal
	// * MenderConfiguration for configuration of devices
	// * MenderMonitoring for the stats and monitoring
	Proto ProtoType `msgpack:"proto"`

	// MsgType is an optional content type header describing
	// the protocol specific content type of the message. E.g.: for MenderShell:
	// * CommandData for the stdout / stderr coming from the terminal and keystrokes from the UI
	// * StartSession for the session start
	// * StopSession for stopping of the session
	MsgType string `msgpack:"typ,omitempty"`

	// Status is the indicator of whether the response contains error or other
	// conditions (required).
	// E.g.: for the remote terminal case a StartSession MsgType maye result in
	//       an error due to mas session limit reached.
	// * Normal message carries a data from successful processing
	// * Error the given message returned error
	// * ProtocolError higher level protocol error
	Status int `msgpack:"state"`

	// SessionID is used to identify one ProtoMsg stream for
	// multiplexing multiple ProtoMsg sessions over the same connection.
	SessionID string `msgpack:"sid,omitempty"`

	// Properties provide a map of optional prototype specific
	// properties (such as http headers or other meta-data).
	Properties map[string]interface{} `msgpack:"props,omitempty"`
}

// ProtoMsg is a wrapper to messages communicated on bidirectional interfaces
// such as websockets to wrap data from other application protocols.
type ProtoMsg struct {
	// Header contains a protocol specific header with a single
	// fixed ProtoType ("typ") field and optional hints for decoding
	// the payload.
	Header ProtoHdr `msgpack:"hdr"`
	// Body contains the raw protocol data. The data contained in Body
	// can be arbitrary and must be decoded according to the protocol
	// defined in the header.
	Body []byte `msgpack:"body,omitempty"`
}
