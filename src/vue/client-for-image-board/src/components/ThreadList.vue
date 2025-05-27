<template>
  <div>
    <div v-for="thread in threads" :key="thread.thread_id" class="thread-item">
      <router-link 
        :to="{
          path: '/thread/' + thread.thread_id,
          query: { 
            title: thread.title,
            // Можно передать другие параметры, если нужно
            is_closed: thread.is_closed 
          }
        }"
        class="thread-link"
      >
        <span class="thread-title">{{ thread.title }}</span>
        <span v-if="thread.is_closed" class="closed-badge">(закрыт)</span>
        <span class="thread-date">{{ formatDate(thread.created_at) }}</span>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.thread-title {
  text-decoration: none;
  color: #2c3e50;
  font-weight: 500;
}
.thread-title:hover {
  color: #42b983;
  text-decoration: underline;
}
</style>

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
