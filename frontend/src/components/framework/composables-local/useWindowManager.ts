import { ref, reactive, computed, Component, Ref } from "vue";
import { v4 as uuidv4 } from "uuid";
import {WindowConfig, WindowState, WindowType} from "../types/window";

export interface CloseOptions {
  isReplaced?: boolean; // 表示窗体是否因为被其他窗体替换而关闭，默认为 undefined 或 false
  replacement?: WindowState; // 如果是替换，则提供新窗体的配置
}
export interface MinimizeOptions {
  isReplaced?: boolean; // 表示窗体是否因为被其他窗体替换而关闭，默认为 undefined 或 false
  replacement?: WindowState; // 如果是替换，则提供新窗体的配置
}

export interface RestoreOptions {
  isReplacing?: boolean; // 是否是替换操作（新窗体替换了旧窗体）
  replacedWindow?: WindowState; // 如果是替换，这里是被替换的旧窗体配置
}
export interface ResizeOptions {
  targetWidth: number;
  targetHeight: number;
  targetX: number;
  targetY: number;
}

export interface CreateOrRestoreWindowOptions {
  data?: Record<string, any>;
  forceCreate?: boolean; // 是否强制创建新窗口，而不是恢复旧窗口
}

// 全局窗口状态存储
const windows = reactive<Record<string, WindowState>>({});
const activeWindowId = ref<string | null>(null);
const highestZIndex = ref(1000);

export interface Options {
  beforeWindowClose?: (w: WindowState, options: CloseOptions) => Promise<void>;
  afterWindowClose?: (w: WindowState, options: CloseOptions) => Promise<void>;
  beforeWindowMinimize?: (w: WindowState, options: MinimizeOptions) => Promise<void>;
  afterWindowMinimize?: (w: WindowState, options: MinimizeOptions) => Promise<void>;
  beforeWindowMaximize?: (w: WindowState) => Promise<void>;
  afterWindowMaximize?: (w: WindowState) => Promise<void>;
  beforeWindowRestore?: (w: WindowState, options: RestoreOptions) => Promise<void>;
  afterWindowRestore?: (w: WindowState, options: RestoreOptions) => Promise<void>;
  beforeWindowResize?: (w: WindowState, options: ResizeOptions) => Promise<void>;
  afterWindowResize?: (w: WindowState, options: ResizeOptions) => Promise<void>;
}

let defaultOptions: Options = {
  beforeWindowClose: async () => {},
  afterWindowClose: async () => {},
  beforeWindowMinimize: async () => {},
  afterWindowMinimize: async () => {},
  beforeWindowMaximize: async () => {},
  afterWindowMaximize: async () => {},
  beforeWindowRestore: async () => {},
  afterWindowRestore: async () => {},
};

export type WindowResizeObserver = {
  resize(w: WindowState, detail: ResizeOptions): void;
}

let windowResizeObserver: WindowResizeObserver | null = null;

