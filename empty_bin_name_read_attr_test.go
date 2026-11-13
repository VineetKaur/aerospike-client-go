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

import "testing"

// The server treats INFO1_GET_ALL as "return every bin whole" and discards the
// operation list. An empty bin name only means "all bins" for a plain read;
// "" is also a valid bin name, so a CDT, bit, HLL or expression read on it must
// not set the flag. These check the attributes without a server.

func emptyBinReadOps() map[string]*Operation {
	return map[string]*Operation{
		"MapGetByKeyOp": MapGetByKeyOp("", "k", MapReturnType.VALUE),
		"MapSizeOp":     MapSizeOp(""),
		"ListSizeOp":    ListSizeOp(""),
		"ExpReadOp":     ExpReadOp("", ExpIntVal(1), ExpReadFlagDefault),
		"BitCountOp":    BitCountOp("", 0, 8),
	}
}

func TestEmptyBinNameCDTReadDoesNotSetGetAll(t *testing.T) {
	for name, op := range emptyBinReadOps() {
		args, err := newOperateArgs(nil, NewWritePolicy(0, 0), nil, []*Operation{op})
		if err != nil {
			t.Fatalf("%s: newOperateArgs: %v", name, err)
		}
		if args.readAttr&_INFO1_GET_ALL != 0 {
			t.Errorf("%s on bin \"\": operate sets INFO1_GET_ALL", name)
		}
		if args.readAttr&_INFO1_READ == 0 {
			t.Errorf("%s on bin \"\": operate does not set INFO1_READ", name)
		}

		if ba := newBatchAttrOps(NewBatchPolicy(), NewBatchWritePolicy(), []*Operation{op}); ba.readAttr&_INFO1_GET_ALL != 0 {
			t.Errorf("%s on bin \"\": batch read sets INFO1_GET_ALL", name)
		}

		ops := []*Operation{PutOp(NewBin("b", 1)), op}
		if ba := newBatchAttrOps(NewBatchPolicy(), NewBatchWritePolicy(), ops); ba.readAttr&_INFO1_GET_ALL != 0 {
			t.Errorf("%s on bin \"\": batch write sets INFO1_GET_ALL", name)
		}

		ba := &batchAttr{}
		ba.adjustWrite(ops)
		if ba.readAttr&_INFO1_GET_ALL != 0 {
			t.Errorf("%s on bin \"\": batchAttr.adjustWrite sets INFO1_GET_ALL", name)
		}
	}
}

func TestGetOpStillSetsGetAll(t *testing.T) {
	args, err := newOperateArgs(nil, NewWritePolicy(0, 0), nil, []*Operation{GetOp()})
	if err != nil {
		t.Fatalf("newOperateArgs: %v", err)
	}
	if args.readAttr&_INFO1_GET_ALL == 0 {
		t.Error("GetOp: operate does not set INFO1_GET_ALL")
	}

	if ba := newBatchAttrOps(NewBatchPolicy(), NewBatchWritePolicy(), []*Operation{GetOp()}); ba.readAttr&_INFO1_GET_ALL == 0 {
		t.Error("GetOp: batch read does not set INFO1_GET_ALL")
	}

	ops := []*Operation{PutOp(NewBin("b", 1)), GetOp()}
	if ba := newBatchAttrOps(NewBatchPolicy(), NewBatchWritePolicy(), ops); ba.readAttr&_INFO1_GET_ALL == 0 {
		t.Error("GetOp: batch write does not set INFO1_GET_ALL")
	}

	ba := &batchAttr{}
	ba.adjustWrite(ops)
	if ba.readAttr&_INFO1_GET_ALL == 0 {
		t.Error("GetOp: batchAttr.adjustWrite does not set INFO1_GET_ALL")
	}
}
