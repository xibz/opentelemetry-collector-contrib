// Copyright 2019, OpenTelemetry Authors
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

package translator

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.5.0"
	"go.uber.org/zap"
)

func TestClientSpanWithURLAttribute(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPURL] = "https://api.example.com/users/junit"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://api.example.com/users/junit"))
}

func TestClientSpanWithSchemeHostTargetAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "https"
	attributes[conventions.AttributeHTTPHost] = "api.example.com"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	attributes["user.id"] = "junit"
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://api.example.com/users/junit"))
}

func TestClientSpanWithPeerAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "http"
	attributes[conventions.AttributeNetPeerName] = "kb234.example.com"
	attributes[conventions.AttributeNetPeerPort] = 8080
	attributes[conventions.AttributeNetPeerIP] = "10.8.17.36"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)

	assert.Equal(t, "10.8.17.36", *httpData.Request.ClientIP)

	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "http://kb234.example.com:8080/users/junit"))
}

func TestClientSpanWithHttpPeerAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPClientIP] = "1.2.3.4"
	attributes[conventions.AttributeNetPeerIP] = "10.8.17.36"
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)

	assert.Equal(t, "1.2.3.4", *httpData.Request.ClientIP)
}

func TestClientSpanWithPeerIp4Attributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "http"
	attributes[conventions.AttributeNetPeerIP] = "10.8.17.36"
	attributes[conventions.AttributeNetPeerPort] = "8080"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)
	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "http://10.8.17.36:8080/users/junit"))
}

func TestClientSpanWithPeerIp6Attributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "https"
	attributes[conventions.AttributeNetPeerIP] = "2001:db8:85a3::8a2e:370:7334"
	attributes[conventions.AttributeNetPeerPort] = "443"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	span := constructHTTPClientSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)
	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://2001:db8:85a3::8a2e:370:7334/users/junit"))
}

func TestServerSpanWithURLAttribute(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPURL] = "https://api.example.com/users/junit"
	attributes[conventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[conventions.AttributeHTTPUserAgent] = "PostmanRuntime/7.21.0"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPServerSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://api.example.com/users/junit"))
}

func TestServerSpanWithSchemeHostTargetAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "https"
	attributes[conventions.AttributeHTTPHost] = "api.example.com"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPServerSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://api.example.com/users/junit"))
}

func TestServerSpanWithSchemeServernamePortTargetAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "https"
	attributes[conventions.AttributeHTTPServerName] = "api.example.com"
	attributes[conventions.AttributeNetHostPort] = 443
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPServerSpan(attributes)

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "https://api.example.com/users/junit"))
}

func TestServerSpanWithSchemeNamePortTargetAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "http"
	attributes[conventions.AttributeHostName] = "kb234.example.com"
	attributes[conventions.AttributeNetHostPort] = 8080
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPServerSpan(attributes)
	timeEvents := constructTimedEventsWithReceivedMessageEvent(span.EndTimestamp())
	timeEvents.CopyTo(span.Events())

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.NotNil(t, httpData)
	assert.NotNil(t, filtered)
	w := testWriters.borrow()
	if err := w.Encode(httpData); err != nil {
		assert.Fail(t, "invalid json")
	}
	jsonStr := w.String()
	testWriters.release(w)
	assert.True(t, strings.Contains(jsonStr, "http://kb234.example.com:8080/users/junit"))
}

func TestSpanWithNotEnoughHTTPRequestURLAttributes(t *testing.T) {
	attributes := make(map[string]interface{})
	attributes[conventions.AttributeHTTPMethod] = "GET"
	attributes[conventions.AttributeHTTPScheme] = "http"
	attributes[conventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[conventions.AttributeHTTPUserAgent] = "PostmanRuntime/7.21.0"
	attributes[conventions.AttributeHTTPTarget] = "/users/junit"
	attributes[conventions.AttributeNetHostPort] = 443
	attributes[conventions.AttributeNetPeerPort] = 8080
	attributes[conventions.AttributeHTTPStatusCode] = 200
	span := constructHTTPServerSpan(attributes)
	timeEvents := constructTimedEventsWithReceivedMessageEvent(span.EndTimestamp())
	timeEvents.CopyTo(span.Events())

	filtered, httpData := makeHTTP(zap.L(), span)

	assert.Nil(t, httpData.Request.URL)
	assert.Equal(t, "192.168.15.32", *httpData.Request.ClientIP)
	assert.Equal(t, "GET", *httpData.Request.Method)
	assert.Equal(t, "PostmanRuntime/7.21.0", *httpData.Request.UserAgent)
	contentLength := *httpData.Response.ContentLength.(*int64)
	assert.Equal(t, int64(12452), contentLength)
	assert.Equal(t, int64(200), *httpData.Response.Status)
	assert.NotNil(t, filtered)
}

func constructHTTPClientSpan(attributes map[string]interface{}) pdata.Span {
	endTime := time.Now().Round(time.Second)
	startTime := endTime.Add(-90 * time.Second)
	spanAttributes := constructSpanAttributes(attributes)

	span := pdata.NewSpan()
	span.SetTraceID(newTraceID())
	span.SetSpanID(newSegmentID())
	span.SetParentSpanID(newSegmentID())
	span.SetName("/users/junit")
	span.SetKind(pdata.SpanKindClient)
	span.SetStartTimestamp(pdata.NewTimestampFromTime(startTime))
	span.SetEndTimestamp(pdata.NewTimestampFromTime(endTime))

	status := pdata.NewSpanStatus()
	status.SetCode(0)
	status.SetMessage("OK")
	status.CopyTo(span.Status())

	spanAttributes.CopyTo(span.Attributes())
	return span
}

func constructHTTPServerSpan(attributes map[string]interface{}) pdata.Span {
	endTime := time.Now().Round(time.Second)
	startTime := endTime.Add(-90 * time.Second)
	spanAttributes := constructSpanAttributes(attributes)

	span := pdata.NewSpan()
	span.SetTraceID(newTraceID())
	span.SetSpanID(newSegmentID())
	span.SetParentSpanID(newSegmentID())
	span.SetName("/users/junit")
	span.SetKind(pdata.SpanKindServer)
	span.SetStartTimestamp(pdata.NewTimestampFromTime(startTime))
	span.SetEndTimestamp(pdata.NewTimestampFromTime(endTime))

	status := pdata.NewSpanStatus()
	status.SetCode(0)
	status.SetMessage("OK")
	status.CopyTo(span.Status())

	spanAttributes.CopyTo(span.Attributes())
	return span
}
