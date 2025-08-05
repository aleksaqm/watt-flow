<script setup lang="ts">
import DatePicker from '@/components/schedule/DatePicker.vue';
import Meeting from '@/components/schedule/Meeting.vue';
import axios from 'axios';
import { useForm } from 'vee-validate';
import { ref, watch } from 'vue';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shad/components/ui/form'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shad/components/ui/popover'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shad/components/ui/command'
import { Button } from '@/shad/components/ui/button'

const formattedTime = (hour: number, minute: number) =>
  `${String(hour).padStart(2, "0")}:${String(minute).padStart(2, "0")}`;

const { handleSubmit, setFieldValue, values, errors } = useForm({
  initialValues: {
    username: '',
    user_id: 0,
  },
});

const onSubmit = handleSubmit((values) => {
  console.log(values)
})

interface User {
  id: number;
  username: string;
}

const users = ref<User[]>([]);
const searchInput = ref("");
const updateSearchTerm = (search: any) => {
  searchInput.value = search;
};

watch(searchInput, (newValue) => {
  fetchUsers(newValue);
});

const fetchUsers = async (searchTerm: string) => {
  const searchQuery = { Role: "Clerk", Username: searchTerm, Status: "Active" };
  const params = { sortBy: "username", page: 1, pageSize: 20  };
  try {
    const result = await axios.post("/api/user/query", searchQuery, { params });
    users.value = result.data.users;
  } catch (error) {
    console.log(error);
  }
};
</script>

<template>
  <main class="flex flex-col items-center min-h-screen pt-10">
    <div class="flex flex-col items-center space-y-10">
      <span class="text-xl">Create new meeting</span>
        <div class="flex items-center space-x-4 w-full">
          <FormField name="username" v-slot="{ field }">
            <FormItem class="flex items-center space-x-4 w-full">
              <FormLabel class="w-1/3 text-center">Meeting with</FormLabel>
              <Popover>
                <PopoverTrigger as-child>
                  <FormControl>
                    <Button
                      variant="outline"
                      role="combobox"
                      class="w-full justify-between !text-left"
                      :class="[!values.username && 'text-muted-foreground']"
                    >
                      {{
                        values.username
                          ? users.find((user) => user.username === values.username)?.username
                          : 'Type to search'
                      }}
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent class="w-[200px] p-0">
                  <Command v-on:update:search-term="updateSearchTerm">
                    <CommandInput placeholder="Search users..." />
                    <CommandEmpty>Nothing found.</CommandEmpty>
                    <CommandList>
                      <CommandGroup>
                        <CommandItem
                          v-for="user in users"
                          :key="user.username"
                          :value="user.username"
                          @select="() => {
                            setFieldValue('username', user.username);
                            setFieldValue('user_id', user.id);
                          }"
                        >
                          {{ user.username }}
                        </CommandItem>
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
              <FormMessage class="text-xs" v-if="errors.username">
                {{ errors.username }}
              </FormMessage>
            </FormItem>
          </FormField>
        </div>

      <DatePicker :user-id="values.user_id" :username="values.username" class="w-full max-w-lg mt-10"></DatePicker>
    </div>
  </main>
</template>
