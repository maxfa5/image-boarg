import { createRouter, createWebHistory } from 'vue-router'
import ThreadList from './components/ThreadList.vue'
import ThreadView from './components/ThreadView.vue'
import LoginForm from './components/LoginForm.vue'
import Registration from './components/Registration.vue'

const routes = [
  { path: '/', component: ThreadList },
  { path: '/thread/:id',  component: ThreadView, props: route => ({ title: route.query.title }) },
  { path: '/login',  component: LoginForm },
  { path: '/registration',  component: Registration }

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router