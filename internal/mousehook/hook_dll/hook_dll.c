#include <windows.h>
#include <stdio.h>
#include <time.h>
#include <commctrl.h>

// 自定义消息定义
#define WM_INSTALL_MOUSE_HOOK   0x0401
#define WM_UNINSTALL_MOUSE_HOOK 0x0402
#define WM_ENABLE_MOUSE_LEAVE   0x0403
#define WM_DISABLE_MOUSE_LEAVE  0x0404

// 全局变量
HWND g_TargetHwnd = NULL;
BOOL g_IgnoreMouseLeave = FALSE;
HMODULE g_hModule = NULL;
UINT_PTR g_SubclassId = 11; // 子类化ID
DWORD g_TargetThreadId = 0; // 目标窗口的线程ID
WNDPROC g_OriginalWndProc = NULL; // 原始窗口过程（用于SetWindowLongPtr方法）
BOOL g_UsingSetWindowLongPtr = FALSE; // 是否使用SetWindowLongPtr方法

// 共享内存结构
typedef struct {
    HWND targetHwnd;
    DWORD processId;
    BOOL isValid;
} SharedData;

HANDLE g_SharedMemoryHandle = NULL;
SharedData* g_SharedData = NULL;

// 日志文件路径
#define LOG_FILE_PATH "C:\\temp\\hook_log.txt"

// 写入日志的函数
void WriteLog(const char* message) {
    FILE* file = fopen(LOG_FILE_PATH, "a");
    if (file) {
        time_t now;
        time(&now);
        struct tm* timeinfo = localtime(&now);
        
        fprintf(file, "[%04d-%02d-%02d %02d:%02d:%02d] %s\n",
                timeinfo->tm_year + 1900, timeinfo->tm_mon + 1, timeinfo->tm_mday,
                timeinfo->tm_hour, timeinfo->tm_min, timeinfo->tm_sec,
                message);
        fclose(file);
    }
}

// 从共享内存读取目标窗口句柄
BOOL ReadTargetHwndFromSharedMemory() {
    // 尝试打开共享内存
    g_SharedMemoryHandle = OpenFileMappingA(FILE_MAP_READ, FALSE, "PigAssistantSharedMemory");
    if (g_SharedMemoryHandle == NULL) {
        WriteLog("Failed to open shared memory");
        return FALSE;
    }
    
    // 映射共享内存
    g_SharedData = (SharedData*)MapViewOfFile(g_SharedMemoryHandle, FILE_MAP_READ, 0, 0, sizeof(SharedData));
    if (g_SharedData == NULL) {
        WriteLog("Failed to map shared memory");
        CloseHandle(g_SharedMemoryHandle);
        g_SharedMemoryHandle = NULL;
        return FALSE;
    }
    
    // 检查数据是否有效
    if (g_SharedData->isValid && IsWindow(g_SharedData->targetHwnd)) {
        g_TargetHwnd = g_SharedData->targetHwnd;
        char logMsg[256];
        sprintf(logMsg, "Target HWND from shared memory: 0x%p", g_TargetHwnd);
        WriteLog(logMsg);
        return TRUE;
    }
    
    WriteLog("Invalid data in shared memory");
    UnmapViewOfFile(g_SharedData);
    CloseHandle(g_SharedMemoryHandle);
    g_SharedData = NULL;
    g_SharedMemoryHandle = NULL;
    return FALSE;
}

// 简单的窗口过程（用于SetWindowLongPtr方法）
LRESULT CALLBACK SimpleHookWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    // 处理自定义消息
    switch (uMsg) {
    case WM_ENABLE_MOUSE_LEAVE:
        // 启用mouseleave事件
        g_IgnoreMouseLeave = FALSE;
        WriteLog("WM_ENABLE_MOUSE_LEAVE received - Mouse leave events enabled");
        return 0;
        
    case WM_DISABLE_MOUSE_LEAVE:
        // 禁用mouseleave事件
        g_IgnoreMouseLeave = TRUE;
        WriteLog("WM_DISABLE_MOUSE_LEAVE received - Mouse leave events disabled");
        return 0;
        
    case WM_MOUSELEAVE:
        // 根据设置决定是否忽略mouseleave事件
        if (g_IgnoreMouseLeave) {
            return 0; // 忽略mouseleave事件
        }
        break;
    }
    
    // 其他消息传递给原始窗口过程
    if (g_OriginalWndProc) {
        return CallWindowProc(g_OriginalWndProc, hwnd, uMsg, wParam, lParam);
    }
    
    return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

