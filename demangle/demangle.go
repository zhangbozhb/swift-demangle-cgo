package demangle

//go:generate stringer -type=returnCode -output return_string.go

/*
#cgo CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/demangle/include
#cgo CPPFLAGS: -I/opt/homebrew/Cellar/zig/0.11.0/lib/zig/libcxx/include -I${SRCDIR}/include -I${SRCDIR}/demangle/include
#cgo LDFLAGS: -L${SRCDIR}/lib -L${SRCDIR}/demangle/lib

#include <stdlib.h>

#include "swift/Demangling/Demangle.h"

#define SYMBOLIC_SWIFT_FEATURE_RETURN_TYPE 0x1
#define SYMBOLIC_SWIFT_FEATURE_PARAMETERS 0x2
#define SYMBOLIC_SWIFT_FEATURE_ALL 0x3

extern "C" int symbolic_demangle_swift(const char *symbol,
                                       char *buffer,
                                       size_t buffer_length,
                                       int features) {
    swift::Demangle::DemangleOptions opts;

    if (features < SYMBOLIC_SWIFT_FEATURE_ALL) {
        opts = swift::Demangle::DemangleOptions::SimplifiedUIDemangleOptions();
        bool return_type = features & SYMBOLIC_SWIFT_FEATURE_RETURN_TYPE;
        bool argument_types = features & SYMBOLIC_SWIFT_FEATURE_PARAMETERS;

        opts.ShowFunctionReturnType = return_type;
        opts.ShowFunctionArgumentTypes = argument_types;
    }

    std::string demangled =
        swift::Demangle::demangleSymbolAsString(llvm::StringRef(symbol), opts);

    if (demangled.size() == 0 || demangled.size() >= buffer_length) {
        return false;
    }

    memcpy(buffer, demangled.c_str(), demangled.size());
    buffer[demangled.size()] = '\0';
    return true;
}

extern "C" int symbolic_demangle_is_swift_symbol(const char *symbol) {
    return swift::Demangle::isSwiftSymbol(symbol);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	SYMBOLIC_SWIFT_FEATURE_RETURN_TYPE = 0x1
	SYMBOLIC_SWIFT_FEATURE_PARAMETERS  = 0x2
	SYMBOLIC_SWIFT_FEATURE_ALL         = 0x3
)

func Demangle(symbol string) (string, error) {
	results := make([]byte, 1024)
	if rc := C.symbolic_demangle_swift(
		(*C.char)(unsafe.Pointer(symbol)),  // const char *symbol
		(*C.char)(unsafe.Pointer(results)), // char *buffer
		C.size_t(len(results)),             // size_t buffer_length
		C.int(SYMBOLIC_SWIFT_FEATURE_ALL),  // features
	); rc != 1 {
		return "", fmt.Errorf("failed to demangle symbol")
	}

	return C.GoString((*C.char)(unsafe.Pointer(results))), nil
}
