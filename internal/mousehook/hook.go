package mousehook

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	comctl32                 = syscall.NewLazyDLL("comctl32.dll")
	procSetWindowsHookEx     = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx  = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx       = user32.NewProc("CallNextHookEx")
	procGetModuleHandle      = kernel32.NewProc("GetModuleHandleW")
	procPostMessage          = user32.NewProc("PostMessageW")
	procSendMessage          = user32.NewProc("SendMessageW")
	procGetClientRect        = user32.NewProc("GetClientRect")
	procScreenToClient       = user32.NewProc("ScreenToClient")
	procPtInRect             = user32.NewProc("PtInRect")
	procGetMessage           = user32.NewProc("GetMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessage      = user32.NewProc("DispatchMessageW")
	procGetWindow            = user32.NewProc("GetWindow")
	procGetWindowText        = user32.NewProc("GetWindowTextW")
	procGetClassName         = user32.NewProc("GetClassNameW")
	procSetWindowSubclass    = comctl32.NewProc("SetWindowSubclass")
	procRemoveWindowSubclass = comctl32.NewProc("RemoveWindowSubclass")
	procDefSubclassProc      = comctl32.NewProc("DefSubclassProc")

	// 新增的API用于跨进程窗口操作
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
	procVirtualAllocEx           = kernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx            = kernel32.NewProc("VirtualFreeEx")
	procWriteProcessMemory       = kernel32.NewProc("WriteProcessMemory")
	procReadProcessMemory        = kernel32.NewProc("ReadProcessMemory")
	procCreateRemoteThread       = kernel32.NewProc("CreateRemoteThread")
	procWaitForSingleObject      = kernel32.NewProc("WaitForSingleObject")
	procGetExitCodeThread        = kernel32.NewProc("GetExitCodeThread")

	// DLL注入相关API
	procGetProcAddress        = kernel32.NewProc("GetProcAddress")
	procLoadLibraryW          = kernel32.NewProc("LoadLibraryW")
	procGetModuleFileNameW    = kernel32.NewProc("GetModuleFileNameW")
	procCreateFileW           = kernel32.NewProc("CreateFileW")
	procWriteFile             = kernel32.NewProc("WriteFile")
	procGetFileSize           = kernel32.NewProc("GetFileSize")
	procReadFile              = kernel32.NewProc("ReadFile")
	procSetFilePointer        = kernel32.NewProc("SetFilePointer")
	procGetTempPathW          = kernel32.NewProc("GetTempPathW")
	procGetTempFileNameW      = kernel32.NewProc("GetTempFileNameW")
	procDeleteFileW           = kernel32.NewProc("DeleteFileW")
	procVirtualProtectEx      = kernel32.NewProc("VirtualProtectEx")
	procFlushInstructionCache = kernel32.NewProc("FlushInstructionCache")

	// 窗口焦点相关API
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procSetFocus            = user32.NewProc("SetFocus")
	procSetActiveWindow     = user32.NewProc("SetActiveWindow")
	procBringWindowToTop    = user32.NewProc("BringWindowToTop")
	procShowWindow          = user32.NewProc("ShowWindow")
	procIsIconic            = user32.NewProc("IsIconic")
	procGetFocus            = user32.NewProc("GetFocus")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")

	// 窗口样式相关API
	procGetWindowLongPtr = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
)

const (
	WM_DESTROY     = 0x0002
	WH_MOUSE_LL    = 14
	WM_ACTIVATE    = 0x0006
	WM_CREATE      = 0x0001
	WM_NCCREATE    = 0x0081
	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
	WM_MOUSEHWHEEL = 0x020E
	WM_MOUSELEAVE  = 0x02A3

	WA_CLICKACTIVE = 2

	// 自定义消息
	WM_INSTALL_MOUSE_HOOK   = 0x0400 + 1
	WM_UNINSTALL_MOUSE_HOOK = 0x0400 + 2
	WM_SET_FOCUS            = 0x0400 + 3

	// DLL控制消息
	WM_DLL_INSTALL_HOOK        = 0x0401
	WM_DLL_UNINSTALL_HOOK      = 0x0402
	WM_DLL_ENABLE_MOUSE_LEAVE  = 0x0403
	WM_DLL_DISABLE_MOUSE_LEAVE = 0x0404

	// 进程访问权限
	PROCESS_ALL_ACCESS        = 0x1F0FFF
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_READ           = 0x0010
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400

	// 内存分配类型
	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
	MEM_RELEASE = 0x8000

	// 等待对象状态
	WAIT_OBJECT_0 = 0x00000000
	WAIT_TIMEOUT  = 0x00000102
	WAIT_FAILED   = 0xFFFFFFFF

	// 文件操作常量
	GENERIC_READ          = 0x80000000
	GENERIC_WRITE         = 0x40000000
	OPEN_EXISTING         = 3
	CREATE_ALWAYS         = 2
	FILE_ATTRIBUTE_NORMAL = 0x80
	INVALID_HANDLE_VALUE  = ^uintptr(0)

	// 内存保护常量
	PAGE_EXECUTE_READWRITE = 0x40
	PAGE_READWRITE         = 0x04
	PAGE_EXECUTE_READ      = 0x20

	// 窗口显示状态常量
	SW_RESTORE = 9
	SW_SHOW    = 5

	// 鼠标按键常量
	VK_LBUTTON = 0x01

	// 窗口样式常量
	GWLP_STYLE     = -16
	WS_CAPTION     = 0x00C00000
	WS_THICKFRAME  = 0x00040000
	WS_MINIMIZEBOX = 0x00020000
	WS_MAXIMIZEBOX = 0x00010000
	WS_SYSMENU     = 0x00080000
	WS_BORDER      = 0x00800000
	WS_DLGFRAME    = 0x00400000
	WS_POPUP       = 0x80000000

	// SetWindowPos常量
	SWP_FRAMECHANGED = 0x0020
	SWP_NOMOVE       = 0x0002
	SWP_NOSIZE       = 0x0001
	SWP_NOZORDER     = 0x0004
	SWP_SHOWWINDOW   = 0x0040
)

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  WPARAM
	LParam  LPARAM
	Time    uint32
	Pt      POINT
}

