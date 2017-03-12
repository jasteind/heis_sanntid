package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import (
	"C"
)

func io_init() bool {

	return int(C.io_init()) != 0
}

func io_set_bit(channel int) {

	C.io_set_bit(C.int(channel))
}

func io_clear_bit(channel int) {

	C.io_clear_bit(C.int(channel))
}

func io_write_analog(channel int, value int) {

	C.io_write_analog(C.int(channel), C.int(value))
}

func io_read_bit(channel int) bool {

	return (C.io_read_bit(C.int(channel))) != 0
}

func io_read_analog(channel int) bool {

	return (C.io_read_analog(C.int(channel))) != 0
}
