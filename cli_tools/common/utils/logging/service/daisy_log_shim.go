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
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/logging"
	"github.com/GoogleCloudPlatform/compute-image-tools/daisy"
	"github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pb"
)

// NewOutputInfoReaderFromWorkflow provides a OutputInfoReader from a daisy workflow.
func NewOutputInfoReaderFromWorkflow(wf *daisy.Workflow) logging.OutputInfoReader {
	if wf == nil {
		return nil
	}
	return workflowOutputInfoReader{wf: wf}
}

type workflowOutputInfoReader struct {
	wf *daisy.Workflow
}

func (w workflowOutputInfoReader) ReadOutputInfo() *pb.OutputInfo {
	return &pb.OutputInfo{
		TargetsSizeGb:         w.getValueAsInt64Slice(targetSizeGb),
		SourcesSizeGb:         w.getValueAsInt64Slice(sourceSizeGb),
		ImportFileFormat:      w.getValue(importFileFormat),
		InflationType:         w.getValue(inflationType),
		InflationTimeMs:       w.getValueAsInt64Slice(inflationTime),
		ShadowInflationTimeMs: w.getValueAsInt64Slice(shadowInflationTime),
		ShadowDiskMatchResult: w.getValue(shadowDiskMatchResult),
		IsUefiCompatibleImage: w.getValueAsBool(isUEFICompatibleImage),
		IsUefiDetected:        w.getValueAsBool(isUEFIDetected),
		SerialOutputs:         w.readSerialPortLogs(),
	}
}

func (w workflowOutputInfoReader) getValue(key string) string {
	return w.wf.GetSerialConsoleOutputValue(key)
}

func (w workflowOutputInfoReader) getValueAsBool(key string) bool {
	v, err := strconv.ParseBool(w.wf.GetSerialConsoleOutputValue(key))
	if err != nil {
		return false
	}
	return v
}

func (w workflowOutputInfoReader) getValueAsInt64Slice(key string) []int64 {
	return getInt64Values(w.wf.GetSerialConsoleOutputValue(key))
}

func (w workflowOutputInfoReader) readSerialPortLogs() []string {
	if w.wf.Logger != nil {
		logs := w.wf.Logger.ReadSerialPortLogs()
		view := make([]string, len(logs))
		copy(view, logs)
		return view
	}
	return nil
}

func getInt64Values(s string) []int64 {
	strs := strings.Split(s, ",")
	var r []int64
	for _, str := range strs {
		i, err := strconv.ParseInt(str, 0, 64)
		if err == nil {
			r = append(r, i)
		}
	}
	return r
}
