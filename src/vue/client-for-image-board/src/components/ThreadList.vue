<template>
  <div>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-else-if="loading">Загрузка тредов...</div>
    <div v-else>
      <div v-for="thread in threads" :key="thread.thread_id" class="thread-link">
        <router-link :to="'/thread/' + thread.thread_id" @click.native="fetchThreadMessages(thread.thread_id)">
          Тред #{{ thread.thread_id.slice(0, 8) }}...
          <span v-if="thread.is_closed" class="closed-badge">(закрыт)</span>
          <span class="created-at">{{ formatDate(thread.created_at) }}</span>
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
        const response = await fetch('http://localhost/api/threads');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        this.threads = await response.json();
      } catch (err) {
        console.error('Ошибка:', err);
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    async fetchThreadMessages(threadId) {
      this.$router.push(`/thread/${threadId}`);
    },
    formatDate(isoString) {
      return new Date(isoString).toLocaleString();
    }
  }
};
</script>

<style scoped>
/* Стили остаются такими же */
</style>