package asm

import (
	"fmt"

	"github.com/consensys/bavard"
)

func (b *Builder) mul(asm *bavard.Assembly) error {
	asm.FuncHeader("_mulADX"+b.elementName, 0, 24)

	asm.WriteLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// however, to benefit from the ADCX and ADOX carry chains
	// we split the inner loops in 2:
	// for i=0 to N-1
	// 		for j=0 to N-1
	// 		    (A,t[j])  := t[j] + x[j]*y[i] + A
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C + A
	`)

	// registers
	t := asm.PopRegisters(b.nbWords)
	x := asm.PopRegister()
	y := asm.PopRegister()
	A := asm.PopRegister()
	tmp := asm.PopRegister()

	// dereference x and y
	asm.MOVQ("x+8(FP)", x)
	asm.MOVQ("y+16(FP)", y)

	for i := 0; i < b.nbWords; i++ {
		asm.XORQ(bavard.DX, bavard.DX)

		asm.MOVQ(y.At(i), bavard.DX)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < b.nbWords; j++ {
			xj := x.At(j)

			reg := A
			if i == 0 {
				if j == 0 {
					asm.MULXQ(xj, t[j], t[j+1])
				} else if j != b.nbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				asm.ADCXQ(A, t[j])
			}

			if !(i == 0 && j == 0) {
				asm.MULXQ(xj, bavard.AX, reg)
				asm.ADOXQ(bavard.AX, t[j])
			}
		}

		asm.Comment("add the last carries to " + string(A))
		asm.MOVQ(0, bavard.DX)
		asm.ADCXQ(bavard.DX, A)
		asm.ADOXQ(bavard.DX, A)

		// m := t[0]*q'[0] mod W
		regM := bavard.DX
		asm.MOVQ(t[0], bavard.DX)
		asm.MULXQ(qInv0(b.elementName), regM, bavard.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		asm.XORQ(bavard.AX, bavard.AX)

		// C,_ := t[0] + m*q[0]
		asm.Comment("C,_ := t[0] + m*q[0]")

		asm.MULXQ(qAt(0, b.elementName), bavard.AX, tmp)
		asm.ADCXQ(t[0], bavard.AX)
		asm.MOVQ(tmp, t[0])

		asm.Comment("for j=1 to N-1")
		asm.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < b.nbWords; j++ {
			asm.ADCXQ(t[j], t[j-1])
			asm.MULXQ(qAt(j, b.elementName), bavard.AX, t[j])
			asm.ADOXQ(bavard.AX, t[j-1])
		}
		asm.MOVQ(0, bavard.AX)
		asm.ADCXQ(bavard.AX, t[b.nbWordsLastIndex])
		asm.ADOXQ(A, t[b.nbWordsLastIndex])
	}

	// free registers
	asm.PushRegister(y, A, tmp)

	// ---------------------------------------------------------------------------------------------
	// reduce
	asm.MOVQ("res+0(FP)", x)
	b.reduce(asm, t, x)

	asm.RET()
	return nil
}

func (b *Builder) mulLarge(asm *bavard.Assembly) error {
	argSize := 8 + 2*b.nbWords*8 // 8 for res ptr, then 8 for each word for x and y
	asm.FuncHeader("_mulLargeADX"+b.elementName, b.nbWords*8, argSize)

	asm.WriteLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// however, to benefit from the ADCX and ADOX carry chains
	// we split the inner loops in 2:
	// for i=0 to N-1
	// 		for j=0 to N-1
	// 		    (A,t[j])  := t[j] + x[j]*y[i] + A
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C + A
	`)

	// registers
	t := asm.PopRegisters(b.nbWords)
	A := asm.PopRegister()

	for i := 0; i < b.nbWords; i++ {

		asm.XORQ(bavard.DX, bavard.DX)

		yi := fmt.Sprintf("y%d+%d(FP)", i, 8+8*b.nbWords+i*8)
		asm.MOVQ(yi, bavard.DX)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < b.nbWords; j++ {
			xj := fmt.Sprintf("x%d+%d(FP)", j, 8+j*8)

			reg := A
			if i == 0 {
				if j == 0 {
					asm.MULXQ(xj, t[j], t[j+1])
				} else if j != b.nbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				asm.ADCXQ(A, t[j])
			}

			if !(i == 0 && j == 0) {
				asm.MULXQ(xj, bavard.AX, reg)
				asm.ADOXQ(bavard.AX, t[j])
			}
		}

		asm.Comment("add the last carries to " + string(A))
		asm.MOVQ(0, bavard.DX)
		asm.ADCXQ(bavard.DX, A)
		asm.ADOXQ(bavard.DX, A)
		asm.PUSHQ(A)

		// m := t[0]*q'[0] mod W
		regM := bavard.DX
		asm.MOVQ(t[0], bavard.DX)
		asm.MULXQ(qInv0(b.elementName), regM, bavard.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		asm.XORQ(bavard.AX, bavard.AX)

		// C,_ := t[0] + m*q[0]
		asm.Comment("C,_ := t[0] + m*q[0]")
		asm.MULXQ(qAt(0, b.elementName), bavard.AX, A)
		asm.ADCXQ(t[0], bavard.AX)
		asm.MOVQ(A, t[0])

		asm.Comment("for j=1 to N-1")
		asm.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < b.nbWords; j++ {
			asm.ADCXQ(t[j], t[j-1])
			asm.MULXQ(qAt(j, b.elementName), bavard.AX, t[j])
			asm.ADOXQ(bavard.AX, t[j-1])
		}

		asm.POPQ(A)
		asm.MOVQ(0, bavard.AX)
		asm.ADCXQ(bavard.AX, t[b.nbWordsLastIndex])
		asm.ADOXQ(A, t[b.nbWordsLastIndex])
	}

	// free registers
	asm.PushRegister(A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := asm.PopRegister()
	asm.MOVQ("res+0(FP)", r)
	b.reduceLarge(asm, t, r)
	asm.RET()
	return nil
}