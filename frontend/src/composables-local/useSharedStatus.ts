import { ref } from "vue";

const isBallMoving = ref(false); // 跟踪悬浮球是否正在移动
const isBallDragging = ref(false);
const lastBallPosition = ref<{ x: number; y: number } | null>(null); // 记录悬浮球的最后位置
const ballSize = ref<number>(40);
const animationDuration = ref<number>(300);

/**
 * 共享状态管理
 */
export const useSharedStatus = () => {
  return {
    isBallMoving,
    isBallDragging,
    lastBallPosition,
    ballSize,
    animationDuration
  };
};
