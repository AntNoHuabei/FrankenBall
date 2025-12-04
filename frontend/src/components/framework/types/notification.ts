import { WindowConfig } from "./window";

/**
 * 通知类型
 */
export type NotificationType = "";
/**
 * 通知配置
 */
export interface NotificationConfig {
  label: string;
  type: NotificationType;
  windowConfig: WindowConfig;
}

/**
 * 通知信息
 */
export interface Notification {
  id: string;
  config: NotificationConfig;
  createTime: Date;
  data?: any;
}
