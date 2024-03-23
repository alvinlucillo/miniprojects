<template>
  <q-page>
    <q-table
      title="Publishers"
      :rows="publishers"
      :columns="columns"
      :visible-columns="visibleColumns"
      row-key="name"
      :filter="searchFilter"
    >
      <template v-slot:top>
        <q-toolbar>
          <q-toolbar-title> Publishers </q-toolbar-title>
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
          label="Create new publisher"
          @click="createNewPublisher"
        />
      </template>

      <template v-slot:body-cell-name="props">
        <q-td :props="props">
          <router-link :to="`/publishers/${props.row.id}`">
            {{ props.row.name }}
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
            @click="editPublisherAction(props.row)"
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
            @click="deletePublisherAction(props.row)"
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
import { useBookStore, Publisher } from 'src/stores/books';

const $q = useQuasar();
const router = useRouter();
const { getPublishers, deletePublisher } = useBookStore();
const { publishers } = toRefs(useBookStore());

const searchFilter = ref('');

onBeforeMount(async () => {
  await reloadPublishers();
});

const reloadPublishers = async () => {
  await getPublishers();
};

const columns = [
  { name: 'id', label: '', field: 'id' },
  { name: 'name', label: 'Publisher Name', field: 'name', align: 'left' },
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

const visibleColumns = ['name', 'edit', 'delete'];

const createNewPublisher = () => {
  router.push('/publishers/new');
};

const editPublisherAction = (publisher: Publisher) => {
  router.push(`/publishers/${publisher.id}/edit`);
};

const deletePublisherAction = (publisher: Publisher) => {
  $q.dialog({
    title: 'Delete Publisher',
    message: `Are you sure you want to delete ${publisher.name}?`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    await deletePublisher(publisher);
    await reloadPublishers();
  });
};
</script>
