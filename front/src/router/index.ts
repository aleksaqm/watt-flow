import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import HouseholdView from '@/views/HouseholdView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import HomeView from '@/views/HomeView.vue'
import SuperAdminChangePasswordView from '@/views/SuperAdminChangePasswordView.vue'
import ManageAdminsView from '@/views/ManageAdminsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/home',
      name: 'home',
      component: HomeView
    },
    {
      path: '/household',
      name: 'household',
      component: HouseholdView
    },
    {
      path: '/superadmin',
      name: 'superadmin',
      component: SuperAdminChangePasswordView
    },
    {
      path: '/manage/admins',
      name: 'manageAdmin',
      component: ManageAdminsView
    },
    {
      path: '/:catchAll(.*)',
      name: 'not-found',
      component: NotFoundView
    },
  ]
})

export default router
