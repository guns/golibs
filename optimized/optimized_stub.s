// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !amd64

#include "textflag.h"

// func Mul64(a, b uint64) (upper, lower uint64)
TEXT ·Mul64(SB), NOSPLIT, $0
	JMP ·mul64(SB)
