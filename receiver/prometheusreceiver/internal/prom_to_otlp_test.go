// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"testing"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/model/pdata"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus"
)

// Parity test to ensure that createNodeAndResource produces identical results to createNodeAndResourcePdata.
func TestCreateNodeAndResourceEquivalence(t *testing.T) {
	job, instance, scheme := "converter", "ocmetrics", "http"
	ocNode, ocResource := createNodeAndResource(job, instance, scheme)
	mdFromOC := opencensus.OCToMetrics(ocNode, ocResource,
		// We need to pass in a dummy set of metrics
		// just to populate and allow for full conversion.
		[]*metricspb.Metric{
			{
				MetricDescriptor: &metricspb.MetricDescriptor{
					Name:        "m1",
					Description: "d1",
					Unit:        "By",
				},
			},
		},
	)

	fromOCResource := mdFromOC.ResourceMetrics().At(0).Resource().Attributes().Sort()
	byDirectOTLPResource := CreateNodeAndResourcePdata(job, instance, scheme).Attributes().Sort()

	require.Equal(t, byDirectOTLPResource, fromOCResource)
}

type jobInstanceDefinition struct {
	job, instance, host, scheme, port string
}

func makeResourceWithJobInstanceScheme(def *jobInstanceDefinition, hasHost bool) *pdata.Resource {
	resource := pdata.NewResource()
	attrs := resource.Attributes()
	// Using hardcoded values to assert on outward expectations so that
	// when variables change, these tests will fail and we'll have reports.
	attrs.UpsertString("service.name", def.job)
	if hasHost {
		attrs.UpsertString("host.name", def.host)
	}
	attrs.UpsertString("job", def.job)
	attrs.UpsertString("instance", def.instance)
	attrs.UpsertString("port", def.port)
	attrs.UpsertString("scheme", def.scheme)
	return &resource
}

func TestCreateNodeAndResourcePromToOTLP(t *testing.T) {
	tests := []struct {
		name, job string
		instance  string
		scheme    string
		want      *pdata.Resource
	}{
		{
			name: "all attributes proper",
			job:  "job", instance: "hostname:8888", scheme: "http",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "hostname:8888", "hostname", "http", "8888",
			}, true),
		},
		{
			name: "missing port",
			job:  "job", instance: "myinstance", scheme: "https",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "myinstance", "myinstance", "https", "",
			}, true),
		},
		{
			name: "blank scheme",
			job:  "job", instance: "myinstance:443", scheme: "",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "myinstance:443", "myinstance", "", "443",
			}, true),
		},
		{
			name: "blank instance, blank scheme",
			job:  "job", instance: "", scheme: "",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "", "", "", "",
			}, true),
		},
		{
			name: "blank instance, non-blank scheme",
			job:  "job", instance: "", scheme: "http",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "", "", "http", "",
			}, true),
		},
		{
			name: "0.0.0.0 address",
			job:  "job", instance: "0.0.0.0:8888", scheme: "http",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "0.0.0.0:8888", "", "http", "8888",
			}, false),
		},
		{
			name: "localhost",
			job:  "job", instance: "localhost:8888", scheme: "http",
			want: makeResourceWithJobInstanceScheme(&jobInstanceDefinition{
				"job", "localhost:8888", "", "http", "8888",
			}, false),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := CreateNodeAndResourcePdata(tt.job, tt.instance, tt.scheme)
			require.Equal(t, got, tt.want)
		})
	}
}