type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type (
	HHOOK   uintptr
	HWND    uintptr
	WPARAM  uintptr
	LPARAM  uintptr
	LRESULT uintptr
	HANDLE  uintptr
	DWORD   uint32
	BOOL    int32
)

var (
	hHook              HHOOK
	targetHWND         HWND
	mainHWND           HWND
	forwardToWindow    bool = true // 是否转发给窗口
	messageLoopRunning bool = false
	quitChan           chan bool
	subclassInstalled  bool = false // 是否已安装子类化
)

func SetWindowsHookEx(idHook int, lpfn uintptr, hMod uintptr, dwThreadId uint32) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		lpfn,
		hMod,
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(uintptr(hhk))
	return ret != 0
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func GetModuleHandle(lpModuleName *uint16) uintptr {
	ret, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(lpModuleName)))
	return ret
}

func PostMessage(hWnd HWND, Msg uint32, wParam WPARAM, lParam LPARAM) bool {
	ret, _, _ := procPostMessage.Call(
		uintptr(hWnd),
		uintptr(Msg),
		uintptr(wParam),
		uintptr(lParam),
	)
	return ret != 0
}

func SendMessage(hWnd HWND, Msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procSendMessage.Call(
		uintptr(hWnd),
		uintptr(Msg),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func GetClientRect(hWnd HWND, lpRect *RECT) bool {
	ret, _, _ := procGetClientRect.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lpRect)),
	)
	return ret != 0
}

func ScreenToClient(hWnd HWND, lpPoint *POINT) bool {
	ret, _, _ := procScreenToClient.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lpPoint)),
	)
	return ret != 0
}

func PtInRect(lprc *RECT, pt POINT) bool {
	ret, _, _ := procPtInRect.Call(
		uintptr(unsafe.Pointer(lprc)),
		uintptr(pt.X),
		uintptr(pt.Y),
	)
	return ret != 0
}

