//  Copyright 2020 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package service

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/compute-image-tools/daisy"
	"github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pb"
	"github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pbtesting"
	"github.com/stretchr/testify/assert"
)

func Test_NewOutputInfoReaderFromWorkflow_ReturnsNilWhenWorkflowNil(t *testing.T) {
	assert.Nil(t, NewOutputInfoReaderFromWorkflow(nil))
}

func Test_WorkflowOutputInfoReader_MapsSingleFieldsCorrectly(t *testing.T) {
	wf := daisy.Workflow{}
	wf.AddSerialConsoleOutputValue(importFileFormat, "vmdk")
	wf.AddSerialConsoleOutputValue(inflationType, "api")
	wf.AddSerialConsoleOutputValue(shadowDiskMatchResult, "shadow-disk-match-result")
	wf.AddSerialConsoleOutputValue(isUEFICompatibleImage, "false")
	wf.AddSerialConsoleOutputValue(isUEFIDetected, "true")
	outputInfoReader := NewOutputInfoReaderFromWorkflow(&wf)
	expected := &pb.OutputInfo{
		ImportFileFormat:      "vmdk",
		InflationType:         "api",
		ShadowDiskMatchResult: "shadow-disk-match-result",
		IsUefiCompatibleImage: false,
		IsUefiDetected:        true,
	}

	pbtesting.AssertEqual(t, expected, outputInfoReader.ReadOutputInfo())
}

func Test_WorkflowOutputInfoReader_MapsArrayFieldsCorrectly(t *testing.T) {
	wf := daisy.Workflow{}
	wf.AddSerialConsoleOutputValue(targetSizeGb, "100,0,5")
	wf.AddSerialConsoleOutputValue(sourceSizeGb, "20,-10,8")
	wf.AddSerialConsoleOutputValue(inflationTime, "20")
	wf.AddSerialConsoleOutputValue(shadowInflationTime, "8")
	outputInfoReader := NewOutputInfoReaderFromWorkflow(&wf)
	expected := &pb.OutputInfo{
		TargetsSizeGb:         []int64{100, 0, 5},
		SourcesSizeGb:         []int64{20, -10, 8},
		InflationTimeMs:       []int64{20},
		ShadowInflationTimeMs: []int64{8},
	}

	pbtesting.AssertEqual(t, expected, outputInfoReader.ReadOutputInfo())
}

func TestWorkflowToOutputInfoReader_ReadSerialPortLogs(t *testing.T) {
	wf := daisy.Workflow{
		Logger: daisyLogger{serialLogs: []string{
			"log-a", "log-b",
		}},
	}
	outputInfoReader := NewOutputInfoReaderFromWorkflow(&wf)

	assert.Equal(t, []string{"log-a", "log-b"}, outputInfoReader.ReadOutputInfo().SerialOutputs)
}

func TestWorkflowToOutputInfoReader_ReadSerialPortLogs_SupportsMissingDaisyLogger(t *testing.T) {
	wf := daisy.Workflow{}
	outputInfoReader := NewOutputInfoReaderFromWorkflow(&wf)

	assert.Empty(t, outputInfoReader.ReadOutputInfo().SerialOutputs)
}

type daisyLogger struct {
	serialLogs []string
}

func (d daisyLogger) WriteLogEntry(e *daisy.LogEntry)                                          {}
func (d daisyLogger) WriteSerialPortLogs(w *daisy.Workflow, instance string, buf bytes.Buffer) {}
func (d daisyLogger) Flush()                                                                   {}
func (d daisyLogger) ReadSerialPortLogs() []string {
	return d.serialLogs
}
