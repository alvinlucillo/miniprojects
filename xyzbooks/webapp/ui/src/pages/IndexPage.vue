<template>
  <q-page>
    <q-table
      title="Books"
      :rows="books"
      :columns="columns"
      :visible-columns="visibleColumns"
      row-key="name"
      :filter="searchFilter"
    >
      <template v-slot:top>
        <q-toolbar>
          <q-toolbar-title> Books </q-toolbar-title>
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
          label="Create new book"
          @click="createNewBookAction"
        />
      </template>

      <template v-slot:body-cell-title="props">
        <q-td :props="props">
          <router-link :to="`/books/${props.row.isbn13}`">
            {{ props.row.title }}
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
            @click="editBookAction(props.row)"
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
            @click="deleteBookAction(props.row)"
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
import { useBookStore, Book } from 'src/stores/books';

const $q = useQuasar();
const router = useRouter();
const { getBooks, deleteBook } = useBookStore();
const { books } = toRefs(useBookStore());

const searchFilter = ref('');

onBeforeMount(async () => {
  await reloadBooks();
});

const reloadBooks = async () => {
  await getBooks();
};

const columns = [
  { name: 'id', label: '', field: 'id' },
  { name: 'title', label: 'Title', field: 'title', align: 'left' },
  {
    name: 'author',
    label: 'Author(s)',
    field: (row) => {
      return row?.authors
        ?.map(
          (author) =>
            `${author.first_name} ${author.middle_name} ${author.last_name}`
        )
        .join(', ');
    },
    align: 'left',
  },
  { name: 'isbn13', label: 'ISBN 13', field: 'isbn13', align: 'left' },
  { name: 'isbn10', label: 'ISBN 10', field: 'isbn10', align: 'left' },
  {
    name: 'publication_year',
    label: 'Publication Year',
    field: 'publication_year',
    align: 'left',
  },
  {
    name: 'publisher',
    label: 'Publisher',
    field: (row) => row?.publisher?.name ?? '',
    align: 'left',
  },
  { name: 'edition', label: 'Edition', field: 'edition', align: 'left' },
  { name: 'price', label: 'Price', field: 'price', align: 'left' },
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
  'title',
  'author',
  'isbn13',
  'isbn10',
  'price',
  'publisher',
  'publication_year',
  'edition',
  'price',
  'edit',
  'delete',
];

const createNewBookAction = () => {
  router.push('/books/new');
};

const editBookAction = (book: Book) => {
  router.push(`/books/${book.isbn13}/edit`);
};

const deleteBookAction = (book: Book) => {
  $q.dialog({
    title: 'Delete Book',
    message: `Are you sure you want to delete ${book.title}?`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    await deleteBook(book);
    await reloadBooks();
  });
};
</script>
