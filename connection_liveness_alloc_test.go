//go:build (linux || darwin || windows) && !app_engine

// Copyright 2014-2026 Aerospike, Inc.
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

package aerospike

import (
	"runtime"
	"testing"
)

// Temporary: reports the probe's allocations for review. Remove with the liveness workflow.
func TestLivenessProbeAllocations(t *testing.T) {
	server := newLivenessServer(t)
	conn, _ := server.connect(t)

	// On Windows the "nothing to read" error of the peek is boxed by the syscall package.
	allowed := 0.0
	if runtime.GOOS == "windows" {
		allowed = 1
	}

	allocations := testing.AllocsPerRun(1000, func() {
		if conn.closedByPeer() {
			t.Fatal("an open idle connection was reported closed")
		}
	})
	if allocations > allowed {
		t.Fatalf("probe allocated %v objects per call, want at most %v", allocations, allowed)
	}
	t.Logf("probe allocated %v objects per call", allocations)
}
