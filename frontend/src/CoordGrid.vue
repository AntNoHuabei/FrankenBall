<script setup lang="ts">
import { ref } from 'vue'

const activeIndex = ref<number | null>(null)
const rootEl = ref<HTMLElement | null>(null)

const positionClass = (i: number) => `p${i}`

const setActiveByLocalPoint = (x: number, y: number) => {
  const el = rootEl.value
  if (!el) return
  const width = el.clientWidth
  const height = el.clientHeight
  if (width <= 0 || height <= 0) return
  const nx = Math.min(Math.max(x, 0), width - 1e-6)
  const ny = Math.min(Math.max(y, 0), height - 1e-6)
  const col = Math.min(3, Math.floor((nx / width) * 3) + 1) // 1..3
  const row = Math.min(3, Math.floor((ny / height) * 3) + 1) // 1..3
  activeIndex.value = (row - 1) * 3 + col
}

const setActiveByClientPoint = (clientX: number, clientY: number) => {
  const el = rootEl.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const x = Math.min(Math.max(clientX - rect.left, 0), rect.width)
  const y = Math.min(Math.max(clientY - rect.top, 0), rect.height)
  setActiveByLocalPoint(x, y)
}

defineExpose({ activeIndex, setActiveByLocalPoint, setActiveByClientPoint })
</script>

<template>
  <table class="coord-grid" ref="rootEl">
    <tbody>
      <tr v-for="row in 3" :key="row">
        <td v-for="col in 3" :key="col">
          <div
            class="dot"
            :class="[positionClass((row - 1) * 3 + col), { 'is-active': activeIndex === (row - 1) * 3 + col }]"
            :data-index="(row - 1) * 3 + col"
          />
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style lang="less" scoped>
.coord-grid {
  width: 100%;
  height: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  background-color: rgba(0, 0, 0, 0.5);
  table-layout: fixed;

  td {
    width: 33.3333%;
    height: 33.3333%;
    padding: 0;
    background-color: transparent;
    position: relative;
  }

  tr:not(:last-child) td {
    border-bottom: 2px dashed rgba(0, 0, 0, 0.35);
  }

  td:not(:last-child) {
    border-right: 2px dashed rgba(0, 0, 0, 0.35);
  }

  .dot {
    position: absolute;
    width: 17%;
    height: 13%;
    border: 2px dashed rgba(255, 255, 255, 0.6);
    background: rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    --tx: 0;
    --ty: 0;
    --scale: 1;
    transform: translate(var(--tx), var(--ty)) scale(var(--scale));
    transition: transform 120ms ease;
    will-change: transform;
  }

  .dot.is-active {
    --scale: 1.2;
  }

  .dot.p1 { top: 0; left: 0; }
  .dot.p2 { top: 0; left: 50%; --tx: -50%; }
  .dot.p3 { top: 0; right: 0; }
  .dot.p4 { top: 50%; left: 0; --ty: -50%; }
  .dot.p5 { top: 50%; left: 50%; --tx: -50%; --ty: -50%; }
  .dot.p6 { top: 50%; right: 0; --ty: -50%; }
  .dot.p7 { bottom: 0; left: 0; }
  .dot.p8 { bottom: 0; left: 50%; --tx: -50%; }
  .dot.p9 { bottom: 0; right: 0; }

  .dot.p1, .dot.p2, .dot.p3 {
    border-top-left-radius: 0;
    border-top-right-radius: 0;
    border-top:none;
  }
  .dot.p7, .dot.p8, .dot.p9 {
    border-bottom-left-radius: 0;
    border-bottom-right-radius: 0;
    border-bottom:none;
  }
  .dot.p1, .dot.p4, .dot.p7 {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    border-left:none;
  }
  .dot.p3, .dot.p6, .dot.p9 {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    border-right:none;
  }
}
</style>
