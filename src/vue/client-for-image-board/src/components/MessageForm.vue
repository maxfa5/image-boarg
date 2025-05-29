<template>
  <div class="message-form-container">
    <h3 class="form-title">{{ isNewThread ? 'Новый тред' : 'Ответить' }}</h3>
    <form @submit.prevent="submitForm" class="message-form">
      
      <div class="form-group">
        <!-- <label for="messageContent" class="form-label">Сообщение</label> -->
        <textarea
          id="messageContent"
          v-model="formData.content"
          class="form-textarea"
          rows="5"
          placeholder="Напишите ваше сообщение..."
          required
        ></textarea>
      </div>

      <div class="form-group-2">
        <div class="image-upload-container">
          <input
            type="file"
            ref="fileInput"
            accept="image/*"
            multiple
            @change="handleFileUpload"
            class="file-input"
            style="display: none"
          >
          <button
            type="button"
            @click="$refs.fileInput.click()"
            class="upload-button"
            :disabled="uploadedImages.length >= 3"
          > Добавить изображения
          </button>
          <div v-if="uploadedImages.length > 0" class="image-preview-container">
            <div v-for="(image, index) in uploadedImages" :key="index" class="image-preview">
              <img :src="image.preview" class="preview-image">
              <button @click="removeImage(index)" class="remove-image-btn">
                &times;
              </button>
            </div>
          </div>
        </div>
        <div class="form-actions">
        <button type="submit" class="submit-btn" :disabled="isSubmitting">
          <span v-if="isSubmitting">Отправка...</span>
          <span v-else>{{ isNewThread ? 'Создать тред' : 'Добавить' }}</span>
        </button>
      </div>
      </div>

      
    </form>
  </div>
</template>

<script>
import {getCookie} from '../utils/cookies.js'
import {isAuthenticated} from '../utils/cookies.js'

export default {
  props: {
    threadId: {
      type: String,
      default: null
    },
    isNewThread: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      formData: {
        threadTitle: '',
        content: ''
      },
      uploadedImages: [],
      isSubmitting: false,
      error: null
    }
  },
  methods: {
    handleFileUpload(event) {
      const files = event.target.files
      if (files.length + this.uploadedImages.length > 3) {
        alert('Максимальное количество изображений - 3')
        return
      }

      for (let i = 0; i < files.length; i++) {
        if (this.uploadedImages.length >= 3) break
        
        const reader = new FileReader()
        reader.onload = (e) => {
          this.uploadedImages.push({
            file: files[i],
            preview: e.target.result
          })
        }
        reader.readAsDataURL(files[i])
      }
    },
    removeImage(index) {
      this.uploadedImages.splice(index, 1)
    },
    async submitForm() {
      console.log(isAuthenticated())
      if(isAuthenticated()==false){
        alert("Пользователь не аутентифицирован!")
        return;
      }
  this.isSubmitting = true;
  this.error = null;

  try {
    const token = getCookie("auth_token");
    // Подготовка данных в новом формате
    let username = getCookie("user_data");
    console.log(username);
    const requestData = {
      action: "create",
      model: "messages",
      data: {
        post_id: "",
        thread_id: this.threadId || "",
        author_id: username, // Здесь нужно использовать реальный ID пользователя
        content: this.formData.content,
        images: this.uploadedImages.map(img => ({
          url: URL.createObjectURL(img.file), // Или загрузить файл на сервер и получить URL
          hash: this.generateFileHash(img.file) // Нужно реализовать функцию хеширования
        })),
        is_thread_root: this.isNewThread
      }
    };

    const endpoint = '/api/push';

    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json' // Указываем тип контента
      },
      body: JSON.stringify(requestData) // Преобразуем объект в JSON
    });

    if (!response.ok) {
      throw new Error('Ошибка при отправке сообщения');
    }

    const result = await response.json();
    this.resetForm();
    this.$emit('message-created', result); // Раскомментируйте, если нужно

  } catch (err) {
    console.error('Ошибка:', err);
    this.error = err.message;
  } finally {
    this.isSubmitting = false;
  }
},
    resetForm() {
      this.formData = {
        threadTitle: '',
        content: ''
      }
      this.uploadedImages = []
    },
    cancelReply() {
      this.$emit('cancel')
    }
  }
}
</script>

<style scoped>
.message-form-container {
  background-color: #E3e3ee;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.form-title {
  margin: 0 auto;
  margin-bottom: 10px;
  color: #2c3e50;
  font-size: 1.25rem;
  text-align: center;
  font-size: 35px;
}

.message-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.form-group-2 {
  display: flex;
  flex-direction: row;
  align-items:center ;
  justify-content: space-between;
}
.form-label {
  font-weight: 500;
  color: #4a5568;
}

.form-input,
.form-textarea {
  padding: 10px 12px;
  background-color: #E3e3ee;
  border: 1px solid #8C8B8B;
  border-radius: 15px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.image-upload-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.upload-button {
  padding: 10px 15px;
  background-color: #E3e3ee;
  color: #8C8B8B;
  border: 1px dashed #8C8B8B;
  border-radius: 15px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.upload-button:hover {
  background-color: #edf2f7;
  border-color: #a0aec0;
}

.upload-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.upload-icon {
  font-size: 1.2rem;
}

.image-preview-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}

.image-preview {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 15px;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-image-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 24px;
  height: 24px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
}

.remove-image-btn:hover {
  background-color: rgba(0, 0, 0, 0.8);
}

.form-actions {
  display: flex;
  gap: 10px;
}

.submit-btn {
  padding: 10px 20px;
  background-color: #E3e3ee;
  color: #733671;
  border: #8C8B8B;
  border-style:  solid;
  border-radius: 15px;
  border-width: 1px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover {
  background-color: #3182ce;
}

.submit-btn:disabled {
  background-color: #a0aec0;
  cursor: not-allowed;
}

.cancel-btn {
  padding: 10px 20px;
  background-color: #f7fafc;
  color: #4a5568;
  border: 1px solid #e2e8f0;
  border-radius: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.cancel-btn:hover {
  background-color: #edf2f7;
}

.error-message {
  color: #e53e3e;
  margin-top: 10px;
  font-size: 0.9rem;
}

@media (max-width: 600px) {
  .message-form-container {
    padding: 15px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .submit-btn,
  .cancel-btn {
    width: 100%;
  }
}
</style>