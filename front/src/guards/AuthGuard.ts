

    import { useUserStore } from '@/stores/user';
    import type { NavigationGuard } from 'vue-router';

    export const authGuard: NavigationGuard = async (to, from, next) => {
        const userStore = useUserStore();
        const authMap = { 
            "Admin": ['/household', '/household/search', '/home', '/properties/requests-manage'], 
            "SuperAdmin": ['/superadmin', '/manage/admins', '/household/search', '/home', '/properties/requests-manage'], 
            "Clerk": ["/home"],
            "Regular": ['/profile', '/home']
        };
        const role = userStore.role;
        let allowedPaths: string[];
        if (role != null && role in authMap) {
            allowedPaths = authMap[role as keyof typeof authMap];
        } else {
            allowedPaths = [];
        }
        if (allowedPaths.includes(to.path)) {
            next();
        } else {
            userStore.clearRole()
            localStorage.removeItem('authToken')
            next('/');
        }
    };