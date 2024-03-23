<template>
  <q-page>
    <div class="q-pa-md">
      <div class="q-gutter-y-md column q-mb-md">
        <div class="row">
          <div class="col-8 q-mr-xs">
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.title"
              label="Title"
              :hide-bottom-space="true"
              :error="v$.title.$error"
              :error-message="v$.title.$errors[0]?.$message"
              @blur="v$.title.$touch"
              :readonly="!isEditMode"
            />
            <q-select
              class="q-mb-sm"
              filled
              v-model="bookPage.authors"
              use-input
              use-chips
              multiple
              input-debounce="0"
              label="Authors"
              option-value="id"
              :option-label="
                (author) =>
                  `${author.first_name} ${author.middle_name}  ${author.last_name}`
              "
              :options="authors"
              style="width: 250px"
              :error="v$.authors.$error"
              :error-message="v$.authors.$errors[0]?.$message"
              @blur="v$.authors.$touch"
              :hide-bottom-space="true"
              :readonly="!isEditMode"
            >
            </q-select>
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.isbn13"
              label="ISBN 13"
              :hide-bottom-space="true"
              :error="v$.isbn13.$error"
              @blur="v$.isbn13.$touch"
              :error-message="v$.isbn13.$errors[0]?.$message"
              :readonly="!isEditMode"
            />
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.isbn10"
              label="ISBN 10"
              :hide-bottom-space="true"
              :error="v$.isbn10.$error"
              @blur="v$.isbn10.$touch"
              :error-message="v$.isbn10.$errors[0]?.$message"
              :readonly="!isEditMode"
            />
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.publication_year"
              label="Publication Year"
              :hide-bottom-space="true"
              :error="v$.publication_year.$error"
              :error-message="v$.publication_year.$errors[0]?.$message"
              @blur="v$.publication_year.$touch"
              :readonly="!isEditMode"
            />
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.edition"
              label="Edition"
              :hide-bottom-space="true"
              :error="v$.edition.$error"
              :error-message="v$.edition.$errors[0]?.$message"
              @blur="v$.edition.$touch"
              :readonly="!isEditMode"
            />
            <q-select
              class="q-mb-sm"
              filled
              use-input
              v-model="bookPage.publisher"
              input-debounce="0"
              label="Publisher"
              option-value="id"
              option-label="name"
              :options="publishers"
              style="width: 250px"
              :error="v$.publisher.$error"
              :error-message="v$.publisher.$errors[0]?.$message"
              @blur="v$.publisher.$touch"
              :hide-bottom-space="true"
              :readonly="!isEditMode"
            >
              <template v-slot:no-option>
                <q-item>
                  <q-item-section class="text-grey">
                    No results
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
            <q-input
              class="q-mb-sm"
              filled
              v-model="bookPage.price"
              label="Price"
              :hide-bottom-space="true"
              :error="v$.price.$error"
              :error-message="v$.price.$errors[0]?.$message"
              @blur="v$.price.$touch"
              :readonly="!isEditMode"
            />
            <q-input
              filled
              v-model="bookPage.image_url"
              label="Image URL"
              :hide-bottom-space="true"
              :error="v$.image_url.$error"
              :error-message="v$.image_url.$errors[0]?.$message"
              @blur="v$.image_url.$touch"
              :readonly="!isEditMode"
            />
          </div>
          <div class="col-3">
            <q-img
              :src="
                bookPage?.image_url === undefined || bookPage?.image_url === ''
                  ? defaultImage
                  : bookPage?.image_url
              "
              spinner-color="white"
              style="height: 140px; max-width: 150px"
              @error="
                () => {
                  bookPage.image_url = defaultImage;
                }
              "
            />
          </div>
        </div>
      </div>
      <!-- buttons -->
      <q-btn
        v-if="isEditBookMode"
        color="primary"
        label="Save Book"
        @click="submitUpdateBook()"
      />
      <q-btn
        v-if="!isEditBookMode && !isNewBookMode"
        color="primary"
        label="Edit Book"
        @click="editBook()"
      />
      <q-btn
        v-if="isNewBookMode"
        color="primary"
        label="Create Book"
        @click="submitCreateBook()"
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

