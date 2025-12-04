<script setup lang="ts">
import Vditor from 'vditor'
import "vditor/src/assets/less/index.less"
import {onBeforeUnmount, onMounted, ref} from "vue";
const editor = ref<Vditor|null>(null)
onMounted(()=>{
   editor.value = new Vditor("memo-editor", {
    height: "100%",
    toolbarConfig: {
      pin: true,
    },
    cache: {
      enable: false,
    },
  })
})
onBeforeUnmount(()=>{
  if(editor.value){
    editor.value.destroy();
  }
})
</script>

<template>
  <div id="memo-editor" class="editor-root">
  </div>

</template>

<style scoped>
.editor-root{
  height: 100%;
  width: 100%;
}

</style>