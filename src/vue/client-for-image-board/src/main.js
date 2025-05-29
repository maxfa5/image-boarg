import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Правильный порядок:
const app = createApp(App)  // 1. Создаем приложение


app.use(router)  // 4. Устанавливаем роутер

app.mount('#app')  // 5. Монтируем приложение