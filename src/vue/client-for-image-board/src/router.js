import { createRouter, createWebHistory } from 'vue-router'
import ThreadList from './components/ThreadList.vue'
import ThreadView from './components/ThreadView.vue'

const routes = [
  { path: '/', component: ThreadList },
  { path: '/thread/:id',  component: ThreadView, props: route => ({ title: route.query.title }) }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router