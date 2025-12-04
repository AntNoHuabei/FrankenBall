# DLL注入功能使用说明

## 概述

本模块实现了完整的Go版本DLL注入功能，用于跨进程窗口消息拦截。虽然鼠标hook不需要跨进程，但DLL注入功能被保留以供其他用途。

## 文件说明

- `hook.go` - 主要的Go实现文件，包含DLL注入逻辑
- `hook_dll.c` - C源代码文件，需要编译成DLL
- `build_dll.bat` - 编译脚本

## 编译DLL

**注意**：由于使用了 `SetWindowSubclass` 等函数，需要链接 `comctl32` 库。

### 方法1：使用MinGW-w64
```bash
gcc -shared -o hook_dll.dll hook_dll.c -lkernel32 -luser32 -lcomctl32
```

### 方法2：使用Visual Studio
```cmd
cl /LD hook_dll.c /Fe:hook_dll.dll /link kernel32.lib user32.lib comctl32.lib
```

### 方法3：使用提供的批处理文件
```cmd
build_dll.bat
```

## 主要功能

### 1. DLL注入
- `InjectDLLAndSubclass()` - 主要的DLL注入函数
- `InjectDLL()` - 执行DLL注入
- `CreateTempDLL()` - 创建临时DLL文件

### 2. 进程操作
- `OpenProcess()` - 打开目标进程
- `VirtualAllocEx()` - 在远程进程中分配内存
- `WriteProcessMemory()` - 写入远程进程内存
- `CreateRemoteThread()` - 创建远程线程

### 3. 模块枚举
- `FindInjectedDLL()` - 查找注入的DLL模块
- `GetRemoteProcAddress()` - 获取远程函数地址
- `CreateToolhelp32Snapshot()` - 创建模块快照

### 4. PE文件解析
- 完整的PE头解析
- 导出表解析
- 函数地址计算

## 使用示例

```go
// 注入DLL到目标窗口的进程
hwnd := mousehook.HWND(windowHandle)
success := mousehook.InjectDLLAndSubclass(hwnd)
if success {
    fmt.Println("DLL注入成功！")
} else {
    fmt.Println("DLL注入失败！")
}
```

## 注意事项

1. **权限要求**：DLL注入需要管理员权限
2. **目标进程**：确保目标进程允许注入
3. **DLL文件**：需要预先编译好DLL文件
4. **错误处理**：所有函数都有详细的错误日志

## 技术细节

### DLL注入流程
1. 打开目标进程
2. 在目标进程中分配内存
3. 写入DLL路径
4. 创建远程线程执行LoadLibraryW
5. 等待DLL加载完成
6. 查找DLL模块基址
7. 获取函数地址
8. 调用DLL中的函数

### PE文件解析
- 读取DOS头验证签名
- 解析PE头获取导出表RVA
- 遍历导出表查找目标函数
- 计算函数实际地址

## 限制和已知问题

1. Go的回调函数无法跨进程执行
2. 某些系统可能阻止DLL注入
3. 需要目标进程的完整访问权限
4. 防病毒软件可能误报

## 扩展建议

1. 添加更多PE文件格式支持
2. 实现更复杂的错误处理
3. 添加注入状态监控
4. 支持64位进程注入
