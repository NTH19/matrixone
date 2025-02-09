// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package date_sub

import (
	"github.com/matrixorigin/matrixone/pkg/container/nulls"
	"github.com/matrixorigin/matrixone/pkg/container/types"
)

var (
	DateSub       func([]types.Date, []int64, []int64, []types.Date) []types.Date
	DatetimeSub   func([]types.Datetime, []int64, []int64, []types.Datetime) []types.Datetime
	DateStringSub func(*types.Bytes, []int64, []int64, *nulls.Nulls, *types.Bytes) *types.Bytes
)

func init() {
	DateSub = dateSub
	DatetimeSub = datetimeSub
	DateStringSub = dateStringSub
}

func dateSub(xs []types.Date, ys []int64, zs []int64, rs []types.Date) []types.Date {
	for i, d := range xs {
		rs[i] = d.ToTime().AddInterval(-ys[0], types.IntervalType(zs[0])).ToDate()
	}
	return rs
}

func datetimeSub(xs []types.Datetime, ys []int64, zs []int64, rs []types.Datetime) []types.Datetime {
	for i, d := range xs {
		rs[i] = d.AddInterval(-ys[0], types.IntervalType(zs[0]))
	}
	return rs
}

func dateStringSub(xs *types.Bytes, ys []int64, zs []int64, ns *nulls.Nulls, rs *types.Bytes) *types.Bytes {
	for i := range xs.Lengths {
		str := string(xs.Get(int64(i)))
		if types.UnitIsDayOrLarger(types.IntervalType(zs[0])) {
			d, e := types.ParseDate(str)
			if e == nil {
				rs.AppendOnce([]byte(d.ToTime().AddInterval(-ys[0], types.IntervalType(zs[0])).ToDate().String()))
				continue
			}
		}
		d, e := types.ParseDatetime(str)
		if e != nil {
			// set null
			nulls.Add(ns, uint64(i))
			rs.AppendOnce([]byte(""))
			continue
		}
		rs.AppendOnce([]byte(d.AddInterval(-ys[0], types.IntervalType(zs[0])).String()))
	}
	return rs
}
