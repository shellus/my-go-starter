package main

import (
	"syscall"
	"unsafe"
	"fmt"
)

func main() {

	//首先,准备输入参数, GetDiskFreeSpaceEx需要4个参数, 可查MSDN
	dir := "C:"
	lpFreeBytesAvailable := int64(0) //注意类型需要跟API的类型相符
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	//获取方法的引用
	kernel32, _ := syscall.LoadLibrary("Kernel32.dll")
	defer syscall.FreeLibrary(kernel32)

	GetDiskFreeSpaceEx, _ := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	r, _, _ := syscall.Syscall6(
		uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dir))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)),
		0,
		0)

	if r != 0 {
		fmt.Printf("Free %.2fGB", float64(lpTotalNumberOfFreeBytes) / 1024 / 1024 / 1024)
	}
}
