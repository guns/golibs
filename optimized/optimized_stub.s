// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !amd64

#include "textflag.h"

// func Mul(x, y uint) (lower, upper uint)
TEXT ·Mul(SB), NOSPLIT, $0
	JMP ·mul(SB)

// func Mul64(a, b uint64) (upper, lower uint64)
TEXT ·Mul64(SB), NOSPLIT, $0
	JMP ·mul64(SB)

// func Mul32(x, y uint32) (lower, upper uint32)
TEXT ·Mul32(SB), NOSPLIT, $0
	JMP ·mul32(SB)

// func RandIntn(n int) int
TEXT ·RandIntn(SB), NOSPLIT, $0
	JMP ·randIntn(SB)
