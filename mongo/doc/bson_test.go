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

package doc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDocumentFromStruct(t *testing.T) {
	testCases := []struct {
		Name string

		Input          interface{}
		AppendElements []bson.E
		Expected       bson.D
	}{
		{
			Name: "Simple success",

			Input: struct {
				Field1 string
				Field2 int
			}{
				Field1: "foo",
				Field2: 321,
			},
			Expected: bson.D{
				{Key: "field1", Value: "foo"},
				{Key: "field2", Value: 321},
			},
		},
		{
			Name: "Bson tags and appends",

			Input: struct {
				Field1 string `bson:"foo"`
				Field2 string `bson:"bar,omitempty"`
			}{
				Field1: "baz",
			},
			AppendElements: []bson.E{
				{Key: "a1", Value: 123},
				{Key: "a2", Value: "foobarbaz"},
			},
			Expected: bson.D{
				{Key: "foo", Value: "baz"},
				{Key: "a1", Value: 123},
				{Key: "a2", Value: "foobarbaz"},
			},
		},
		{
			Name: "Not a struct",

			Input:    "Panic attack!",
			Expected: nil,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			doc := DocumentFromStruct(tc.Input, tc.AppendElements...)
			assert.Equal(t, tc.Expected, doc)
		})
	}

}
