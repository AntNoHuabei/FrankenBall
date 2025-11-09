package main

import (
	"embed"
	"fmt"
	"log"
	"unsafe"

	"github.com/AntNoHuabei/Remo/internal/mousehook"
	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "Remo",
		Description: "一个基于Wails的AI悬浮球应用",
		Services: []application.Service{
			application.NewService(&MouseEventService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			WndProcInterceptor: func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool) {

				// 处理自定义消息
				switch msg {
				//
				//case mousehook.WM_CREATE:
				//	log.Printf("Received WM_CREATE message for window 0x%x", hwnd)
				//	// 窗体创建时设置为无边框
				//	success := mousehook.SetWindowBorderless(mousehook.HWND(hwnd))
				//	if success {
				//		log.Printf("Successfully set window 0x%x to borderless on create", hwnd)
				//	} else {
				//		log.Printf("Failed to set window 0x%x to borderless on create", hwnd)
				//	}
				//
				//case mousehook.WM_NCCREATE:
				//	log.Printf("Received WM_NCCREATE message for window 0x%x", hwnd)
				//	// 非客户端区域创建时也设置为无边框
				//	success := mousehook.SetWindowBorderless(mousehook.HWND(hwnd))
				//	if success {
				//		log.Printf("Successfully set window 0x%x to borderless on nc create", hwnd)
				//	} else {
				//		log.Printf("Failed to set window 0x%x to borderless on nc create", hwnd)
				//	}

				case mousehook.WM_DESTROY:
					mousehook.UninstallMouseHook()

				//case mousehook.WM_SET_FOCUS:
				//	fmt.Println("Received WM_SET_FOCUS message")
				//	w, _ := application.Get().Window.Get("WinMain")
				//	if !w.IsIgnoreMouseEvents() {
				//
				//		//	mousehook.SetFocus(mousehook.HWND(hwnd))
				//	}

				case mousehook.WM_INSTALL_MOUSE_HOOK:
					log.Printf("Received WM_INSTALL_MOUSE_HOOK message (0x%x)", msg)

					// 查找 WebView 子窗口
					webViewHwnd := mousehook.FindWebViewChild(mousehook.HWND(hwnd))
					if webViewHwnd == 0 {
						log.Println("Failed to find WebView child window, using main window")
						webViewHwnd = mousehook.HWND(hwnd)
					}

					success := mousehook.InstallMouseHookSimple(unsafe.Pointer(webViewHwnd))
					if success {
						log.Printf("Mouse hook installed successfully in message pump thread, target hwnd=0x%x", webViewHwnd)
						mousehook.TestHook()
					} else {
						log.Println("Failed to install mouse hook in message pump thread")
					}
					return 0, true // 消息已处理，不继续传递

				case mousehook.WM_UNINSTALL_MOUSE_HOOK:
					log.Printf("Received WM_UNINSTALL_MOUSE_HOOK message (0x%x)", msg)
					success := mousehook.UninstallMouseHook()
					if success {
						log.Println("Mouse hook uninstalled successfully in message pump thread")
					} else {
						log.Println("Failed to uninstall mouse hook in message pump thread")
					}
					return 0, true // 消息已处理，不继续传递

					//default:
					//fmt.Printf("Unhandled WndProcInterceptor message: hwnd=0x%x, msg=0x%x, wParam=0x%x, lParam=0x%x", hwnd, msg, wParam, lParam)
				}

				return 0, false
			},
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.

	app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
		ps := application.Get().Screen.GetPrimary()

		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Name: "WinMain",
			Windows: application.WindowsWindow{
				DisableFramelessWindowDecorations: true,
				HiddenOnTaskbar:                   false,
				//ExStyle:                           w32.WS_EX_LAYERED | w32.WS_POPUP,
			},
			BackgroundColour: application.NewRGBA(0, 0, 0, 0),
			URL:              "/",
			Frameless:        false,
			AlwaysOnTop:      true,
			DevToolsEnabled:  true,
			InitialPosition:  application.WindowCentered,
			BackgroundType:   application.BackgroundTypeTransparent,
			StartState:       application.WindowStateNormal,
			DisableResize:    false,
			//IgnoreMouseEvents: true,
			OpenInspectorOnStartup: true,
			X:                      0,
			Y:                      0,
			Width:                  ps.PhysicalWorkArea.Width,
			Height:                 ps.PhysicalWorkArea.Height,
		})
		w.OnWindowEvent(events.Windows.WindowActive, func(event *application.WindowEvent) {
			fmt.Println("========================")

			success := mousehook.SetWindowBorderless(mousehook.HWND(w.NativeWindow()))
			if success {
				log.Printf("Successfully set window to borderless on create")
			} else {
				log.Printf("Failed to set window  to borderless on create")
			}
			w.SetBounds(application.Rect{
				X:      0,
				Y:      0,
				Width:  ps.WorkArea.Width,
				Height: ps.WorkArea.Height, //避免窗口达到最大化效果导致 焦点无法获取和下层窗口无法重绘
			})
		})
	})

	//w.OnWindowEvent(events.Windows.WindowActive, func(event *application.WindowEvent) {
	//	w.SetIgnoreMouseEvents(true)
	//	mousehook.InstallMouseHookSimple(w.NativeWindow())
	//})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