// MAKELPARAM combines two 16-bit values (lo and hi) into a single 32-bit LPARAM.
// Typically used to encode coordinates like x (low) and y (high) for mouse messages.
func MAKELPARAM(lo, hi uint16) uintptr {
	return uintptr(uint32(lo) | uint32(hi)<<16)
}
func mouseHookProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	//fmt.Printf("Mouse hook: %d, %d, %d\n", nCode, wParam, lParam)
	if nCode < 0 {
		return CallNextHookEx(hHook, nCode, wParam, lParam)
	}

	if targetHWND == 0 {
		return CallNextHookEx(hHook, nCode, wParam, lParam)
	}

	msg := uint32(wParam)
	hookStruct := (*MSLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))
	//
	//// 调试信息 - 显示所有鼠标事件
	//switch msg {
	////case WM_MOUSEMOVE:
	////	// 鼠标移动事件太频繁，只在窗口内时打印
	//case WM_LBUTTONDOWN:
	//	//将targetHWND设置为前台应用并获取焦点
	//	if mainHWND != 0 {
	//
	//		var clientRect RECT
	//		if !GetClientRect(targetHWND, &clientRect) {
	//			return CallNextHookEx(hHook, nCode, wParam, lParam)
	//		}
	//
	//		// 将屏幕坐标转换为客户端坐标
	//		clientPoint := hookStruct.Pt
	//		if !ScreenToClient(targetHWND, &clientPoint) {
	//			return CallNextHookEx(hHook, nCode, wParam, lParam)
	//		}
	//
	//		//fmt.Printf("Mouse in window: screen(%d,%d) -> client(%d,%d)\n", hookStruct.Pt.X, hookStruct.Pt.Y, clientPoint.X, clientPoint.Y)
	//		//fmt.Printf("Client Rect")
	//
	//		// 检查鼠标是否在窗口客户端区域内
	//		if PtInRect(&clientRect, clientPoint) {
	//			PostMessage(mainHWND, WM_SET_FOCUS, 0, 0)
	//		}
	//	}

	//	fmt.Printf("Left button down at (%d, %d)\n", hookStruct.Pt.X, hookStruct.Pt.Y)
	//case WM_LBUTTONUP:
	//	fmt.Printf("Left button up at (%d, %d)\n", hookStruct.Pt.X, hookStruct.Pt.Y)
	//case WM_RBUTTONDOWN:
	//	fmt.Printf("Right button down at (%d, %d)\n", hookStruct.Pt.X, hookStruct.Pt.Y)
	//case WM_RBUTTONUP:
	//	fmt.Printf("Right button up at (%d, %d)\n", hookStruct.Pt.X, hookStruct.Pt.Y)
	//case WM_MOUSEWHEEL:
	//	wheelDelta := int16(hookStruct.MouseData >> 16)
	//	fmt.Printf("Mouse wheel at (%d, %d), delta: %d\n", hookStruct.Pt.X, hookStruct.Pt.Y, wheelDelta)
	//default:
	//	fmt.Printf("Mouse event: %d, at (%d, %d)\n", msg, hookStruct.Pt.X, hookStruct.Pt.Y)

	//}

	if msg == WM_MOUSEMOVE {

		// 转发鼠标事件到窗口
		if forwardToWindow {
			// 获取窗口客户端区域
			var clientRect RECT
			if !GetClientRect(targetHWND, &clientRect) {
				return CallNextHookEx(hHook, nCode, wParam, lParam)
			}

			// 将屏幕坐标转换为客户端坐标
			clientPoint := hookStruct.Pt
			if !ScreenToClient(targetHWND, &clientPoint) {
				return CallNextHookEx(hHook, nCode, wParam, lParam)
			}

			//fmt.Printf("Mouse in window: screen(%d,%d) -> client(%d,%d)\n", hookStruct.Pt.X, hookStruct.Pt.Y, clientPoint.X, clientPoint.Y)
			//fmt.Printf("Client Rect")

			// 检查鼠标是否在窗口客户端区域内
			if PtInRect(&clientRect, clientPoint) {
				// 构造lParam (低16位是x坐标，高16位是y坐标)
				lParamNew := LPARAM((uint32(clientPoint.Y) << 16) | (uint32(clientPoint.X) & 0xFFFF))
				//lParamNew := MAKELPARAM(uint16(clientPoint.X), uint16(clientPoint.Y))
				// 调试信息
				//fmt.Printf("Mouse in window: screen(%d,%d) -> client(%d,%d), posting to hwnd=0x%x\n",
				//	hookStruct.Pt.X, hookStruct.Pt.Y, clientPoint.X, clientPoint.Y, targetHWND)

				// 发送消息到窗口
				PostMessage(targetHWND, msg, 0, lParamNew)

			}
		} else {
			//fmt.Println("Not Forwarding.........")
			//fmt.Printf("Mouse out of window: screen(%d,%d)\n", hookStruct.Pt.X, hookStruct.Pt.Y)
		}
	}

	r := CallNextHookEx(hHook, nCode, wParam, lParam)

	return r
}

func InstallMouseHook(hwnd unsafe.Pointer) bool {
	if hHook != 0 {
		fmt.Println("Hook already installed")
		return false
	}

	targetHWND = HWND(hwnd)
	hMod := GetModuleHandle(nil)

	fmt.Printf("Installing mouse hook: hwnd=0x%x, hMod=0x%x\n", targetHWND, hMod)

	// 关键修复：使用当前模块句柄，而不是0
	hHook = SetWindowsHookEx(WH_MOUSE_LL, syscall.NewCallback(mouseHookProc), hMod, 0)

	if hHook != 0 {
		fmt.Printf("Mouse hook installed successfully: hHook=0x%x\n", hHook)

		fmt.Println("Hook is now active, try moving your mouse...")
		return true
	} else {
		fmt.Println("Failed to install mouse hook")
		// 获取错误信息
		return false
	}
}

// 简化版本：不使用消息循环，让Wails应用的消息循环处理
func InstallMouseHookSimple(hwnd unsafe.Pointer) bool {
	if hHook != 0 {
		fmt.Println("Hook already installed")
		return false
	}

	targetHWND = HWND(hwnd)
	hMod := GetModuleHandle(nil)

	fmt.Printf("Installing mouse hook (simple): hwnd=0x%x, hMod=0x%x\n", targetHWND, hMod)

	// 使用原来的逻辑，不检查跨进程

	// 安装窗口子类化来拦截 WM_MOUSELEAVE 消息
	if !subclassInstalled {
		// 检查 comctl32.dll 是否可用
		fmt.Printf("Checking comctl32.dll availability...\n")
		if err := comctl32.Load(); err != nil {
			fmt.Printf("Failed to load comctl32.dll: %v\n", err)
		} else {
			fmt.Printf("comctl32.dll loaded successfully\n")
		}

		// 检查窗口句柄是否有效
		className := GetWindowClassName(targetHWND)
		windowText := GetWindowText(targetHWND)
		fmt.Printf("Target window info: hwnd=0x%x, class='%s', text='%s'\n", targetHWND, className, windowText)

		// 检查窗口是否存在
		var clientRect RECT
		if GetClientRect(targetHWND, &clientRect) {
			fmt.Printf("Window client rect: left=%d, top=%d, right=%d, bottom=%d\n",
				clientRect.Left, clientRect.Top, clientRect.Right, clientRect.Bottom)
		} else {
			fmt.Printf("Failed to get client rect for window 0x%x\n", targetHWND)
		}

		success := InjectCrossProcessDLL(targetHWND)

		//success := SetWindowSubclass(targetHWND, syscall.NewCallback(subclassProc), 111, 0)
		if success {
			subclassInstalled = true
			ControlDLLMouseLeave(targetHWND, false)
			fmt.Printf("Window subclass installed successfully for hwnd=0x%x\n", targetHWND)
		} else {
			fmt.Printf("Failed to install window subclass for hwnd=0x%x\n", targetHWND)
		}
	}

	// 使用当前模块句柄
	hHook = SetWindowsHookEx(WH_MOUSE_LL, syscall.NewCallback(mouseHookProc), hMod, 0)

	if hHook != 0 {
		fmt.Printf("Mouse hook installed successfully: hHook=0x%x\n", hHook)
		fmt.Println("Hook is now active, try moving your mouse...")
		return true
	} else {
		fmt.Println("Failed to install mouse hook")
		err := syscall.GetLastError()
		fmt.Printf("Last error: %d\n", err)
		return false
	}
}

