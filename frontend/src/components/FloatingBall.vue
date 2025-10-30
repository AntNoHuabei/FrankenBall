<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

</script>

<template>
  <div
    ref="ballRef"
    class="floating-ball mouse-interactive_drag"
  >
    <div class="ball-content">
      <div class="ball-icon">🐷</div>
    </div>
  </div>
</template>

<style scoped>
.floating-ball {
  position: absolute;
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  /* cursor: move; */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
  z-index: 9999;
  user-select: none;
}

.floating-ball:hover {
  transform: scale(1.1);
}

.floating-ball:active {
  transform: scale(0.95);
}

.ball-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.ball-icon {
  font-size: 24px;
  color: white;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}
</style>

