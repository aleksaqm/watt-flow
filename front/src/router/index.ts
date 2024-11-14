import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import HouseholdView from '@/views/HouseholdView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import HomeView from '@/views/HomeView.vue'
import HouseholdSearch from '@/views/HouseholdSearch.vue'
import PropertyRequestView from '@/views/PropertyRequestView.vue'
import UsersPropertyesTable from '@/components/property/UsersPropertyesTable.vue'

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
      path: '/household/search',
      name: 'household-search',
      component: HouseholdSearch
    },
    {
      path: '/property-request',
      name: 'property-request',
      component: PropertyRequestView
    },
    {
      path: '/my-property-request',
      name: 'my-property-request',
      component: UsersPropertyesTable
    },
    {
      path: '/:catchAll(.*)',
      name: 'not-found',
      component: NotFoundView
    },
  ]
})

export default router
