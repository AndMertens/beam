package dataflow

import (
	"encoding/json"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam/graph"
	"google.golang.org/api/googleapi"
)

// newMsg creates a json-encoded RawMessage. Panics if encoding fails.
func newMsg(msg interface{}) googleapi.RawMessage {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return googleapi.RawMessage(data)
}

// pipelineOptions models Job/Environment/SdkPipelineOptions
type pipelineOptions struct {
	DisplayData []*displayData `json:"display_data,omitempty"`
	// Options interface{} `json:"options,omitempty"`
}

// NOTE(herohde) 2/9/2017: most of the v1b3 messages are weakly-typed json
// blobs. We manually add them here for convenient and safer use.

// userAgent models Job/Environment/UserAgent. Example value:
//    "userAgent": {
//        "name": "Apache Beam SDK for Python",
//        "version": "0.6.0.dev"
//    },
type userAgent struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// version models Job/Environment/Version. Example value:
//    "version": {
//       "job_type": "PYTHON_BATCH",
//       "major": "5"
//    },
// See: dataflow_job_settings.gcl for range.
type version struct {
	JobType string `json:"job_type,omitempty"`
	Major   string `json:"major,omitempty"`
}

// properties models Step/Properties. Note that the valid subset of fields
// depend on the step kind.
type properties struct {
	UserName    string        `json:"user_name,omitempty"`
	DisplayData []displayData `json:"display_data,omitempty"`

	// Element []string  `json:"element,omitempty"`
	// UserFn string  `json:"user_fn,omitempty"`

	CustomSourceInputStep *customSourceInputStep      `json:"custom_source_step_input,omitempty"`
	ParallelInput         *outputReference            `json:"parallel_input,omitempty"`
	NonParallelInputs     map[string]*outputReference `json:"non_parallel_inputs,omitempty"`
	Format                string                      `json:"format,omitempty"`
	SerializedFn          string                      `json:"serialized_fn,omitempty"`
	OutputInfo            []output                    `json:"output_info,omitempty"`
}

type output struct {
	UserName   string          `json:"user_name,omitempty"`
	OutputName string          `json:"output_name,omitempty"`
	Encoding   *graph.CoderRef `json:"encoding,omitempty"`
}

type integer struct {
	Type  string `json:"@type,omitempty"` // "http://schema.org/Integer"
	Value int    `json:"value,omitempty"`
}

func newInteger(value int) *integer {
	return &integer{
		Type:  "http://schema.org/Integer",
		Value: value,
	}
}

type customSourceInputStep struct {
	Spec     customSourceInputStepSpec      `json:"spec"`
	Metadata *customSourceInputStepMetadata `json:"metadata,omitempty"`
}

type customSourceInputStepSpec struct {
	Type             string `json:"@type,omitempty"` // "CustomSourcesType"
	SerializedSource string `json:"serialized_source,omitempty"`
}

type customSourceInputStepMetadata struct {
	EstimatedSizeBytes *integer `json:"estimated_size_bytes,omitempty"`
}

func newCustomSourceInputStep(serializedSource string) *customSourceInputStep {
	return &customSourceInputStep{
		Spec: customSourceInputStepSpec{
			Type:             "CustomSourcesType",
			SerializedSource: serializedSource,
		},
		Metadata: &customSourceInputStepMetadata{
			EstimatedSizeBytes: newInteger(5 << 20), // 5 MB
		},
	}
}

type outputReference struct {
	Type       string `json:"@type,omitempty"` // "OutputReference"
	StepName   string `json:"step_name,omitempty"`
	OutputName string `json:"output_name,omitempty"`
}

func newOutputReference(step, output string) *outputReference {
	return &outputReference{
		Type:       "OutputReference",
		StepName:   step,
		OutputName: output,
	}
}

type displayData struct {
	Key        string      `json:"key,omitempty"`
	Label      string      `json:"label,omitempty"`
	Namespace  string      `json:"namespace,omitempty"`
	ShortValue string      `json:"shortValue,omitempty"`
	Type       string      `json:"type,omitempty"`
	Value      interface{} `json:"value,omitempty"`
}

func findDisplayDataType(value interface{}) (string, interface{}) {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "INTEGER", value
	case bool:
		return "BOOLEAN", value
	case string:
		return "STRING", value
	default:
		return "STRING", fmt.Sprintf("%v", value)
	}
}

func newDisplayData(key, label, namespace string, value interface{}) *displayData {
	t, v := findDisplayDataType(value)

	return &displayData{
		Key:       key,
		Label:     label,
		Namespace: namespace,
		Type:      t,
		Value:     v,
	}
}
