//go:build !as_performance && !app_engine

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

package aerospike_test

import (
	gg "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"

	as "github.com/aerospike/aerospike-client-go/v8"
)

// "" is a valid bin name (namespaces migrated from single-bin keep their data
// in it). A CDT read on that bin must return the operation's result, not every
// bin of the record: an empty bin name only means "read all bins" for a plain
// read.
var _ = gg.Describe("CDT read on the empty bin name", func() {

	var ns = *namespace
	var set = randString(50)
	var key *as.Key

	data := map[interface{}]interface{}{
		"a": map[interface{}]interface{}{"x": 1},
		"b": map[interface{}]interface{}{"x": 2},
		"c": map[interface{}]interface{}{"x": 3},
	}
	want := map[interface{}]interface{}{"x": 2}

	gg.BeforeEach(func() {
		var err error
		key, err = as.NewKey(ns, set, randString(50))
		gm.Expect(err).ToNot(gm.HaveOccurred())

		err = client.PutBins(nil, key, as.NewBin("", data), as.NewBin("other", 7))
		gm.Expect(err).ToNot(gm.HaveOccurred())
	})

	gg.AfterEach(func() {
		_, err := client.Delete(nil, key)
		gm.Expect(err).ToNot(gm.HaveOccurred())
	})

	gg.It("Operate MapGetByKeyOp returns only the requested entry", func() {
		rec, err := client.Operate(nil, key, as.MapGetByKeyOp("", "b", as.MapReturnType.VALUE))
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(rec.Bins).To(gm.Equal(as.BinMap{"": want}))
	})

	gg.It("Operate MapSizeOp returns the size, not the map", func() {
		rec, err := client.Operate(nil, key, as.MapSizeOp(""))
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(rec.Bins).To(gm.Equal(as.BinMap{"": 3}))
	})

	gg.It("Operate with a write and a CDT read on bin \"\" returns the read result", func() {
		rec, err := client.Operate(nil, key,
			as.PutOp(as.NewBin("other", 8)),
			as.MapGetByKeyOp("", "b", as.MapReturnType.VALUE),
		)
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(rec.Bins[""]).To(gm.Equal(want))
	})

	gg.It("BatchOperate read ops return only the requested entry", func() {
		br := as.NewBatchReadOps(nil, key, as.MapGetByKeyOp("", "b", as.MapReturnType.VALUE))
		err := client.BatchOperate(nil, []as.BatchRecordIfc{br})
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(br.Err).ToNot(gm.HaveOccurred())
		gm.Expect(br.Record.Bins).To(gm.Equal(as.BinMap{"": want}))
	})

	gg.It("BatchOperate write with a CDT read on bin \"\" returns the read result", func() {
		bw := as.NewBatchWrite(nil, key,
			as.PutOp(as.NewBin("other", 9)),
			as.MapGetByKeyOp("", "b", as.MapReturnType.VALUE),
		)
		err := client.BatchOperate(nil, []as.BatchRecordIfc{bw})
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(bw.Err).ToNot(gm.HaveOccurred())
		gm.Expect(bw.Record.Bins[""]).To(gm.Equal(want))
	})

	gg.It("GetOp still reads all bins", func() {
		rec, err := client.Operate(nil, key, as.GetOp())
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(rec.Bins).To(gm.Equal(as.BinMap{"": data, "other": 7}))
	})
})
