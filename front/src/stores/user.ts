import { defineStore } from 'pinia';

interface UserState {
  role: string | null;
  id: number | null | undefined;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    role: null,
    id: null,
  }),
  actions: {
    setRole(role: string | null) {
      this.role = role;
    },
    clearRole() {
      this.role = null;
    },
    setId(id: number | null| undefined){
      this.id = id;
    },
    clearId(){
      this.id = null;
    }
  },
});