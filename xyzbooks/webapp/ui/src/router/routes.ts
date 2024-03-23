import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { path: '/books', component: () => import('pages/IndexPage.vue') },
      { path: '/books/:isbn13', component: () => import('pages/books/BookPage.vue') },
      { path: '/books/:isbn13/edit', component: () => import('pages/books/BookPage.vue') },
      { path: '/books/new', component: () => import('pages/books/BookPage.vue') },
    ],
  },
  {
    path: '/authors',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/authors/IndexPage.vue') },
      { path: '/authors/new', component: () => import('pages/authors/AuthorPage.vue') },
      { path: '/authors/:id', component: () => import('pages/authors/AuthorPage.vue') },
      { path: '/authors/:id/edit', component: () => import('pages/authors/AuthorPage.vue') },
    ],
  },
  {
    path: '/publishers',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/publishers/IndexPage.vue') },
      { path: '/publishers/new', component: () => import('pages/publishers/PublisherPage.vue') },
      { path: '/publishers/:id', component: () => import('pages/publishers/PublisherPage.vue') },
      { path: '/publishers/:id/edit', component: () => import('pages/publishers/PublisherPage.vue') },
    ],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
