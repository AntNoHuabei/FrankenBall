import { computed, ref } from "vue";

const sharedToolbarRef = ref<Record<string, HTMLElement>>({});

export const useSharedToolbarRef = (winId: string) => {
  const toolbarRef = computed(() => sharedToolbarRef.value[winId]);

  const setToolbarRef = (ref: HTMLElement) => {
    sharedToolbarRef.value[winId] = ref;
  };

  const removeToolbarRef = () => {
    delete sharedToolbarRef.value[winId];
  };

  return {
    toolbarRef,
    setToolbarRef,
    removeToolbarRef,
  };
};
