package main

import (
	"unsafe"
	cli "github.com/skycoin/skycoin/src/cli"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cli_GenerateWallet
func SKY_cli_GenerateWallet(_walletFile string, _opts *C.wallet__Options, _numAddrs uint64, _arg3 *C.Wallet__Handle) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	opts := *(*wallet.Options)(unsafe.Pointer(_opts))
	numAddrs := _numAddrs
	__arg3, ____return_err := cli.GenerateWallet(walletFile, opts, numAddrs)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg3 = registerWalletHandle(__arg3)
	}
	return
}

//export SKY_cli_MakeAlphanumericSeed
func SKY_cli_MakeAlphanumericSeed(_arg0 *C.GoString_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	__arg0 := cli.MakeAlphanumericSeed()
	copyString(__arg0, _arg0)
	return
}