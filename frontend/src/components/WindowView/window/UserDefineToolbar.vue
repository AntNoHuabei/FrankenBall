<script setup lang="ts">
import { inject, onMounted, ref, Ref } from "vue";
import { useSharedToolbarRef } from "../composables/useSharedToolbarRef";

const toolbarDom = ref<HTMLElement | null>(null);

onMounted(() => {
  const winId = inject<string>("winId");
  const { toolbarRef } = useSharedToolbarRef(winId!);
  if (toolbarRef.value) {
    toolbarDom.value = toolbarRef.value;
  }
});
</script>

<template>
  <Teleport v-if="toolbarDom" :to="toolbarDom">
    <slot></slot>
  </Teleport>
</template>

<style scoped></style>
