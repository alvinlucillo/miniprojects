<template>
  <q-page>
    <div class="q-pa-md">
      <div class="q-gutter-y-md column q-mb-md">
        <q-input
          class="q-mb-sm"
          filled
          v-model="publisherPage.name"
          label="Name"
          :hide-bottom-space="true"
          :error="v$.name.$error"
          :error-message="v$.name.$errors[0]?.$message"
          @blur="v$.name.$touch"
          :readonly="!isEditMode"
        />
      </div>
      <!-- buttons -->
      <q-btn
        v-if="isEditPublisherMode"
        color="primary"
        label="Save Publisher"
        @click="submitUpdatePublisher()"
      />
      <q-btn
        v-if="!isEditPublisherMode && !isNewPublisherMode"
        color="primary"
        label="Edit Publisher"
        @click="editPublisher()"
      />
      <q-btn
        v-if="isNewPublisherMode"
        color="primary"
        label="Create Publisher"
        @click="submitCreatePublisher()"
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

import { useBookStore, Publisher } from 'src/stores/books';

const $q = useQuasar();
const route = useRoute();
const router = useRouter();

const { getPublisherById, createPublisher, updatePublisher } = useBookStore();
const { publisher } = toRefs(useBookStore());

const publisherID = ref('');
const publisherPage: ref<Publisher> = ref({
  name: '',
});

// vuelidate rules
const rules = {
  name: {
    required: helpers.withMessage('Name is required', required),
    maximumLength: helpers.withMessage(
      'Name must not exceed 100 characters',
      (value: string) => value.length <= 100
    ),
  },
};

const v$ = useVuelidate(rules, publisherPage);

// methods
const submitUpdatePublisher = async () => {
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

  await updatePublisher(publisherPage.value);
};

const submitCreatePublisher = async () => {
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

  await createPublisher(publisherPage.value);
};

const editPublisher = () => {
  router.push(`/publishers/${publisherID.value}/edit`);
};

// computed variables
const isEditPublisherMode = computed(() => {
  return route.path.endsWith('/edit');
});

const isNewPublisherMode = computed(() => {
  return route.path.endsWith('/new');
});

const isEditMode = computed(() => {
  return isEditPublisherMode.value || isNewPublisherMode.value;
});

// lifecycle hooks
onBeforeMount(async () => {
  if (!isNewPublisherMode.value) {
    publisherID.value = route.params.id;
    await getPublisherById(publisherID.value);

    publisherPage.value = publisher.value;
  }
});
</script>
