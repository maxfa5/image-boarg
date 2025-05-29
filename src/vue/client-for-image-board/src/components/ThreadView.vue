<template>
  <div>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-else-if="loadingMessages">Загрузка сообщений...</div>
    <div v-else>
      <h2 class="thread_msg">Сообщения треда </h2>
      <div v-for="message in messages" :key="message.post_id" class="message">
        <div class="message-author" v-if="message.author_id">
          Автор: {{ message.author_id }}
        </div>
        <div class="message-content">{{ message.content }}</div>
        <div class="message-meta">
          <span class="message-date">{{ formatDate(message.timestamp) }}</span>
        </div>
      </div>
    </div>
    <MessageForm 
    :thread-id=threadId 
    @message-created="handleNewMessage"
    @cancel="closeReplyForm"
  />
  </div>
</template>


<script>
import MessageForm from './MessageForm.vue';
export default {
  components: {
    MessageForm
  },
  data() {
    return {
      messages: [],
      title: '',
      threadId: '',
      error: null,
      loadingMessages: true
    };
  },
  created() {
    this.threadId = this.$route.params.id;
    this.title = this.$route.query.title || 'Без названия';
    this.fetchMessages();
  },
  methods: {
    async fetchMessages() {
      try {
        const response = await fetch(`/api/messages/${this.threadId}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        this.messages = await response.json();
      } catch (err) {
        console.error('Ошибка:', err);
        this.error = err.message;
      } finally {
        this.loadingMessages = false;
      }
    },
    formatDate(isoString) {
      if (typeof isoString !== 'string') return '';
      try {
        const normalizedDate = isoString.replace(/(\.\d{3})\d+Z$/, '$1Z');
        const date = new Date(normalizedDate);
        return date.toLocaleString('ru-RU', {
          day: '2-digit',
          month: '2-digit',
          year: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit'
        });
      } catch (e) {
        console.error('Ошибка форматирования даты:', e);
        return isoString.split('.')[0].replace('T', ' ');
      }
    }
  }
};
</script>

<style scoped>
.error {
  color: red;
  padding: 10px;
  margin-bottom: 15px;
}
.message {
  margin: 15px 0;
  padding: 15px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #fff;
}
.message-author {
  font-weight: bold;
  margin-bottom: 8px;
  color: #8B5E88;
}
.message-content {
  margin-bottom: 10px;
  white-space: pre-line;
}
.thread_msg{
  font-size: 30px;
  text-align: center;
}
.message-meta {
  font-size: 0.8em;
  color: #7f8c8d;
  display: flex;
  justify-content: space-between;
}
</style>