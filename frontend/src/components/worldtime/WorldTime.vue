<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

type Zone = { id: string; label: string; flag?: string };

const zones = ref<Zone[]>([
  { id: "Asia/Shanghai", label: "北京", flag: "🇨🇳" },
  { id: "Asia/Tokyo", label: "东京", flag: "🇯🇵" },
  { id: "Europe/London", label: "伦敦", flag: "🇬🇧" },
  { id: "America/New_York", label: "纽约", flag: "🇺🇸" },
  { id: "Australia/Sydney", label: "悉尼", flag: "🇦🇺" },
  { id: "UTC", label: "UTC", flag: "🌐" },
]);

const now = ref<Date>(new Date());
const hour12 = ref(false);
let timer: number | null = null;

const formatTime = (tz: string) =>
  new Intl.DateTimeFormat("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: hour12.value,
    timeZone: tz,
  }).format(now.value);

const formatDate = (tz: string) =>
  new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    weekday: "short",
    timeZone: tz,
  }).format(now.value);

onMounted(() => {
  timer = window.setInterval(() => {
    now.value = new Date();
  }, 1000);
});

onUnmounted(() => {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
});
</script>

<template>
  <div class="world-time">
    <div class="header">
      <div class="title">世界时间</div>
      <div class="controls">
        <span class="label">24小时制</span>
        <label class="switch">
          <input type="checkbox" :checked="!hour12" @change="hour12 = !(($event.target as HTMLInputElement).checked)" />
          <span class="slider"></span>
        </label>
      </div>
    </div>
    <div class="grid">
      <div class="card" v-for="z in zones" :key="z.id">
        <div class="card-header">
          <span class="flag">{{ z.flag }}</span>
          <span class="city">{{ z.label }}</span>
        </div>
        <div class="time">{{ formatTime(z.id) }}</div>
        <div class="date">{{ formatDate(z.id) }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.world-time {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  height: 100%;
  padding: 12px;
  box-sizing: border-box;
  overflow: hidden;
  color: rgba(255, 255, 255, 0.92);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 0 0 auto;
}

.title {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label {
  font-size: 12px;
  opacity: 0.85;
}

.switch {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.25);
  transition: 0.2s;
  border-radius: 10px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  top: 2px;
  background-color: var(--primary-color);
  border-radius: 50%;
  transition: 0.2s;
}

.switch input:checked + .slider {
  background: rgba(255, 255, 255, 0.35);
}

.switch input:checked + .slider:before {
  transform: translateX(16px);
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(0, 1fr));
  gap: 12px;
  flex: 1 1 auto;
  overflow: hidden;
}

.card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  backdrop-filter: blur(8px);
  min-width: 0;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.flag {
  font-size: 16px;
}

.city {
  font-size: 14px;
  font-weight: 500;
}

.time {
  font-size: 22px;
  font-weight: 600;
  color: #ffffff;
}

.date {
  font-size: 12px;
  opacity: 0.85;
}
</style>