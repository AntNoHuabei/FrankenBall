import {defineStore} from "pinia";
import {KEY_TODO} from "./keys";
import {ref} from "vue";
import type {Todo} from "../types/todo";

export const useTodo = defineStore(KEY_TODO, () => {
    const todos = ref<Todo[]>([]);

    function addTodo(reminder: Omit<Todo, 'id' | 'createdTime'>) {
        todos.value.push({
            ...reminder,
            id: Math.random().toString(36).substring(2, 9),
            createdTime: new Date()
        });
    }

    function removeTodo(id: string) {
        const index = todos.value.findIndex(r => r.id === id);
        if (index !== -1) {
            todos.value.splice(index, 1);
        }
    }

    function updateTodo(id: string, updates: Partial<Todo>) {
        const reminder = todos.value.find(r => r.id === id);
        if (reminder) {
            Object.assign(reminder, updates);
        }
    }

    return {
        todos,
        addTodo,
        removeTodo,
        updateTodo
    };
});