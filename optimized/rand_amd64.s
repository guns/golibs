// Copyright (c) 2018 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

#include "textflag.h"

// func RandIntn(n int) int
TEXT ·RandIntn(SB), NOSPLIT, $0
	JMP ·RandInt63n(SB)
