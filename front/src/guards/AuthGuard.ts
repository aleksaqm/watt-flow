

import { useUserStore } from '@/stores/user';
import type { NavigationGuard } from 'vue-router';

export const authGuard: NavigationGuard = async (to, from, next) => {
  const userStore = useUserStore();
  const authMap = {
    "Admin": ['/household/info/:id', '/household/search', '/home', '/properties/requests-manage', '/manage/clerks', '/manage/clerks/new'],
    "SuperAdmin": ['/manage/admins', '/household/info/:id', '/household/search', '/home', '/properties/requests-manage', '/manage/clerks', '/manage/clerks/new'],
    "Clerk": ["/home"],
    "Regular": ['/profile', '/home', '/my-property-request', '/property-request']
  };
  const role = userStore.role;
  let allowedPaths: string[];
  if (role != null && role in authMap) {
    allowedPaths = authMap[role as keyof typeof authMap];
  } else {
    allowedPaths = [];
  }
  console.log(to)
  if (allowedPaths.includes(to.matched[0].path)) {
    next();
  } else {
    // userStore.clearRole()
    // localStorage.removeItem('authToken')
    next('/home');
  }
};
