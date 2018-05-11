package main

import (
	wallet "github.com/skycoin/skycoin/src/wallet"
	"reflect"
	"unsafe"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_wallet_LoadWallets
func SKY_wallet_LoadWallets(_dir string, _arg1 *C.wallet__Wallets) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	dir := _dir
	__arg1, ____return_err := wallet.LoadWallets(dir)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg1 = *(*C.wallet__Wallets)(unsafe.Pointer(&__arg1))
	}
	return
}

//export SKY_wallet_Wallets_ToReadable
func SKY_wallet_Wallets_ToReadable(_wlts *C.wallet__Wallets, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	wlts := *(*wallet.Wallets)(unsafe.Pointer(_wlts))
	__arg0 := wlts.ToReadable()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}