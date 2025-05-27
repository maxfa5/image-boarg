
<template>
  <div>
    <button @click="fetchMessages">Получить сообщения</button>
    <div v-if="loading">Загрузка...</div>
    <div v-if="error" style="color: red;">Ошибка: {{ error }}</div>
    <div v-if="messages">
      <h3>Сообщения:</h3>
      <pre>{{ messages }}</pre>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    async fetchMessages() {
      try {
        const response = await fetch('http://localhost/api/messages');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        this.messages = await response.json();
        console.log(this.messages)
      } catch (err) {
        console.error('Ошибка:', err);
        this.error = err.message;
      }
    }
  }
}
</script>