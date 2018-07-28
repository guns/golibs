// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

#include "textflag.h"

// func Mul64(x, y uint64) (lower, upper uint64)
TEXT ·Mul64(SB), NOSPLIT, $0
	MOVQ x+0(FP), AX
	MOVQ y+8(FP), DX
	MULQ DX               // DX*AX -> DX:AX
	MOVQ AX, lower+16(FP)
	MOVQ DX, upper+24(FP)
	RET
