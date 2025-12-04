<script setup lang="ts">
import {ref} from "vue";
import {Todo} from "../../types/todo";
import {useTodo} from "../../store/useTodo";
import {Bell, Plus, Trash2,X} from 'lucide-vue-next'
const todoStore = useTodo()
const todo = defineModel<Todo>("todo");

const toggleTodoCompleted = () => {
  console.log('toggleTodoCompleted', todo.value?.completed)
  if (todo.value) {
    todo.value.completed = !todo.value.completed
  }
}

// 工具: 格式化提醒文本
const formatReminderLabel = (date: Date | null | undefined) => {
  if (!date) return '未设置提醒'
  const d = new Date(date)
  const now = new Date()
  const isToday = d.toDateString() === now.toDateString()
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)
  const hhmm = `${pad(d.getHours())}:${pad(d.getMinutes())}`
  if (isToday) return `今天 ${hhmm} 提醒`
  const y = d.getFullYear()
  const m = pad(d.getMonth() + 1)
  const day = pad(d.getDate())
  return `${y}-${m}-${day} ${hhmm} 提醒`
}

// 删除任务
const deleteTodo = () => {
  if(todo.value){

    todoStore.removeTodo(todo.value.id)
  }
}

// 设置/更新提醒时间（列表中行内设置）
const onItemReminderChange = (value: string) => {
  if(todo.value){
    const hasValue = !!value
    todoStore.updateTodo(todo.value.id, {
      remind: hasValue,
      remindTime: hasValue ? new Date(value) : null,
    })
  }

}


// 工具: 将 Date 转为 datetime-local 可用值
const toInputValue = (date: Date | null|undefined) => {
  if (!date) return ''
  const d = new Date(date)
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)
  const yyyy = d.getFullYear()
  const mm = pad(d.getMonth() + 1)
  const dd = pad(d.getDate())
  const hh = pad(d.getHours())
  const mi = pad(d.getMinutes())
  return `${yyyy}-${mm}-${dd}T${hh}:${mi}`
}

const dtPicker = ref<HTMLInputElement | null>(null);

const showDtPicker = () => {
  if (dtPicker.value) {
    dtPicker.value.value = toInputValue(todo.value?.remindTime)
    dtPicker.value.showPicker ()
  }
}

const clearRemindTime = () => {
  if(todo.value){
    todoStore.updateTodo(todo.value.id, {
      remind: false,
      remindTime: null,
    })
  }
}
</script>

<template>

  <div class="todo-item">
    <div class="todo-content">
      <div class="todo-checkbox" :class="{checked: todo?.completed}" @click="toggleTodoCompleted"></div>
      <div class="todo-text">{{ todo?.title }}</div>
      <div class="actions">
        <Bell v-show="!todo?.remind" :size="16" @click="showDtPicker"/>
        <X v-show="todo?.remind" :size="16" @click="clearRemindTime"/>
        <Trash2 @click="deleteTodo" :size="16" class="delete-icon" />
      </div>
    </div>
    <div class="action-section"
    v-show="todo?.remind"
    >
      <div class="reminder-info">

        <Bell  :size="16" class="reminder-icon" @click="showDtPicker"/>
        <span @click="showDtPicker">{{ formatReminderLabel(todo?.remindTime) }}</span>
        <input ref="dtPicker" type="datetime-local" style="display: none"
               @change="onItemReminderChange( ($event.target as HTMLInputElement).value)"
               :value="toInputValue(todo?.remindTime)">
      </div>
    </div>
  </div>

</template>


<style lang="less" scoped>

.todo-item {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  margin-bottom: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  position: relative;
  cursor: pointer;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.25);

    .todo-content{
      .actions{
        display: flex;
        flex-direction: row;
        justify-content: space-around;
        align-items: center;
        gap: 8px;
      }
    }
  }

  .action-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);

    .reminder-info {
      font-size: 12px;
      color: #c7c7c7;
      display: flex;
      align-items: center;
      gap: 4px;
    }


    .reminder-icon {
      color: #FAAD14;
      cursor: pointer;
    }
    .delete-icon{
      color: #c7c7c7;
      &:hover{
        color: red;
      }
    }

  }
  .todo-content {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 12px;

    .todo-text {
      flex: 1;
      font-size: 16px;
      line-height: 1.4;
      color: #ffffff;
      text-align: left;
    }


    .todo-checkbox {
      width: 20px;
      height: 20px;
      border: 2px solid rgba(255, 255, 255, 0.2);
      border-radius: 50%;
      cursor: pointer;
      flex-shrink: 0;
      position: relative;
      transition: all 0.3s ease;
      background: transparent;
    }

    .todo-checkbox.checked {
      background: #5c7cfa;
      border-color: #5c7cfa;

      &::after {
        content: '✓';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: #ffffff;
        font-size: 12px;
        font-weight: bold;
      }
    }

    .actions{
    
      display: none;
    }
  }


}

.todo-item.completed {
  opacity: 0.7;
  background: rgba(255, 255, 255, 0.05);
  
  .todo-text {
    text-decoration: line-through;
    color: #aaaaaa;
  }
  
}
</style>