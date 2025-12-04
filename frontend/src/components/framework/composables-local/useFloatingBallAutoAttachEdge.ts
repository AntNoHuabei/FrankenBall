import {onUnmounted, ref, Ref, watch} from "vue";
import {FloatingArea, useFloatingCoord} from "../composables/useFloatingCoord";


const distanceForEdge = 0;
interface Position {
    x: number;
    y: number;
}
export function useFloatingBallAutoAttachEdge(currentElement: Ref<HTMLElement|null>) {

    const {activeFloatingArea} = useFloatingCoord(ref(null))

    const acquireMovementOrigin = (toAttachArea: FloatingArea,el: HTMLElement)=>{

        const rect = el.getBoundingClientRect()

        switch (toAttachArea){

            case 1:
                // 左上
                return {x: rect.left + distanceForEdge,y: rect.top + distanceForEdge}
            case 2:
                // 中上
                return {x: rect.left + rect.width / 2,y: rect.top + distanceForEdge}
            case 3:
                // 右上
                return {x: rect.right - distanceForEdge,y: rect.top + distanceForEdge}
            case 4:
                // 左中
                return {x: rect.left + distanceForEdge,y: rect.top + rect.height / 2}
            case 5:
                // 中中
                return {x: rect.left + rect.width / 2,y: rect.top + rect.height / 2}
            case 6:
                // 右中
                return {x: rect.right - distanceForEdge,y: rect.top + rect.height / 2}
            case 7:
                // 左下
                return {x: rect.left + distanceForEdge,y: rect.bottom - distanceForEdge}
            case 8:
                // 中下
                return {x: rect.left + rect.width / 2,y: rect.bottom - distanceForEdge}
            case 9:
                // 右下
                return {x: rect.right - distanceForEdge,y: rect.bottom - distanceForEdge}
            default:
                return {x: rect.left + rect.width / 2,y: rect.top + rect.height / 2}


        }
    }
    const acquireMovementTarget = (toAttachArea: FloatingArea)=>{
        switch (toAttachArea){
            case 1:
                // 左上
                return {x: distanceForEdge,y: distanceForEdge}
            case 2:
                // 中上

                return {x: window.innerWidth / 2,y: distanceForEdge}
            case 3:
                // 右上
                return {x: window.innerWidth - distanceForEdge,y: distanceForEdge}
            case 4:
                // 左中
                return {x: distanceForEdge,y: window.innerHeight / 2}
            case 5:
                // 中中
                return {x: window.innerWidth / 2,y: window.innerHeight / 2}
            case 6:
                // 右中
                return {x: window.innerWidth - distanceForEdge,y: window.innerHeight / 2}
            case 7:
                // 左下
                return {x: distanceForEdge,y: window.innerHeight - distanceForEdge}
            case 8:
                // 中下
                return {x: window.innerWidth / 2,y: window.innerHeight - distanceForEdge}
            case 9:
                // 右下
                return {x: window.innerWidth - distanceForEdge,y: window.innerHeight - distanceForEdge}
            default:
                return {x: window.innerWidth / 2,y: window.innerHeight / 2}

        }
    }

    const attachSide = (toAttachArea: FloatingArea)=>{

        if(!currentElement.value){
             return
        }

        //获取元素移动起始点
        const movementOrigin = acquireMovementOrigin(toAttachArea, currentElement.value)



    }

    watch(()=>activeFloatingArea,(nv)=>{
        attachSide(nv.value)
    })

    //初始化贴边吸贴
    attachSide(activeFloatingArea.value)


    //监听元素尺寸变化
    const resizeObserver = new ResizeObserver(()=>{
        attachSide(activeFloatingArea.value)
    })

    if(currentElement.value){
        resizeObserver.observe(currentElement.value)
    }else{
        const stopWatch = watch(()=>currentElement.value,(nv)=>{
            if(nv){
                resizeObserver.observe(nv)
                stopWatch()
            }
        })
    }

    onUnmounted(()=>{
        resizeObserver.disconnect()
    })
}