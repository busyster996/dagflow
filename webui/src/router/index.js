import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    redirect: '/tasks',
    children: [
      {
        path: '/tasks',
        name: 'Tasks',
        component: () => import('@/views/TaskManagement.vue'),
      },
      {
        path: '/pipelines',
        name: 'Pipelines',
        component: () => import('@/views/PipelineManagement.vue'),
      },
      {
        path: '/upload',
        name: 'FileUpload',
        component: () => import('@/views/FileUpload.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router