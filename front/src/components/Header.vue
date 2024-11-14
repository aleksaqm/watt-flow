<script setup lang="ts">
import { ref } from 'vue'
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


interface MenuItem {
  title: string;
  href: string;
  children?: {
    title: string;
    href: string;
    description?: string;
  }[];
}

const router = useRouter();

const menuItems = ref<MenuItem[]>([
  {
    title: 'Home',
    href: '/',
  },
  {
    title: 'Users',
    href: '/products',
    children: [
      {
        title: 'Manage users',
        href: '/products/feature-2',
        description: 'Add new or manage existing users',
      },
      {
        title: 'Manage admins',
        href: '/manage/admins',
        description: 'Add new admins to the system',
      },
    ],
  },
  {
    title: 'Properties',
    href: '/properties',
    children: [
      {
        title: 'Requests',
        href: '/properties/requests-manage',
        description: 'Manage new property requests',
      },
    ],
  },
  {
    title: 'Households',
    href: '/household',
    children: [
      {
        title: 'Search',
        href: '/household/search',
        description: 'Search for households',
      },
      {
        title: 'Requests',
        href: '/households/requests-manage',
        description: 'Manage household ownership requests',
      },
    ],
  },
  {
    title: 'Bills',
    href: '/bills',
    children: [
      {
        title: 'Price Management',
        href: '/bills/prices',
        description: 'Manage active and create new price lists',
      },
      {
        title: 'Send bills',
        href: '/bills/send',
        description: 'Send bills to users',
      },
    ],
  },
  {
    title: 'Profile',
    href: '/profile',
  },
  {
    title: 'Logout',
    href: '/logout',
  },
])


const isMenuOpen = ref(false)

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleNavigation = (item: MenuItem) => {
  if (item.title === 'Logout') {
    localStorage.removeItem('authToken')
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
              <NavigationMenuItem v-for="item in menuItems" :key="item.title">
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
