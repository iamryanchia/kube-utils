// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonextracter

import (
	"encoding/json"
	"testing"
)

func TestFieldPath(t *testing.T) {
	type args struct {
		obj              map[string]interface{}
		allowMissingKeys bool
		fields           []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"kind", args{obj: podData, allowMissingKeys: true, fields: []string{"kind"}}, `{"kind":"Pod"}`, false},
		{"lables", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels"}}, `{"metadata":{"labels":{"app":"pause","name":"pause"}}}`, false},
		{"label name", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "name"}}, `{"metadata":{"labels":{"name":"pause"}}}`, false},
		{"containers", args{obj: podData, allowMissingKeys: true, fields: []string{"spec", "containers"}}, `{"spec":{"containers":[{"image":"registry.k8s.io/pause:3.8","imagePullPolicy":"IfNotPresent","name":"pause1","resources":{"limits":{"cpu":"100m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}},{"image":"registry.k8s.io/pause:3.8","imagePullPolicy":"IfNotPresent","name":"pause2","resources":{"limits":{"cpu":"10m","memory":"64Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}}]}}`, false},
		{"test wrong type", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "name", "xx"}}, "null", true},
		{"test not allow miss key", args{obj: podData, allowMissingKeys: false, fields: []string{"metadata", "labels", "xx"}}, "null", true},
		{"test allow miss key", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "xx"}}, `{"metadata":{"labels":{}}}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NestedFieldNoCopy(tt.args.obj, tt.args.allowMissingKeys, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NestedFieldNoCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			data, _ := json.Marshal(got)
			if string(data) != tt.want {
				t.Errorf("NestedFieldNoCopy() = %v, want %v", string(data), tt.want)
			}
		})
	}
}

func BenchmarkFieldPath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NestedFieldNoCopy(podData, false, "kind")
	}
}

func BenchmarkRecurFieldPath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RecurNestedFieldNoCopy(podData, false, "kind")
	}
}