import { useBookStore, Book } from 'src/stores/books';

const $q = useQuasar();
const route = useRoute();
const router = useRouter();

const { getBookByISBN13, getPublishers, getAuthors, updateBook, createBook } =
  useBookStore();
const { book, publishers, authors } = toRefs(useBookStore());

const defaultImage =
  'https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg';
const isbn13 = ref('');
const bookPage: ref<Book> = ref({
  id: '',
  title: '',
  isbn13: '',
  isbn10: '',
  edition: '',
  authors: [],
  publication_year: 0,
  publisher: null,
  price: 0,
});

// vuelidate rules
const rules = {
  title: {
    required: helpers.withMessage('Title is required', required),
    maximumLength: helpers.withMessage(
      'Title must not exceed 100 characters',
      (value: string) => value.length <= 100
    ),
  },
  isbn13: {
    required: helpers.withMessage('ISBN 13 is required', required),
    length: helpers.withMessage(
      'ISBN 13 must be exactly 13 digits',
      (value: string) => {
        console.log(value, value.length);
        return value.length === 13;
      }
    ),
    numeric: helpers.withMessage(
      'ISBN 13 must be all numbers',
      (value: string) => {
        // check value if all numeric via regex
        const regex = new RegExp('^[0-9]+$');
        return regex.test(value);
      }
    ),
  },
  publication_year: {
    required: helpers.withMessage('Publication Year is required', required),
    numeric: helpers.withMessage(
      'Publication Year must be all numbers',
      (value: string) => {
        // check value if all numeric via regex
        const regex = new RegExp('^[0-9]+$');
        return regex.test(value);
      }
    ),
  },
  price: {
    required: helpers.withMessage('Price is required', required),
    numeric: helpers.withMessage(
      'Price must a valid number',
      (value: string) => {
        // create regex for a valid price
        const regex = new RegExp('^[0-9]+(\.[0-9]{1,2})?$');
        return regex.test(value);
      }
    ),
  },
  authors: {
    required: helpers.withMessage(
      'At least one author is required',
      (value) => Array.isArray(value) && value.length > 0
    ),
  },
  publisher: {
    required: helpers.withMessage('Publisher is required', required),
  },
  isbn10: {
    maxLength: helpers.withMessage(
      'ISBN 10 must be exactly 10 digits',
      (value: string) => {
        if (value) {
          return value.length === 10;
        }
        return true;
      }
    ),
    alphanumeric: helpers.withMessage(
      'ISBN 10 must be alphanumeric',
      (value: string) => {
        if (value) {
          // check value if all alphanumeric via regex
          const regex = new RegExp('^[a-zA-Z0-9]+$');
          return regex.test(value);
        }
        return true;
      }
    ),
  },
  edition: {
    maxLength: helpers.withMessage(
      'Edition must not exceed 100 characters',
      (value: string) => {
        if (value) {
          return value.length <= 100;
        }
        return true;
      }
    ),
  },
  image_url: {
    maxLength: helpers.withMessage(
      'Image URL must not exceed 1000 characters',
      (value: string) => {
        if (value) {
          return value.length <= 1000;
        }
        return true;
      }
    ),
  },
};

const v$ = useVuelidate(rules, bookPage);

// methods
const submitUpdateBook = async () => {
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

  await updateBook(bookPage.value);
};

const submitCreateBook = async () => {
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

  await createBook(bookPage.value);
};

const editBook = () => {
  router.push(`/books/${isbn13.value}/edit`);
};

// computed variables
const isEditBookMode = computed(() => {
  return route.path.endsWith('/edit');
});

const isNewBookMode = computed(() => {
  return route.path.endsWith('/new');
});

const isEditMode = computed(() => {
  return isEditBookMode.value || isNewBookMode.value;
});

// lifecycle hooks
onBeforeMount(async () => {
  if (!isNewBookMode.value) {
    isbn13.value = route.params.isbn13;

    await getBookByISBN13(isbn13.value);
  }

  await getPublishers();
  await getAuthors();

  bookPage.value = book.value;
});
</script>
