<script setup lang="ts">
import {computed, ref} from 'vue'
import {useTodo} from '../../store/useTodo'
import {storeToRefs} from 'pinia'
import {Bell, Plus, Trash2} from 'lucide-vue-next'
import Item from "./Item.vue";

const todoStore = useTodo()
const {todos} = storeToRefs(todoStore)

// 新增任务输入
const newTodoTitle = ref('')
const newReminderTime = ref<string>('')

// 任务统计文案
const taskCountText = computed(() => `${todos.value.length} 个任务`)

// 添加任务
const addTodo = () => {
  const title = newTodoTitle.value.trim()
  if (!title) return

  todoStore.addTodo({
    title,
    detail: '',
    priority: 'medium',
    remind: !!newReminderTime.value,
    remindTime: newReminderTime.value ? new Date(newReminderTime.value) : null,
    completed: false,
  })
  newTodoTitle.value = ''
  newReminderTime.value = ''
}

// 切换完成状态
const toggleCompleted = (id: string, value: boolean | undefined) => {
  todoStore.updateTodo(id, {completed: !value})
}

</script>

<template>
  <div class="todo-list-container">

    <!-- 待办列表区域 -->
    <div class="reminder-list">
      <div class="list-header">
        <h2 class="list-title">待办列表</h2>
        <span class="task-count">{{ taskCountText }}</span>
      </div>

      <div v-if="todos.length > 0" class="list">
        <Item v-for="todo in todos" :todo="todo" :key="todo.id">
        </Item>
      </div>

      <!-- 空状态提示 -->
      <div v-else class="empty-state">
        <div class="empty-icon">📥</div>
        <p>暂无任务，开始添加您的第一个任务吧</p>
      </div>
    </div>

    <!-- 添加待办 -->
    <div class="add-task-container">
      <form @submit.prevent="addTodo" class="add-form">
        <input
            type="text"
            v-model="newTodoTitle"
            class="add-input"
            placeholder="添加新任务..."
            required
        />
        <Plus class="add-btn"/>
      </form>
    </div>
  </div>
</template>

<style lang="less" scoped>
.todo-list-container {
  width: 100%;
  height: 100%;
  display: flex;
  box-sizing: border-box;
  flex-direction: column;
  color: #e0e0e0;
  font-family: 'Inter', sans-serif;
}

.header {
  padding: 16px 20px 0 20px;
  text-align: center;
}

.title {
  margin: 0;
  color: #ffffff;
  font-size: 22px;
  font-weight: 700;
}

.subtitle {
  margin: 6px 0 0 0;
  color: #aaaaaa;
  font-size: 13px;
}

.reminder-list {
  overflow: auto;
  padding: 16px 20px 20px 20px;
  flex: 1;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.list-title {
  margin: 0;
  color: #ffffff;
  font-weight: 600;
}

.task-count {
  color: #c7c7c7;
  font-size: 12px;
  background: #333F51;
  padding: 4px 10px;
  border-radius: 999px;
}

.list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.empty-state {
  text-align: center;
  color: #aaaaaa;
  font-style: italic;
  padding: 30px 0;
}

.empty-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.add-task-container {
  padding: 12px 12px 16px 12px;
  margin-top: auto;
}

.add-form {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px;
}

.add-input {
  flex: 1;
  background-color: rgba(0, 0, 0, 0.2);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  padding: 8px 10px;
}

.add-input:focus {
  outline: none;
  border-color: #5c7cfa;
  box-shadow: 0 0 0 2px rgba(92, 124, 250, 0.3);
}

.toggle-bell {
  width: 14px;
  height: 14px;
  color: #FAAD14;
  margin-right: 6px;
}

.toggle-datetime {
  background: transparent;
  border: none;
  outline: none;
  color: #fff;
  font-size: 12px;
}

.add-btn {
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
}

.add-btn:hover {


}

/* 滚动条样式优化 */
.reminder-list::-webkit-scrollbar {
  width: 6px;
}

.reminder-list::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.reminder-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.reminder-list::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>