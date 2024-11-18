import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import HouseholdView from '@/views/household/HouseholdView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import HomeView from '@/views/HomeView.vue'
import HouseholdSearch from '@/views/household/HouseholdSearch.vue'
import PropertyRequestView from '@/views/property/PropertyRequestView.vue'
import UsersPropertyesTable from '@/components/property/UsersPropertyesTable.vue'
import SuperAdminChangePasswordView from '@/views/auth/SuperAdminChangePasswordView.vue'
import ManageAdminsView from '@/views/user/ManageAdminsView.vue'
import UserProfileView from '@/views/user/UserProfileView.vue'
import AdminPropertyRequestsView from '@/views/property/AdminPropertyRequestsView.vue'
import OwnersPropertyRequestsView from '@/views/property/OwnersPropertyRequestsView.vue'

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
      component: OwnersPropertyRequestsView
    },
    {
      path: '/properties/requests-manage',
      name: 'property-request-manage',
      component: AdminPropertyRequestsView
    },
    {
      path: '/superadmin',
      name: 'superadmin',
      component: SuperAdminChangePasswordView
    },
    {
      path: '/manage/admins',
      name: 'manageAdmins',
      component: ManageAdminsView
    },
    {
      path: '/profile',
      name: 'profile',
      component: UserProfileView
    },
    {
      path: '/:catchAll(.*)',
      name: 'not-found',
      component: NotFoundView
    },
  ]
})

export default router