func UninstallMouseHook() bool {
	if hHook == 0 {
		return false
	}

	result := UnhookWindowsHookEx(hHook)
	if result {
		hHook = 0
		targetHWND = 0
		fmt.Println("Mouse hook uninstalled successfully")
	}

	ControlDLLMouseLeave(targetHWND, true)
	return result
}

func IsHookInstalled() bool {
	return hHook != 0
}

func SetForwardToWindow(forward bool) {
	forwardToWindow = forward
	if IsHookInstalled() && targetHWND != 0 {
		ControlDLLMouseLeave(targetHWND, !forward)
	}
	if !forward {
		fmt.Printf("Not forwarding to window\n")
	} else {
		fmt.Printf("Forwarding to window\n")
	}
}

// 发送安装hook的自定义消息
func SendInstallHookMessage(hwnd unsafe.Pointer) bool {
	fmt.Printf("Sending install hook message to hwnd=0x%x, message=0x%x\n", hwnd, WM_INSTALL_MOUSE_HOOK)
	result := PostMessage(HWND(hwnd), WM_INSTALL_MOUSE_HOOK, 0, 0)
	fmt.Printf("PostMessage result: %v\n", result)
	if !result {
		err := syscall.GetLastError()
		fmt.Printf("PostMessage failed with error: %d\n", err)
	}
	mainHWND = HWND(hwnd)
	return result
}

// 发送卸载hook的自定义消息
func SendUninstallHookMessage(hwnd unsafe.Pointer) bool {
	fmt.Printf("Sending uninstall hook message to hwnd=0x%x, message=0x%x\n", hwnd, WM_UNINSTALL_MOUSE_HOOK)
	result := PostMessage(HWND(hwnd), WM_UNINSTALL_MOUSE_HOOK, 0, 0)
	fmt.Printf("PostMessage result: %v\n", result)
	if !result {
		err := syscall.GetLastError()
		fmt.Printf("PostMessage failed with error: %d\n", err)
	}
	return result
}

// 使用 SendMessage 同步发送安装hook消息
func SendInstallHookMessageSync(hwnd unsafe.Pointer) bool {
	fmt.Printf("Sending install hook message (sync) to hwnd=0x%x, message=0x%x\n", hwnd, WM_INSTALL_MOUSE_HOOK)
	result := SendMessage(HWND(hwnd), WM_INSTALL_MOUSE_HOOK, 0, 0)
	fmt.Printf("SendMessage result: 0x%x\n", result)
	return result != 0
}

// 使用 SendMessage 同步发送卸载hook消息
func SendUninstallHookMessageSync(hwnd unsafe.Pointer) bool {
	fmt.Printf("Sending uninstall hook message (sync) to hwnd=0x%x, message=0x%x\n", hwnd, WM_UNINSTALL_MOUSE_HOOK)
	result := SendMessage(HWND(hwnd), WM_UNINSTALL_MOUSE_HOOK, 0, 0)
	fmt.Printf("SendMessage result: 0x%x\n", result)
	return result != 0
}

// 获取窗口的第一个子窗口句柄
func GetFirstChildWindow(hwnd HWND) HWND {
	ret, _, _ := procGetWindow.Call(uintptr(hwnd), 5) // GW_CHILD = 5
	return HWND(ret)
}

// 获取下一个兄弟窗口句柄
func GetNextWindow(hwnd HWND) HWND {
	ret, _, _ := procGetWindow.Call(uintptr(hwnd), 2) // GW_HWNDNEXT = 2
	return HWND(ret)
}

// 获取窗口的类名
func GetWindowClassName(hwnd HWND) string {
	var className [256]uint16
	ret, _, _ := procGetClassName.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&className[0])), 256)
	if ret > 0 {
		return syscall.UTF16ToString(className[:ret])
	}
	return ""
}

// 获取窗口的标题
func GetWindowText(hwnd HWND) string {
	var text [256]uint16
	ret, _, _ := procGetWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&text[0])), 256)
	if ret > 0 {
		return syscall.UTF16ToString(text[:ret])
	}
	return ""
}

// 递归查找 WebView 子窗口
func FindWebViewChild(hwnd HWND) HWND {
	return findWebViewChildRecursive(hwnd, 0)
}

