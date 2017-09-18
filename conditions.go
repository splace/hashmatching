package main

func leadingZeroBits(c uint) func([]byte) bool {
	if c%8 == 0 {
		return leadingZeroBytes(c >> 3)
	}
	return func(t []byte) bool {
		return leadingZeroBytes(c>>3)(t) && t[c>>3]>>(8-c%8) == 0
	}
}

func leadingZeroBytes(c uint) func([]byte) bool {
	return func(t []byte) bool {
		switch c {
		case 16:
			if t[15] != 0 {
				return false
			}
			fallthrough
		case 15:
			if t[14] != 0 {
				return false
			}
			fallthrough
		case 14:
			if t[13] != 0 {
				return false
			}
			fallthrough
		case 13:
			if t[12] != 0 {
				return false
			}
			fallthrough
		case 12:
			if t[11] != 0 {
				return false
			}
			fallthrough
		case 11:
			if t[10] != 0 {
				return false
			}
			fallthrough
		case 10:
			if t[9] != 0 {
				return false
			}
			fallthrough
		case 9:
			if t[8] != 0 {
				return false
			}
			fallthrough
		case 8:
			if t[7] != 0 {
				return false
			}
			fallthrough
		case 7:
			if t[6] != 0 {
				return false
			}
			fallthrough
		case 6:
			if t[5] != 0 {
				return false
			}
			fallthrough
		case 5:
			if t[4] != 0 {
				return false
			}
			fallthrough
		case 4:
			if t[3] != 0 {
				return false
			}
			fallthrough
		case 3:
			if t[2] != 0 {
				return false
			}
			fallthrough
		case 2:
			if t[1] != 0 {
				return false
			}
			fallthrough
		case 1:
			if t[0] != 0 {
				return false
			}
			fallthrough
		default:
			return true
		}
	}
}


func leadingSetBits(c uint) func([]byte) bool {
	if c%8 == 0 {
		return leadingSetBytes(c >> 3)
	}
	return func(t []byte) bool {
		return leadingSetBytes(c>>3)(t) && t[c>>3]>>(8-c%8) >0
	}
}

func leadingSetBytes(c uint) func([]byte) bool {
	return func(t []byte) bool {
		switch c {
		case 16:
			if t[15] != 0xff {
				return false
			}
			fallthrough
		case 15:
			if t[14] != 0xff {
				return false
			}
			fallthrough
		case 14:
			if t[13] != 0xff {
				return false
			}
			fallthrough
		case 13:
			if t[12] != 0xff {
				return false
			}
			fallthrough
		case 12:
			if t[11] != 0xff {
				return false
			}
			fallthrough
		case 11:
			if t[10] != 0xff {
				return false
			}
			fallthrough
		case 10:
			if t[9] != 0xff {
				return false
			}
			fallthrough
		case 9:
			if t[8] != 0xff {
				return false
			}
			fallthrough
		case 8:
			if t[7] != 0xff {
				return false
			}
			fallthrough
		case 7:
			if t[6] != 0xff {
				return false
			}
			fallthrough
		case 6:
			if t[5] != 0xff {
				return false
			}
			fallthrough
		case 5:
			if t[4] != 0xff {
				return false
			}
			fallthrough
		case 4:
			if t[3] != 0xff {
				return false
			}
			fallthrough
		case 3:
			if t[2] != 0xff {
				return false
			}
			fallthrough
		case 2:
			if t[1] != 0xff {
				return false
			}
			fallthrough
		case 1:
			if t[0] != 0xff {
				return false
			}
			fallthrough
		default:
			return true
		}
	}
}

