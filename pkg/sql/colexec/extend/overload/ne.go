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

package overload

import (
	"matrixone/pkg/container/types"
	"matrixone/pkg/container/vector"
	"matrixone/pkg/encoding"
	"matrixone/pkg/vectorize/ne"
	"matrixone/pkg/vm/process"
	"matrixone/pkg/vm/register"

	roaring "github.com/RoaringBitmap/roaring/roaring64"
)

func init() {
	BinOps[NE] = []*BinOp{
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]int8), rv.Col.([]int8)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Int8NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int8NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Int8NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int8NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int8NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int8NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Int8NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Int8Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]int16), rv.Col.([]int16)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Int16NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int16NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Int16NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int16NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int16NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int16NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Int16NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Int16Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]int32), rv.Col.([]int32)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Int32NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int32NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Int32NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int32NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int32NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int32NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Int32NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Int32Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]int64), rv.Col.([]int64)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Int64NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int64NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Int64NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Int64NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int64NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Int64NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Int64NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Int64Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]uint8), rv.Col.([]uint8)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Uint8NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint8NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Uint8NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint8NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint8NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint8NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Uint8NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Uint8Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]uint16), rv.Col.([]uint16)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Uint16NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint16NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Uint16NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint16NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint16NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint16NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Uint16NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Uint16Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]uint32), rv.Col.([]uint32)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Uint32NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint32NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Uint32NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint32NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint32NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint32NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Uint32NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Uint32Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]uint64), rv.Col.([]uint64)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Uint64NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint64NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Uint64NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Uint64NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint64NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Uint64NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Uint64NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Uint64Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]float32), rv.Col.([]float32)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Float32NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Float32NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Float32NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Float32NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Float32NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Float32NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Float32NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Float32Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.([]float64), rv.Col.([]float64)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.Float64NeNullableScalar(lvs[0], rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Float64NeScalar(lvs[0], rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.Float64NeNullableScalar(rvs[0], lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.Float64NeScalar(rvs[0], lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Float64NeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.Float64NeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.Float64NeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.Float64Ne(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_char,
			RightType:  types.T_char,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.(*types.Bytes), rv.Col.(*types.Bytes)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs.Lengths))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.StrNeNullableScalar(lvs.Data, rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.StrNeScalar(lvs.Data, rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs.Lengths)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.StrNeNullableScalar(rvs.Data, lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.StrNeScalar(rvs.Data, lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs.Lengths))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs.Lengths)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.StrNe(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
		&BinOp{
			LeftType:   types.T_varchar,
			RightType:  types.T_varchar,
			ReturnType: types.T_sel,
			Fn: func(lv, rv *vector.Vector, proc *process.Process, lc, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.(*types.Bytes), rv.Col.(*types.Bytes)
				switch {
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs.Lengths))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() {
						vec.SetCol(ne.StrNeNullableScalar(lvs.Data, rvs, rv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.StrNeScalar(lvs.Data, rvs, rs))
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths))*8, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs.Lengths)]
					if lv.Nsp.Any() {
						vec.SetCol(ne.StrNeNullableScalar(rvs.Data, lvs, lv.Nsp.Np, rs))
					} else {
						vec.SetCol(ne.StrNeScalar(rvs.Data, lvs, rs))
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				}
				vec, err := register.Get(proc, int64(len(lvs.Lengths))*8, SelsType)
				if err != nil {
					return nil, err
				}
				rs := encoding.DecodeInt64Slice(vec.Data)
				rs = rs[:len(lvs.Lengths)]
				switch {
				case lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, roaring.Or(lv.Nsp.Np, rv.Nsp.Np), rs))
				case !lv.Nsp.Any() && rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, rv.Nsp.Np, rs))
				case lv.Nsp.Any() && !rv.Nsp.Any():
					vec.SetCol(ne.StrNeNullable(lvs, rvs, lv.Nsp.Np, rs))
				default:
					vec.SetCol(ne.StrNe(lvs, rvs, rs))
				}
				if lv.Ref == 0 {
					register.Put(proc, lv)
				}
				if rv.Ref == 0 {
					register.Put(proc, rv)
				}
				return vec, nil
			},
		},
	}
}
