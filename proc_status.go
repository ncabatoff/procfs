// Copyright 2018 The Prometheus Authors
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

package procfs

import (
	"bufio"
	"fmt"
	"os"
)

// ProcStatus provides status information about the process,
// read from /proc/[pid]/status.
type (
	ProcStatus struct {
		TID                      int
		TracerPid                int
		UIDReal                  int
		UIDEffective             int
		UIDSavedSet              int
		UIDFileSystem            int
		GIDReal                  int
		GIDEffective             int
		GIDSavedSet              int
		GIDFileSystem            int
		FDSize                   int
		VmPeakKB                 int
		VmSizeKB                 int
		VmLckKB                  int
		VmHWMKB                  int
		VmRSSKB                  int
		VmDataKB                 int
		VmStkKB                  int
		VmExeKB                  int
		VmLibKB                  int
		VmPTEKB                  int
		VmSwapKB                 int
		VoluntaryCtxtSwitches    int
		NonvoluntaryCtxtSwitches int
	}

	procStatusScanner struct {
		format string
		args   []interface{}
	}

	procStatusBuilder struct {
		ps       ProcStatus
		scanners []procStatusScanner
	}
)

func newProcStatusBuilder() *procStatusBuilder {
	var b procStatusBuilder
	b.scanners = []procStatusScanner{
		{"Pid: %d", []interface{}{&b.ps.TID}},
		{"TracerPid: %d", []interface{}{&b.ps.TracerPid}},
		{"Uid: %d %d %d %d", []interface{}{
			&b.ps.UIDReal,
			&b.ps.UIDEffective,
			&b.ps.UIDSavedSet,
			&b.ps.UIDFileSystem,
		}},
		{"Gid: %d %d %d %d", []interface{}{
			&b.ps.GIDReal,
			&b.ps.GIDEffective,
			&b.ps.GIDSavedSet,
			&b.ps.GIDFileSystem,
		}},
		{"FDSize: %d", []interface{}{&b.ps.FDSize}},
		{"VmPeak: %d kB", []interface{}{&b.ps.VmPeakKB}},
		{"VmSize: %d kB", []interface{}{&b.ps.VmSizeKB}},
		{"VmLck:  %d kB", []interface{}{&b.ps.VmLckKB}},
		{"VmHWM:  %d kB", []interface{}{&b.ps.VmHWMKB}},
		{"VmRSS:  %d kB", []interface{}{&b.ps.VmRSSKB}},
		{"VmData: %d kB", []interface{}{&b.ps.VmDataKB}},
		{"VmStk:  %d kB", []interface{}{&b.ps.VmStkKB}},
		{"VmExe:  %d kB", []interface{}{&b.ps.VmExeKB}},
		{"VmLib:  %d kB", []interface{}{&b.ps.VmLibKB}},
		{"VmPTE:  %d kB", []interface{}{&b.ps.VmPTEKB}},
		{"VmSwap: %d kB", []interface{}{&b.ps.VmSwapKB}},
		{"voluntary_ctxt_switches:    %d", []interface{}{&b.ps.VoluntaryCtxtSwitches}},
		{"nonvoluntary_ctxt_switches: %d", []interface{}{&b.ps.NonvoluntaryCtxtSwitches}},
	}
	return &b
}

func (b *procStatusBuilder) readStatus(r *bufio.Reader) (ProcStatus, error) {
	for _, s := range b.scanners {
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return ProcStatus{}, err
			}

			_, err = fmt.Sscanf(line, s.format, s.args...)
			if err == nil {
				break
			}
		}
	}
	return b.ps, nil
}

// NewStatus returns the current status information of the process.
func (p Proc) NewStatus() (ProcStatus, error) {
	f, err := os.Open(p.path("status"))
	if err != nil {
		return ProcStatus{}, err
	}
	defer f.Close()

	return newProcStatusBuilder().readStatus(bufio.NewReader(f))
}
