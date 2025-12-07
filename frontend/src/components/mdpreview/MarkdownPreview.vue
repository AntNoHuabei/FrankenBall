<script setup lang="ts">
import {ref, onMounted, onBeforeUnmount, nextTick} from 'vue';
import  'table-of-contents-element';
import {Button,message} from "ant-design-vue";
import { createMarkdownExit } from 'markdown-exit'
import { codeToHtml } from 'shiki'
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'
import {BookOpen} from 'lucide-vue-next'
// factory helper
const md = createMarkdownExit({
  html:true,
  linkify:true,
  async highlight(code, lang) {
    return codeToHtml(code, {lang, theme: 'nord'});
  }
})

//TODO 处理链接点击事件,从默认浏览器打开

const mdContent = ref<string>('');
const isFileOpened = ref<boolean>(false);


const openFile =async ()=>{
  try {
    // 打开原生文件选择器
    const [fileHandle] = await showOpenFilePicker({
      types: [{
        description: 'Markdown文件',
        accept: {
          'text/markdown': ['.md'],
        },
      }],
    });
    // 获取文件对象
    const file = await fileHandle.getFile();
    // 读取文件内容为文本
    const content = await file.text();
    mdContent.value = await md.renderAsync( content)
    isFileOpened.value = true;
    message.success('文件已打开')

    await nextTick(() => {
      setSlugs();
      // @ts-ignore
      toc.value?.render()
    })
  } catch (err:any) {
    console.error('无法读取文件:', err);
    // 用户取消选择文件时不显示错误消息
    if (err.name !== 'AbortError') {
      message.error('无法打开文件');
    }
  }
}
import 'github-markdown-css'


const toc = ref<any>(null);

function slugify(s:string) {
  return encodeURIComponent(String(s).trim().toLowerCase().replace(/\s+/g, '-'));
};
function setSlugs() {
  document.querySelectorAll('.markdown-body :is(h1, h2, h3, h4, h5, h6)').forEach(item => {
    item.id = slugify(item.textContent || item.innerHTML)
  });
}

</script>

<template>
  <div class="markdown-preview-container">
    <div v-if="!isFileOpened" class="empty-state">
      <div class="empty-content">
        <h3>欢迎使用 Markdown 预览器</h3>
        <p>请打开一个 Markdown 文件开始预览</p>
        <Button @click="openFile" type="primary">选择文件</Button>
      </div>
    </div>
    <splitpanes v-else class="default-theme markdown-content-container">
      <pane size="20" style="overflow-y: auto">
        <table-of-contents ref="toc" selector=".markdown-body :is(h1, h2, h3, h4, h5, h6)">
          <header>目录</header>
          <div data-toc-render-target></div>
        </table-of-contents>
      </pane>
      <pane style="overflow:auto;">
        <div v-html="mdContent" class="markdown-body"></div>
      </pane>
    </splitpanes>
    
    <!-- 悬浮按钮 -->
    <div 
      v-if="isFileOpened" 
      class="floating-button-container"
      @click="openFile"
    >
      <BookOpen :size="40" class="floating-toggle">
      </BookOpen>
    </div>
  </div>
</template>

<style lang="less" scoped>
.markdown-preview-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  padding: 10px;
  box-sizing: border-box;
  position: relative;
  :deep(.default-theme.splitpanes .splitpanes__pane){
    background-color: #00000000;
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 15px;
    margin-top: 10px;
    border: 2px dashed rgba(255, 255, 255, 0.1);
    
    .empty-content {
      text-align: center;
      color: #ffffffcc;
      
      h3 {
        margin-bottom: 10px;
        font-size: 1.5em;
        color: #ffffff;
      }
      
      p {
        margin-bottom: 20px;
        font-size: 1.1em;
      }
    }
  }
  
  .markdown-content-container{
    width: 100%;
    flex: 1;
    overflow: auto;
    display: flex;
    flex-direction: row;
    border-bottom-left-radius: 15px;
    border-bottom-right-radius: 15px;

    /* 基础链接样式 - 近白色文字 */
    :deep(a) {
      color: #ffffff; /* 纯白 */
      text-decoration: none;
      position: relative;
      transition: all 0.3s ease;
      font-weight: 500;
      padding: 4px 8px;
    }

    /* 鼠标悬停效果 - 轻微透明度变化 */
    :deep(a:hover) {
      color: #f7fafc; /* 极浅灰白 */
      background-color: rgba(255, 255, 255, 0.1);
      border-radius: 4px;
      transform: translateY(-1px);
    }

    /* 下划线动画效果 */
    :deep(a::after) {
      content: '';
      position: absolute;
      width: 0;
      height: 2px;
      bottom: -2px;
      left: 50%;
      background-color: #ffffff;
      transition: all 0.3s ease;
      transform: translateX(-50%);
    }

    :deep(a:hover::after) {
      width: 80%;
    }
  }
  
  // 悬浮按钮样式
  .floating-button-container {
    position: absolute;
    bottom: 20px;
    right: 20px;
    z-index: 11000;
    
    .floating-toggle {
      padding: 5px;
      box-sizing: border-box;
      border-radius: 50%;
      background-color: rgba(255, 255, 255, 0.1);
      backdrop-filter: blur(10px);
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      color: var(--primary-color);
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
      transition: all 0.3s ease;
      
      &:hover {
        background-color: rgba(255, 255, 255, 0.2);
        transform: scale(1.1);
      }
    }
    
    .floating-menu {
      position: absolute;
      top: 50px;
      right: 0;
      width: 150px;
      background-color: rgba(30, 30, 30, 0.9);
      border-radius: 8px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
      padding: 10px;
      backdrop-filter: blur(10px);
      animation: fadeIn 0.2s ease;
      
      .floating-menu-item {
        margin-bottom: 5px;
        
        &:last-child {
          margin-bottom: 0;
        }
      }
    }
    
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(-10px); }
      to { opacity: 1; transform: translateY(0); }
    }
  }
}
</style>