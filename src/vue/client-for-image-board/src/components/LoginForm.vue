<template>
  <div class="login-container">
    <h2>Авторизация</h2>
    <form @submit.prevent="login">
      <div class="form-group">
        <label for="username">Имя пользователя:</label>
        <input 
          type="text" 
          id="username" 
          v-model="credentials.username" 
          required
          class="form-control"
        />
      </div>
      
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input 
          type="password" 
          id="password" 
          v-model="credentials.password" 
          required
          class="form-control"
        />
      </div>
      
      <button 
        type="submit" 
        class="btn-submit"
        :disabled="loading"
      >
        <span v-if="loading">Вход...</span>
        <span v-else>Войти</span>
      </button>
      
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </form>
  </div>
</template>

<script>
import axios from 'axios';
import { setCookie } from '../utils/cookies.js';

export default {
  name: 'LoginForm',
  data() {
    return {
      credentials: {
        username: '',
        password: ''
      },
      loading: false,
      error: null
    };
  },
  methods: {
    async login() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await axios.post(
          '/api/auth/login', 
          {
            username: this.credentials.username,
            password: this.credentials.password
          },
          {
            headers: {
              'Content-Type': 'application/json'
            }
          }
        );
        
        // Сохраняем данные в куки
        setCookie('auth_token', response.data.token, 10); 
        setCookie('user_data', (response.data.username), 10);
        
        // Перенаправляем пользователя
        this.$router.push('/');
        
      } catch (err) {
        console.error('Ошибка авторизации:', err);
        this.error = err.response?.data?.message || 'Неверные учетные данные';
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.btn-submit {
  width: 100%;
  padding: 0.75rem;
  background-color: #4a76a8;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-submit:hover {
  background-color: #3a6598;
}

.btn-submit:disabled {
  background-color: #a0aec0;
  cursor: not-allowed;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #fee2e2;
  color: #ef4444;
  border-radius: 4px;
  text-align: center;
}
</style>