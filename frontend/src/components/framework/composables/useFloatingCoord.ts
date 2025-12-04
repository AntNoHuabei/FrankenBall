import {computed, Ref, ref} from "vue";

export type FloatingArea =  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;


const _activeFloatingArea = ref<FloatingArea>(3);

export function useFloatingCoord(currentElement: Ref<HTMLElement|null>) {

    const activeFloatingArea = computed(() => {
        return _activeFloatingArea.value
    })

    const setActiveByLocalPoint = (x: number, y: number) => {
        const el = currentElement
        if (!el?.value) return
        const width = el.value.clientWidth
        const height = el.value.clientHeight
        if (width <= 0 || height <= 0) return
        const nx = Math.min(Math.max(x, 0), width - 1e-6)
        const ny = Math.min(Math.max(y, 0), height - 1e-6)
        const col = Math.min(3, Math.floor((nx / width) * 3) + 1) // 1..3
        const row = Math.min(3, Math.floor((ny / height) * 3) + 1) // 1..3
        let activeIndex = (row - 1) * 3 + col
        _activeFloatingArea.value = activeIndex as FloatingArea
    }

    const setActiveByClientPoint = (clientX: number, clientY: number) => {
        const el = currentElement
        if (!el?.value) return
        const rect = el.value.getBoundingClientRect()
        const x = Math.min(Math.max(clientX - rect.left, 0), rect.width)
        const y = Math.min(Math.max(clientY - rect.top, 0), rect.height)
        setActiveByLocalPoint(x, y)
    }

    return {
        setActiveByClientPoint,
        activeFloatingArea
    }
}