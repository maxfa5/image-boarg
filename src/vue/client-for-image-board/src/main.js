import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Правильный порядок:
const app = createApp(App)  // 1. Создаем приложение

const pinia = createPinia()  // 2. Создаем экземпляр Pinia

app.use(pinia)  // 3. Устанавливаем Pinia
app.use(router)  // 4. Устанавливаем роутер

app.mount('#app')  // 5. Монтируем приложение