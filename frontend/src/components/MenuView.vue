<template>
  <div
      ref="menuElement"
      id="absolute-menu-container-id"
      class="absolute-menu-container"
      :style="{height:`${ballSize}px`,'padding-right':`${ballSize * 0.2}px`}"
  >
    <div class="menu" :style="{height:`${ballSize}px`}">

      <div class="menu-item" :style="{width:`${ballSize}px`,height:`${ballSize}px`}">
        <!--悬浮球占位-->
      </div>
      <div class="menu-item"
           :style="{width:`${ballSize * 0.8}px`,height:`${ballSize * 0.8}px`}"
           v-for="item in visibleMenuItems"
           :key="item.id"
           @click="handleMenuClick(item)"
      >
        <component :size="ballSize * 0.5" :is="item.icon" class="menu-item-icon"
                   v-if="'function' === typeof item.icon"/>
        <img :style="{width:`${ballSize * 0.6}px`,height:`${ballSize * 0.8}px`}" :src="item.icon"
             v-else-if="typeof item.icon === 'string'" class="menu-item-icon" :alt="item.label"/>
        <Sparkle :size="ballSize * 0.6" v-else/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {Sparkle} from "lucide-vue-next";
import {computed} from "vue";
import {useMenuItems} from "../composables-local/useMenuItems";
import {useSharedStatus} from "../composables-local/useSharedStatus";
import {MenuItem} from "../types/menu";

const {menuItems} = useMenuItems();

const visibleMenuItems = computed(() => {
  return menuItems.value.filter((item) => item.visible);
});

const {ballSize} = useSharedStatus();

const handleMenuClick = (menu: MenuItem) => {
  // console.log("Menu item clicked:", menuId, context);
};

</script>

<style scoped lang="less">
.absolute-menu-container {
  position: absolute;
  left: 1px;
  top: 1px;
  /* 移除固定尺寸，由JavaScript动态控制 */
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12),
  0 2px 8px rgba(0, 0, 0, 0.08);
  box-shadow: 0px 0px 1px 0px rgba(95, 108, 144, 0.3);
  pointer-events: auto;
  z-index: 9998;
  background-size: cover;

  .menu {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }
}

.menu-item {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 50%;
  //padding: 0;
  justify-content: center;
  color: #ffffff;
  box-sizing: border-box;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
}

.item-meeting {
  &:hover {
    background: rgba(52, 199, 89, 0.2);
    transform: translateX(2px);
  }
}

.item-note {
  &:hover {
    background: rgba(255, 141, 40, 0.2);
    transform: translateX(2px);
  }
}

.item-todo {
  &:hover {
    background: rgba(255, 45, 85, 0.2);
    transform: translateX(2px);
  }
}

/* .item-chat 保留默认样式，无 hover 效果 */
.item-chat {
  flex: 1;
}

.menu-item-icon {
}

.menu-item-icon img {
  width: 50%;
  height: 50%;
  border-radius: 6px;
  object-fit: cover;
}

.menu-item-label {
  flex: 1;
  font-size: 10px;
  color: #333;
  font-weight: 500;
}

.menu-item-arrow {
  width: 16px;
  height: 16px;
  color: #999;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.menu-item:hover .menu-item-arrow {
  opacity: 1;
}

.menu-empty {
  padding: 32px 16px;
  text-align: center;
  color: #999;
  font-size: 14px;
}

/* 菜单样式保持不变，变换动画由 transitionUtils 处理 */
</style>
