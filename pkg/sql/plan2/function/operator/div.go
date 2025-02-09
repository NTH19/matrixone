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

package operator

import (
	"github.com/matrixorigin/matrixone/pkg/container/nulls"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/encoding"
	"github.com/matrixorigin/matrixone/pkg/vectorize/div"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
	"golang.org/x/exp/constraints"
)

func Div[T constraints.Float](vectors []*vector.Vector, proc *process.Process) (*vector.Vector, error) {
	lv, rv := vectors[0], vectors[1]
	lvs, rvs := lv.Col.([]T), rv.Col.([]T)
	rtl := lv.Typ.Oid.FixedLength()

	if lv.IsScalarNull() || rv.IsScalarNull() {
		return proc.AllocScalarNullVector(lv.Typ), nil
	}

	switch {
	case lv.IsScalar() && rv.IsScalar():
		vec := proc.AllocScalarVector(lv.Typ)
		rs := make([]T, 1)
		nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if v == 0 {
					return nil, ErrDivByZero
				}
			}
			vector.SetCol(vec, div.NumericDiv[T](lvs, rvs, rs))
			return vec, nil
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if rvs[i] == 0 {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vector.SetCol(vec, div.NumericDivSels[T](lvs, rvs, rs, sels))
		return vec, nil
	case lv.IsScalar() && !rv.IsScalar():
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if v == 0 {
					return nil, ErrDivByZero
				}
			}
			vec, err := proc.AllocVector(lv.Typ, int64(rtl)*int64(len(rvs)))
			if err != nil {
				return nil, err
			}
			rs := encoding.DecodeFixedSlice[T](vec.Data, rtl)
			nulls.Set(vec.Nsp, rv.Nsp)
			vector.SetCol(vec, div.NumericDivScalar[T](lvs[0], rvs, rs))
			return vec, nil
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if rvs[i] == 0 {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vec, err := proc.AllocVector(lv.Typ, int64(rtl)*int64(len(rvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeFixedSlice[T](vec.Data, rtl)
		nulls.Set(vec.Nsp, rv.Nsp)
		vector.SetCol(vec, div.NumericDivScalarSels[T](lvs[0], rvs, rs, sels))
		return vec, nil
	case !lv.IsScalar() && rv.IsScalar():
		if rvs[0] == 0 {
			return nil, ErrDivByZero
		}
		vec, err := proc.AllocVector(lv.Typ, int64(rtl)*int64(len(lvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeFixedSlice[T](vec.Data, rtl)
		nulls.Set(vec.Nsp, lv.Nsp)
		vector.SetCol(vec, div.NumericDivByScalar[T](rvs[0], lvs, rs))
		return vec, nil
	}
	vec, err := proc.AllocVector(lv.Typ, int64(rtl)*int64(len(lvs)))
	if err != nil {
		return nil, err
	}
	rs := encoding.DecodeFixedSlice[T](vec.Data, rtl)
	nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
	if !nulls.Any(rv.Nsp) {
		for _, v := range rvs {
			if v == 0 {
				return nil, ErrDivByZero
			}
		}
		vector.SetCol(vec, div.NumericDiv[T](lvs, rvs, rs))
		return vec, nil
	}
	sels := process.GetSels(proc)
	defer process.PutSels(sels, proc)
	for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
		if nulls.Contains(rv.Nsp, i) {
			continue
		}
		if rvs[i] == 0 {
			return nil, ErrDivByZero
		}
		sels = append(sels, int64(i))
	}
	vector.SetCol(vec, div.NumericDivSels[T](lvs, rvs, rs, sels))
	return vec, nil
}

func DivDecimal64(vectors []*vector.Vector, proc *process.Process) (*vector.Vector, error) {
	lv, rv := vectors[0], vectors[1]
	lvs, rvs := lv.Col.([]types.Decimal64), rv.Col.([]types.Decimal64)
	lvScale, rvScale := lv.Typ.Scale, rv.Typ.Scale
	resultScale := lv.Typ.Scale
	resultTyp := types.Type{Oid: types.T_decimal128, Size: 16, Width: 38, Scale: resultScale}

	if lv.IsScalarNull() || rv.IsScalarNull() {
		return proc.AllocScalarNullVector(lv.Typ), nil
	}

	switch {
	case lv.IsScalar() && rv.IsScalar():
		vec := proc.AllocScalarVector(resultTyp)
		rs := make([]types.Decimal128, 1)
		nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if int64(v) == 0 {
					return nil, ErrDivByZero
				}
			}
			vector.SetCol(vec, div.Decimal64Div(lvs, rvs, lvScale, rvScale, rs))
			vec.Typ = resultTyp
			return vec, nil
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if int64(rvs[i]) == 0 {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vector.SetCol(vec, div.Decimal64DivSels(lvs, rvs, lvScale, rvScale, rs, sels))
		vec.Typ = resultTyp
		return vec, nil
	case lv.IsScalar() && !rv.IsScalar():
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if int64(v) == 0 {
					return nil, ErrDivByZero
				}
			}
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if int64(rvs[i]) == 0 {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vec, err := proc.AllocVector(resultTyp, int64(resultTyp.Size)*int64(len(rvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeDecimal128Slice(vec.Data)
		rs = rs[:len(rvs)]
		nulls.Set(vec.Nsp, rv.Nsp)
		vector.SetCol(vec, div.Decimal64DivScalarSels(lvs[0], rvs, lvScale, rvScale, rs, sels))
		vec.Typ = resultTyp
		return vec, nil
	case !lv.IsScalar() && rv.IsScalar():
		if int64(rvs[0]) == 0 {
			return nil, ErrDivByZero
		}
		vec, err := proc.AllocVector(resultTyp, int64(resultTyp.Size)*int64(len(lvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeDecimal128Slice(vec.Data)
		rs = rs[:len(lvs)]
		nulls.Set(vec.Nsp, lv.Nsp)
		vector.SetCol(vec, div.Decimal64DivByScalar(rvs[0], lvs, rvScale, lvScale, rs))
		vec.Typ = resultTyp
		return vec, nil
	}
	vec, err := proc.AllocVector(resultTyp, int64(resultTyp.Size)*int64(len(lvs)))
	if err != nil {
		return nil, err
	}
	rs := encoding.DecodeDecimal128Slice(vec.Data)
	rs = rs[:len(rvs)]
	nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
	if !nulls.Any(rv.Nsp) {
		for _, v := range rvs {
			if int64(v) == 0 {
				return nil, ErrDivByZero
			}
		}
		vector.SetCol(vec, div.Decimal64Div(lvs, rvs, lvScale, rvScale, rs))
		vec.Typ = resultTyp
		return vec, nil
	}
	sels := process.GetSels(proc)
	defer process.PutSels(sels, proc)
	for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
		if nulls.Contains(rv.Nsp, i) {
			continue
		}
		if int64(rvs[i]) == 0 {
			return nil, ErrDivByZero
		}
		sels = append(sels, int64(i))
	}
	vector.SetCol(vec, div.Decimal64DivSels(lvs, rvs, lvScale, rvScale, rs, sels))
	vec.Typ = resultTyp
	return vec, nil
}

func DivDecimal128(vectors []*vector.Vector, proc *process.Process) (*vector.Vector, error) {
	lv, rv := vectors[0], vectors[1]
	lvs, rvs := lv.Col.([]types.Decimal128), rv.Col.([]types.Decimal128)
	lvScale, rvScale := lv.Typ.Scale, rv.Typ.Scale
	resultScale := lv.Typ.Scale
	resultTyp := types.Type{Oid: types.T_decimal128, Size: 16, Width: 38, Scale: resultScale}
	if lv.IsScalarNull() || rv.IsScalarNull() {
		return proc.AllocScalarNullVector(lv.Typ), nil
	}

	switch {
	case lv.IsScalar() && rv.IsScalar():
		vec := proc.AllocScalarVector(resultTyp)
		rs := make([]types.Decimal128, 1)
		nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if types.Decimal128IsZero(v) {
					return nil, ErrDivByZero
				}
			}
			vector.SetCol(vec, div.Decimal128Div(lvs, rvs, lvScale, rvScale, rs))
			vec.Typ = resultTyp
			return vec, nil
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if types.Decimal128IsZero(rvs[i]) {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vector.SetCol(vec, div.Decimal128DivSels(lvs, rvs, lvScale, rvScale, rs, sels))
		vec.Typ = resultTyp
		return vec, nil
	case lv.IsScalar() && !rv.IsScalar():
		if !nulls.Any(rv.Nsp) {
			for _, v := range rvs {
				if types.Decimal128IsZero(v) {
					return nil, ErrDivByZero
				}
			}
			vec, err := proc.AllocVector(lv.Typ, int64(resultTyp.Size)*int64(len(rvs)))
			if err != nil {
				return nil, err
			}
			rs := encoding.DecodeDecimal128Slice(vec.Data)
			rs = rs[:len(rvs)]
			nulls.Set(vec.Nsp, rv.Nsp)
			vector.SetCol(vec, div.Decimal128DivScalar(lvs[0], rvs, lvScale, rvScale, rs))
			vec.Typ = resultTyp
			return vec, nil
		}
		sels := process.GetSels(proc)
		defer process.PutSels(sels, proc)
		for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
			if nulls.Contains(rv.Nsp, i) {
				continue
			}
			if types.Decimal128IsZero(rvs[i]) {
				return nil, ErrDivByZero
			}
			sels = append(sels, int64(i))
		}
		vec, err := proc.AllocVector(lv.Typ, int64(resultTyp.Size)*int64(len(rvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeDecimal128Slice(vec.Data)
		rs = rs[:len(rvs)]
		nulls.Set(vec.Nsp, rv.Nsp)
		vector.SetCol(vec, div.Decimal128DivScalarSels(lvs[0], rvs, lvScale, rvScale, rs, sels))
		vec.Typ = resultTyp
		return vec, nil
	case !lv.IsScalar() && rv.IsScalar():
		if types.Decimal128IsZero(rvs[0]) {
			return nil, ErrDivByZero
		}
		vec, err := proc.AllocVector(lv.Typ, int64(resultTyp.Size)*int64(len(lvs)))
		if err != nil {
			return nil, err
		}
		rs := encoding.DecodeDecimal128Slice(vec.Data)
		rs = rs[:len(lvs)]
		nulls.Set(vec.Nsp, lv.Nsp)
		vector.SetCol(vec, div.Decimal128DivByScalar(rvs[0], lvs, rvScale, lvScale, rs))
		vec.Typ = resultTyp
		return vec, nil
	}
	vec, err := proc.AllocVector(lv.Typ, int64(resultTyp.Size)*int64(len(lvs)))
	if err != nil {
		return nil, err
	}
	rs := encoding.DecodeDecimal128Slice(vec.Data)
	rs = rs[:len(rvs)]
	nulls.Or(lv.Nsp, rv.Nsp, vec.Nsp)
	if !nulls.Any(rv.Nsp) {
		for _, v := range rvs {
			if types.Decimal128IsZero(v) {
				return nil, ErrDivByZero
			}
		}
		vector.SetCol(vec, div.Decimal128Div(lvs, rvs, lvScale, rvScale, rs))
		vec.Typ = resultTyp
		return vec, nil
	}
	sels := process.GetSels(proc)
	defer process.PutSels(sels, proc)
	for i, j := uint64(0), uint64(len(rvs)); i < j; i++ {
		if nulls.Contains(rv.Nsp, i) {
			continue
		}
		if types.Decimal128IsZero(rvs[i]) {
			return nil, ErrDivByZero
		}
		sels = append(sels, int64(i))
	}
	vector.SetCol(vec, div.Decimal128DivSels(lvs, rvs, lvScale, rvScale, rs, sels))
	vec.Typ = resultTyp
	return vec, nil
}
