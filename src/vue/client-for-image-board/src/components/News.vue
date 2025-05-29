<template>
  <div >
    <h2 class="h2_News">Новости</h2>
    <div class="news-container">
        <div v-for="(news, index) in newsList" :key="index" class="news-item">
            <p class="news-text">{{ news.text }}</p>
        <span class="news-time">{{ news.time }}</span>
        </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      newsList: [
        { text: 'Ксюша не выпила пиво и от горя...', time: '18.04.25' },
        { text: 'Ксюша не выпила пиво и от горя...', time: '18.04.25' },
        { text: 'Погода сегодня солнечная и теплая.', time: '18.04.25' },
        { text: 'Новый фильм выходит в прокат завтра.', time: '18.04.25' },
        { text: 'Спортсмены готовятся к чемпионату.', time: '18.04.25' },
        { text: 'В парке открылась новая выставка.', time: '18.04.25' }
      ]
    };
  }
};
</script>

<style>
.news-container {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16px; /* Отступы между новостями */
  padding: 10px;
  justify-content: space-between;
}


.news-item {
  background-color: #D9D9D9;
  border-radius: 8px;
  padding: 12px 16px;
  box-shadow: 0 2px 4px rgba(44, 62, 80, 0.1);
  max-width: 250px;
  flex: 1 1 250px; /* Элементы растут и сжимаются, базовая ширина 250px */
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.h2_News{
    font-size: 30px;
    text-align: center;
}

.news-text {
  font-size: 1rem;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.news-time {
  font-size: 0.85rem;
  color: #7f8c8d;
  align-self: flex-end;
}
</style>