// 递归查找 WebView 子窗口的辅助函数
func findWebViewChildRecursive(hwnd HWND, depth int) HWND {
	if depth > 10 { // 防止无限递归
		fmt.Printf("Max depth reached, stopping search\n")
		return 0
	}

	child := GetFirstChildWindow(hwnd)
	if child == 0 {
		if depth == 0 {
			fmt.Printf("No child window found for hwnd=0x%x\n", hwnd)
		}
		return 0
	}

	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	// 遍历所有兄弟窗口
	current := child
	for current != 0 {
		className := GetWindowClassName(current)
		windowText := GetWindowText(current)
		fmt.Printf("%sDepth %d: Found child window: hwnd=0x%x, class='%s', text='%s'\n", indent, depth, current, className, windowText)

		// 查找 WebView 相关的子窗口
		// 常见的 WebView 类名包括: "Chrome_WidgetWin_1", "WebView2", "CefBrowserWindow" 等
		if className == "Chrome_RenderWidgetHostHWND" {
			fmt.Printf("%sFound WebView child window: hwnd=0x%x, class='%s'\n", indent, current, className)
			return current
		}

		// 递归查找下一层
		deepChild := findWebViewChildRecursive(current, depth+1)
		if deepChild != 0 {
			return deepChild
		}

		// 获取下一个兄弟窗口
		current = GetNextWindow(current)
	}

	// 如果所有层都没有找到 WebView 窗口，返回第一个子窗口
	if depth == 0 {
		fmt.Printf("Using first child window as target: hwnd=0x%x\n", child)
	}
	return child
}

// 测试hook是否工作
func TestHook() {
	fmt.Printf("Hook status: installed=%v, hHook=0x%x, targetHWND=0x%x, messageLoop=%v\n",
		IsHookInstalled(), hHook, targetHWND, messageLoopRunning)
}

// 获取窗口所属的进程ID和线程ID
func GetWindowThreadProcessId(hwnd HWND) (threadId, processId uint32) {
	ret, _, _ := procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processId)))
	threadId = uint32(ret)
	return
}

// 打开进程句柄
func OpenProcess(processId uint32) HANDLE {
	ret, _, _ := procOpenProcess.Call(
		uintptr(PROCESS_ALL_ACCESS),
		0, // bInheritHandle
		uintptr(processId),
	)
	return HANDLE(ret)
}

// 关闭句柄
func CloseHandle(handle HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(handle))
	return ret != 0
}

// 在远程进程中分配内存
func VirtualAllocEx(hProcess HANDLE, size uintptr) uintptr {
	ret, _, _ := procVirtualAllocEx.Call(
		uintptr(hProcess),
		0, // lpAddress
		size,
		uintptr(MEM_COMMIT|MEM_RESERVE),
		0x40, // PAGE_EXECUTE_READWRITE
	)
	return ret
}

// 释放远程进程中的内存
func VirtualFreeEx(hProcess HANDLE, lpAddress uintptr) bool {
	ret, _, _ := procVirtualFreeEx.Call(
		uintptr(hProcess),
		lpAddress,
		0, // dwSize
		uintptr(MEM_RELEASE),
	)
	return ret != 0
}

// 写入远程进程内存
func WriteProcessMemory(hProcess HANDLE, lpBaseAddress uintptr, lpBuffer []byte) bool {
	var bytesWritten uintptr
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&lpBuffer[0])),
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	return ret != 0 && bytesWritten == uintptr(len(lpBuffer))
}

// 读取远程进程内存
func ReadProcessMemory(hProcess HANDLE, lpBaseAddress uintptr, size uintptr) ([]byte, bool) {
	buffer := make([]byte, size)
	var bytesRead uintptr
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&buffer[0])),
		size,
		uintptr(unsafe.Pointer(&bytesRead)),
	)
	return buffer, ret != 0 && bytesRead == size
}

// 读取远程进程内存到现有缓冲区
func ReadProcessMemoryToBuffer(hProcess HANDLE, lpBaseAddress uintptr, buffer []byte) bool {
	var bytesRead uintptr
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&bytesRead)),
	)
	return ret != 0 && bytesRead == uintptr(len(buffer))
}

// 创建远程线程
func CreateRemoteThread(hProcess HANDLE, lpStartAddress uintptr, lpParameter uintptr) HANDLE {
	ret, _, _ := procCreateRemoteThread.Call(
		uintptr(hProcess),
		0, // lpThreadAttributes
		0, // dwStackSize
		lpStartAddress,
		lpParameter,
		0, // dwCreationFlags
		0, // lpThreadId
	)
	return HANDLE(ret)
}

// 等待单个对象
func WaitForSingleObject(hHandle HANDLE, dwMilliseconds uint32) uint32 {
	ret, _, _ := procWaitForSingleObject.Call(uintptr(hHandle), uintptr(dwMilliseconds))
	return uint32(ret)
}

// 获取线程退出代码
func GetExitCodeThread(hThread HANDLE) (uint32, bool) {
	var exitCode uint32
	ret, _, _ := procGetExitCodeThread.Call(uintptr(hThread), uintptr(unsafe.Pointer(&exitCode)))
	return exitCode, ret != 0
}

