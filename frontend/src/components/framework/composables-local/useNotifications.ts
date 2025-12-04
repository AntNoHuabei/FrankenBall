import { computed, markRaw, reactive, ref } from "vue";
import { v4 as uuidv4 } from "uuid";
import TodoReminderTab from "~/components/floatingball/TodoReminderTab.vue";
import MeetingDialog from "~/components/floatingball/MeetingDialogTab.vue";
import MeetingEndDialog from "~/components/floatingball/MeetingEndDialogTab.vue";
import {NotificationConfig, NotificationType,Notification} from "../types/notification";
import {CreateOrRestoreWindowOptions, useWindowManager} from "./useWindowManager";
import {WindowConfig} from "../types/window";

export interface NotifyParams {
  /**
   * 通知类型,注册的时候填写的type
   */
  type: NotificationType;
  /**
   * 是否覆盖相同type的通知
   */
  override: boolean;
  /**
   * 通知数据
   */
  data?: any;
}

export const useNotifications = () => {
  const { createOrRestoreWindow } = useWindowManager();

  /**
   * 注册通知
   */
  const configs = ref<NotificationConfig[]>([]);

  /**
   * 发送通知
   * @param params
   */
  const notify = (params: NotifyParams) => {
    const config = configs.value.find((item) => item.type === params.type);
    if (config) {
      const id = uuidv4();
      const nf: Notification = {
        id,
        config,
        data: params.data,
        createTime: new Date(),
      };
      let windowOption: CreateOrRestoreWindowOptions = {
        data: params.data,
        forceCreate: !params.override,
      };
      createOrRestoreWindow(nf.config.windowConfig, windowOption);
      return id;
    }
  };

  /**
   * 注册目前已知的所有通知类型
   */

  const todoWindowConfig: WindowConfig = {
    name: "TodoReminderTab",
    title: "待办提醒",
    type: "NotificationWindow",
    component: markRaw(TodoReminderTab),
    layout: {
      width: 360,
      height: 170,
      resizable: false,
    },
    toolbar: {
      show: false,
    },
  };

  const meetingWindowConfig: WindowConfig = {
    name: "MeetingDialog",
    title: "会议弹框",
    type: "NotificationWindow",
    component: markRaw(MeetingDialog),
    layout: {
      width: 360,
      height: 170,
      resizable: false,
    },
    toolbar: {
      show: false,
    },
  };
  
  // 会议结束提醒弹框
  const meetingEndWindowConfig: WindowConfig = {
    name: "MeetingEndDialog",
    title: "结束录音",
    type: "NotificationWindow",
    component: markRaw(MeetingEndDialog),
    layout: {
      width: 330,
      height: 150,
      resizable: false,
    },
    toolbar: {
      show: false,
    },
  };



  return {
    notify,
  };
};
