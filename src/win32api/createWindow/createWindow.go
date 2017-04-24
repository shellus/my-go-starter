package main

import (
	"github.com/AllenDang/w32"
	"syscall"
	"unsafe"
)


//A simple message box example
func main() {
	windowMain()

}

func MakeIntResource(id uint16) (*uint16) {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}
func WndProc(hWnd w32.HWND, msg uint32, wParam, lParam uintptr) (uintptr) {
	switch msg {
	case w32.WM_CREATE:
	_ = addButton(10,10,100,40, hWnd, "发送")
	case w32.WM_DESTROY:
		w32.PostQuitMessage(0)
	default:
		return w32.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}

func windowMain() int {

	lpszClassName := syscall.StringToUTF16Ptr("WNDclass")
	hInstance := w32.GetModuleHandle("")
	var wcex w32.WNDCLASSEX

	wcex.Size = uint32(unsafe.Sizeof(wcex))
	wcex.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	wcex.WndProc = syscall.NewCallback(WndProc)
	wcex.ClsExtra = 0
	wcex.WndExtra = 0
	wcex.Instance = hInstance
	wcex.Icon = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))
	wcex.Cursor = w32.LoadCursor(0, MakeIntResource(w32.IDC_ARROW))
	wcex.Background = w32.COLOR_WINDOW + 11

	wcex.MenuName = nil

	wcex.ClassName = lpszClassName
	wcex.IconSm = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))

	wcex.MenuName = nil

	wcex.ClassName = lpszClassName
	wcex.IconSm = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))

	w32.RegisterClassEx(&wcex)

	hWnd := w32.CreateWindowEx(
		w32.WS_EX_CLIENTEDGE,
		lpszClassName,
		syscall.StringToUTF16Ptr("Simple Go Window!"),
		w32.WS_OVERLAPPEDWINDOW | w32.WS_VISIBLE,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		600,
		400,
		0,
		0,
		hInstance,
		nil)

	w32.ShowWindow(hWnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(hWnd)

	var msg w32.MSG
	for {
		if w32.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}
	return int(msg.WParam)
}
func addButton(x, y, w, h int, parent w32.HWND, text string) (btnHwnd w32.HWND) {
	hWnd := w32.CreateWindowEx(
		w32.WS_EX_CLIENTEDGE,
		syscall.StringToUTF16Ptr("button"),
		syscall.StringToUTF16Ptr(text),
		w32.WS_CHILD | w32.WS_VISIBLE | w32.BS_PUSHBUTTON,
		x,
		y,
		w,
		h,
		parent,
		0,
		0,
		nil)

	return hWnd
}