// 我们的子类化窗口过程
LRESULT CALLBACK HookSubclassProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR uIdSubclass, DWORD_PTR dwRefData) {
    // 处理自定义消息
    switch (uMsg) {
    case WM_INSTALL_MOUSE_HOOK:
        // 安装hook的消息，现在主要用于调试
        WriteLog("WM_INSTALL_MOUSE_HOOK received (for debugging)");
        return 0;
        
    case WM_UNINSTALL_MOUSE_HOOK:
        // 卸载hook的消息
        WriteLog("WM_UNINSTALL_MOUSE_HOOK received");
        if (g_TargetHwnd) {
            if (g_UsingSetWindowLongPtr) {
                // 使用SetWindowLongPtr方法，恢复原始窗口过程
                if (g_OriginalWndProc) {
                    SetWindowLongPtr(g_TargetHwnd, GWLP_WNDPROC, (LONG_PTR)g_OriginalWndProc);
                    WriteLog("Window procedure restored using SetWindowLongPtr");
                }
                g_UsingSetWindowLongPtr = FALSE;
                g_OriginalWndProc = NULL;
            } else {
                // 使用SetWindowSubclass方法
                if (!RemoveWindowSubclass(g_TargetHwnd, HookSubclassProc, g_SubclassId)) {
                    DWORD error = GetLastError();
                    char errorMsg[256];
                    sprintf(errorMsg, "Failed to remove window subclass, GetLastError: %lu (0x%lX)", error, error);
                    WriteLog(errorMsg);
                } else {
                    WriteLog("Window subclass removed successfully");
                }
            }
            g_TargetHwnd = NULL;
        }
        return 0;
        
    case WM_ENABLE_MOUSE_LEAVE:
        // 启用mouseleave事件
        g_IgnoreMouseLeave = FALSE;
        WriteLog("WM_ENABLE_MOUSE_LEAVE received - Mouse leave events enabled");
        return 0;
        
    case WM_DISABLE_MOUSE_LEAVE:
        // 禁用mouseleave事件
        g_IgnoreMouseLeave = TRUE;
        WriteLog("WM_DISABLE_MOUSE_LEAVE received - Mouse leave events disabled");
        return 0;
        
    case WM_MOUSELEAVE:
        // 根据设置决定是否忽略mouseleave事件
        if (g_IgnoreMouseLeave) {
            return 0; // 忽略mouseleave事件
        }
        break;
    }
    
    // 其他消息传递给默认窗口过程
    return DefSubclassProc(hwnd, uMsg, wParam, lParam);
}


// 在目标线程中安装子类化的函数
BOOL InstallSubclassInTargetThread() {
    if (!g_TargetHwnd) {
        WriteLog("No target window available for subclassing");
        return FALSE;
    }

    // 检查当前线程是否是目标窗口的线程
    DWORD currentThreadId = GetCurrentThreadId();
    char logMsg[256];
    sprintf(logMsg, "Current thread ID: %lu, Target thread ID: %lu", currentThreadId, g_TargetThreadId);
    WriteLog(logMsg);

    if (currentThreadId == g_TargetThreadId) {
        WriteLog("Already in target thread, installing subclass directly");
        // 已经在正确的线程中，直接安装
        if (!SetWindowSubclass(g_TargetHwnd, HookSubclassProc, g_SubclassId, 0)) {
            DWORD error = GetLastError();
            char errorMsg[256];
            sprintf(errorMsg, "Failed to install window subclass in target thread, GetLastError: %lu (0x%lX)", error, error);
            WriteLog(errorMsg);
            return FALSE;
        }
        WriteLog("Window subclass installed successfully in target thread");
        return TRUE;
    }

    // 不在目标线程中，尝试使用SendMessage同步调用
    WriteLog("Not in target thread, trying to install subclass anyway");

    // 直接尝试安装，看看是否真的需要在线程中
    if (!SetWindowSubclass(g_TargetHwnd, HookSubclassProc, g_SubclassId, 0)) {
        DWORD error = GetLastError();
        char errorMsg[256];
        sprintf(errorMsg, "Failed to install window subclass from different thread, GetLastError: %lu (0x%lX)", error, error);
        WriteLog(errorMsg);

        // 如果失败，尝试使用SetWindowLongPtr作为备选方案
        WriteLog("Trying SetWindowLongPtr as fallback method");
        g_OriginalWndProc = (WNDPROC)GetWindowLongPtr(g_TargetHwnd, GWLP_WNDPROC);
        if (g_OriginalWndProc) {
            if (SetWindowLongPtr(g_TargetHwnd, GWLP_WNDPROC, (LONG_PTR)SimpleHookWndProc)) {
                g_UsingSetWindowLongPtr = TRUE;
                WriteLog("SetWindowLongPtr method succeeded");
                return TRUE;
            } else {
                DWORD error = GetLastError();
                char errorMsg[256];
                sprintf(errorMsg, "SetWindowLongPtr failed, GetLastError: %lu (0x%lX)", error, error);
                WriteLog(errorMsg);
                return FALSE;
            }
        } else {
            WriteLog("Failed to get original window procedure");
            return FALSE;
        }
    } else {
        WriteLog("Window subclass installed successfully from different thread");
        return TRUE;
    }
}

