/* eslint-disable unicorn/prefer-single-call */
/* eslint-disable require-await */
import { computed, nextTick, onBeforeUnmount, onMounted, ref, Ref,ComponentPublicInstance } from "vue";
import {useSharedStatus} from "./useSharedStatus";
import {WindowState} from "../types/window";
import {Options as WindowManagerOptions,CloseOptions, MinimizeOptions, ResizeOptions, RestoreOptions, useWindowManager} from "./useWindowManager";
import {ViewMode} from "../types/view";

export interface Options {
  rootElement: Ref<HTMLElement | null>;
  ballViewEle: Ref<ComponentPublicInstance | null>;
  menuViewEle: Ref<ComponentPublicInstance | null>;
  windowViewEle: Ref<ComponentPublicInstance | null>;
}

const { lastBallPosition, isBallDragging, isBallMoving, ballSize,animationDuration } = useSharedStatus();
const snapConfig = {
  threshold: 100,
  edgeDistance: 10,
  ballSize: ballSize.value,
};

export const useViewManager = (option: Options) => {


  const windowManagerOptions: WindowManagerOptions = {
    beforeWindowResize: async (w: WindowState, detail: ResizeOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
    },
    afterWindowResize: async (w: WindowState, detail: ResizeOptions) => {
      if (!w.el) {
        return;
      }
      w.el.style.width = `${detail.targetWidth}px`;
      w.el.style.height = `${detail.targetHeight}px`;

      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      //如果当前是组件视图,而且当前窗体是可见的,则需要更新根组件的大小
      if (w.isVisible && viewMode.value === "component") {
        const { needsRepositioning, newPosition } = checkBoundaryAndAdjustPosition(
          detail.targetWidth,
          detail.targetHeight,
        );
        if (needsRepositioning && newPosition) {
          //调整根组件的尺寸,保持同步
          adjustRootSizeWithAnim(detail.targetWidth, detail.targetHeight, newPosition.x, newPosition.y);
        } else {
          //如果不需要调整位置,则需要更新根组件的大小
          adjustRootSizeWithAnim(detail.targetWidth, detail.targetHeight);
        }
      }
    },
    beforeWindowRestore: async (w: WindowState, detail: RestoreOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      console.log("beforeWindowRestore", w, detail);

      await nextTick(() => {
        if (!w.el) {
          return;
        }
        if (viewMode.value === "menu") {
          //因为不排除有悬浮球直接切换到组件的情况,所以这里需要判断一下
          //如果是从菜单切换到组件,那么就需要切换到组件视图
          inTransformView.value = ["menu", "component"];
        } else if (viewMode.value === "ball") {
          // 记录当前悬浮球位置
          recordBallPosition();
          inTransformView.value = ["ball", "component"];
        } else {
          inTransformView.value = ["component"];
        }

        viewMode.value = "component";

        //先将透明度设置为0
        adjustEleOpacity(w.el!, 0);
      });
    },
    afterWindowRestore: async (w: WindowState, detail: RestoreOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      console.log("afterWindowRestore", w, detail);
      await nextTick(() => {
        let todo: Promise<void>[]  = [];
        if (detail.isReplacing && detail.replacedWindow?.el) {
          todo.push(adjustEleOpacityWithAnim(detail.replacedWindow.el!, 0, false));
        } else if (inTransformView.value.includes("menu")) {
          //如果是从菜单切换过来,则菜单需要隐藏
        //  todo.push(adjustEleOpacityWithAnim(menuViewEle.value!.$el, 0, false));
        } else if (inTransformView.value.includes("ball")) {
          //如果是从悬浮球切换过来,则悬浮球需要隐藏
        //  todo.push(adjustEleOpacityWithAnim(ballViewEle.value!.$el, 0, false));
        }
        //再将透明度设置为1
        todo.push(adjustEleOpacityWithAnim(w.el!, 1));
        if (w.windowConfig?.layout?.x && w.windowConfig?.layout?.y) {
          todo.push(
            adjustRootSizeWithAnim(w.size.width, w.size.height, w.windowConfig?.layout?.x, w.windowConfig?.layout?.y),
          );
        } else {
          const { needsRepositioning, newPosition } = checkBoundaryAndAdjustPosition(w.size.width, w.size.height);
          if (needsRepositioning && newPosition) {
            todo.push(adjustRootSizeWithAnim(w.size.width, w.size.height, newPosition.x, newPosition.y));
          } else {
            todo.push(adjustRootSizeWithAnim(w.size.width, w.size.height));
          }
        }

        Promise.all(todo).then(() => {
          // 动画完成后的清理工作
          inTransformView.value = [];
        });
      });
    },
    beforeWindowMinimize: async (w: WindowState, detail: MinimizeOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      return new Promise<void>((resolve, reject) => {
        nextTick(() => {
          if (!w.el) {
            reject();
            return;
          }

          if (detail.isReplaced) {
            //如果是被替换的窗口,则需要将替换窗口的透明度设置为0
            adjustEleOpacityWithAnim(w.el!, 0, false);
            resolve();
          } else {
            viewMode.value = "ball";
            inTransformView.value = ["ball", "component"];
            //将窗口视图元素透明度设置为0
            Promise.all([
              //最小化 先不要透明度过度....要眼睁睁的看着窗口视图元素尺寸变成一个球
              //adjustEleOpacityWithAnim(w.el!, 0, true),
             // adjustEleOpacityWithAnim(ballViewEle.value!.$el, 1),
              adjustRootSizeWithAnim(
                snapConfig.ballSize,
                snapConfig.ballSize,
                lastBallPosition.value?.x,
                lastBallPosition.value?.y,
              ),
            ])
              .then(() => {
                // 动画完成后的清理工作
                inTransformView.value = [];
                resolve();
              })
              .catch((error) => {
                console.error("beforeWindowMinimize: 清理过程中出错", error);
                reject(error);
              });
          }
        });
      });
    },
    beforeWindowClose: async (w: WindowState, detail: CloseOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      return new Promise<void>((resolve, reject) => {
        nextTick(() => {
          if (!w.el) {
            reject();
            return;
          }
          if (detail.isReplaced) {
            adjustEleOpacityWithAnim(w.el!, 0, false);
            resolve();
            return;
          } else {
            //视图状态切换回悬浮球
            viewMode.value = "ball";
            inTransformView.value = ["ball", "component"];

            Promise.all([
              adjustEleOpacityWithAnim(w.el!, 0, false),
              //adjustEleOpacityWithAnim(ballViewEle.value!.$el, 1),
              adjustRootSizeWithAnim(
                snapConfig.ballSize,
                snapConfig.ballSize,
                lastBallPosition.value?.x,
                lastBallPosition.value?.y,
              ),
            ])
              .then(() => {
                // 动画完成后的清理工作
                inTransformView.value = [];
                console.log("afterWindowClose: 清理完成");
                resolve();
              })
              .catch((error) => {
                console.error("afterWindowClose: 清理过程中出错", error);
                reject(error);
              });
          }
        });
      });
    },
    afterWindowClose: async (w: WindowState, detail: CloseOptions) => {
      if (w.windowConfig.type !== "AttachWindow") {
        return;
      }
      console.log("afterWindowClose", w, detail);
      // 动画完成后的清理工作
    },
  };

  const { createOrRestoreWindow, setDefaultOptions ,setWindowResizeObserver} = useWindowManager(windowManagerOptions);

  // 窗口管理器默认选项
  setDefaultOptions(windowManagerOptions);


  setWindowResizeObserver({
    resize(w: WindowState, detail: ResizeOptions) {

      if(viewMode.value === "component"){
        if(!detail.targetX ){
          detail.targetX = lastBallPosition.value.x
        }
        if(!detail.targetY){
          detail.targetY = lastBallPosition.value.y
        }
        adjustRootSizeWithAnim(detail.targetWidth, detail.targetHeight, detail.targetX, detail.targetY)
      }
    }
  })

  const {
    /**
     * 悬浮球视图元素(FloatingBall.vue)
     */
    ballViewEle,
    /**
     * 菜单视图元素(MenuView.vue)
     */
    menuViewEle,
    /**
     * 窗体视图元素(WindowView.vue)
     */
    windowViewEle,
    /**
     * 根元素(根元素是body)
     */
    rootElement,
  } = option;

  const inTransformView = ref<ViewMode[]>([]);

  let hideMenuTimer: number | null = null;
  let showMenuTimer: number | null = null;

  // 菜单尺寸（固定展开方向为 right-bottom）
  //const menuWidth = menuSize.value.width;

  /**
   * 当前视图模式(ball,menu,component)
   */
  const viewMode = ref<ViewMode>("ball");

  const isAnimating = ref(false);

  const adjustEleOpacity = (ele: HTMLElement, opacity: number) => {
    ele.style.opacity = opacity.toString();
  };

  const adjustEleOpacityWithAnim = async (ele: HTMLElement, opacity: number, hideOnZero = false): Promise<void> => {
    if (!ele) {
      console.warn("adjustEleOpacityWithAnim: 元素不存在或无效");
      return;
    }

    return new Promise<void>((resolve) => {
      // 监听动画结束事件
      const handleTransitionEnd = (event: TransitionEvent) => {
        // 确保是透明度过渡结束
        if (event.propertyName === "opacity") {
          console.log("透明度变换结束事件", event);
          // 移除事件监听器
          ele.removeEventListener("transitionend", handleTransitionEnd);

          // 当透明度为0且设置了hideOnZero时，隐藏元素
          if (hideOnZero && opacity === 0) {
            ele.style.display = "none";
          }

          // 清除过渡样式，避免影响其他操作

          // 恢复原始的transition设置
          // ele.style.transition = "";
          resolve();
        }
      };

      // 添加过渡结束事件监听
      ele.addEventListener("transitionend", handleTransitionEnd);

      // 设置CSS过渡动画
      ele.style.transition = `opacity ${animationDuration.value}ms  cubic-bezier(0.4, 0, 1, 1)`;

      console.log("原始透明度,目标透明度", ele.style.opacity, opacity);
      setTimeout(() => {
        // 设置目标透明度
        ele.style.opacity = opacity.toString();
      }, 50);
    });
  };

  /**
   * 调整根元素尺寸以适应当前视图模式，并在需要时附加位移动画
   */
  const adjustRootSizeWithAnim = async (width: number, height: number, x?: number, y?: number): Promise<void> => {
    const rootElement = option.rootElement.value;
    if (!rootElement) {
      return;
    }

    return new Promise<void>((resolve) => {


      // 监听动画结束事件
      const handleTransitionEnd = (event: TransitionEvent) => {
        // 确保是我们关心的过渡属性结束
        if (event.propertyName === "width" || event.propertyName === "height") {
          console.log("尺寸过渡结束", event);
          rootElement.removeEventListener("transitionend", handleTransitionEnd);
          // 清除过渡样式，避免影响其他操作
          rootElement.style.transition = "";
          resolve();
        }
      };

      // 添加过渡结束事件监听
      rootElement.addEventListener("transitionend", handleTransitionEnd);

      // 设置CSS过渡动画
      rootElement.style.transition = `
        width ${animationDuration.value}ms ease-in-out,
        height ${animationDuration.value}ms ease-in-out,
        left ${animationDuration.value}ms ease-in-out,
        top ${animationDuration.value}ms ease-in-out
      `;

      setTimeout(() => {

        const rect = rootElement.getBoundingClientRect();
        // 如果提供了位置信息，则同时进行位移动画
        if (x !== undefined && y !== undefined && rect.left !== x || rect.top !== y) {
          const screenWidth = window.innerWidth;
          const screenHeight = window.innerHeight;

          // 确保目标位置在屏幕范围内
          const safeTargetX = Math.max(0, Math.min(x!, screenWidth - width));
          const safeTargetY = Math.max(0, Math.min(y!, screenHeight - height));
          rootElement.style.left = `${Math.round(safeTargetX)}px`;
          rootElement.style.top = `${Math.round(safeTargetY)}px`;
        }
        // 设置目标尺寸
        rootElement.style.width = `${Math.round(width)}px`;
        rootElement.style.height = `${Math.round(height)}px`;

      }, 50);
    });
  };

  /**
   * 绑定/解绑悬浮球视图事件
   */
  const doBallViewEventBinding = (action: "bind" | "unbind") => {
    if (action === "bind") {
      if (ballViewEle.value) {
        ballViewEle.value.$el.addEventListener("mouseover", handleBallHover);
        ballViewEle.value.$el.addEventListener("mouseleave", handleBallLeave);
      }
    } else if (ballViewEle.value) {
      ballViewEle.value.$el.removeEventListener("mouseover", handleBallHover);
      ballViewEle.value.$el.removeEventListener("mouseleave", handleBallLeave);
    }
  };

  /**
   * 绑定/解绑菜单视图事件
   */
  const doMenuViewEventBinding = (action: "bind" | "unbind") => {
    if (action === "bind") {
      if (menuViewEle.value) {
        menuViewEle.value.$el.addEventListener("mouseenter", handleMenuEnter);
        menuViewEle.value.$el.addEventListener("mouseleave", handleMenuLeave);
      }
    } else if (menuViewEle.value) {
      menuViewEle.value.$el.removeEventListener("mouseenter", handleMenuEnter);
      menuViewEle.value.$el.removeEventListener("mouseleave", handleMenuLeave);
    }
  };

  /**
   * 绑定/解绑窗体视图事件
   */
  const doWindowViewEventBinding = (action: "bind" | "unbind") => {};

  const handleBallHover = () => {
    if (isBallDragging.value) {
      return;
    }

    // 取消隐藏菜单的定时器（鼠标重新进入悬浮球区域）
    if (hideMenuTimer) {
      clearTimeout(hideMenuTimer);
      hideMenuTimer = null;
    }

    // 如果已经在显示菜单，则不重复触发
    if (viewMode.value === "menu") {
      return;
    }

    // 设置1秒延迟后显示菜单
    if (showMenuTimer) {
      clearTimeout(showMenuTimer);
    }
    showMenuTimer = window.setTimeout(() => {
      showMenuOnHover();
      showMenuTimer = null;
    }, 500);
  };

  const handleBallLeave = () => {
    // 如果悬浮球正在移动，忽略 mouseleave 事件
    if (isBallMoving.value) {
      // console.log("悬浮球正在移动，忽略 mouseleave 事件");
      return;
    }
    if (isBallDragging.value) {
      return;
    }

    // 取消显示菜单的定时器
    if (showMenuTimer) {
      clearTimeout(showMenuTimer);
      showMenuTimer = null;
    }

    // 鼠标离开悬浮球，如果菜单已显示，则延迟隐藏
    if (viewMode.value === "menu") {
      // console.log("Ball leave detected, scheduling hide");
      if (hideMenuTimer) {
        clearTimeout(hideMenuTimer);
      }
      hideMenuTimer = window.setTimeout(() => {
        hideMenuWithAnimation();
        hideMenuTimer = null;
      }, 200);
    }
  };

  const hideMenuWithAnimation = async () => {
    console.log("Starting menu collapse animation");
    // 切换到ball视图
    viewMode.value = "ball";
    inTransformView.value = ["ball", "menu"];

    await Promise.all([
     // adjustEleOpacityWithAnim(menuViewEle.value!.$el, 0),
      adjustRootSizeWithAnim(
        snapConfig.ballSize,
        snapConfig.ballSize,
        lastBallPosition.value?.x,
        lastBallPosition.value?.y,
      ),
      //adjustEleOpacityWithAnim(ballViewEle.value!.$el, 1),
    ]);
    // 动画完成后的清理工作
    inTransformView.value = [];
    // 同时执行菜单到悬浮球的变换和回到最后位置
    // await Promise.all([performViewTransition("menu", "ball", 300), returnToLastPosition()]);
  };
  /**
   *记录悬浮球当前位置
   */
  const recordBallPosition = () => {
    const ballBounds = rootElement.value?.getBoundingClientRect();
    if (ballBounds) {
      lastBallPosition.value = {
        x: ballBounds.x,
        y: ballBounds.y,
      };

      // console.log("记录悬浮球位置:", lastBallPosition.value);
    }
  };

  /**
   * 鼠标悬浮在悬浮球上显示菜单
   */
  const showMenuOnHover = () => {
    if (isBallDragging.value) {
      return;
    }

    // 记录当前悬浮球位置
    recordBallPosition();

    // 使用新的变换方法
    viewMode.value = "menu";
    inTransformView.value = ["ball", "menu"];

    //先将菜单透明度设置为0
   // adjustEleOpacity(menuViewEle.value!.$el, 0);

    let todo: Promise<void>[] = [];

    let menuWidth = menuViewEle.value!.$el.getBoundingClientRect().width;
    let menuHeight = menuViewEle.value!.$el.getBoundingClientRect().height;

    const { needsRepositioning, newPosition } = checkBoundaryAndAdjustPosition(menuWidth, menuHeight);
    // 如果需要重新定位悬浮球，同时执行位置动画和菜单展开动画
    if (needsRepositioning && newPosition) {
      // console.log("开始悬浮球位置调整和菜单展开同步动画");

      // 需要移动悬浮球到新位置
      //todo.push(animateToPosition(rootElement.value!, newBallPosition.x, newBallPosition.y));
      //需要调整根元素到菜单的尺寸,同时需要新的位置

      todo.push(adjustRootSizeWithAnim(menuWidth, menuHeight, newPosition.x, newPosition.y));
    } else {
      //不需要移动悬浮球,就不需要新的位置,只需要调整尺寸
      todo.push(adjustRootSizeWithAnim(menuWidth, menuHeight));
    }

    // 最后设置菜单透明度为1渐变动效
    //todo.push(adjustEleOpacityWithAnim(menuViewEle.value!.$el, 1));
    //todo.push(adjustEleOpacityWithAnim(ballViewEle.value!.$el, 0));

    console.info("菜单展开动画开始", todo);
    Promise.all(todo).then(() => {
      inTransformView.value = [];
      // 动画完成后，重置移动状态
      console.log("菜单展开动画完成");
    });
  };

  const handleMenuEnter = () => {
    // 鼠标进入菜单区域，取消隐藏定时器
    if (hideMenuTimer) {
      clearTimeout(hideMenuTimer);
      hideMenuTimer = null;
    }
  };

  const handleMenuLeave = () => {
    // 如果悬浮球正在移动，忽略 menuleave 事件
    if (isBallMoving.value) {
      // console.log("悬浮球正在移动，忽略 menuleave 事件");
      return;
    }

    // console.log("Menu leave detected, scheduling hide");
    // 延迟200ms隐藏菜单，给用户时间移回悬浮球或菜单
    if (hideMenuTimer) {
      clearTimeout(hideMenuTimer);
    }
    if (viewMode.value === "menu") {
      hideMenuTimer = window.setTimeout(() => {
        hideMenuWithAnimation();
        hideMenuTimer = null;
      }, 200);
    }
  };


  /**
   * 创建一个CSS动画，将元素从当前位置移动到指定位置
   * @param ele 元素
   * @param targetX
   * @param targetY
   */
  const animateToPosition = (ele: HTMLElement, targetX: number, targetY: number) => {
    // 验证输入参数
    if (!Number.isFinite(targetX) || !Number.isFinite(targetY)) {
      console.error("Invalid target position:", { targetX, targetY });
      return;
    }

    return new Promise((resolve) => {
      const screenWidth = window.innerWidth;
      const screenHeight = window.innerHeight;

      // 限制目标位置在窗口范围内
      const safeTargetX = Math.max(0, Math.min(targetX, screenWidth - snapConfig.ballSize));
      const safeTargetY = Math.max(0, Math.min(targetY, screenHeight - snapConfig.ballSize));

      // console.log("CSS Animation start:", {
      //   startPos: { x: position.x, y: position.y },
      //   targetPos: { x: safeTargetX, y: safeTargetY },
      //   delta: { x: safeTargetX - position.x, y: safeTargetY - position.y },
      // });

      // 设置CSS过渡动画
      ele.style.transition = `left ${animationDuration.value}ms cubic-bezier(0.4, 0, 0.2, 1), top ${animationDuration.value}ms cubic-bezier(0.4, 0, 0.2, 1)`;

      // 直接设置DOM位置，绕过Vue响应式系统
      const setElementPosition = (x: number, y: number) => {
        ele.style.left = `${Math.round(x)}px`;
        ele.style.top = `${Math.round(y)}px`;
      };

      // 直接设置目标位置，让CSS处理动画
      setElementPosition(safeTargetX, safeTargetY);

      // 监听动画结束事件
      const handleTransitionEnd = () => {
        ele.removeEventListener("transitionend", handleTransitionEnd);
        ele.style.transition = ""; // 清除过渡样式
        resolve({ x: safeTargetX, y: safeTargetY });
      };

      ele.addEventListener("transitionend", handleTransitionEnd);
    });
  };

  // 通用的边界检查和位置调整函数
  const checkBoundaryAndAdjustPosition = (targetWidth: number, targetHeight: number) => {
    const screenWidth = window.innerWidth;
    const screenHeight = window.innerHeight;

    let ballBounds = rootElement.value?.getBoundingClientRect();
    if (!ballBounds) {
      return { needsRepositioning: false, newPosition: { x: 0, y: 0 } };
    }

    // 计算目标组件展开后的边界位置（固定从右下角展开）
    // 注意：这里应该使用悬浮球的左上角作为基准，而不是直接使用 ballBounds
    const targetBounds = {
      left: ballBounds.x,
      top: ballBounds.y,
      right: ballBounds.x + targetWidth,
      bottom: ballBounds.y + targetHeight,
    };

    // 检查是否会超出屏幕边界
    const isOutOfBounds =
      targetBounds.left < 0 ||
      targetBounds.top < 0 ||
      targetBounds.right > screenWidth ||
      targetBounds.bottom > screenHeight;

    console.log("组件边界检查详情:", {
      组件边界: targetBounds,
      屏幕边界: { width: screenWidth, height: screenHeight },
      是否超出边界: isOutOfBounds,
      组件尺寸: { width: targetWidth, height: targetHeight },
      悬浮球位置: { x: ballBounds.x, y: ballBounds.y },
      距离右边: screenWidth - ballBounds.x,
      距离下边: screenHeight - ballBounds.y,
    });

    let needsRepositioning = false;
    let newPosition = { x: ballBounds.x, y: ballBounds.y };

    // 如果组件超出边界，计算新的悬浮球位置
    if (isOutOfBounds) {
      // console.log("组件将超出屏幕边界，需要调整悬浮球位置");
      needsRepositioning = true;

      // 计算调整后的悬浮球位置，确保组件完全在屏幕内（针对 right-bottom 展开）
      let adjustedX = ballBounds.x;
      let adjustedY = ballBounds.y;

      // 针对 right-bottom 展开的边界调整
      // 如果组件右边超出屏幕，向左移动悬浮球
      if (targetBounds.right > screenWidth) {
        adjustedX = screenWidth - targetWidth - 10; // 右边留10px边距
      }

      // 如果组件下边超出屏幕，向上移动悬浮球
      if (targetBounds.bottom > screenHeight) {
        adjustedY = screenHeight - targetHeight - 10; // 下边留10px边距
      }

      const halfBallSize = snapConfig.ballSize / 2;
      // 确保悬浮球不会超出屏幕边界（悬浮球尺寸32x32）
      adjustedX = Math.max(halfBallSize, Math.min(adjustedX, screenWidth - halfBallSize));
      adjustedY = Math.max(halfBallSize, Math.min(adjustedY, screenHeight - halfBallSize));

      newPosition = {
        x: adjustedX,
        y: adjustedY,
      };

      console.log("悬浮球位置调整:", {
        原始位置: { x: ballBounds.x, y: ballBounds.y },
        新位置: newPosition,
        组件边界: targetBounds,
        屏幕边界: { width: screenWidth, height: screenHeight },
      });
    }

    return { needsRepositioning, newPosition };
  };

  /**
   * 组件挂载时绑定三个视图的事件
   */
  onMounted(() => {
    doBallViewEventBinding("bind");
    doMenuViewEventBinding("bind");
    doWindowViewEventBinding("bind");
  });
  /**
   * 组件销毁时解绑三个视图的事件
   */
  onBeforeUnmount(() => {
    doBallViewEventBinding("unbind");
    doMenuViewEventBinding("unbind");
    doWindowViewEventBinding("unbind");
  });

  const viewVisibility = computed(() => {
    return {
      ball: viewMode.value === "ball" || inTransformView.value.includes("ball"),
      menu: viewMode.value === "menu" || inTransformView.value.includes("menu"),
      window: viewMode.value === "component" || inTransformView.value.includes("component"),
    };
  });

  
  
  const restoreLastBounds = () => {
    const lastPos = localStorage.getItem("lastBallPosition");

    const moveToDefaultPos = () => {
      const screenWidth = window.innerWidth;
      const screenHeight = window.innerHeight;


      console.log("screenWidth:", screenWidth)
      console.log("screenHeight:", screenHeight)


      const initPos = {
        x: screenWidth - snapConfig.ballSize / 2 - snapConfig.edgeDistance,
        y: screenHeight - snapConfig.ballSize / 2 - snapConfig.edgeDistance,
      };
      if (rootElement.value) {
        adjustRootSizeWithAnim(snapConfig.ballSize, snapConfig.ballSize,initPos.x, initPos.y)
      }
    };

    if (rootElement.value) {
      if (lastPos) {
        try {
          const pos = JSON.parse(lastPos);
          adjustRootSizeWithAnim(snapConfig.ballSize, snapConfig.ballSize,pos.x, pos.y)
        } catch (error) {
          console.error("解析最后位置失败:", error);
          moveToDefaultPos();
        }
      } else {
        moveToDefaultPos();
      }


      window.onresize = () => {
        const screenWidth = window.innerWidth;
        const screenHeight = window.innerHeight;
        console.log("screenWidth:", screenWidth)
        console.log("screenHeight:", screenHeight)
        moveToDefaultPos();
      };
    }

  };

  // 组件挂载时恢复最后位置
  onMounted(() => {
    restoreLastBounds();
  });

  return {
    inTransformView,
    viewVisibility,
  };
};
