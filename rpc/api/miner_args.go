// Copyright 2015 The go-krypton Authors
// This file is part of the go-krypton library.
//
// The go-krypton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-krypton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-krypton library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"encoding/json"

	"math/big"

	"github.com/krypton/go-krypton/common"
	"github.com/krypton/go-krypton/rpc/shared"
)

type StartMinerArgs struct {
	Threads int
}

func (args *StartMinerArgs) UnmarshalJSON(b []byte) (err error) {
	var obj []interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return shared.NewDecodeParamError(err.Error())
	}

	if len(obj) == 0 || obj[0] == nil {
		args.Threads = -1
		return nil
	}

	var num *big.Int
	if num, err = numString(obj[0]); err != nil {
		return err
	}
	args.Threads = int(num.Int64())
	return nil
}

type SetExtraArgs struct {
	Data string
}

func (args *SetExtraArgs) UnmarshalJSON(b []byte) (err error) {
	var obj []interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return shared.NewDecodeParamError(err.Error())
	}

	if len(obj) < 1 {
		return shared.NewInsufficientParamsError(len(obj), 1)
	}

	extrastr, ok := obj[0].(string)
	if !ok {
		return shared.NewInvalidTypeError("Price", "not a string")
	}
	args.Data = extrastr

	return nil
}

type GasPriceArgs struct {
	Price string
}

func (args *GasPriceArgs) UnmarshalJSON(b []byte) (err error) {
	var obj []interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return shared.NewDecodeParamError(err.Error())
	}

	if len(obj) < 1 {
		return shared.NewInsufficientParamsError(len(obj), 1)
	}

	if pricestr, ok := obj[0].(string); ok {
		args.Price = pricestr
		return nil
	}

	return shared.NewInvalidTypeError("Price", "not a string")
}

type SetKrbaseArgs struct {
	Krbase common.Address
}

func (args *SetKrbaseArgs) UnmarshalJSON(b []byte) (err error) {
	var obj []interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return shared.NewDecodeParamError(err.Error())
	}

	if len(obj) < 1 {
		return shared.NewInsufficientParamsError(len(obj), 1)
	}

	if addr, ok := obj[0].(string); ok {
		args.Krbase = common.HexToAddress(addr)
		if (args.Krbase == common.Address{}) {
			return shared.NewInvalidTypeError("Krbase", "not a valid address")
		}
		return nil
	}

	return shared.NewInvalidTypeError("Krbase", "not a string")
}

type MakeDAGArgs struct {
	BlockNumber int64
}

func (args *MakeDAGArgs) UnmarshalJSON(b []byte) (err error) {
	args.BlockNumber = -1
	var obj []interface{}

	if err := json.Unmarshal(b, &obj); err != nil {
		return shared.NewDecodeParamError(err.Error())
	}

	if len(obj) < 1 {
		return shared.NewInsufficientParamsError(len(obj), 1)
	}

	if err := blockHeight(obj[0], &args.BlockNumber); err != nil {
		return err
	}

	return nil
}