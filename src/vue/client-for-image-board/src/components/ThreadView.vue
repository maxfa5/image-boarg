<template>
  <div>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-else-if="loadingMessages">Загрузка сообщений...</div>
    <div v-else>
      <h2>Сообщения треда #{{ threadId.slice(0, 8) }}...</h2>
      <div v-for="message in messages" :key="message.id" class="message">
        <div class="message-content">{{ message.content }}</div>
        <div class="message-meta">
          <span class="message-date">{{ formatDate(message.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      messages: [],
      threadId: '',
      error: null,
      loadingMessages: true
    };
  },
  created() {
    this.threadId = this.$route.params.id;
    this.fetchMessages();
  },
  methods: {
    async fetchMessages() {
      try {
        const response = await fetch(`http://localhost/api/messages/${this.threadId}`);
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
      return new Date(isoString).toLocaleString();
    }
  }
};
</script>

<style scoped>
.message {
  margin: 15px 0;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 5px;
}
.message-content {
  margin-bottom: 5px;
}
.message-meta {
  font-size: 0.8em;
  color: #7f8c8d;
}
</style>