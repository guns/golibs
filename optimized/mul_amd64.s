// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

#include "textflag.h"

// func Mul(x, y uint) (lower, upper uint)
TEXT ·Mul(SB), NOSPLIT, $0
	MOVQ x+0(FP), AX
	MOVQ y+8(FP), DX
	MULQ DX               // DX*AX -> DX:AX
	MOVQ AX, lower+16(FP)
	MOVQ DX, upper+24(FP)
	RET

// func Mul64(x, y uint64) (lower, upper uint64)
TEXT ·Mul64(SB), NOSPLIT, $0
	MOVQ x+0(FP), AX
	MOVQ y+8(FP), DX
	MULQ DX               // DX*AX -> DX:AX
	MOVQ AX, lower+16(FP)
	MOVQ DX, upper+24(FP)
	RET

// func Mul32(x, y uint32) (lower, upper uint32)
TEXT ·Mul32(SB), NOSPLIT, $0
	MOVL x+0(FP), AX
	MOVL y+4(FP), DX
	MULL DX               // DX*AX -> DX:AX
	MOVL AX, lower+8(FP)
	MOVL DX, upper+12(FP)
	RET
