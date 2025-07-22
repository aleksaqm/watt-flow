import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import HouseholdView from '@/views/household/HouseholdView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import HomeView from '@/views/HomeView.vue'
import HouseholdSearch from '@/views/household/HouseholdSearch.vue'
import PropertyRequestView from '@/views/property/PropertyRequestView.vue'
import SuperAdminChangePasswordView from '@/views/auth/SuperAdminChangePasswordView.vue'
import ManageAdminsView from '@/views/user/ManageAdminsView.vue'
import UserProfileView from '@/views/user/UserProfileView.vue'
import AdminPropertyRequestsView from '@/views/property/AdminPropertyRequestsView.vue'
import OwnersPropertyRequestsView from '@/views/property/OwnersPropertyRequestsView.vue'
import ManageClerksView from '@/views/user/ManageClerksView.vue'
import { authGuard } from '@/guards/AuthGuard'
import RegisterClerkView from '@/views/user/RegisterClerkView.vue'
import ClerkMeetingScheduleView from '@/views/schedule/ClerkMeetingScheduleView.vue'
import OwnershipRequestsView from '@/views/household/OwnershipRequestsView.vue'
import RegularCreateMeetingView from '@/views/schedule/RegularCreateMeetingView.vue'
import RegularMeetingScheduleView from '@/views/schedule/RegularMeetingScheduleView.vue'
import PricelistManagement from '@/views/bills/PricelistManagement.vue'
import IssueBills from '@/views/bills/IssueBills.vue'
import CityConsumption from '@/views/household/CityConsumption.vue'
import ClerkProfileView from '@/views/user/ClerkProfileView.vue'
import OwnerBills from '@/views/bills/OwnerBills.vue'
import PaymentPage from '@/views/bills/PaymentPage.vue'
import PaymentSuccessPage from '@/views/bills/PaymentSuccessPage.vue'

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
      path: '/manage/clerks',
      name: 'manageClerks',
      component: ManageClerksView,
      meta: { requiresAuth: true }
    },
    {
      path: '/manage/clerks/new',
      name: 'new-clerk',
      component: RegisterClerkView,
      meta: { requiresAuth: true }
    },
    {
      path: '/meeting/clerk',
      name: 'clerk-schedule',
      component: ClerkMeetingScheduleView,
      meta: { requiresAuth: true }
    },
    {
      path: '/meeting/schedule',
      name: 'regular-schedule',
      component: RegularMeetingScheduleView,
      meta: { requiresAuth: true }
    },
    {
      path: '/meeting/user',
      name: 'regular-meeting',
      component: RegularCreateMeetingView,
      meta: { requiresAuth: true }
    },
    {
      path: '/bills/prices',
      name: 'manage-prices',
      component: PricelistManagement,
      meta: { requiresAuth: true }
    },
    {
      path: '/bills/send',
      name: 'send-bills',
      component: IssueBills,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: UserProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/clerk/profile/:id',
      name: 'clerk-profile',
      component: ClerkProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ownership/requests',
      name: 'ownership-requests',
      component: OwnershipRequestsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/consumption',
      name: 'city-consumption',
      component: CityConsumption,
      meta : { requiresAuth: true }
    },
    {
      path: '/bills/owner',
      name: 'owner-bills',
      component: OwnerBills,
      meta : { requiresAuth: true}
    },
    {
      path: '/bills/pay/:billId',
      name: 'pay-bill',
      component: PaymentPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/bills/pay/success',
      name: 'pay-success',
      component: PaymentSuccessPage,
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
