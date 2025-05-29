<template>
  <div class="header-section">
      <div class="loginAndRegister">
      <button 
        v-if="!isAuthenticated"
        @click="goToLogin"
        class="login-button"
      >
        <i class="icon-login"></i> Войти
      </button>
            
      <div v-if="isAuthenticated">
        <button @click="logout" class="logout-button">Выйти</button>
      </div>
      <button 
        @click="goToRegister"
        class="login-button"
      >
        <i class="icon-login"></i> Регистрация
      </button>
    </div>
    </div>
  <div class="threads-container">
    
    <h1 class="page-title">Треды</h1>
    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-else-if="loading" class="loading-indicator">
      <div class="spinner"></div>
      <span>Загрузка тредов...</span>
    </div>
    
    <div v-else class="threads-list">

      <div v-for="thread in threads" :key="thread.thread_id" class="thread-item">
        <router-link 
          :to="{
            path: '/thread/' + thread.thread_id,
            query: { 
              title: thread.title,
              is_closed: thread.is_closed 
            }
          }"
          class="thread-link"
        >
          <div class="thread-content">
            <span class="thread-title">{{ thread.title }}</span>
            <div class="thread-meta">
              <span class="thread-date">
                <i class="icon-clock"></i>
                {{ formatDate(thread.created_at) }}
              </span>
              <span v-if="thread.is_closed" class="closed-badge">
                <i class="icon-lock"></i>
                Закрыт
              </span>
            </div>
          </div>
        </router-link>
      </div>
    </div>
    <MessageForm 
    :thread-id=threadId
    :isNewThread=true 
    @message-created="handleNewMessage"
    @cancel="closeReplyForm"
  />
    <News></News>
</div>
  <div class="footer"></div>
</template>

<script>
import { isAuthenticated } from '../utils/cookies';
import { deleteCookie } from '../utils/cookies';

import MessageForm from './MessageForm.vue';
import News from './News.vue';

export default {
  components: {
    MessageForm,
    News
  },
  data() {
    return {
      threads: [],
      error: null,
      loading: true,
      isAuthenticated: false
    };
  },
  async created() {
    this.checkAuthStatus();
    await this.fetchThreads();
  },
  methods: {
    checkAuthStatus() {
      this.isAuthenticated = isAuthenticated();
    },
    async fetchThreads() {
      try {
        const response = await fetch('/api/threads');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        this.threads = await response.json();
      } catch (err) {
        console.error('Ошибка:', err);
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    formatDate(isoString) {
      return new Date(isoString).toLocaleString('ru-RU', {
        day: 'numeric',
        month: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }, goToLogin() {
      this.$router.push('/login');
    }, goToRegister() {
      this.$router.push('/registration');
    }, logout() {
      this.isAuthenticated = false; // Обновите состояние аутентификации'
      deleteCookie('user_data');
    }
  }
};
</script>
<style scoped>
.threads-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}
html, body{
  height: 100%;
  margin: 0; 
  box-sizing: border-box;
}
.footer{
  background-image: url('@/assets/footer.png'); /* Укажите путь к изображению */
  background-size: cover; /* Масштабирование изображения для заполнения контейнера */
  background-position: center; /* Центрирование изображения */
  height: 250px; /* Установите высоту для футера */
  width: 100%; /* Занимаем всю ширину */
}
.header-section {
  background-image: url("@/assets/header.png");
  background-size: cover; /* Масштабирование изображения для заполнения контейнера */
  background-position: center; /* Центрирование изображения */
  height: 331px; /* Установите высоту для футера */
  width: 100%; /* Занимаем всю ширину */
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  text-align: center;
  font-size: 30px;
  color: #333;
}

.login-button {
  background-color: #4a76a8;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  margin-left: 10px;
  align-items: center;
  gap: 5px;
  transition: background-color 0.3s;
}

.login-button:hover {
  background-color: #3a6598;
}

.create-thread-section {
  margin-top: 30px;
  text-align: center;
}
.loginAndRegister{
  margin: 10px;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  width: 100%;
}

.create-button {
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.3s;
}

.create-button:hover {
  background-color: #218838;
}

.threads-list {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  padding: 20px;
}

.thread-item {
  align-items: center;
  background-color: #D9D9D9;
  width: 140px;
  margin: 20px;
  border: 1px solid #e0e0e0;
  border-radius: 15px;
  overflow: hidden;
  transition: box-shadow 0.3s;
  
}

.thread-item:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.thread-link {
  display: block;
  text-decoration: none;
  color: inherit;
  text-align: center;
  padding: 15px;
}

.thread-content {
  display: flex;
  flex-direction: column;
}

.thread-title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 8px;
   display: -webkit-box;
  -webkit-line-clamp: 2; /* Количество строк */
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.thread-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #666;
}
.logout-button {
  background-color: #f44336; /* Красный цвет для кнопки выхода */
  color: white;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  border-radius: 5px;
}

.logout-button:hover {
  background-color: #d32f2f; /* Темнее при наведении */
}
.thread-date, .closed-badge {
  display: flex;
  align-items: center;
  gap: 4px;
}

.closed-badge {
  color: #dc3545;
}

.icon-clock, .icon-lock, .icon-login, .icon-plus {
  font-size: 14px;
}

.loading-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 30px;
}

.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-left-color: #4a76a8;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  background-color: #fee2e2;
  color: #ef4444;
  padding: 15px;
  border-radius: 4px;
  text-align: center;
  margin: 20px 0;
}
</style>