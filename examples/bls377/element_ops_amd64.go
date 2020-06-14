// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goff (v0.2.2) DO NOT EDIT

// Package bls377 contains field arithmetic operations
package bls377

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

//go:noescape
func mulAssignElement(res, y *Element)

//go:noescape
func mulElement(res, x, y *Element)

//go:noescape
func addAssignElement(res, y *Element)

//go:noescape
func addElement(res, x, y *Element)

//go:noescape
func subAssignElement(res, y *Element)

//go:noescape
func subElement(res, x, y *Element)

//go:noescape
func doubleElement(res, y *Element)

//go:noescape
func fromMontElement(res *Element)

//go:noescape
func reduceElement(res *Element) // for test purposes

//go:noescape
func squareElement(res, y *Element)

// modulus
var modulusElement = Element{
	9586122913090633729,
	1660523435060625408,
	2230234197602682880,
	1883307231910630287,
	14284016967150029115,
	121098312706494698,
}

var modulusElementInv0 uint64 = 9586122913090633727

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *Element) FromMont() *Element {
	fromMontElement(z)
	return z
}

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Mul(x, y *Element) *Element {
	mulElement(z, x, y)
	return z
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) MulAssign(x *Element) *Element {
	mulAssignElement(z, x)
	return z
}

// Add z = x + y mod q
func (z *Element) Add(x, y *Element) *Element {
	addElement(z, x, y)
	return z
}

// AddAssign z = z + x mod q
func (z *Element) AddAssign(x *Element) *Element {
	addAssignElement(z, x)
	return z
}

// Double z = x + x mod q, aka Lsh 1
func (z *Element) Double(x *Element) *Element {
	doubleElement(z, x)
	return z
}

// Sub  z = x - y mod q
func (z *Element) Sub(x, y *Element) *Element {
	subElement(z, x, y)
	return z
}

// SubAssign  z = z - x mod q
func (z *Element) SubAssign(x *Element) *Element {
	subAssignElement(z, x)
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Square(x *Element) *Element {
	squareElement(z, x)
	return z
}
