import { useUserStore } from '@/stores/user';
import type { NavigationGuard } from 'vue-router';

export const authGuard: NavigationGuard = async (to, from, next) => {
  if (localStorage.getItem("authToken") === null){
    next('/login');
    return;
  }
  const userStore = useUserStore();
  const authMap = {
    "Admin": ['/household/info/:id', '/household/search', '/home', '/properties/requests-manage', '/manage/clerks', '/manage/clerks/new', '/ownership/requests', '/bills/prices', '/bills/send', '/consumption', '/clerk/profile/:id', '/bills/owner',  '/bills/pay/:billId', '/bills/pay/success'],
    "SuperAdmin": ['/manage/admins', '/household/info/:id', '/household/search', '/home', '/properties/requests-manage', '/manage/clerks', '/manage/clerks/new', '/ownership/requests', '/bills/prices', '/bills/send', '/consumption', '/clerk/profile/:id', '/bills/owner', '/bills/pay/:billId', '/bills/pay/success'],
    "Clerk": ["/home", "/meeting/clerk"],
    "Regular": ['/profile', '/home', '/my-property-request', '/property-request', '/household/search', '/ownership/requests', '/meeting/user', '/meeting/schedule', '/household/my', '/household/my/:id', '/bills/owner', '/bills/pay/:billId', '/bills/pay/success']
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


  // TODO must add invalid token check and redirect to login page
};