// 递归查找子窗口的函数
HWND FindChromeRenderWidget(HWND parentHwnd, int depth) {
    if (depth > 10) return NULL; // 防止无限递归
    
    HWND childHwnd = GetWindow(parentHwnd, GW_CHILD);
    while (childHwnd) {
        char className[256];
        char windowText[256];
        GetClassNameA(childHwnd, className, sizeof(className));
        GetWindowTextA(childHwnd, windowText, sizeof(windowText));
        
        char indent[32] = "";
        for (int i = 0; i < depth; i++) {
            strcat(indent, "  ");
        }
        
        char logMsg[512];
        sprintf(logMsg, "%sChecking child HWND: 0x%p, Class: '%s', Text: '%s'", 
                indent, childHwnd, className, windowText);
        WriteLog(logMsg);
        
        // 精确匹配 Chrome_RenderWidgetHostHWND
        if (strcmp(className, "Chrome_RenderWidgetHostHWND") == 0) {
            sprintf(logMsg, "%sFound Chrome_RenderWidgetHostHWND: HWND: 0x%p, Class: '%s', Text: '%s'", 
                    indent, childHwnd, className, windowText);
            WriteLog(logMsg);
            return childHwnd;
        }
        
        // 递归查找子窗口的子窗口
        HWND found = FindChromeRenderWidget(childHwnd, depth + 1);
        if (found) {
            return found;
        }
        
        childHwnd = GetWindow(childHwnd, GW_HWNDNEXT);
    }
    
    return NULL;
}

// 查找主窗口的函数
HWND FindMainWindow() {
    WriteLog("Searching for Chrome_RenderWidgetHostHWND window...");
    
    // 先获取前台窗口
    HWND foregroundHwnd = GetForegroundWindow();
    if (foregroundHwnd) {
        char windowText[256];
        char className[256];
        GetWindowTextA(foregroundHwnd, windowText, sizeof(windowText));
        GetClassNameA(foregroundHwnd, className, sizeof(className));
        
        // 获取前台窗口的进程ID和线程ID
        DWORD foregroundProcessId = 0;
        DWORD foregroundThreadId = GetWindowThreadProcessId(foregroundHwnd, &foregroundProcessId);
        
        char logMsg[512];
        sprintf(logMsg, "Found foreground window: HWND: 0x%p, Class: '%s', Text: '%s', ProcessID: %lu (0x%lX)", 
                foregroundHwnd, className, windowText, foregroundProcessId, foregroundProcessId);
        WriteLog(logMsg);
        
        // 在前台窗口的子窗口中查找 Chrome_RenderWidgetHostHWND
        WriteLog("Searching for Chrome_RenderWidgetHostHWND in foreground window children...");
        HWND chromeHwnd = FindChromeRenderWidget(foregroundHwnd, 0);
        if (chromeHwnd) {
            return chromeHwnd;
        }
    }
    
    // 如果前台窗口没有找到，尝试直接查找
    WriteLog("Chrome_RenderWidgetHostHWND not found in foreground window, trying direct search...");
    HWND directHwnd = FindWindowA("Chrome_RenderWidgetHostHWND", NULL);
    if (directHwnd) {
        char windowText[256];
        char className[256];
        GetWindowTextA(directHwnd, windowText, sizeof(windowText));
        GetClassNameA(directHwnd, className, sizeof(className));
        
        // 获取直接查找窗口的进程ID
        DWORD directProcessId = 0;
        GetWindowThreadProcessId(directHwnd, &directProcessId);
        
        char logMsg[512];
        sprintf(logMsg, "Found Chrome_RenderWidgetHostHWND via direct search: HWND: 0x%p, Class: '%s', Text: '%s', ProcessID: %lu (0x%lX)", 
                directHwnd, className, windowText, directProcessId, directProcessId);
        WriteLog(logMsg);
        return directHwnd;
    }
    
    WriteLog("No Chrome_RenderWidgetHostHWND window found");
    return NULL;
}