// 获取当前进程ID
func GetCurrentProcessId() uint32 {
	procGetCurrentProcessId := kernel32.NewProc("GetCurrentProcessId")
	ret, _, _ := procGetCurrentProcessId.Call()
	return uint32(ret)
}

// 跨进程DLL注入并安装窗口子类化
func InjectCrossProcessDLL(hwnd HWND) bool {
	_, processId := GetWindowThreadProcessId(hwnd)
	fmt.Printf("Injecting DLL for hwnd=0x%x, process=%d\n", hwnd, processId)

	// 设置环境变量传递hwnd
	hwndStr := fmt.Sprintf("0x%x", hwnd)
	os.Setenv("PIG_ASSISTANT_HWND", hwndStr)
	fmt.Printf("Set environment variable PIG_ASSISTANT_HWND=%s\n", hwndStr)

	// 同时创建共享内存传递hwnd
	createSharedMemory(hwnd, processId)

	// 打开目标进程
	hProcess := OpenProcess(processId)
	if hProcess == 0 {
		fmt.Printf("Failed to open process %d\n", processId)
		return false
	}
	defer CloseHandle(hProcess)

	// 获取DLL文件路径
	dllPath, err := GetDLLPath()
	if err != nil {
		fmt.Printf("Failed to get DLL path: %v\n", err)
		return false
	}

	fmt.Printf("Using DLL: %s\n", dllPath)

	// 执行DLL注入
	success := InjectDLL(hProcess, dllPath)
	if success {
		fmt.Printf("DLL injection successful!\n")

		// 等待DLL加载完成并自动安装子类化
		time.Sleep(200 * time.Millisecond)

		fmt.Printf("DLL automatically installed subclassing in target process\n")
		return true
	} else {
		fmt.Printf("DLL injection failed\n")
		return false
	}
}

// ==================== DLL注入相关函数 ====================

// 执行DLL注入
func InjectDLL(hProcess HANDLE, dllPath string) bool {
	fmt.Printf("Injecting DLL: %s\n", dllPath)

	// 转换DLL路径为UTF16
	_, err := syscall.UTF16PtrFromString(dllPath)
	if err != nil {
		fmt.Printf("Failed to convert DLL path: %v\n", err)
		return false
	}

	// 在目标进程中分配内存
	dllPathSize := (len(dllPath) + 1) * 2 // UTF16字符大小
	remoteMemory := VirtualAllocEx(hProcess, uintptr(dllPathSize))
	if remoteMemory == 0 {
		fmt.Printf("Failed to allocate memory in target process\n")
		return false
	}
	defer VirtualFreeEx(hProcess, remoteMemory)

	// 写入DLL路径到目标进程
	pathBytes := syscall.StringToUTF16(dllPath)
	// 转换[]uint16为[]byte
	pathBytesAsByte := make([]byte, len(pathBytes)*2)
	for i, v := range pathBytes {
		pathBytesAsByte[i*2] = byte(v & 0xFF)
		pathBytesAsByte[i*2+1] = byte((v >> 8) & 0xFF)
	}
	if !WriteProcessMemory(hProcess, remoteMemory, pathBytesAsByte) {
		fmt.Printf("Failed to write DLL path to target process\n")
		return false
	}

	// 获取LoadLibraryW的地址
	hKernel32 := GetModuleHandle(syscall.StringToUTF16Ptr("kernel32.dll"))
	if hKernel32 == 0 {
		fmt.Printf("Failed to get kernel32.dll handle\n")
		return false
	}

	loadLibraryAddr, _, _ := procGetProcAddress.Call(hKernel32, uintptr(unsafe.Pointer(syscall.StringBytePtr("LoadLibraryW"))))
	if loadLibraryAddr == 0 {
		fmt.Printf("Failed to get LoadLibraryW address\n")
		return false
	}

	// 创建远程线程执行LoadLibraryW
	hThread := CreateRemoteThread(hProcess, loadLibraryAddr, remoteMemory)
	if hThread == 0 {
		fmt.Printf("Failed to create remote thread\n")
		return false
	}
	defer CloseHandle(hThread)

	// 等待线程完成
	waitResult := WaitForSingleObject(hThread, 5000) // 5秒超时
	if waitResult != WAIT_OBJECT_0 {
		fmt.Printf("Remote thread did not complete in time\n")
		return false
	}

	// 获取线程退出代码（LoadLibraryW的返回值）
	exitCode, success := GetExitCodeThread(hThread)
	if !success {
		fmt.Printf("Failed to get thread exit code\n")
		return false
	}

	if exitCode == 0 {
		fmt.Printf("LoadLibraryW failed in target process\n")
		return false
	}

	fmt.Printf("DLL loaded successfully, module handle: 0x%x\n", exitCode)
	return true
}

