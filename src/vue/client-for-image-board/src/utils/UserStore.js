// Утилиты для работы с куками
const setCookie = (name, value) => {
  const date = new Date();
  document.cookie = `${name}=${value};path=/`;
};

const getCookie = (name) => {
  const nameEQ = `${name}=`;
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
};

const deleteCookie = (name) => {
  setCookie(name, '', -1);
};

// Хранилище пользователя
const UserStore = {
  // Сохранить данные пользователя
  setUser(userData) {
    setCookie('user', JSON.stringify(userData));
  },

  // Получить данные пользователя
  getUser() {
    const user = getCookie('user');
    return user ? JSON.parse(user) : null;
  },

  // Получить никнейм
  getUsername() {
    const user = this.getUser();
    return user ? user.username : null;
  },

  // Проверить авторизацию
  isAuthenticated() {
    return !!this.getUser();
  },

  // Выход
  logout() {
    deleteCookie('user');
  }
};

export default UserStore;