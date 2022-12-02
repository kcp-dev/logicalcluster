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
)

func TestPath_Split(t *testing.T) {
	tests := []struct {
		cn     Path
		parent Path
		name   string
	}{
		{NewPath(""), NewPath(""), ""},
		{NewPath("foo"), NewPath(""), "foo"},
		{NewPath("foo:bar"), NewPath("foo"), "bar"},
		{NewPath("foo:bar:baz"), NewPath("foo:bar"), "baz"},
		{NewPath("foo::baz"), NewPath("foo:"), "baz"},
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

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"", false},
		{"*", true},

		{"elephant", true},
		{"elephant:foo", true},
		{"elephant:foo:bar", true},

		{"system", true},
		{"system:foo", true},
		{"system:foo:bar", true},
		{"elephant:0a", true},
		{"elephant:0bar", true},

		// the plugin does not decide about segment length, the server does
		{"elephant:b1234567890123456789012345678912", true},
		{"elephant:test-8827a131-f796-4473-8904-a0fa527696eb:b1234567890123456789012345678912", true},
		{"elephant:test-too-long-org-0020-4473-0030-a0fa-0040-5276-0050-sdg2-0060:b1234567890123456789012345678912", true},

		{"elephant:", false},
		{":elephant", false},
		{"elephant::foo", false},
		{"elephant:föö:bär", false},
		{"elephant:bar_bar", false},
		{"elephant/bar", false},
		{"elephant:bar-", false},
		{"elephant:-bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPath(tt.name).IsValid(); got != tt.valid {
				t.Errorf("isValid(%q) = %v, want %v", tt.name, got, tt.valid)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	type container struct {
		ClusterPath Path `json:"clusterPath"`
	}

	initial := container{
		ClusterPath: NewPath("foo:bar"),
	}

	raw, err := json.Marshal(initial)
	if err != nil {
		t.Fatal(err)
	}
	if actual, expected := string(raw), `{"clusterPath":"foo:bar"}`; actual != expected {
		t.Fatalf("incorrect marshalled bytes, expected %s, got %s", expected, actual)
	}

	var final container
	if err := json.Unmarshal(raw, &final); err != nil {
		t.Fatal(err)
	}
	if actual, expected := initial.ClusterPath, final.ClusterPath; actual != expected {
		t.Fatalf("incorrect unmarshalled name, expected %s, got %s", expected, actual)
	}
}
