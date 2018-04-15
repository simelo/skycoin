package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

import (
	"unsafe"
	webrpc "github.com/skycoin/skycoin/src/api/webrpc"
	wallet "github.com/skycoin/skycoin/src/wallet"
	cli "github.com/skycoin/skycoin/src/api/cli"
	gcli "github.com/urfave/cli"

)

type Handle uint64

var (
	handleMap = make(map[Handle]interface{})
)

func registerHandle(obj interface{}) Handle {
	ptr := &obj
	handle := *(*Handle)(unsafe.Pointer(&ptr))
	handleMap[handle] = obj
	return handle
}

func lookupHandleObj(handle Handle) (interface{}, bool) {
	obj, ok := handleMap[handle]
	return obj, ok
}

func registerWebRpcClientHandle(obj *webrpc.Client) C.WebrpcClient__Handle{
	return (C.WebrpcClient__Handle)(registerHandle(obj))
}

func lookupWebRpcClientHandle(handle C.WebrpcClient__Handle) (*webrpc.Client, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*webrpc.Client); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerWalletHandle(obj *wallet.Wallet) C.Wallet__Handle{
	return (C.Wallet__Handle)(registerHandle(obj))
}

func lookupWalletHandle(handle C.Wallet__Handle) (*wallet.Wallet, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*wallet.Wallet); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerConfigHandle(obj *cli.Config) C.Config__Handle{
	return (C.Config__Handle)(registerHandle(obj))
}

func lookupConfigHandle(handle C.Config__Handle) (*cli.Config, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.Config); isOK {
			return obj, true
		}
	}
	return nil, false
}

func registerAppHandle(obj *cli.App) C.App__Handle{
	return (C.App__Handle)(registerHandle(obj))
}

func lookupAppHandle(handle C.App__Handle) (*cli.App, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*cli.App); isOK {
			return obj, true
		}
	}
	return nil, false
}


func registerContextHandle(obj *gcli.Context) C.GcliContext__Handle{
	return (C.GcliContext__Handle)(registerHandle(obj))
}

func lookupContextHandle(handle C.GcliContext__Handle) (*gcli.Context, bool){
	obj, ok := lookupHandleObj(Handle(handle))
	if ok {
		if obj, isOK := (obj).(*gcli.Context); isOK {
			return obj, true
		}
	}
	return nil, false
}

func closeHandle(handle Handle) {
	delete(handleMap, handle)
}

//export SKY_handle_close
func SKY_handle_close(handle C.Handle){
	closeHandle(Handle(handle))
}