// 获取DLL文件路径
func GetDLLPath() (string, error) {
	// 首先尝试在当前目录查找预编译的DLL
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// 尝试多个可能的DLL路径
	possiblePaths := []string{
		filepath.Join(currentDir, "hook_dll.dll"),
		filepath.Join(currentDir, "internal", "mousehook", "hook_dll.dll"),
		filepath.Join(currentDir, "..", "internal", "mousehook", "hook_dll.dll"),
		filepath.Join(currentDir, "..", "..", "internal", "mousehook", "hook_dll.dll"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Found DLL at: %s\n", path)
			return path, nil
		}
	}

	// 如果没有找到预编译的DLL，返回错误而不是创建临时DLL
	return "", fmt.Errorf("pre-compiled DLL not found in any of the expected locations: %v", possiblePaths)
}

// 加载库
func LoadLibrary(lpLibFileName string) uintptr {
	procLoadLibrary := kernel32.NewProc("LoadLibraryW")
	libName, err := syscall.UTF16PtrFromString(lpLibFileName)
	if err != nil {
		return 0
	}
	ret, _, _ := procLoadLibrary.Call(uintptr(unsafe.Pointer(libName)))
	return ret
}

// 释放库
func FreeLibrary(hLibModule uintptr) bool {
	procFreeLibrary := kernel32.NewProc("FreeLibrary")
	ret, _, _ := procFreeLibrary.Call(hLibModule)
	return ret != 0
}

// 共享内存结构
type SharedData struct {
	TargetHwnd HWND
	ProcessId  uint32
	IsValid    bool
}

// 全局变量保存共享内存句柄
var g_SharedMemoryHandle HANDLE
var g_SharedData *SharedData

// 创建共享内存传递hwnd
func createSharedMemory(hwnd HWND, processId uint32) {
	// 创建共享内存
	procCreateFileMapping := kernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile := kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile := kernel32.NewProc("UnmapViewOfFile")
	procCloseHandle := kernel32.NewProc("CloseHandle")

	// 创建文件映射对象
	hMapping, _, _ := procCreateFileMapping.Call(
		^uintptr(0), // INVALID_HANDLE_VALUE
		0,           // lpFileMappingAttributes
		uintptr(PAGE_READWRITE),
		0,                                    // dwMaximumSizeHigh
		uintptr(unsafe.Sizeof(SharedData{})), // dwMaximumSizeLow
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("PigAssistantSharedMemory"))),
	)

	if hMapping == 0 {
		fmt.Printf("Failed to create shared memory mapping\n")
		return
	}

	// 保存句柄，不要立即关闭
	g_SharedMemoryHandle = HANDLE(hMapping)

	// 映射共享内存
	pSharedData, _, _ := procMapViewOfFile.Call(
		hMapping,
		uintptr(0x0002), // FILE_MAP_WRITE
		0, 0, 0,
	)
	if pSharedData == 0 {
		fmt.Printf("Failed to map shared memory view\n")
		procCloseHandle.Call(hMapping)
		return
	}

	// 保存映射指针
	g_SharedData = (*SharedData)(unsafe.Pointer(pSharedData))

	// 写入数据到共享内存
	g_SharedData.TargetHwnd = hwnd
	g_SharedData.ProcessId = processId
	g_SharedData.IsValid = true

	fmt.Printf("Created shared memory with hwnd=0x%x, processId=%d\n", hwnd, processId)

	// 等待DLL加载完成后再清理
	go func() {
		time.Sleep(5 * time.Second) // 等待5秒让DLL读取
		if g_SharedData != nil {
			procUnmapViewOfFile.Call(uintptr(unsafe.Pointer(g_SharedData)))
			g_SharedData = nil
		}
		if g_SharedMemoryHandle != 0 {
			procCloseHandle.Call(uintptr(g_SharedMemoryHandle))
			g_SharedMemoryHandle = 0
		}
		fmt.Printf("Cleaned up shared memory\n")
	}()
}

// 控制DLL中的mouseleave行为
func ControlDLLMouseLeave(hwnd HWND, enable bool) bool {
	var message uint32
	if enable {
		message = WM_DLL_ENABLE_MOUSE_LEAVE
		fmt.Printf("Enabling mouseleave events for hwnd=0x%x\n", hwnd)
	} else {
		message = WM_DLL_DISABLE_MOUSE_LEAVE
		fmt.Printf("Disabling mouseleave events for hwnd=0x%x\n", hwnd)
	}

	// 发送消息到窗口
	result := PostMessage(hwnd, message, 0, 0)
	if result {
		fmt.Printf("Successfully sent message 0x%x to hwnd=0x%x\n", message, hwnd)
	} else {
		fmt.Printf("Failed to send message 0x%x to hwnd=0x%x\n", message, hwnd)
	}

	return result
}

// 卸载DLL中的hook
func UninstallDLLHook(hwnd HWND) bool {
	fmt.Printf("Uninstalling DLL hook for hwnd=0x%x\n", hwnd)

	result := PostMessage(hwnd, WM_DLL_UNINSTALL_HOOK, 0, 0)
	if result {
		fmt.Printf("Successfully sent uninstall message to hwnd=0x%x\n", hwnd)
	} else {
		fmt.Printf("Failed to send uninstall message to hwnd=0x%x\n", hwnd)
	}

	return result
}

// GetFocus 获取当前具有键盘焦点的窗口句柄
func GetFocus() HWND {
	ret, _, _ := procGetFocus.Call()
	return HWND(ret)
}

