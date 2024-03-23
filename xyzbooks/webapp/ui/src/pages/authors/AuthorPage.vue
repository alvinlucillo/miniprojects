<template>
  <q-page>
    <div class="q-pa-md">
      <div class="q-gutter-y-md column q-mb-md">
        <q-input
          class="q-mb-sm"
          filled
          v-model="authorPage.first_name"
          label="First Name"
          :hide-bottom-space="true"
          :error="v$.first_name.$error"
          :error-message="v$.first_name.$errors[0]?.$message"
          @blur="v$.first_name.$touch"
          :readonly="!isEditMode"
        />
        <q-input
          v-model="authorPage.middle_name"
          filled
          label="Middle Name"
          :hide-bottom-space="true"
          :error="v$.middle_name.$error"
          :error-message="v$.middle_name.$errors[0]?.$message"
          @blur="v$.middle_name.$touch"
          :readonly="!isEditMode"
        />
        <q-input
          class="q-mb-sm"
          filled
          v-model="authorPage.last_name"
          label="Last Name"
          :hide-bottom-space="true"
          :error="v$.last_name.$error"
          :error-message="v$.last_name.$errors[0]?.$message"
          @blur="v$.last_name.$touch"
          :readonly="!isEditMode"
        />
      </div>
      <!-- buttons -->
      <q-btn
        v-if="isEditAuthorMode"
        color="primary"
        label="Save Book"
        @click="submitUpdateAuthor()"
      />
      <q-btn
        v-if="!isEditAuthorMode && !isNewAuthorMode"
        color="primary"
        label="Edit Author"
        @click="editAuthor()"
      />
      <q-btn
        v-if="isNewAuthorMode"
        color="primary"
        label="Create Author"
        @click="submitCreateAuthor()"
      />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, toRefs, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, helpers } from '@vuelidate/validators';
import { useQuasar } from 'quasar';

import { useBookStore, Author } from 'src/stores/books';

const $q = useQuasar();
const route = useRoute();
const router = useRouter();

const { getAuthorById, createAuthor, updateAuthor } = useBookStore();
const { author } = toRefs(useBookStore());

const authorID = ref('');
const authorPage: ref<Author> = ref({
  first_name: '',
  last_name: '',
  middle_name: '',
});

// vuelidate rules
const rules = {
  first_name: {
    required: helpers.withMessage('First Name is required', required),
    maximumLength: helpers.withMessage(
      'First Name must not exceed 100 characters',
      (value: string) => value.length <= 100
    ),
  },
  last_name: {
    required: helpers.withMessage('Last Name is required', required),
    maximumLength: helpers.withMessage(
      'Last Name must not exceed 100 characters',
      (value: string) => value.length <= 100
    ),
  },
  middle_name: {
    maximumLength: helpers.withMessage(
      'Middle Name must not exceed 100 characters',
      (value: string) => value.length <= 100
    ),
  },
};

const v$ = useVuelidate(rules, authorPage);

// methods
const submitUpdateAuthor = async () => {
  v$.value.$touch();
  if (v$.value.$error) {
    console.log(v$.value.$errors);
    $q.notify({
      color: 'negative',
      message: 'Please fix the errors',
      icon: 'report_problem',
      actions: [
        {
          icon: 'close',
          color: 'white',
          round: true,
          handler: () => {
            /* ... */
          },
        },
      ],
    });
    return;
  }

  await updateAuthor(authorPage.value);
};

const submitCreateAuthor = async () => {
  v$.value.$touch();
  if (v$.value.$error) {
    console.log(v$.value.$errors);
    $q.notify({
      color: 'negative',
      message: 'Please fix the errors',
      icon: 'report_problem',
      actions: [
        {
          icon: 'close',
          color: 'white',
          round: true,
          handler: () => {
            /* ... */
          },
        },
      ],
    });
    return;
  }

  await createAuthor(authorPage.value);
};

const editAuthor = () => {
  router.push(`/authors/${authorID.value}/edit`);
};

// computed variables
const isEditAuthorMode = computed(() => {
  return route.path.endsWith('/edit');
});

const isNewAuthorMode = computed(() => {
  return route.path.endsWith('/new');
});

const isEditMode = computed(() => {
  return isEditAuthorMode.value || isNewAuthorMode.value;
});

// lifecycle hooks
onBeforeMount(async () => {
  if (!isNewAuthorMode.value) {
    authorID.value = route.params.id;
    await getAuthorById(authorID.value);

    authorPage.value = author.value;
  }
});
</script>