// 安装窗口子类化的函数
BOOL InstallWindowSubclass() {
    WriteLog("InstallWindowSubclass called");
    
    // 查找主窗口
    if(!g_TargetHwnd) {
        g_TargetHwnd = FindMainWindow();
    }
    
    if (!g_TargetHwnd) {
        DWORD error = GetLastError();
        char errorMsg[256];
        sprintf(errorMsg, "Failed to find main window, GetLastError: %lu (0x%lX)", error, error);
        WriteLog(errorMsg);
        return FALSE;
    }
    
    // 输出目标窗口的详细信息
    char windowText[256];
    char className[256];
    GetWindowTextA(g_TargetHwnd, windowText, sizeof(windowText));
    GetClassNameA(g_TargetHwnd, className, sizeof(className));
    
    // 获取窗口所在进程ID和线程ID
    DWORD processId = 0;
    g_TargetThreadId = GetWindowThreadProcessId(g_TargetHwnd, &processId);
    
    char logMsg[512];
    sprintf(logMsg, "Target window details: HWND: 0x%p, Class: '%s', Text: '%s', ProcessID: %lu (0x%lX), ThreadID: %lu (0x%lX)", 
            g_TargetHwnd, className, windowText, processId, processId, g_TargetThreadId, g_TargetThreadId);
    WriteLog(logMsg);
    
    // 使用线程安全的方法安装子类化
    if (!InstallSubclassInTargetThread()) {
        WriteLog("Failed to install window subclass using thread-safe method");
        return FALSE;
    }
    
    return TRUE;
}

// DLL入口点
BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved) {
    switch (ul_reason_for_call) {
    case DLL_PROCESS_ATTACH:
        g_hModule = hModule;
        WriteLog("DLL_PROCESS_ATTACH - DLL loaded");       
        // 如果环境变量也没有找到，尝试从共享内存获取
        if (!g_TargetHwnd) {
            if (ReadTargetHwndFromSharedMemory()) {
                WriteLog("Successfully read target HWND from shared memory");
            } else {
                WriteLog("Failed to read target HWND from shared memory");
            }
        }
        
        // DLL被加载时自动安装窗口子类化
        InstallWindowSubclass();
        break;
        
    case DLL_THREAD_ATTACH:
        WriteLog("DLL_THREAD_ATTACH");
        break;
        
    case DLL_THREAD_DETACH:
        WriteLog("DLL_THREAD_DETACH");
        break;
        
    case DLL_PROCESS_DETACH:
        WriteLog("DLL_PROCESS_DETACH - DLL unloading");
        // DLL被卸载时移除子类化
        if (g_TargetHwnd) {
            if (g_UsingSetWindowLongPtr) {
                // 使用SetWindowLongPtr方法，恢复原始窗口过程
                if (g_OriginalWndProc) {
                    SetWindowLongPtr(g_TargetHwnd, GWLP_WNDPROC, (LONG_PTR)g_OriginalWndProc);
                    WriteLog("Window procedure restored during DLL unload using SetWindowLongPtr");
                }
            } else {
                // 使用SetWindowSubclass方法
                if (!RemoveWindowSubclass(g_TargetHwnd, HookSubclassProc, g_SubclassId)) {
                    DWORD error = GetLastError();
                    char errorMsg[256];
                    sprintf(errorMsg, "Failed to remove window subclass during DLL unload, GetLastError: %lu (0x%lX)", error, error);
                    WriteLog(errorMsg);
                } else {
                    WriteLog("Window subclass removed during DLL unload");
                }
            }
        }
        break;
    }
    return TRUE;
}
