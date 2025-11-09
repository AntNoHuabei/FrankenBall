import {computed, Ref, ref} from "vue";
import {useMousePosition} from "../composables/useMousePosition";
import {useFloatingCoord} from "../composables/useFloatingCoord";

const dragState = ref({
    isDragging: false,
    startTime: 0,
    canShowCoordGrid: false
});

export function useCoordGrid(rootElement: Ref<HTMLElement|null>){

    const {setActiveByClientPoint, activeFloatingArea} = useFloatingCoord(rootElement)

    const showCoordGrid = computed(()=>{
        return dragState.value.canShowCoordGrid
    })

// 使用拖拽composable - 自动检测多个app-region拖拽样式元素
    useMousePosition({
        autoDetect: true,
        selector: [".mouse-drag", ".mouse-interactive_drag"],
        // 自定义拖拽行为 - 实现app-container拖拽
        onCustomDrag: (e: MouseEvent, element: HTMLElement) => {
            // console.log("自定义拖拽开始", element);

            // 检查是否是app-container元素

            // 初始化拖拽状态
            dragState.value.isDragging = true;
            dragState.value.startTime = Date.now();

            // 禁用默认拖拽行为，使用自定义实现
            return false;

            // 其他元素使用默认行为
        },
        onCustomDragMove: (e: MouseEvent, element: HTMLElement) => {
            // 只对app-container元素进行自定义拖拽处理
            if ( dragState.value.isDragging) {

                if (!dragState.value.canShowCoordGrid) {
                    if (Date.now() - dragState.value.startTime > 500) {

                        dragState.value.canShowCoordGrid = true
                    }
                } else {
                    setActiveByClientPoint(e.clientX, e.clientY);
                }

                //console.log("自定义拖拽移动", newX, newY)

                // 阻止默认拖拽行为
                return false;
            }

            // 其他元素使用默认行为
        },
        onCustomDragEnd: (e: MouseEvent, element: HTMLElement) => {
            // console.log("自定义拖拽结束", element);

                // 重置拖拽状态
                dragState.value.isDragging = false;
                dragState.value.canShowCoordGrid = false;
                dragState.value.startTime = 0;
        },
    });

    return {activeFloatingArea,showCoordGrid}
}
