<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import axios from "axios";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shad/components/ui/table";
import { Button } from "@/shad/components/ui/button";
import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationListItem,
  PaginationNext,
  PaginationPrev,
} from "@/shad/components/ui/pagination";
import Input from "@/shad/components/ui/input/Input.vue";
import router from "@/router";
import type { Clerk } from "./clerk";

import { useToast } from "@/shad/components/ui/toast/use-toast";
import Toaster from "@/shad/components/ui/toast/Toaster.vue";

const props = defineProps({
  query: {
    type: Object,
    default: () => ({}),
  },
  triggerSearch: {
    type: Number,
    default: 0,
  },
});

// Watch for changes in triggerSearch to run fetch
watch(
  () => props.triggerSearch,
  () => {
    pagination.value.page = 1;
    fetchClerks();
  },
);
const clerks = ref<Clerk[]>([]);

const pagination = ref({ page: 1, total: 0, perPage: 10 });

const { toast } = useToast();

const sortBy = ref("username");
const sortOrder = ref<{ [key: string]: "asc" | "desc" | "" }>({
  firstName: "",
  lastName: "",
  username: "",
  status: "",
});
const totalPages = computed(() =>
  Math.ceil(pagination.value.total / pagination.value.perPage),
);

async function fetchClerks() {
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.perPage,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value[sortBy.value],
    };
    console.log(props.query);

    const response = await axios.post("/api/user/query", props.query, {
      params: params,
    });

    if (response.data && response.data.users) {
      clerks.value = response.data.users.map((user: any) => mapToUser(user));

      pagination.value.total = response.data.total;
    }
  } catch (error) {
    console.error("Failed to fetch clerks:", error);
  }
}

function mapToUser(data: any): Clerk {
  return {
    id: data.id,
    firstName: data.first_name,
    lastName: data.last_name,
    username: data.username,
    status: data.status,
  };
}

function onPageChange(page: number) {
  pagination.value.page = page;
  fetchClerks();
}

function onSortChange(field: string) {
  let temp = sortOrder.value[field];
  sortOrder.value.firstName = "";
  sortOrder.value.lastName = "";
  sortOrder.value.username = "";
  sortOrder.value.statsu = "";
  sortOrder.value[field] = temp === "asc" ? "desc" : "asc";
  sortBy.value = field;
  fetchClerks();
}

function getButtonStyle(isSelected: boolean) {
  return isSelected ? ["bg-indigo-500"] : [];
}

function viewClerk(id: number) {
  router.push({ name: "clerk", params: { id: id } });
}

const changeAccountStatus = async (id: number, status: string) => {
  try {
    const action = status === "Active" ? "suspend" : "unsuspend";
    const response = await axios.get("/api/user/" + action + "/" + id);
    if (response.status == 200) {
      if (action === "suspend") {
        toast({
          title: "Account suspended!",
          description:
            "Account is suspended and user will not be able to login!",
          variant: "default",
        });
      } else {
        toast({
          title: "Account unsuspended!",
          description: "Account is unsuspended and user will be able to login!",
          variant: "default",
        });
      }
      fetchClerks();
    }
  } catch (error) {
    console.error("Failed to fetch users:", error);
  }
};
</script>

<template>
  <div class="p-7 flex flex-col bg-white w-10/12 shadow-lg">
    <Table
      class="gap-5 items-center border rounded-2xl border-gray-300 shadow-gray-500 p-10 mb-10"
    >
      <TableHeader>
        <TableRow>
          <TableHead
            @click="onSortChange('firstName')"
            :orientation="sortOrder.firstName"
            >First name</TableHead
          >
          <TableHead
            @click="onSortChange('lastName')"
            :orientation="sortOrder.lastName"
            >Last name</TableHead
          >
          <TableHead
            @click="onSortChange('username')"
            :orientation="sortOrder.username"
            >Username</TableHead
          >
          <TableHead
            @click="onSortChange('status')"
            :orientation="sortOrder.status"
            >Status</TableHead
          >
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow
          v-for="clerk in clerks"
          :key="clerk.id"
          @click="viewClerk(clerk.id)"
        >
          <TableCell>{{ clerk.firstName }}</TableCell>
          <TableCell>{{ clerk.lastName }}</TableCell>
          <TableCell>{{ clerk.username }}</TableCell>
          <TableCell>{{ clerk.status }}</TableCell>
          <TableCell>
            <Button
              class="bg-indigo-500 text-white mr-2 hover:bg-indigo-300"
              @click.stop="changeAccountStatus(clerk.id, clerk.status)"
              >{{ clerk.status === "Active" ? "Suspend" : "Unsuspend" }}</Button
            >
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
    <div class="flex gap-20 pt-10">
      <Pagination
        v-slot="{ page }"
        :total="pagination.total"
        :sibling-count="1"
        show-edges
        :default-page="pagination.page"
        :items-per-page="pagination.perPage"
      >
        <PaginationList v-slot="{ items }" class="flex items-center gap-1">
          <PaginationFirst
            @click="onPageChange(1)"
            :disabled="pagination.page === 1"
          />
          <PaginationPrev
            @click="onPageChange(pagination.page - 1)"
            :disabled="pagination.page === 1"
          />
          <template v-for="(item, index) in items">
            <PaginationListItem
              v-if="item.type === 'page'"
              :key="index"
              :value="item.value"
              as-child
            >
              <Button
                class="w-10 h-10 p-0 hover:bg-indigo-300"
                :class="getButtonStyle(item.value === page)"
                :variant="item.value === page ? 'default' : 'outline'"
                @click="onPageChange(item.value)"
              >
                {{ item.value }}
              </Button>
            </PaginationListItem>
            <PaginationEllipsis v-else :key="item.type" :index="index" />
          </template>

          <PaginationNext
            @click="onPageChange(pagination.page + 1)"
            :disabled="pagination.page === totalPages"
          />
          <PaginationLast
            @click="onPageChange(totalPages)"
            :disabled="pagination.page === totalPages"
          />
        </PaginationList>
      </Pagination>
      <div class="flex items-center gap-2">
        <span>Rows per page:</span>
        <Input
          v-model="pagination.perPage"
          type="number"
          class="w-20"
          min="1"
          placeholder="Rows per page"
        />
      </div>
    </div>
  </div>
  <Toaster />
</template>
