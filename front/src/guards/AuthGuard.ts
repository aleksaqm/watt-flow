

    import { useUserStore } from '@/stores/user';
    import type { NavigationGuard } from 'vue-router';

    export const authGuard: NavigationGuard = async (to, from, next) => {
        const userStore = useUserStore();
        const authMap = { 
            "Admin": ['/household', '/household/search, /home'], 
            "SuperAdmin": ['/superadmin', '/manage/admins', '/home'], 
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
            next('/');
        }
    };