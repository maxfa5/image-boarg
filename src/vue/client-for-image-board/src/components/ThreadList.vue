<template>
  <div class="threads-container">
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
  </div>
</template>

<script>
export default {
  data() {
    return {
      threads: [],
      error: null,
      loading: true
    };
  },
  async created() {
    await this.fetchThreads();
  },
  methods: {
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
    }
  }
};
</script>

<style scoped>
.threads-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.error-message {
  padding: 15px;
  background-color: #ffebee;
  color: #c62828;
  border-radius: 4px;
  margin-bottom: 20px;
  border-left: 4px solid #c62828;
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 20px;
  color: #666;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top-color: #42b983;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.threads-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.thread-item {
  transition: transform 0.2s;
}

.thread-item:hover {
  transform: translateX(5px);
}

.thread-link {
  display: block;
  text-decoration: none;
  color: inherit;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.2s;
}

.thread-link:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background-color: #f8f9fa;
}

.thread-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.thread-title {
  font-size: 1.1rem;
  font-weight: 500;
  color: #2c3e50;
}

.thread-meta {
  display: flex;
  gap: 15px;
  font-size: 0.85rem;
  color: #7f8c8d;
}

.thread-date {
  display: flex;
  align-items: center;
  gap: 4px;
}

.closed-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #e53935;
  font-weight: 500;
}

.icon-clock::before {
  content: "🕒";
}

.icon-lock::before {
  content: "🔒";
}

@media (max-width: 600px) {
  .threads-container {
    padding: 10px;
  }
  
  .thread-link {
    padding: 12px;
  }
  
  .thread-meta {
    flex-direction: column;
    gap: 5px;
  }
}
</style>