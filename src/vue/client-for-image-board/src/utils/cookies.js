// Установка куки
export const setCookie = (name, value, days = 365) => {
  const date = new Date();
  date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
  const expires = `expires=${date.toUTCString()}`;
  document.cookie = `${name}=${encodeURIComponent(value)};${expires};path=/;SameSite=Lax`;
};

// Получение куки
export const getCookie = (name) => {
  const nameEQ = `${name}=`;
  const ca = document.cookie.split(';');
  
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1);
    if (c.indexOf(nameEQ) === 0) {
      return decodeURIComponent(c.substring(nameEQ.length));
    }
  }
  
  return null;
};

// Удаление куки
export const deleteCookie = (name) => {
  document.cookie = `${name}=; Max-Age=0; path=/;`;
};

// Получение данных пользователя
export const getUserData = () => {
  const userData = getCookie('user_data');
  return userData ? JSON.parse(userData) : null;
};

// Получение токена
export const getAuthToken = () => {
  return getCookie('auth_token');
};

// Проверка авторизации
export const isAuthenticated = () => {
  console.log(getCookie('user_data'));
  return getCookie('user_data')!=null
};

// Выход из системы
export const logout = () => {
  deleteCookie('auth_token');
  deleteCookie('user_data');
};