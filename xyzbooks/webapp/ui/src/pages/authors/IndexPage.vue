<template>
  <q-page>
    <q-table
      title="Authors"
      :rows="authors"
      :columns="columns"
      :visible-columns="visibleColumns"
      row-key="name"
      :filter="searchFilter"
    >
      <template v-slot:top>
        <q-toolbar>
          <q-toolbar-title> Authors </q-toolbar-title>
        </q-toolbar>
        <q-input
          debounce="300"
          color="primary"
          placeholder="Search"
          v-model="searchFilter"
        >
          <template v-slot:prepend>
            <q-icon name="search" />
          </template>
        </q-input>
        <q-space />

        <q-btn
          color="primary"
          label="Create new author"
          @click="createNewBookAuthor"
        />
      </template>

      <template v-slot:body-cell-first_name="props">
        <q-td :props="props">
          <router-link :to="`/authors/${props.row.id}`">
            {{ props.row.first_name }}
          </router-link>
        </q-td>
      </template>

      <template v-slot:body-cell-edit="props">
        <q-td auto-width :props="props">
          <q-btn
            flat
            round
            color="primary"
            icon="edit"
            @click="editAuthorAction(props.row)"
          />
        </q-td>
      </template>

      <template v-slot:body-cell-delete="props">
        <q-td auto-width :props="props">
          <q-btn
            flat
            round
            color="negative"
            icon="delete"
            @click="deleteAuthorAction(props.row)"
          />
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup lang="ts">
import { onBeforeMount, toRefs, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useBookStore, Book, Author } from 'src/stores/books';

const $q = useQuasar();
const router = useRouter();
const { getAuthors, deleteAuthor } = useBookStore();
const { authors } = toRefs(useBookStore());

const searchFilter = ref('');

onBeforeMount(async () => {
  await reloadAuthors();
});

const reloadAuthors = async () => {
  await getAuthors();
};

const columns = [
  { name: 'id', label: '', field: 'id' },
  {
    name: 'first_name',
    label: 'First Name',
    field: 'first_name',
    align: 'left',
  },
  {
    name: 'middle_name',
    label: 'Middle Name',
    field: 'middle_name',
    align: 'left',
  },
  { name: 'last_name', label: 'Last Name', field: 'last_name', align: 'left' },
  {
    name: 'edit',
    label: 'Edit',
    align: 'left',
    button: true,
  },
  {
    name: 'delete',
    label: 'Delete',
    align: 'left',
    button: true,
  },
];

const visibleColumns = [
  'first_name',
  'middle_name',
  'last_name',
  'edit',
  'delete',
];

const createNewBookAuthor = () => {
  router.push('/authors/new');
};

const editAuthorAction = (author: Author) => {
  router.push(`/authors/${author.id}/edit`);
};

const deleteAuthorAction = (author: Author) => {
  $q.dialog({
    title: 'Delete Author',
    message: `Are you sure you want to delete ${author.first_name} ${author.middle_name} ${author.last_name}?`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    await deleteAuthor(author);
    await reloadAuthors();
  });
};
</script>
