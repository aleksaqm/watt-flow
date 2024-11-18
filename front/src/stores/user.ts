import { defineStore } from 'pinia';

interface UserState {
  role: string | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    role: null,
  }),
  actions: {
    setRole(role: string | null) {
      this.role = role;
    },
    clearRole() {
      this.role = null;
    },
  },
});