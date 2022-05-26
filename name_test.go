/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logicalcluster

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName_Split(t *testing.T) {
	tests := []struct {
		cn     Name
		parent Name
		name   string
	}{
		{New(""), New(""), ""},
		{New("foo"), New(""), "foo"},
		{New("foo:bar"), New("foo"), "bar"},
		{New("foo:bar:baz"), New("foo:bar"), "baz"},
		{New("foo::baz"), New("foo:"), "baz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParent, gotName := tt.cn.Split()
			if gotParent != tt.parent {
				t.Errorf("Split() gotParent = %v, want %v", gotParent, tt.parent)
			}
			if gotName != tt.name {
				t.Errorf("Split() gotName = %v, want %v", gotName, tt.name)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	type JWT struct {
		I  int  `json:"i"`
		CN Name `json:"cn"`
	}

	jwt := JWT{
		I:  1,
		CN: New("foo:bar"),
	}

	bs, err := json.Marshal(jwt)
	require.NoError(t, err)
	require.Equal(t, `{"i":1,"cn":"foo:bar"}`, string(bs))

	var jwt2 JWT
	err = json.Unmarshal(bs, &jwt2)
	require.NoError(t, err)
	require.Equal(t, jwt, jwt2)
}

type testObject struct {
	annotations map[string]string
}

func (t *testObject) GetAnnotations() map[string]string {
	return t.annotations
}

func (t *testObject) SetAnnotations(a map[string]string) {
	t.annotations = a
}

func TestSet(t *testing.T) {

	testCases := []struct {
		name          string
		LogicaCluster Name
		Object        *testObject
	}{
		{
			name:          "Set cluster if not set",
			LogicaCluster: Name{value: "root:org:ws-1"},
			Object: &testObject{
				annotations: map[string]string{},
			},
		},
		{
			name:          "overwrite cluster if set",
			LogicaCluster: Name{value: "root:org:ws-1"},
			Object: &testObject{
				annotations: map[string]string{
					LogicalClusterAnnotationKey: "root:org:ws-2",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.LogicaCluster.Set(tt.Object)
			require.Equal(t, tt.LogicaCluster.value, tt.Object.annotations[LogicalClusterAnnotationKey])
		})
	}
}
