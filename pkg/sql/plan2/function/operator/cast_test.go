// Copyright 2022 Matrix Origin
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

package operator

import (
	roaring "github.com/RoaringBitmap/roaring/roaring64"
	"github.com/matrixorigin/matrixone/pkg/container/nulls"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCastSameType(t *testing.T) {
	makeTempVectors := func(src interface{}, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeVector(src, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors(int8(-23), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{-23},
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors(int16(-23), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{-23},
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors(int32(-23), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{-23},
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors(int64(-23), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{-23},
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors(uint8(23), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{23},
			wantScalar: true,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors(uint16(23), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{23},
			wantScalar: true,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors(uint32(23), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{23},
			wantScalar: true,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors(uint64(23), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{23},
			wantScalar: true,
		},
		{
			name:       "Test09",
			vecs:       makeTempVectors(float32(23.5), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{23.5},
			wantScalar: true,
		},
		{
			name:       "Test10",
			vecs:       makeTempVectors(float64(23.5), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{23.5},
			wantScalar: true,
		},
		{
			name:       "Test11",
			vecs:       makeTempVectors(int8(-23), types.T_int8, false),
			proc:       procs,
			wantValues: []int8{-23},
			wantScalar: false,
		},
		{
			name:       "Test12",
			vecs:       makeTempVectors(int16(-23), types.T_int16, false),
			proc:       procs,
			wantValues: []int16{-23},
			wantScalar: false,
		},
		{
			name:       "Test13",
			vecs:       makeTempVectors(int32(-23), types.T_int32, false),
			proc:       procs,
			wantValues: []int32{-23},
			wantScalar: false,
		},
		{
			name:       "Test14",
			vecs:       makeTempVectors(int64(-23), types.T_int64, false),
			proc:       procs,
			wantValues: []int64{-23},
			wantScalar: false,
		},
		{
			name:       "Test15",
			vecs:       makeTempVectors(uint8(23), types.T_uint8, false),
			proc:       procs,
			wantValues: []uint8{23},
			wantScalar: false,
		},
		{
			name:       "Test16",
			vecs:       makeTempVectors(uint16(23), types.T_uint16, false),
			proc:       procs,
			wantValues: []uint16{23},
			wantScalar: false,
		},
		{
			name:       "Test17",
			vecs:       makeTempVectors(uint32(23), types.T_uint32, false),
			proc:       procs,
			wantValues: []uint32{23},
			wantScalar: false,
		},
		{
			name:       "Test18",
			vecs:       makeTempVectors(uint64(23), types.T_uint64, false),
			proc:       procs,
			wantValues: []uint64{23},
			wantScalar: false,
		},
		{
			name:       "Test19",
			vecs:       makeTempVectors(float32(23.5), types.T_float32, false),
			proc:       procs,
			wantValues: []float32{23.5},
			wantScalar: false,
		},
		{
			name:       "Test20",
			vecs:       makeTempVectors(float64(23.5), types.T_float64, false),
			proc:       procs,
			wantValues: []float64{23.5},
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func TestCastSameType2(t *testing.T) {
	makeTempVectors := func(src interface{}, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeVector(src, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	//types.Date | types.Datetime | types.Timestamp
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors(types.Date(729848), types.T_date, true),
			proc:       procs,
			wantValues: []types.Date{729848},
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors(types.Datetime(66122056321728512), types.T_datetime, true),
			proc:       procs,
			wantValues: []types.Datetime{66122056321728512},
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors(types.Timestamp(66122026122739712), types.T_timestamp, true),
			proc:       procs,
			wantValues: []types.Timestamp{66122026122739712},
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors(types.Date(729848), types.T_date, false),
			proc:       procs,
			wantValues: []types.Date{729848},
			wantScalar: false,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors(types.Datetime(66122056321728512), types.T_datetime, false),
			proc:       procs,
			wantValues: []types.Datetime{66122056321728512},
			wantScalar: false,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors(types.Timestamp(66122026122739712), types.T_timestamp, false),
			proc:       procs,
			wantValues: []types.Timestamp{66122026122739712},
			wantScalar: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func TestCastLeftToRight(t *testing.T) {
	// int8 -> (int16/int32/int64/uint8/uint16/uint32/uint64/float32/float64)
	// int16 -> (int8/int32/int64/uint8/uint16/uint32/uint64/float32/float64)
	// int32 -> (int8/int16/int64/uint8/uint16/uint32/uint64/float32/float64)
	// int64 -> (int8/int16/int32/uint8/uint16/uint32/uint64/float32/float64)
	// uint8 -> (int8/int16/int32/int64/uint16/uint32/uint64/float32/float64)
	// uint16 -> (int8/int16/int32/int64/uint8/uint32/uint64/float32/float64)
	// uint32 -> (int8/int16/int32/int64/uint8/uint16/uint64/float32/float64)
	// uint64 -> (int8/int16/int32/int64/uint8/uint16/uint32/float32/float64)
	// float32 -> (int8/int16/int32/int64/uint8/uint16/uint32/uint64/float64)
	// float64 -> (int8/int16/int32/int64/uint8/uint16/uint32/uint64/float32)

	makeTempVectors := func(src interface{}, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeVector(src, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors(int8(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors(int8(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors(int8(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors(int8(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors(int8(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors(int8(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors(int8(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors(int8(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test09",
			vecs:       makeTempVectors(int8(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test10",
			vecs:       makeTempVectors(int8(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test11",
			vecs:       makeTempVectors(int16(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test12",
			vecs:       makeTempVectors(int16(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test13",
			vecs:       makeTempVectors(int16(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test14",
			vecs:       makeTempVectors(int16(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test15",
			vecs:       makeTempVectors(int16(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test16",
			vecs:       makeTempVectors(int16(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test17",
			vecs:       makeTempVectors(int16(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test18",
			vecs:       makeTempVectors(int16(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test19",
			vecs:       makeTempVectors(int16(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test20",
			vecs:       makeTempVectors(int16(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test21",
			vecs:       makeTempVectors(int32(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test22",
			vecs:       makeTempVectors(int32(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test23",
			vecs:       makeTempVectors(int32(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test24",
			vecs:       makeTempVectors(int32(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test25",
			vecs:       makeTempVectors(int32(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test26",
			vecs:       makeTempVectors(int32(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test27",
			vecs:       makeTempVectors(int32(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test28",
			vecs:       makeTempVectors(int32(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test29",
			vecs:       makeTempVectors(int32(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test30",
			vecs:       makeTempVectors(int32(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test31",
			vecs:       makeTempVectors(int64(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test32",
			vecs:       makeTempVectors(int64(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test33",
			vecs:       makeTempVectors(int64(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test34",
			vecs:       makeTempVectors(int64(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test35",
			vecs:       makeTempVectors(int64(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test36",
			vecs:       makeTempVectors(int64(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test37",
			vecs:       makeTempVectors(int64(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test38",
			vecs:       makeTempVectors(int64(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test39",
			vecs:       makeTempVectors(int64(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test40",
			vecs:       makeTempVectors(int64(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test41",
			vecs:       makeTempVectors(uint8(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test42",
			vecs:       makeTempVectors(uint8(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test43",
			vecs:       makeTempVectors(uint8(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test44",
			vecs:       makeTempVectors(uint8(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test45",
			vecs:       makeTempVectors(uint8(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test46",
			vecs:       makeTempVectors(uint8(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test47",
			vecs:       makeTempVectors(uint8(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test48",
			vecs:       makeTempVectors(uint8(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test49",
			vecs:       makeTempVectors(uint8(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test50",
			vecs:       makeTempVectors(uint8(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test51",
			vecs:       makeTempVectors(uint16(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test52",
			vecs:       makeTempVectors(uint16(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test53",
			vecs:       makeTempVectors(uint16(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test54",
			vecs:       makeTempVectors(uint16(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test55",
			vecs:       makeTempVectors(uint16(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test56",
			vecs:       makeTempVectors(uint16(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test57",
			vecs:       makeTempVectors(uint16(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test58",
			vecs:       makeTempVectors(uint16(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test59",
			vecs:       makeTempVectors(uint16(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test60",
			vecs:       makeTempVectors(uint16(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test61",
			vecs:       makeTempVectors(uint32(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test62",
			vecs:       makeTempVectors(uint32(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test63",
			vecs:       makeTempVectors(uint32(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test64",
			vecs:       makeTempVectors(uint32(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test65",
			vecs:       makeTempVectors(uint32(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test66",
			vecs:       makeTempVectors(uint32(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test67",
			vecs:       makeTempVectors(uint32(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test68",
			vecs:       makeTempVectors(uint32(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test69",
			vecs:       makeTempVectors(uint32(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test70",
			vecs:       makeTempVectors(uint32(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test71",
			vecs:       makeTempVectors(uint64(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test72",
			vecs:       makeTempVectors(uint64(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test73",
			vecs:       makeTempVectors(uint64(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test74",
			vecs:       makeTempVectors(uint64(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test75",
			vecs:       makeTempVectors(uint64(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test76",
			vecs:       makeTempVectors(uint64(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test77",
			vecs:       makeTempVectors(uint64(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test78",
			vecs:       makeTempVectors(uint64(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test79",
			vecs:       makeTempVectors(uint64(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test80",
			vecs:       makeTempVectors(uint64(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test81",
			vecs:       makeTempVectors(float32(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test82",
			vecs:       makeTempVectors(float32(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test83",
			vecs:       makeTempVectors(float32(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test84",
			vecs:       makeTempVectors(float32(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test85",
			vecs:       makeTempVectors(float32(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test86",
			vecs:       makeTempVectors(float32(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test87",
			vecs:       makeTempVectors(float32(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test88",
			vecs:       makeTempVectors(float32(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test89",
			vecs:       makeTempVectors(float32(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test90",
			vecs:       makeTempVectors(float32(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
		{
			name:       "Test91",
			vecs:       makeTempVectors(float64(125), types.T_int8, true),
			proc:       procs,
			wantValues: []int8{125},
			wantScalar: true,
		},
		{
			name:       "Test92",
			vecs:       makeTempVectors(float64(125), types.T_int16, true),
			proc:       procs,
			wantValues: []int16{125},
			wantScalar: true,
		},
		{
			name:       "Test93",
			vecs:       makeTempVectors(float64(125), types.T_int32, true),
			proc:       procs,
			wantValues: []int32{125},
			wantScalar: true,
		},
		{
			name:       "Test94",
			vecs:       makeTempVectors(float64(125), types.T_int64, true),
			proc:       procs,
			wantValues: []int64{125},
			wantScalar: true,
		},
		{
			name:       "Test95",
			vecs:       makeTempVectors(float64(125), types.T_uint8, true),
			proc:       procs,
			wantValues: []uint8{125},
			wantScalar: true,
		},
		{
			name:       "Test96",
			vecs:       makeTempVectors(float64(125), types.T_uint16, true),
			proc:       procs,
			wantValues: []uint16{125},
			wantScalar: true,
		},
		{
			name:       "Test97",
			vecs:       makeTempVectors(float64(125), types.T_uint32, true),
			proc:       procs,
			wantValues: []uint32{125},
			wantScalar: true,
		},
		{
			name:       "Test98",
			vecs:       makeTempVectors(float64(125), types.T_uint64, true),
			proc:       procs,
			wantValues: []uint64{125},
			wantScalar: true,
		},
		{
			name:       "Test99",
			vecs:       makeTempVectors(float64(125), types.T_float32, true),
			proc:       procs,
			wantValues: []float32{125},
			wantScalar: true,
		},
		{
			name:       "Test100",
			vecs:       makeTempVectors(float64(125), types.T_float64, true),
			proc:       procs,
			wantValues: []float64{125},
			wantScalar: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func TestCastSpecials1Int(t *testing.T) {
	// (char / varhcar) -> (int8 / int16 / int32/ int64 / uint8 / uint16 / uint32 / uint64)

	makeTempVectors := func(src string, srcType types.T, srcIsConst bool, destType types.T) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeStringVector(src, srcType, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_int8),
			proc:       procs,
			wantValues: []int8{15},
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_int16),
			proc:       procs,
			wantValues: []int16{15},
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_int32),
			proc:       procs,
			wantValues: []int32{15},
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_int64),
			proc:       procs,
			wantValues: []int64{15},
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_uint8),
			proc:       procs,
			wantValues: []uint8{15},
			wantScalar: true,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_uint16),
			proc:       procs,
			wantValues: []uint16{15},
			wantScalar: true,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_uint32),
			proc:       procs,
			wantValues: []uint32{15},
			wantScalar: true,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors("15", types.T_varchar, true, types.T_uint64),
			proc:       procs,
			wantValues: []uint64{15},
			wantScalar: true,
		},
		{
			name:       "Test09",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_int8),
			proc:       procs,
			wantValues: []int8{15},
			wantScalar: true,
		},
		{
			name:       "Test10",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_int16),
			proc:       procs,
			wantValues: []int16{15},
			wantScalar: true,
		},
		{
			name:       "Test11",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_int32),
			proc:       procs,
			wantValues: []int32{15},
			wantScalar: true,
		},
		{
			name:       "Test12",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_int64),
			proc:       procs,
			wantValues: []int64{15},
			wantScalar: true,
		},
		{
			name:       "Test13",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_uint8),
			proc:       procs,
			wantValues: []uint8{15},
			wantScalar: true,
		},
		{
			name:       "Test14",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_uint16),
			proc:       procs,
			wantValues: []uint16{15},
			wantScalar: true,
		},
		{
			name:       "Test15",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_uint32),
			proc:       procs,
			wantValues: []uint32{15},
			wantScalar: true,
		},
		{
			name:       "Test16",
			vecs:       makeTempVectors("15", types.T_char, true, types.T_uint64),
			proc:       procs,
			wantValues: []uint64{15},
			wantScalar: true,
		},
		{
			name:       "Test17",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_int8),
			proc:       procs,
			wantValues: []int8{15},
			wantScalar: false,
		},
		{
			name:       "Test18",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_int16),
			proc:       procs,
			wantValues: []int16{15},
			wantScalar: false,
		},
		{
			name:       "Test19",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_int32),
			proc:       procs,
			wantValues: []int32{15},
			wantScalar: false,
		},
		{
			name:       "Test20",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_int64),
			proc:       procs,
			wantValues: []int64{15},
			wantScalar: false,
		},
		{
			name:       "Test21",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_uint8),
			proc:       procs,
			wantValues: []uint8{15},
			wantScalar: false,
		},
		{
			name:       "Test22",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_uint16),
			proc:       procs,
			wantValues: []uint16{15},
			wantScalar: false,
		},
		{
			name:       "Test23",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_uint32),
			proc:       procs,
			wantValues: []uint32{15},
			wantScalar: false,
		},
		{
			name:       "Test24",
			vecs:       makeTempVectors("15", types.T_varchar, false, types.T_uint64),
			proc:       procs,
			wantValues: []uint64{15},
			wantScalar: false,
		},
		{
			name:       "Test25",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_int8),
			proc:       procs,
			wantValues: []int8{15},
			wantScalar: false,
		},
		{
			name:       "Test26",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_int16),
			proc:       procs,
			wantValues: []int16{15},
			wantScalar: false,
		},
		{
			name:       "Test27",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_int32),
			proc:       procs,
			wantValues: []int32{15},
			wantScalar: false,
		},
		{
			name:       "Test28",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_int64),
			proc:       procs,
			wantValues: []int64{15},
			wantScalar: false,
		},
		{
			name:       "Test29",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_uint8),
			proc:       procs,
			wantValues: []uint8{15},
			wantScalar: false,
		},
		{
			name:       "Test30",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_uint16),
			proc:       procs,
			wantValues: []uint16{15},
			wantScalar: false,
		},
		{
			name:       "Test31",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_uint32),
			proc:       procs,
			wantValues: []uint32{15},
			wantScalar: false,
		},
		{
			name:       "Test32",
			vecs:       makeTempVectors("15", types.T_char, false, types.T_uint64),
			proc:       procs,
			wantValues: []uint64{15},
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}

}

func TestCastSpecials1Float(t *testing.T) {
	// (char / varhcar) -> (float32 / float64)
	makeTempVectors := func(src string, srcType types.T, srcIsConst bool, destType types.T) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeStringVector(src, srcType, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors("15.23", types.T_varchar, true, types.T_float32),
			proc:       procs,
			wantValues: []float32{15.23},
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors("15.23", types.T_varchar, true, types.T_float64),
			proc:       procs,
			wantValues: []float64{15.23},
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors("15.23", types.T_char, true, types.T_float32),
			proc:       procs,
			wantValues: []float32{15.23},
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors("15.23", types.T_char, true, types.T_float64),
			proc:       procs,
			wantValues: []float64{15.23},
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors("15.23", types.T_varchar, false, types.T_float32),
			proc:       procs,
			wantValues: []float32{15.23},
			wantScalar: false,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors("15.23", types.T_varchar, false, types.T_float64),
			proc:       procs,
			wantValues: []float64{15.23},
			wantScalar: false,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors("15.23", types.T_char, false, types.T_float32),
			proc:       procs,
			wantValues: []float32{15.23},
			wantScalar: false,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors("15.23", types.T_char, false, types.T_float64),
			proc:       procs,
			wantValues: []float64{15.23},
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}

}

func TestCastSpecials2Float(t *testing.T) {
	//(float32/float64) -> (char / varhcar)
	makeTempVectors := func(src interface{}, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeVector(src, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantBytes  []byte
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors(float32(23.65), types.T_char, true),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors(float64(23.65), types.T_char, true),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors(float32(23.65), types.T_varchar, true),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors(float64(23.65), types.T_varchar, true),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors(float32(23.65), types.T_char, false),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: false,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors(float64(23.65), types.T_char, false),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: false,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors(float32(23.65), types.T_varchar, false),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: false,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors(float64(23.65), types.T_varchar, false),
			proc:       procs,
			wantBytes:  []byte("23.65"),
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			col := castRes.Col.(*types.Bytes)
			offset := col.Offsets[0]
			length := col.Lengths[0]
			resBytes := col.Data[offset:length]
			require.Equal(t, c.wantBytes, resBytes)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func TestCastSpecials3(t *testing.T) {
	// char -> char
	// char -> varhcar
	// varhcar -> char
	// varhcar -> varhcar
	makeTempVectors := func(src string, srcType types.T, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeStringVector(src, srcType, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantBytes  []byte
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors("abcsedn", types.T_char, types.T_char, true),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors("abcsedn", types.T_char, types.T_varchar, true),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors("abcsedn", types.T_varchar, types.T_char, true),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors("abcsedn", types.T_varchar, types.T_varchar, true),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors("abcsedn", types.T_char, types.T_char, false),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: false,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors("abcsedn", types.T_char, types.T_varchar, false),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: false,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors("abcsedn", types.T_varchar, types.T_char, false),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: false,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors("abcsedn", types.T_varchar, types.T_varchar, false),
			proc:       procs,
			wantBytes:  []byte("abcsedn"),
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			col := castRes.Col.(*types.Bytes)
			offset := col.Offsets[0]
			length := col.Lengths[0]
			resBytes := col.Data[offset:length]
			require.Equal(t, c.wantBytes, resBytes)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}

}

func TestCastSpecial4(t *testing.T) {
	//(int8/int16/int32/int64) to decimal128
	// (uint8/uint16/uint32/uint64) to decimal128

	makeTempVectors := func(src interface{}, destType types.T, srcIsConst bool) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeVector(src, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}
	resType := types.Type{Oid: types.T_decimal128, Size: 16, Width: 38, Scale: 0}
	decimal128 := types.Decimal128{Lo: 123, Hi: 0}
	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantType   types.Type
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors(int8(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors(int16(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors(int32(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors(int64(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test05",
			vecs:       makeTempVectors(uint8(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test06",
			vecs:       makeTempVectors(uint16(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test07",
			vecs:       makeTempVectors(uint32(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test08",
			vecs:       makeTempVectors(uint64(123), types.T_decimal128, true),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: true,
		},
		{
			name:       "Test09",
			vecs:       makeTempVectors(int8(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test10",
			vecs:       makeTempVectors(int16(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test11",
			vecs:       makeTempVectors(int32(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test12",
			vecs:       makeTempVectors(int64(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test13",
			vecs:       makeTempVectors(uint8(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test14",
			vecs:       makeTempVectors(uint16(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test15",
			vecs:       makeTempVectors(uint32(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
		{
			name:       "Test16",
			vecs:       makeTempVectors(uint64(123), types.T_decimal128, false),
			proc:       procs,
			wantValues: []types.Decimal128{decimal128},
			wantType:   resType,
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			decimalres, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, decimalres.Col)
			require.Equal(t, c.wantType, decimalres.Typ)
			require.Equal(t, c.wantScalar, decimalres.IsScalar())
		})
	}

}

func TestCastVarcharAsDate(t *testing.T) {
	//Cast converts varchar to date type
	//Cast converts varchar to datetime type
	//Cast converts varchar to timestamp type
	makeTempVectors := func(src string, srcType types.T, srcIsConst bool, destType types.T) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeStringVector(src, srcType, srcIsConst)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantValues interface{}
		wantType   types.T
		wantScalar bool
	}{
		{
			name:       "Test01",
			vecs:       makeTempVectors("1988-03-14", types.T_varchar, true, types.T_date),
			proc:       procs,
			wantValues: []types.Date{725809},
			wantType:   types.T_date,
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVectors("1988-03-14", types.T_varchar, false, types.T_date),
			proc:       procs,
			wantValues: []types.Date{725809},
			wantType:   types.T_date,
			wantScalar: false,
		},
		{
			name:       "Test03",
			vecs:       makeTempVectors("2020-06-14 16:24:15.23", types.T_varchar, true, types.T_datetime),
			proc:       procs,
			wantValues: []types.Datetime{66823387773895280},
			wantType:   types.T_datetime,
			wantScalar: true,
		},
		{
			name:       "Test04",
			vecs:       makeTempVectors("2020-06-14 16:24:15.23", types.T_varchar, false, types.T_datetime),
			proc:       procs,
			wantValues: []types.Datetime{66823387773895280},
			wantType:   types.T_datetime,
			wantScalar: false,
		},
		//{
		//	name:       "Test05",
		//	vecs:       makeTempVectors("2020-06-14 16:24:15.23", types.T_varchar, true, types.T_timestamp),
		//	proc:       procs,
		//	wantValues: []types.Timestamp{66823357574906480},
		//	wantType:   types.T_timestamp,
		//	wantScalar: true,
		//},
		//{
		//	name:       "Test06",
		//	vecs:       makeTempVectors("2020-06-14 16:24:15.23", types.T_varchar, false, types.T_timestamp),
		//	proc:       procs,
		//	wantValues: []types.Timestamp{66823357574906480},
		//	wantType:   types.T_timestamp,
		//	wantScalar: false,
		//},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantValues, castRes.Col)
			require.Equal(t, c.wantType, castRes.Typ.Oid)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}

}

func TestCastDecimal64AsDecimal128(t *testing.T) {
	//Cast converts decimal64 to decimal128
	makeTempVector := func(left types.Decimal64, leftType types.Type, leftScalar bool, destType types.Type) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = &vector.Vector{
			Col:     []types.Decimal64{left},
			Nsp:     &nulls.Nulls{},
			Typ:     leftType,
			IsConst: leftScalar,
			Length:  1,
		}
		vectors[1] = &vector.Vector{
			Col:     nil,
			Nsp:     &nulls.Nulls{},
			Typ:     destType,
			IsConst: false,
			Length:  0,
		}
		return vectors
	}
	// decimal(10,5)
	leftType := types.Type{Oid: types.T_decimal64, Size: 8, Width: 10, Scale: 5}
	//decimal(20, 5)
	destType := types.Type{Oid: types.T_decimal128, Size: 16, Width: 20, Scale: 5}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantBytes  interface{}
		wantType   types.T
		wantScalar bool
	}{
		{
			name:       "TEST01", //cast(333.33300 as decimal(20, 5))
			vecs:       makeTempVector(33333300, leftType, true, destType),
			proc:       procs,
			wantBytes:  []types.Decimal128{{Lo: 33333300, Hi: 0}},
			wantType:   types.T_decimal128,
			wantScalar: true,
		},
		{
			name:       "TEST01", //cast(333.33300 as decimal(20, 5))
			vecs:       makeTempVector(33333300, leftType, false, destType),
			proc:       procs,
			wantBytes:  []types.Decimal128{{Lo: 33333300, Hi: 0}},
			wantType:   types.T_decimal128,
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantBytes, castRes.Col)
			require.Equal(t, c.wantType, castRes.Typ.Oid)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}

}

func TestCastDecimal64AsDecimal64(t *testing.T) {
	//Cast converts decimal64 to decimal64
	makeTempVector := func(left types.Decimal64, leftType types.Type, leftScalar bool, destType types.Type) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = &vector.Vector{
			Col:     []types.Decimal64{left},
			Nsp:     &nulls.Nulls{},
			Typ:     leftType,
			IsConst: leftScalar,
			Length:  1,
		}

		vectors[1] = &vector.Vector{
			Col:     nil,
			Nsp:     &nulls.Nulls{},
			Typ:     destType,
			IsConst: false,
			Length:  0,
		}
		return vectors
	}
	// decimal(10,5)
	leftType := types.Type{Oid: types.T_decimal64, Size: 8, Width: 10, Scale: 5}
	//decimal(10, 4)
	destType := types.Type{Oid: types.T_decimal64, Size: 8, Width: 10, Scale: 4}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantBytes  interface{}
		wantType   types.Type
		wantScalar bool
	}{
		{
			name:       "TEST01", //cast(333.33300 as decimal(10, 4))
			vecs:       makeTempVector(33333300, leftType, true, destType),
			proc:       procs,
			wantBytes:  []types.Decimal64{33333300},
			wantType:   destType,
			wantScalar: true,
		},
		{
			name:       "TEST02", //cast(333.33300 as decimal(10, 4))
			vecs:       makeTempVector(33333300, leftType, false, destType),
			proc:       procs,
			wantBytes:  []types.Decimal64{33333300},
			wantType:   destType,
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantBytes, castRes.Col)
			require.Equal(t, c.wantType.Oid, castRes.Typ.Oid)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func TestCastDecimal128AsDecimal128(t *testing.T) {
	// Cast converts decimal128 to decimal128
	makeTempVector := func(left types.Decimal128, leftType types.Type, leftScalar bool, destType types.Type) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = &vector.Vector{
			Col:     []types.Decimal128{left},
			Nsp:     &nulls.Nulls{},
			Typ:     leftType,
			IsConst: leftScalar,
			Length:  1,
		}

		vectors[1] = &vector.Vector{
			Col:     nil,
			Nsp:     &nulls.Nulls{},
			Typ:     destType,
			IsConst: false,
			Length:  0,
		}
		return vectors
	}

	leftType := types.Type{Oid: types.T_decimal128, Size: 16, Width: 20, Scale: 5}
	destType := types.Type{Oid: types.T_decimal128, Size: 16, Width: 20, Scale: 5}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		wantBytes  interface{}
		wantType   types.Type
		wantScalar bool
	}{
		{
			name:       "TEST01", //cast(333.33300 as decimal(20, 5))
			vecs:       makeTempVector(types.Decimal128{Lo: 33333300, Hi: 0}, leftType, true, destType),
			proc:       procs,
			wantBytes:  []types.Decimal128{{Lo: 33333300, Hi: 0}},
			wantType:   destType,
			wantScalar: true,
		},
		{
			name:       "Test02",
			vecs:       makeTempVector(types.Decimal128{Lo: 33333300, Hi: 0}, leftType, false, destType),
			proc:       procs,
			wantBytes:  []types.Decimal128{{Lo: 33333300, Hi: 0}},
			wantType:   destType,
			wantScalar: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.wantBytes, castRes.Col)
			require.Equal(t, c.wantType.Oid, castRes.Typ.Oid)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

//func TestCastTimeStampAsDatetime(t *testing.T) {
//	//Cast converts timestamp to datetime
//	makeTempVectors := func(src types.Timestamp, srcType types.T, srcIsConst bool, destType types.T) []*vector.Vector {
//		vectors := make([]*vector.Vector, 2)
//		vectors[0] = &vector.Vector{
//			Col:     []types.Timestamp{src},
//			Nsp:     &nulls.Nulls{},
//			Typ:     types.Type{Oid: srcType},
//			IsConst: srcIsConst,
//			Length:  1,
//		}
//		vectors[1] = makeTypeVector(destType)
//		return vectors
//	}
//
//	procs := makeProcess()
//	cases := []struct {
//		name       string
//		vecs       []*vector.Vector
//		proc       *process.Process
//		wantBytes  interface{}
//		wantType   types.T
//		wantScalar bool
//	}{
//		{
//			name:       "TEST01", //cast(c_timestamp as datetime)  c_timestamp:'1999-04-05 11:01:02'
//			vecs:       makeTempVectors(66122026122739712, types.T_timestamp, true, types.T_datetime),
//			proc:       procs,
//			wantBytes:  []types.Datetime{66122056321728512},
//			wantType:   types.T_datetime,
//			wantScalar: true,
//		},
//		{
//			name:       "TEST02", //cast(333.33300 as decimal(20, 5))
//			vecs:       makeTempVectors(66122026122739712, types.T_timestamp, false, types.T_datetime),
//			proc:       procs,
//			wantBytes:  []types.Datetime{66122056321728512},
//			wantType:   types.T_datetime,
//			wantScalar: false,
//		},
//	}
//
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			castRes, err := Cast(c.vecs, c.proc)
//			if err != nil {
//				t.Fatal(err)
//			}
//			require.Equal(t, c.wantBytes, castRes.Col)
//			require.Equal(t, c.wantType, castRes.Typ.Oid)
//			require.Equal(t, c.wantScalar, castRes.IsScalar())
//		})
//	}
//
//}

func TestCastNullAsAllType(t *testing.T) {
	//Cast null as (int8/int16/int32/int64/uint8/uint16/uint32/uint64/float32/float64/date/datetime/timestamp/decimal64/decimal128/char/varchar)
	makeTempVectors := func(srcType types.T, destType types.T) []*vector.Vector {
		vectors := make([]*vector.Vector, 2)
		vectors[0] = makeScalarNullVector(srcType)
		vectors[1] = makeTypeVector(destType)
		return vectors
	}

	procs := makeProcess()
	cases := []struct {
		name       string
		vecs       []*vector.Vector
		proc       *process.Process
		bitMap     *roaring.Bitmap
		wantScalar bool
	}{
		{
			name:       "Test1",
			vecs:       makeTempVectors(types.T_int8, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test2",
			vecs:       makeTempVectors(types.T_int8, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test3",
			vecs:       makeTempVectors(types.T_int8, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test4",
			vecs:       makeTempVectors(types.T_int8, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test5",
			vecs:       makeTempVectors(types.T_int8, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test6",
			vecs:       makeTempVectors(types.T_int8, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test7",
			vecs:       makeTempVectors(types.T_int8, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test8",
			vecs:       makeTempVectors(types.T_int8, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test9",
			vecs:       makeTempVectors(types.T_int8, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test10",
			vecs:       makeTempVectors(types.T_int8, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test11",
			vecs:       makeTempVectors(types.T_int8, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test12",
			vecs:       makeTempVectors(types.T_int8, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test13",
			vecs:       makeTempVectors(types.T_int8, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test14",
			vecs:       makeTempVectors(types.T_int16, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test15",
			vecs:       makeTempVectors(types.T_int16, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test16",
			vecs:       makeTempVectors(types.T_int16, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test17",
			vecs:       makeTempVectors(types.T_int16, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test18",
			vecs:       makeTempVectors(types.T_int16, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test19",
			vecs:       makeTempVectors(types.T_int16, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test20",
			vecs:       makeTempVectors(types.T_int16, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test21",
			vecs:       makeTempVectors(types.T_int16, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test22",
			vecs:       makeTempVectors(types.T_int16, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test23",
			vecs:       makeTempVectors(types.T_int16, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test24",
			vecs:       makeTempVectors(types.T_int16, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test25",
			vecs:       makeTempVectors(types.T_int16, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test26",
			vecs:       makeTempVectors(types.T_int16, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test27",
			vecs:       makeTempVectors(types.T_int32, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test28",
			vecs:       makeTempVectors(types.T_int32, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test29",
			vecs:       makeTempVectors(types.T_int32, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test30",
			vecs:       makeTempVectors(types.T_int32, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test31",
			vecs:       makeTempVectors(types.T_int32, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test32",
			vecs:       makeTempVectors(types.T_int32, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test33",
			vecs:       makeTempVectors(types.T_int32, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test34",
			vecs:       makeTempVectors(types.T_int32, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test35",
			vecs:       makeTempVectors(types.T_int32, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test36",
			vecs:       makeTempVectors(types.T_int32, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test37",
			vecs:       makeTempVectors(types.T_int32, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test38",
			vecs:       makeTempVectors(types.T_int32, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test39",
			vecs:       makeTempVectors(types.T_int32, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test40",
			vecs:       makeTempVectors(types.T_int64, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test41",
			vecs:       makeTempVectors(types.T_int64, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test42",
			vecs:       makeTempVectors(types.T_int64, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test43",
			vecs:       makeTempVectors(types.T_int64, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test44",
			vecs:       makeTempVectors(types.T_int64, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test45",
			vecs:       makeTempVectors(types.T_int64, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test46",
			vecs:       makeTempVectors(types.T_int64, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test47",
			vecs:       makeTempVectors(types.T_int64, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test48",
			vecs:       makeTempVectors(types.T_int64, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test49",
			vecs:       makeTempVectors(types.T_int64, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test50",
			vecs:       makeTempVectors(types.T_int64, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test51",
			vecs:       makeTempVectors(types.T_int64, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test52",
			vecs:       makeTempVectors(types.T_int64, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test53",
			vecs:       makeTempVectors(types.T_uint8, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test54",
			vecs:       makeTempVectors(types.T_uint8, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test55",
			vecs:       makeTempVectors(types.T_uint8, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test56",
			vecs:       makeTempVectors(types.T_uint8, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test57",
			vecs:       makeTempVectors(types.T_uint8, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test58",
			vecs:       makeTempVectors(types.T_uint8, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test59",
			vecs:       makeTempVectors(types.T_uint8, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test60",
			vecs:       makeTempVectors(types.T_uint8, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test61",
			vecs:       makeTempVectors(types.T_uint8, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test62",
			vecs:       makeTempVectors(types.T_uint8, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test63",
			vecs:       makeTempVectors(types.T_uint8, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test64",
			vecs:       makeTempVectors(types.T_uint8, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test65",
			vecs:       makeTempVectors(types.T_uint8, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test66",
			vecs:       makeTempVectors(types.T_uint16, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test67",
			vecs:       makeTempVectors(types.T_uint16, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test68",
			vecs:       makeTempVectors(types.T_uint16, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test69",
			vecs:       makeTempVectors(types.T_uint16, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test70",
			vecs:       makeTempVectors(types.T_uint16, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test71",
			vecs:       makeTempVectors(types.T_uint16, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test72",
			vecs:       makeTempVectors(types.T_uint16, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test73",
			vecs:       makeTempVectors(types.T_uint16, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test74",
			vecs:       makeTempVectors(types.T_uint16, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test75",
			vecs:       makeTempVectors(types.T_uint16, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test76",
			vecs:       makeTempVectors(types.T_uint16, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test77",
			vecs:       makeTempVectors(types.T_uint16, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test78",
			vecs:       makeTempVectors(types.T_uint16, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test79",
			vecs:       makeTempVectors(types.T_uint32, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test80",
			vecs:       makeTempVectors(types.T_uint32, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test81",
			vecs:       makeTempVectors(types.T_uint32, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test82",
			vecs:       makeTempVectors(types.T_uint32, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test83",
			vecs:       makeTempVectors(types.T_uint32, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test84",
			vecs:       makeTempVectors(types.T_uint32, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test85",
			vecs:       makeTempVectors(types.T_uint32, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test86",
			vecs:       makeTempVectors(types.T_uint32, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test87",
			vecs:       makeTempVectors(types.T_uint32, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test88",
			vecs:       makeTempVectors(types.T_uint32, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test89",
			vecs:       makeTempVectors(types.T_uint32, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test90",
			vecs:       makeTempVectors(types.T_uint32, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test91",
			vecs:       makeTempVectors(types.T_uint32, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test92",
			vecs:       makeTempVectors(types.T_uint64, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test93",
			vecs:       makeTempVectors(types.T_uint64, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test94",
			vecs:       makeTempVectors(types.T_uint64, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test95",
			vecs:       makeTempVectors(types.T_uint64, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test96",
			vecs:       makeTempVectors(types.T_uint64, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test97",
			vecs:       makeTempVectors(types.T_uint64, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test98",
			vecs:       makeTempVectors(types.T_uint64, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test99",
			vecs:       makeTempVectors(types.T_uint64, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test100",
			vecs:       makeTempVectors(types.T_uint64, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test101",
			vecs:       makeTempVectors(types.T_uint64, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test102",
			vecs:       makeTempVectors(types.T_uint64, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test103",
			vecs:       makeTempVectors(types.T_uint64, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test104",
			vecs:       makeTempVectors(types.T_uint64, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test105",
			vecs:       makeTempVectors(types.T_float32, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test106",
			vecs:       makeTempVectors(types.T_float32, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test107",
			vecs:       makeTempVectors(types.T_float32, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test108",
			vecs:       makeTempVectors(types.T_float32, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test109",
			vecs:       makeTempVectors(types.T_float32, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test110",
			vecs:       makeTempVectors(types.T_float32, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test111",
			vecs:       makeTempVectors(types.T_float32, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test112",
			vecs:       makeTempVectors(types.T_float32, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test113",
			vecs:       makeTempVectors(types.T_float32, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test114",
			vecs:       makeTempVectors(types.T_float32, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test115",
			vecs:       makeTempVectors(types.T_float32, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test116",
			vecs:       makeTempVectors(types.T_float32, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test117",
			vecs:       makeTempVectors(types.T_float64, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test118",
			vecs:       makeTempVectors(types.T_float64, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test119",
			vecs:       makeTempVectors(types.T_float64, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test120",
			vecs:       makeTempVectors(types.T_float64, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test121",
			vecs:       makeTempVectors(types.T_float64, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test122",
			vecs:       makeTempVectors(types.T_float64, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test123",
			vecs:       makeTempVectors(types.T_float64, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test124",
			vecs:       makeTempVectors(types.T_float64, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test125",
			vecs:       makeTempVectors(types.T_float64, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test126",
			vecs:       makeTempVectors(types.T_float64, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test127",
			vecs:       makeTempVectors(types.T_float64, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test128",
			vecs:       makeTempVectors(types.T_float64, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0129",
			vecs:       makeTempVectors(types.T_char, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0130",
			vecs:       makeTempVectors(types.T_char, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0131",
			vecs:       makeTempVectors(types.T_char, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0132",
			vecs:       makeTempVectors(types.T_char, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0133",
			vecs:       makeTempVectors(types.T_char, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0134",
			vecs:       makeTempVectors(types.T_char, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0135",
			vecs:       makeTempVectors(types.T_char, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0136",
			vecs:       makeTempVectors(types.T_char, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0137",
			vecs:       makeTempVectors(types.T_char, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0138",
			vecs:       makeTempVectors(types.T_char, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0139",
			vecs:       makeTempVectors(types.T_char, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0140",
			vecs:       makeTempVectors(types.T_char, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0141",
			vecs:       makeTempVectors(types.T_varchar, types.T_int8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0142",
			vecs:       makeTempVectors(types.T_varchar, types.T_int16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0143",
			vecs:       makeTempVectors(types.T_varchar, types.T_int32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0144",
			vecs:       makeTempVectors(types.T_varchar, types.T_int64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0145",
			vecs:       makeTempVectors(types.T_varchar, types.T_uint8),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0146",
			vecs:       makeTempVectors(types.T_varchar, types.T_uint16),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0147",
			vecs:       makeTempVectors(types.T_varchar, types.T_uint32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0148",
			vecs:       makeTempVectors(types.T_varchar, types.T_uint64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0149",
			vecs:       makeTempVectors(types.T_varchar, types.T_float32),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0150",
			vecs:       makeTempVectors(types.T_varchar, types.T_float64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0151",
			vecs:       makeTempVectors(types.T_varchar, types.T_date),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0152",
			vecs:       makeTempVectors(types.T_varchar, types.T_datetime),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0153",
			vecs:       makeTempVectors(types.T_varchar, types.T_timestamp),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0154",
			vecs:       makeTempVectors(types.T_varchar, types.T_char),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0155",
			vecs:       makeTempVectors(types.T_varchar, types.T_varchar),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0156",
			vecs:       makeTempVectors(types.T_date, types.T_date),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0157",
			vecs:       makeTempVectors(types.T_datetime, types.T_datetime),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0158",
			vecs:       makeTempVectors(types.T_datetime, types.T_datetime),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0159",
			vecs:       makeTempVectors(types.T_timestamp, types.T_datetime),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0160",
			vecs:       makeTempVectors(types.T_timestamp, types.T_timestamp),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0161",
			vecs:       makeTempVectors(types.T_decimal64, types.T_decimal64),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0162",
			vecs:       makeTempVectors(types.T_decimal64, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
		{
			name:       "Test0163",
			vecs:       makeTempVectors(types.T_decimal128, types.T_decimal128),
			proc:       procs,
			bitMap:     roaring.BitmapOf(0),
			wantScalar: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			castRes, err := Cast(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, c.bitMap, castRes.Nsp.Np)
			require.Equal(t, c.wantScalar, castRes.IsScalar())
		})
	}
}

func makeTypeVector(t types.T) *vector.Vector {
	return &vector.Vector{
		Col:     nil,
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: t},
		IsConst: false,
		Length:  0,
	}
}

// make vector for type of int8,int16,int32,int64,uint8,uint16,uint32,uint64,date,datetime,timestamp
func makeVector(src interface{}, isSrcConst bool) *vector.Vector {
	var typeOid types.T
	var col interface{}
	switch val := src.(type) {
	case int8:
		typeOid = types.T_int8
		col = []int8{val}
	case int16:
		typeOid = types.T_int16
		col = []int16{val}
	case int32:
		typeOid = types.T_int32
		col = []int32{val}
	case int64:
		typeOid = types.T_int64
		col = []int64{val}
	case uint8:
		typeOid = types.T_uint8
		col = []uint8{val}
	case uint16:
		typeOid = types.T_uint16
		col = []uint16{val}
	case uint32:
		typeOid = types.T_uint32
		col = []uint32{val}
	case uint64:
		typeOid = types.T_uint64
		col = []uint64{val}
	case float32:
		typeOid = types.T_float32
		col = []float32{val}
	case float64:
		typeOid = types.T_float64
		col = []float64{val}
	case types.Date:
		typeOid = types.T_date
		col = []types.Date{val}
	case types.Datetime:
		typeOid = types.T_datetime
		col = []types.Datetime{val}
	case types.Timestamp:
		typeOid = types.T_timestamp
		col = []types.Timestamp{val}
	}

	return &vector.Vector{
		Col:     col,
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: typeOid},
		IsConst: isSrcConst,
		Length:  1,
	}
}

func makeScalarNullVector(srcType types.T) *vector.Vector {
	nullVector := &vector.Vector{
		Col:     nil,
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: srcType},
		IsConst: true,
		Length:  1,
	}
	nulls.Add(nullVector.Nsp, 0)
	return nullVector
}

// make vector for type of char and varchar
func makeStringVector(src string, t types.T, isConst bool) *vector.Vector {
	srcBytes := &types.Bytes{
		Data:    []byte(src),
		Offsets: []uint32{0},
		Lengths: []uint32{uint32(len(src))},
	}
	if t == types.T_char {
		return &vector.Vector{
			Col:     srcBytes,
			Nsp:     &nulls.Nulls{},
			Typ:     types.Type{Oid: types.T_char, Size: 24},
			IsConst: isConst,
			Length:  1,
		}
	} else if t == types.T_varchar {
		return &vector.Vector{
			Col:     srcBytes,
			Nsp:     &nulls.Nulls{},
			Typ:     types.Type{Oid: types.T_varchar, Size: 24},
			IsConst: isConst,
			Length:  1,
		}
	}
	return nil
}
