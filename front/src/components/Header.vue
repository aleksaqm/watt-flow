<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/shad/components/ui/navigation-menu"
import { Button } from "@/shad/components/ui/button"
import { BoltIcon } from '@heroicons/vue/16/solid';
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user';


interface MenuItem {
  title: string;
  href: string;
  roles?: string[];
  children?: {
    title: string;
    href: string;
    description?: string;
    roles?: string[];
  }[];
}

const router = useRouter();
const userStore = useUserStore();

const menuItems = ref<MenuItem[]>([
  { title: 'Home', href: '/home', roles: ['Regular', 'Clerk', 'Admin', 'SuperAdmin'] },
  {
    title: 'Users', href: '/users', roles: ['Admin', 'SuperAdmin'], children: [
      { title: 'Manage clerks', href: '/manage/clerks', description: 'Search or add new clerks', roles: ['Admin', 'SuperAdmin'] },
      { title: 'Manage admins', href: '/manage/admins', description: 'Add new admins to the system', roles: ['SuperAdmin'] },
    ],
  },
  {
    title: 'Properties', href: '/properties', roles: ['Admin', 'SuperAdmin', 'Regular'], children: [
      { title: 'Requests', href: '/properties/requests-manage', description: 'Manage new property requests', roles: ['Admin', 'SuperAdmin'] },
      { title: 'Requests', href: '/property-request', description: 'Create new property requests', roles: ['Regular'] },
      { title: 'My requests', href: '/my-property-request', description: 'View my property requests', roles: ['Regular'] },
    ],
  },
  {
    title: 'Households', href: '/household', roles: ['Regular', 'Admin', 'SuperAdmin'], children: [
      { title: 'Search', href: '/household/search', description: 'Search for households', roles: ['Admin', 'SuperAdmin', 'Regular'] },
      { title: 'Ownership Requests', href: '/ownership/requests', description: 'Look at ownership requests', roles: ['Regular', "Admin", "SuperAdmin"] },
    ],
  },
  {
    title: 'Schedule', href: '/schedule', roles: ['Clerk', 'Regular'], children: [
      { title: 'Meetings', href: '/clerk/schedule', description: 'Organize meetings', roles: ['Clerk'] },
      { title: 'New Meeting', href: '/regular/meeting', description: 'Crete new meeting', roles: ['Regular'] },
      { title: 'Meetings', href: '/regular/schedule', description: 'Schedule meetings', roles: ['Regular'] },
    ],
  },
  {
    title: 'Bills', href: '/bills', roles: [], children: [
      { title: 'Price Management', href: '/bills/prices', description: 'Manage active and create new price lists', roles: ['Admin'] },
      { title: 'Send bills', href: '/bills/send', description: 'Send bills to users', roles: ['Admin'] },
    ],
  },
  { title: 'Profile', href: '/profile', roles: ['Regular'] },
  { title: 'Logout', href: '/logout', roles: ['Regular', 'Clerk', 'Admin', 'SuperAdmin'] },
]);

const filteredMenuItems = computed(() => {
  const role = userStore.role;
  if (role !== null) {
    return menuItems.value.filter(item => {
      return !item.roles || item.roles.includes(role);
    }).map(item => ({
      ...item,
      children: item.children?.filter(child => !child.roles || child.roles.includes(role)),
    }));
  }
  return menuItems.value.filter(item => {
    return !item.roles;
  }).map(item => ({
    ...item,
    children: item.children?.filter(child => !child.roles),
  }));
});


const isMenuOpen = ref(false)

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleNavigation = (item: MenuItem) => {
  if (item.title === 'Logout') {
    localStorage.removeItem('authToken')
    const userStore = useUserStore();
    userStore.clearRole();
    userStore.clearId();
    console.log("Logging out...");
    router.push({ name: 'login' });
  } else {
    router.push(item.href);
  }
}
</script>

<template>
  <header class="border-b">
    <div class="container mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <BoltIcon class="w-8 fill-indigo-500"></BoltIcon>
          <a href="/" class="text-xl font-bold">
            Watt-Flow
          </a>
        </div>

        <!-- Desktop Navigation -->
        <div class="hidden md:block">
          <NavigationMenu>
            <NavigationMenuList>
              <NavigationMenuItem v-for="item in filteredMenuItems" :key="item.title">
                <template v-if="item.children">
                  <NavigationMenuTrigger class="nav-item">{{ item.title }}</NavigationMenuTrigger>
                  <NavigationMenuContent>
                    <ul class="grid w-[400px] gap-3 p-4 md:w-[500px] md:grid-cols-2">
                      <li v-for="child in item.children" :key="child.title">
                        <NavigationMenuLink :href="child.href">
                          <div class="text-sm font-medium">{{ child.title }}</div>
                          <p v-if="child.description" class="text-sm text-muted-foreground">
                            {{ child.description }}
                          </p>
                        </NavigationMenuLink>
                      </li>
                    </ul>
                  </NavigationMenuContent>
                </template>
                <NavigationMenuLink v-else href="#" :class="navigationMenuTriggerStyle()"
                  @click.prevent="handleNavigation(item)">
                  {{ item.title }}
                </NavigationMenuLink>
              </NavigationMenuItem>
            </NavigationMenuList>
          </NavigationMenu>
        </div>

        <!-- Mobile Menu Button -->
        <Button variant="ghost" class="md:hidden" @click="toggleMenu">
          <span class="sr-only">Toggle menu</span>
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!isMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </Button>

        <!-- Mobile Menu -->
        <div v-if="isMenuOpen" class="absolute top-16 left-0 right-0 bg-white border-b md:hidden">
          <nav class="px-4 py-2">
            <ul class="space-y-2">
              <li v-for="item in menuItems" :key="item.title">
                <a href="#" @click.prevent="handleNavigation(item)"
                  class="block px-4 py-2 hover:bg-gray-100 rounded-md">
                  {{ item.title }}
                </a>
                <ul v-if="item.children" class="pl-6 space-y-2 mt-2">
                  <li v-for="child in item.children" :key="child.title">
                    <a :href="child.href" class="block px-4 py-2 hover:bg-gray-100 rounded-md">
                      {{ child.title }}
                    </a>
                  </li>
                </ul>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.nav-item {
  font-style: "Inter", sans-serif;
}
</style>