// GetForegroundWindow 获取当前前台窗口句柄
func GetForegroundWindow() HWND {
	ret, _, _ := procGetForegroundWindow.Call()
	return HWND(ret)
}

// SetFocus 设置窗口焦点，获取键盘输入焦点
func SetFocus(hwnd HWND) bool {
	if hwnd == 0 {
		fmt.Printf("Invalid window handle\n")
		return false
	}

	fmt.Printf("Setting focus for hwnd=0x%x\n", hwnd)

	// 检查当前焦点窗口
	currentFocus := GetFocus()
	currentForeground := GetForegroundWindow()

	fmt.Printf("Current focus: 0x%x, current foreground: 0x%x\n", currentFocus, currentForeground)

	// 如果目标窗口已经是焦点窗口，直接返回成功
	if currentFocus == hwnd && currentForeground == hwnd {
		fmt.Printf("Window 0x%x already has focus\n", hwnd)
		return true
	}

	// 检查窗口是否最小化
	isIconic, _, _ := procIsIconic.Call(uintptr(hwnd))
	if isIconic != 0 {
		fmt.Printf("Window is minimized, restoring it first\n")
		// 恢复窗口
		procShowWindow.Call(uintptr(hwnd), SW_RESTORE)
		time.Sleep(100 * time.Millisecond) // 等待窗口恢复
	}

	// 1. 将窗口置于前台
	result1, _, err := procSetForegroundWindow.Call(uintptr(hwnd))
	if result1 == 0 {
		fmt.Printf("Failed to set foreground window:%v\n", err)
	}

	// 2. 激活窗口
	result2, _, err := procSetActiveWindow.Call(uintptr(hwnd))
	if result2 == 0 {
		fmt.Printf("Failed to set active window:%v\n", err)
	}

	// 3. 将窗口置于顶层
	result3, _, err := procBringWindowToTop.Call(uintptr(hwnd))
	if result3 == 0 {
		fmt.Printf("Failed to bring window to top:%v \n", err)
	}

	// 4. 设置键盘焦点
	result4, _, err := procSetFocus.Call(uintptr(hwnd))
	if result4 == 0 {
		fmt.Printf("Failed to set focus\n")
		fmt.Printf("SetFocus error: %v \n", err)
		return false
	}

	// 验证焦点是否设置成功
	newFocus := GetFocus()
	if newFocus == hwnd {
		fmt.Printf("Successfully set focus for hwnd=0x%x\n", hwnd)
		return true
	} else {
		fmt.Printf("Focus verification failed: expected 0x%x, got 0x%x\n", hwnd, newFocus)
		return false
	}
}

// GetAsyncKeyState 获取指定虚拟键的异步状态
func GetAsyncKeyState(vKey int32) int16 {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return int16(ret)
}

// IsLeftMouseButtonDown 检测鼠标左键是否按下
func IsLeftMouseButtonDown() bool {
	keyState := GetAsyncKeyState(VK_LBUTTON)
	// GetAsyncKeyState 返回值的最高位表示键是否被按下
	// 使用 uint16 来避免溢出，然后转换为 int16
	return (uint16(keyState) & 0x8000) != 0
}

// GetWindowLongPtr 获取窗口的扩展信息
func GetWindowLongPtr(hwnd HWND, nIndex int32) uintptr {
	ret, _, _ := procGetWindowLongPtr.Call(uintptr(hwnd), uintptr(nIndex))
	return ret
}

// SetWindowLongPtr 设置窗口的扩展信息
func SetWindowLongPtr(hwnd HWND, nIndex int32, dwNewLong uintptr) uintptr {
	ret, _, _ := procSetWindowLongPtr.Call(uintptr(hwnd), uintptr(nIndex), dwNewLong)
	return ret
}

// SetWindowPos 设置窗口位置和大小
func SetWindowPos(hwnd HWND, hWndInsertAfter HWND, x, y, cx, cy int32, uFlags uint32) bool {
	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags),
	)
	return ret != 0
}

// SetWindowBorderless 设置指定窗口为无边框、无标题栏、无菜单
func SetWindowBorderless(hwnd HWND) bool {
	if hwnd == 0 {
		fmt.Printf("Invalid window handle\n")
		return false
	}

	fmt.Printf("Setting window 0x%x to borderless\n", hwnd)

	// 获取当前窗口样式
	currentStyle := GetWindowLongPtr(hwnd, GWLP_STYLE)
	fmt.Printf("Current window style: 0x%x\n", currentStyle)

	// 移除边框、标题栏、菜单等样式，并添加WS_POPUP样式
	newStyle := (currentStyle &^ (WS_CAPTION | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_SYSMENU | WS_BORDER | WS_DLGFRAME)) | WS_POPUP
	fmt.Printf("New window style: 0x%x\n", newStyle)

	// 设置新的窗口样式
	result := SetWindowLongPtr(hwnd, GWLP_STYLE, newStyle)
	if result == 0 {
		fmt.Printf("Failed to set window style\n")
		return false
	}

	// 强制窗口重绘边框
	//	SetWindowPos(hwnd, 0, 0, 0, 0, 0, SWP_FRAMECHANGED|SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_SHOWWINDOW)

	fmt.Printf("Successfully set window 0x%x to borderless\n", hwnd)
	return true
}
