// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"testing"

	pipelinev1beta1test "github.com/tektoncd/pipeline/test"
	rtesting "knative.dev/pkg/reconciler/testing"
)

func SeedV1beta1TestData(t *testing.T, d pipelinev1beta1test.Data) (pipelinev1beta1test.Clients, pipelinev1beta1test.Informers) {
	ctx, _ := rtesting.SetupFakeContext(t)
	return pipelinev1beta1test.SeedTestData(t, ctx, d)
}
