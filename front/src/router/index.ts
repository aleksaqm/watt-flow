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
import { useUserStore } from '@/stores/user'
import { authGuard } from '@/guards/AuthGuard'

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
      component: RegisterView,
    },
    {
      path: '/home',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/household/info/:id',
      name: 'household',
      component: HouseholdView,
      meta: { requiresAuth: true }
    },
    {
      path: '/household/search',
      name: 'household-search',
      component: HouseholdSearch,
      meta: { requiresAuth: true }
    },
    {
      path: '/property-request',
      name: 'property-request',
      component: PropertyRequestView,
      meta: { requiresAuth: true }
    },
    {
      path: '/my-property-request',
      name: 'my-property-request',
      component: OwnersPropertyRequestsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/properties/requests-manage',
      name: 'property-request-manage',
      component: AdminPropertyRequestsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/superadmin',
      name: 'superadmin',
      component: SuperAdminChangePasswordView,
    },
    {
      path: '/manage/admins',
      name: 'manageAdmins',
      component: ManageAdminsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: UserProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/:catchAll(.*)',
      name: 'not-found',
      component: NotFoundView
    },
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    authGuard(to, from, next);
  } else {
    next();
  }
});

export default router
