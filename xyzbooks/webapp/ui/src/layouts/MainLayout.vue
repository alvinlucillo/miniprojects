<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title> XYZ Books </q-toolbar-title>

        <div>Quasar v{{ $q.version }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header> Navigation Links </q-item-label>

        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
          :is-active="link.title.toLocaleLowerCase() === activePage"
        />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import EssentialLink, { EssentialLinkProps } from 'components/NavLinks.vue';

const route = useRoute();

const activePage = computed(() => {
  // Remove the leading slash and return the first part of the path
  if (route.path === '/') return 'books';
  return route.path.substring(1).split('/')[0].toLocaleLowerCase();
});

const essentialLinks: EssentialLinkProps[] = [
  {
    title: 'Books',
    caption: 'All books',
    icon: 'book',
    link: '/',
  },
  {
    title: 'Authors',
    caption: 'All authors',
    icon: 'people',
    link: '/authors',
  },
  {
    title: 'Publishers',
    caption: 'All publishers',
    icon: 'factory',
    link: '/publishers',
  },
];

const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>