// 窗口管理器
export function useWindowManager(options?: Options) {


  const setWindowResizeObserver = (newObserver: WindowResizeObserver) =>{
    windowResizeObserver = newObserver;
  }

  if (options) {
    options = { ...defaultOptions, ...options };
  } else {
    options = defaultOptions;
  }

  const setDefaultOptions = (newOptions: Options) => {
    defaultOptions = newOptions;
  };

  // 注册窗口
  const createWindow = (windowConfig: WindowConfig): string => {
    const layout = windowConfig.layout || {
      x: undefined,
      y: undefined,
      width: 400,
      height: 300,
      resizable: undefined,
      maxWidth: undefined,
      maxHeight: undefined,
      minWidth: undefined,
      minHeight: undefined,
    };
    //给toolbar默认值
    windowConfig.toolbar = windowConfig.toolbar || {};
    windowConfig.toolbar.show = windowConfig.toolbar.show === undefined || windowConfig.toolbar.show === true;
    windowConfig.toolbar.hoverShow = windowConfig.toolbar.hoverShow === true;
    windowConfig.toolbar.showMinimizeButton =
      windowConfig.toolbar.showMinimizeButton === undefined || windowConfig.toolbar.showMinimizeButton === true;
    windowConfig.toolbar.showMaximizeButton =
      windowConfig.toolbar.showMaximizeButton === undefined || windowConfig.toolbar.showMaximizeButton === true;
    windowConfig.toolbar.showCloseButton =
      windowConfig.toolbar.showCloseButton === undefined || windowConfig.toolbar.showCloseButton === true;

    // 设置默认位置和尺寸
    const position = {
      x: layout.x === undefined ? 100 : layout.x,
      y: layout.y === undefined ? 100 : layout.y,
    };

    const size = {
      width: layout.width || 400,
      height: layout.height || 300,
      maxWidth: layout.maxWidth,
      maxHeight: layout.maxHeight,
      minWidth: layout.minWidth,
      minHeight: layout.minHeight,
    };

    // 设置行为属性
    const behavior = windowConfig.behavior || {};

    const id = uuidv4();
    windows[id] = {
      id,
      isVisible: false,
      position,
      size,
      zIndex: behavior.zIndex || highestZIndex.value,
      resizable: layout.resizable === undefined || layout.resizable === true,
      draggable: behavior.draggable === undefined || behavior.draggable === true,
      modal: behavior.modal === true,
      windowConfig,
    };

    return id;
  };

  const createOrRestoreWindow = async (
    windowConfig: WindowConfig,
    corawOptions?: CreateOrRestoreWindowOptions,
  ): Promise<string> => {
    let existingWindow = Object.values(windows).find(
      (window) => window.windowConfig.name === windowConfig.name && window.windowConfig.type === windowConfig.type,
    );

    if (corawOptions?.forceCreate) {
      existingWindow = undefined;
    }

    let oldWindow = Object.values(windows).find((window) => window.isVisible);

    if (existingWindow) {
      for (const window of Object.values(windows).filter((window) => window.windowConfig.type === windowConfig.type)) {
        if (window.id === existingWindow.id) {
          try {
            if (options?.beforeWindowRestore) {
              await options.beforeWindowRestore(window, {
                isReplacing: oldWindow !== undefined,
                replacedWindow: oldWindow,
              });
            }
            window.isVisible = true;
            window.props = corawOptions?.data || {};
            if (options?.afterWindowRestore) {
              await options.afterWindowRestore(window, {
                isReplacing: oldWindow !== undefined,
                replacedWindow: oldWindow,
              });
            }
          } catch (e) {
            console.error(e);
          }
          if (window.windowConfig.type !== "AttachWindow") {
            break;
          }
        } else if (window.isVisible && window.windowConfig.type === "AttachWindow") {
          try {
            if (options?.beforeWindowMinimize) {
              await options.beforeWindowMinimize(window, {
                isReplaced: true,
                replacement: existingWindow,
              });
            }
            window.isVisible = false;
            if (options?.afterWindowMinimize) {
              await options.afterWindowMinimize(window, {
                isReplaced: true,
                replacement: existingWindow,
              });
            }
          } catch (e) {
            console.error(e);
          }
        }
      }

      console.log(windows);
      return existingWindow.id;
    }

    const layout = windowConfig.layout || {
      x: undefined,
      y: undefined,
      width: 400,
      height: 300,
      resizable: undefined,
      maxWidth: undefined,
      maxHeight: undefined,
      minWidth: undefined,
      minHeight: undefined,
    };

    //给toolbar默认值
    windowConfig.toolbar = windowConfig.toolbar || {};
    windowConfig.toolbar.show = windowConfig.toolbar.show === undefined || windowConfig.toolbar.show === true;
    windowConfig.toolbar.hoverShow = windowConfig.toolbar.hoverShow === true;
    windowConfig.toolbar.showMinimizeButton =
      windowConfig.toolbar.showMinimizeButton === undefined || windowConfig.toolbar.showMinimizeButton === true;
    windowConfig.toolbar.showMaximizeButton =
      windowConfig.toolbar.showMaximizeButton === undefined || windowConfig.toolbar.showMaximizeButton === true;
    windowConfig.toolbar.showCloseButton =
      windowConfig.toolbar.showCloseButton === undefined || windowConfig.toolbar.showCloseButton === true;

    // 设置默认位置和尺寸
    const position = {
      x: layout.x === undefined ? 100 : layout.x,
      y: layout.y === undefined ? 100 : layout.y,
    };

    const size = {
      width: layout.width || 400,
      height: layout.height || 300,
      maxWidth: layout.maxWidth,
      maxHeight: layout.maxHeight,
      minWidth: layout.minWidth,
      minHeight: layout.minHeight,
    };

    // 设置行为属性
    const behavior = windowConfig.behavior || {};

    const id = uuidv4();
    existingWindow = {
      id,
      isVisible: false,
      position,
      size,
      zIndex: behavior.zIndex || highestZIndex.value,
      resizable: layout.resizable === undefined || layout.resizable === true,
      draggable: behavior.draggable === undefined || behavior.draggable === true,
      modal: behavior.modal === true,
      windowConfig,
      props: corawOptions?.data || {},
    };

    windows[id] = existingWindow;

    for (const window of Object.values(windows).filter((window) => window.windowConfig.type === windowConfig.type)) {
      if (window.id === id) {
        try {
          if (options?.beforeWindowRestore) {
            await options.beforeWindowRestore(window, {
              isReplacing: oldWindow !== undefined,
              replacedWindow: oldWindow,
            });
          }
          window.isVisible = true;
          if (options?.afterWindowRestore) {
            await options.afterWindowRestore(window, {
              isReplacing: oldWindow !== undefined,
              replacedWindow: oldWindow,
            });
          }
        } catch (e) {
          console.error(e);
        }
        if (window.windowConfig.type !== "AttachWindow") {
          break;
        }
      } else if (window.isVisible && window.windowConfig.type === "AttachWindow") {
        try {
          if (options?.beforeWindowMinimize) {
            await options.beforeWindowMinimize(window, {
              isReplaced: true,
              replacement: existingWindow,
            });
          }
          window.isVisible = false;
          if (options?.afterWindowMinimize) {
            await options.afterWindowMinimize(window, {
              isReplaced: true,
              replacement: existingWindow,
            });
          }
        } catch (e) {
          console.error(e);
        }
      }
    }

    console.log(windows);
    return id;
  };

  // 显示窗口
  const showWindow = (id: string): void => {
    if (windows[id]) {
      windows[id].isVisible = true;
      bringToFront(id);
    }
  };

  // 切换窗口显示状态
  const toggleWindow = (id: string): void => {
    if (windows[id]) {
      windows[id].isVisible = !windows[id].isVisible;
      if (windows[id].isVisible) {
        bringToFront(id);
      }
    }
  };

  // 将窗口置于最前
  const bringToFront = (id: string): void => {
    if (windows[id]) {
      highestZIndex.value += 1;
      windows[id].zIndex = highestZIndex.value;
      activeWindowId.value = id;
    }
  };

  // 设置窗口位置
  const setWindowPosition = (id: string, x: number, y: number): void => {
    if (windows[id]) {
      windows[id].position = { x, y };
    }
  };

  // 设置窗口大小
  const setWindowSize = (id: string, width: number, height: number): void => {
    if (windows[id]) {
      if (!windows[id].size){
        windows[id].size = { width, height };
      }else{
        windows[id].size.width = width;
        windows[id].size.height = height;
      }
    }
  };

  // 获取窗口状态
  const getWindowState = (id: string) => {
    return computed(() => windows[id] || null);
  };

  // 获取所有窗口
  const getAllWindows = computed(() => windows);

  /**
   * 获取Attach窗口
   */
  const attachWindows = computed(() => {
    return Object.values(windows).filter((window) => window.windowConfig.type === "AttachWindow");
  });

  const notificationWindows = computed(() => {
    return Object.values(windows).filter((window) => window.windowConfig.type === "NotificationWindow");
  });

  const visibleAttachWindow = computed(() => {
    return Object.values(windows).find((window) => window.isVisible && window.windowConfig.type === "AttachWindow");
  });

  /**
   * 获取浮动窗口
   */
  const floatingWindows = computed(() => {
    return Object.values(windows).filter((window) => window.windowConfig.type === "FloatingWindow");
  });

  // 获取活动窗口
  const getActiveWindow = computed(() => (activeWindowId.value ? windows[activeWindowId.value] : null));

  // 获取可见窗口
  const getVisibleWindows = computed(() => {
    return Object.values(windows).filter((window) => window.isVisible);
  });

  // 关闭所有窗口
  const closeAllWindows = (): void => {
    Object.keys(windows).forEach((id) => {
      windows[id].isVisible = false;
    });
    activeWindowId.value = null;
  };

  /**
   * 最小化窗口
   * @param id 窗体id
   */
  const minimizeWindow = async (id: string): Promise<void> => {
    const w = windows[id];

    if (w) {
      try {
        await options?.beforeWindowMinimize?.(w, {
          isReplaced: false,
        });

        // 调用窗口配置中的最小化回调函数
        const onMinimize = w.windowConfig?.toolbar?.onMinimize;
        if (onMinimize && w.windowConfig) {
          const isHandled = onMinimize(w);
          if (!isHandled) {
            w.isVisible = false;
          }
        } else {
          w.isVisible = false;
        }
        await options?.afterWindowMinimize?.(w, {
          isReplaced: false,
        });
      } catch (error) {
        console.error("Error minimizing window:", error);
      } finally {
        /* empty */
      }
    }
  };
  /**
   * 关闭窗口
   * @param id 窗体id
   */
  const closeWindow = async (id: string): Promise<void> => {
    const w = windows[id];
    if (w) {
      try {
        await options?.beforeWindowClose?.(windows[id], {
          isReplaced: false,
        });
        windows[id].isVisible = false;
        await options?.afterWindowClose?.(windows[id], {
          isReplaced: false,
        });
        delete windows[id];
      } catch (error) {
        console.error("Error closing window:", error);
      }
    }
  };

  /**
   * 最大化窗口
   * @param id 窗体id
   */
  const maximizeWindow = async (id: string): Promise<void> => {
    const w = windows[id];
    if (w) {
      try {
        await options?.beforeWindowMaximize?.(w);
        if (w.windowConfig.type === "FloatingWindow" || w.windowConfig.type === "AttachWindow") {
          // 获取屏幕尺寸
          const screenWidth = window.innerWidth;
          const screenHeight = window.innerHeight;

          // 保存原始尺寸（如果需要恢复）
          const originalSize = { ...w.size };
          const originalPosition = { ...w.position };

          const layout = w.windowConfig.layout;
          if(w.size.width != screenWidth || w.size.height != screenHeight){
            // 设置为最大化尺寸
            setWindowPosition(id, 0, 0);
            setWindowSize(id, screenWidth, screenHeight);
            if(windowResizeObserver){
              windowResizeObserver.resize(w,{
                targetWidth: screenWidth,
                targetHeight: screenHeight,
                targetX: 0,
                targetY: 0,
              })
            }
          }else{
            setWindowPosition(id, layout.x, layout.y);
            setWindowSize(id, layout.width, layout.height);
            if(windowResizeObserver){
              windowResizeObserver.resize(w,{
                targetWidth: layout.width,
                targetHeight: layout.height,
                targetX: layout.x,
                targetY: layout.y,
              })
            }
          }


        }

        // 调用窗口配置中的最大化回调函数
        const onMaximize = w.windowConfig?.toolbar?.onMaximize;
        if (onMaximize && w.windowConfig) {
          const isHandled = onMaximize(w);
          if (!isHandled) {
            w.isVisible = true;
          }
        } else {
          w.isVisible = true;
        }
        await options?.afterWindowMaximize?.(w);
      } catch (error) {
        console.error("Error maximizing window:", error);
      }
    }
  };

  const resizeWindow = async (id: string, width: number, height: number): Promise<void> => {
    const w = windows[id];
    if (w) {
      try {
        await options?.beforeWindowResize?.(w, {
          targetWidth: width,
          targetHeight: height,
          targetX: w.position.x,
          targetY: w.position.y,
        });
        setWindowSize(id, width, height);
        await options?.afterWindowResize?.(w, {
          targetWidth: width,
          targetHeight: height,
          targetX: w.position.x,
          targetY: w.position.y,
        });
      } catch (error) {
        console.error("Error resizing window:", error);
      }
    }
  };

  // 恢复窗口
  const restoreWindow = (id: string, originalSize: any, originalPosition: any): void => {
    if (windows[id]) {
      setWindowPosition(id, originalPosition.x, originalPosition.y);
      setWindowSize(id, originalSize.width, originalSize.height);
    }
  };

  const getWindowsByType = (type: WindowType) => {
    return Object.values(windows).filter((window) => window.windowConfig.type === type);
  };

  const bindWindowEl = (id: string, el: HTMLElement) => {
    let w = windows[id];
    if (w) {
      w.el = el;
    }
  };
  return {
    createWindow,
    createOrRestoreWindow,
    showWindow,
    toggleWindow,
    bringToFront,
    setWindowPosition,
    setWindowSize,
    getWindowState,
    getAllWindows,
    getActiveWindow,
    getVisibleWindows,
    getWindowsByType,
    closeAllWindows,
    minimizeWindow,
    maximizeWindow,
    closeWindow,
    restoreWindow,
    resizeWindow,
    attachWindows,
    notificationWindows,
    visibleAttachWindow,
    floatingWindows,
    bindWindowEl,
    setDefaultOptions,
    setWindowResizeObserver
  };
}
