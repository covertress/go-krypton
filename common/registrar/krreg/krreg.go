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

package krreg

import (
	"math/big"

	"github.com/krypton/go-krypton/common/registrar"
	"github.com/krypton/go-krypton/xkr"
)

// implements a versioned Registrar on an archiving full node
type KrReg struct {
	backend  *xkr.XKr
	registry *registrar.Registrar
}

func New(xe *xkr.XKr) (self *KrReg) {
	self = &KrReg{backend: xe}
	self.registry = registrar.New(xe)
	return
}

func (self *KrReg) Registry() *registrar.Registrar {
	return self.registry
}

func (self *KrReg) Resolver(n *big.Int) *registrar.Registrar {
	xe := self.backend
	if n != nil {
		xe = self.backend.AtStateNum(n.Int64())
	}
	return registrar.New(xe)
}
