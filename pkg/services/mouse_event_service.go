package services

import (
	"fmt"

	"github.com/AntNoHuabei/Remo/internal/mousehook"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type MouseEventService struct {
	webViewHwnd mousehook.HWND
}

func (m *MouseEventService) MousePassthroughWithMove(ignore bool) {

	w, _ := application.Get().Window.Get("WinMain")

	// 设置窗口忽略鼠标事件，实现鼠标穿透
	w.SetIgnoreMouseEvents(ignore)
	fmt.Println("IsIgnoreMouseEvents:", w.IsIgnoreMouseEvents())
	if ignore {

		if !mousehook.IsHookInstalled() {
			//application.InvokeSync(func() {
			//	mousehook.InstallMouseHookSimple(w.NativeWindow())
			//})
			mousehook.SendInstallHookMessage(w.NativeWindow())
		}
		mousehook.SetForwardToWindow(true)

	} else {

		mousehook.SetForwardToWindow(false)
		//mousehook.SendUninstallHookMessage(w.NativeWindow())
	}

}

// SetWindowBorderless 设置指定窗口为无边框、无标题栏、无菜单
func (m *MouseEventService) SetWindowBorderless(hwnd uintptr) bool {
	return mousehook.SetWindowBorderless(mousehook.HWND(hwnd))
}
