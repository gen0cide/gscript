// +build darwin

package main

/*

#include <stdio.h>

int          NXArgc = 0;
const char** NXArgv = NULL;
const char** environ = NULL;
const char*  __progname = NULL;

struct __DATA__dyld {
	long			lazy;
	int				(*lookup)(const char*, void**);
	// ProgramVars
	const void*		mh;
	int*			NXArgcPtr;
	const char***	NXArgvPtr;
	const char***	environPtr;
	const char**	__prognamePtr;
};

static volatile struct __DATA__dyld  myDyldSection __attribute__ ((section ("__DATA,__dyld")))
	= { 0, 0, NULL, &NXArgc, &NXArgv, &environ, &__progname };


uintptr_t GetAddr(const char* fnname) {
	void *fnp;
	printf("%p\n", fnp);
	int ec = myDyldSection.lookup(fnname, &fnp);
	printf("%d\n", ec);
	printf("%p\n", fnp);
	return (uintptr_t)fnp;
}

*/
import "C"

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must supply a function name")
		os.Exit(1)
	}
	fnName := os.Args[1]
	foo := C.GetAddr(C.CString(fnName))
	spew.Dump(uintptr(foo))
}
