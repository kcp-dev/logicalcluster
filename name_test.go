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

import "testing"

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"", false},
		{"*", false},

		{"elephant", true},
		{"elephant:foo", false},
		{"elephant:foo:bar", false},

		{"system", true},
		{"system:foo", false},
		{"system-foo", true},
		{"system:foo:bar", false},
		{"system-foo-bar", true},
		{"elephant:0a", false},
		{"elephant-0a", true},
		{"elephant:0bar", false},
		{"elephant-0bar", true},

		{"elephant:b1234567890123456789012345678912", false},
		{"elephant:test-8827a131-f796-4473-8904-a0fa527696eb:b1234567890123456789012345678912", false},
		{"elephant:test-too-long-org-0020-4473-0030-a0fa-0040-5276-0050-sdg2-0060:b1234567890123456789012345678912", false},
		{"elephant-test-8827a131-f796-4473-8904-a0fa527696eb-b1234567890123456789012345678912", false},

		{"elephant:", false},
		{":elephant", false},
		{"-elephant", false},
		{"elephant-", false},
		{"elephant::foo", false},
		{"elephant:föö:bär", false},
		{"elephant:bar_bar", false},
		{"elephant/bar", false},
		{"elephant:bar-", false},
		{"elephant:-bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Name(tt.name).IsValid(); got != tt.valid {
				t.Errorf("isValid(%q) = %v, want %v", tt.name, got, tt.valid)
			}
		})
	}
}
