import { Component } from "vue";

/**
 * AttachWindow 嵌入到悬浮球上的窗体,可以移动,但是与悬浮球坐标同步
 * NotificationWindow 通知窗体,不可移动
 * FloatingWindow 浮动窗体,可以随意移动
 */
export type WindowType = "AttachWindow" | "NotificationWindow" | "FloatingWindow";
export type WindowClassic = "";

export interface WindowConfig {
  name: string;
  title: string;
  component: Component;

  type: WindowType; // 窗口类型：附加窗口、通知窗口或浮动窗口

  /** 布局：位置与尺寸 */
  layout?: {
    x?: number; // 初始 X 坐标
    y?: number; // 初始 Y 坐标
    width: number; // 宽度
    height: number; // 高度
    resizable?: boolean; // 是否可调整大小
    maxWidth?: number; // 最大宽度
    maxHeight?: number; // 最大高度
    minWidth?: number; // 最小宽度
    minHeight?: number; // 最小高度
  };

  /** 工具栏控制按钮配置（扁平化，推荐） */
  toolbar?: {
    show?: boolean; // 是否显示工具栏，默认 true
    hoverShow?: boolean; // 是否悬停显示，默认 true
    showMinimizeButton?: boolean; // 是否显示最小化，默认 true
    showMaximizeButton?: boolean; // 是否显示最大化，默认 true
    showCloseButton?: boolean; // 是否显示关闭，默认 true

    onMinimize?: (w: WindowState) => boolean; //返回true代表事件已经自行处理，默认false
    onMaximize?: (w: WindowState) => boolean; //返回true代表事件已经自行处理，默认false
    onClose?: (w: WindowState) => boolean; //返回true代表事件已经自行处理，默认false
  };

  /** 窗口行为配置 */
  behavior?: {
    draggable?: boolean; // 是否可拖动，默认 true
    zIndex?: number; // 层级，可控制叠放顺序
    modal?: boolean; // 是否模态窗口
  };
}

export interface WindowState {
  id: string;
  isVisible: boolean;
  position: {
    x: number;
    y: number;
  };
  size: {
    width: number;
    height: number;
    maxWidth?: number;
    maxHeight?: number;
    minWidth?: number;
    minHeight?: number;
  };
  zIndex: number;
  resizable: boolean;
  draggable: boolean;
  modal: boolean;
  windowConfig: WindowConfig;
  el?: HTMLElement | null;
  props?: Record<string, any>;
}
