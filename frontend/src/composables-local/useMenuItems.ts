// 获取屏幕宽度的75%
import { computed, markRaw, ref } from "vue";
import {PenTool,ListTodo,MessageCircleMore,Settings} from 'lucide-vue-next';
import {MenuItem} from "../types/menu";

function getScreenWidth75Percent(): number {
  return Math.floor(window.screen.width * 0.75);
}

/**
 * 注册菜单项
 */
export function useMenuItems() {
  const _menuItems = ref<MenuItem[]>([]);


  // Vue 组件菜单 - 随手记
  _menuItems.value.push(
      {
        id: "ball-quick-chat",
        label: "问一问",
        type: "chat",
        icon: MessageCircleMore,
        visible: true,
        windowConfig: {
          name: "ChatWithMe",
          title: "问一问",
          component: markRaw(MessageCircleMore),
          type: "AttachWindow",
          layout: {
            width: 400,
            height: 500,
            resizable: true,
          },
          toolbar: {
            show: true,
          },
        },
      },
    {
      id: "ball-quick-notes",
      label: "备忘录",
      type: "note",
      icon: PenTool,
      visible: true,
      windowConfig: {
        name: "QuickNotes",
        title: "备忘录",
        component: markRaw(PenTool),
        type: "AttachWindow",
        layout: {
          width: 400,
          height: 500,
          resizable: true,
        },
        toolbar: {
          show: true,
        },
      },
    },

    {
      id: "ball-todo",
      type: "todo",
      icon: ListTodo,
      label: "待办提醒",
      visible: true,
      windowConfig: {
        name: "TodoTab",
        title: "待办提醒",
        component: markRaw(ListTodo),
        type: "AttachWindow",
        layout: {
          width: 350,
          height: 400,
          resizable: true,
        },
        toolbar: {
          show: true,
        },
      },
    },
      {
        id: "settings",
        type: "setting",
        icon: Settings,
        label: "设置",
        visible: true,
        windowConfig: {
          name: "Settings",
          title: "配置",
          component: markRaw(Settings),
          type: "AttachWindow",
          layout: {
            width: 350,
            height: 400,
            resizable: true,
          },
          toolbar: {
            show: true,
          },
        },
      }

  );

  const menuItems = computed(() => {
    return _menuItems.value;
  });

  const getMenuItemById = (id: string) => {
    return menuItems.value.find((item) => item.id === id);
  };
  const getMenuItemByType = (type: string) => {
    return menuItems.value.find((item) => item.type === type);
  };

  return {
    menuItems,
    getMenuItemById,
    getMenuItemByType,
  };
}
