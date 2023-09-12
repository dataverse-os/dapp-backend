package lib

/*
#cgo LDFLAGS: -L./ -lrs_binding
#include <stdlib.h>
#include "./rs-binding.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func GenerateDID(key string) (string, error) {
	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))
	o := C.generate_did(c_key)
	if o.err == 0 {
		return C.GoString(o.data), nil
	} else {
		return "", errors.New(C.GoString(o.data))
	}
}

func GetCeramicNodeStatus(ceramic, key string) (bool, error) {
	c_input := C.CString(ceramic)
	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_input))
	defer C.free(unsafe.Pointer(c_key))
	o := C.get_ceramic_node_status(c_input, c_key)
	if o.err == 0 {
		return C.GoString(o.data) == "success", nil
	} else {
		return false, errors.New(C.GoString(o.data))
	}
}
