import { defineStore } from 'pinia';
import { Ref, ref } from 'vue';
import { api } from 'boot/axios';
import { useQuasar } from 'quasar'
import { au } from 'app/dist/spa/assets/index.dbb50748';

export interface Book {
  id: string;
  title: string;
  author: string;
  isbn13: string;
  isbn10: string;
  publication_year: number;
  edition: string;
  price: number;
  image_url?: string;
  publisher: Publisher
  authors: Author[];
}

export interface Author {
  id: string;
  first_name: string;
  middle_name?: string;
  last_name: string;
}

export interface Publisher {
  id: string;
  name: string;
}

export const useBookStore = defineStore('book', () => {
  const books: Ref<Book[]> = ref([]);
  const book: Ref<Book> = ref({} as Book);

  const author: Ref<Author> = ref({} as Author);
  const authors: Ref<string[]> = ref([]);

  const publishers: Ref<string[]> = ref([]);
  const publisher: Ref<Publisher> = ref({} as Publisher);

  const $q = useQuasar()

  const getBooks = async () => {
    $q.loading.show();
    try {
      const response = await api.get('/books');
      books.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const getBookByISBN13 = async (isbn13: string) => {
    $q.loading.show();
    try {
      const response = await api.get(`/books/${isbn13}`);
      book.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const updateBook = async (b: Book) => {
    $q.loading.show();

    // reinforce price and year to be a number
    b.price = Number(b.price);
    b.publication_year = Number(b.publication_year);

    try {
      const response = await api.put(`/books/${b.id}`, b, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Book updated successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const deleteBook = async (b: Book) => {
    $q.loading.show();
    try {
      const response = await api.delete(`/books/${b.id}`);
      showSuccess('Book deleted successfully');

    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const createBook = async (b: Book) => {
    $q.loading.show();

    // reinforce price and year to be a number
    b.price = Number(b.price);
    b.publication_year = Number(b.publication_year);

    try {
      const response = await api.post(`/books`, b, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Book created successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const getAuthors = async () => {
    $q.loading.show();
    try {
      const response = await api.get('/authors');
      authors.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const getAuthorById = async (id: string) => {
    $q.loading.show();
    try {
      const response = await api.get(`/authors/${id}`);
      author.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const deleteAuthor = async (a: Author) => {
    $q.loading.show();
    try {
      const response = await api.delete(`/authors/${a.id}`);
      showSuccess('Author deleted successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const updateAuthor = async (a: Author) => {
    $q.loading.show();
    try {
      const response = await api.put(`/authors/${a.id}`, a, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Author updated successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const createAuthor = async (a: Author) => {
    $q.loading.show();
    try {
      const response = await api.post(`/authors`, a, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Author created successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const getPublishers = async () => {
    $q.loading.show();
    try {
      const response = await api.get('/publishers');
      publishers.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const getPublisherById = async (id: string) => {
    $q.loading.show();
    try {
      const response = await api.get(`/publishers/${id}`);
      publisher.value = response.data;
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const deletePublisher = async (p: Publisher) => {
    $q.loading.show();
    try {
      const response = await api.delete(`/publishers/${p.id}`);
      showSuccess('Publisher deleted successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const updatePublisher = async (p: Publisher) => {
    $q.loading.show();
    try {
      const response = await api.put(`/publishers/${p.id}`, p, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Publisher updated successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const createPublisher = async (p: Publisher) => {
    $q.loading.show();
    try {
      const response = await api.post(`/publishers`, p, { headers: { 'Content-Type': 'application/json' } });
      showSuccess('Publisher created successfully');
    } catch (error) {
      console.log(error);
      showError(error);
    }
    $q.loading.hide();
  }

  const showSuccess = (msg: string) => {
    $q.notify({
      color: 'positive',
      message: msg,
      actions: [
        { icon: 'close', color: 'white', round: true, handler: () => { /* ... */ } }]
    });
  }

  const showError = (error: any) => {
    let msgs: [string] = ['Server error occurred'];
    if (error?.response?.status == 404) {
      msgs[0] = 'Record not found';
    } else if (error?.response?.status == 400) {
      if (error?.response?.data?.errors) {
        msgs = error?.response?.data?.errors.map((e: any) => e.message)
      } else {
        msgs[0] = error?.response?.data ?? 'Error processing request';
      }
    }

    for (const msg of msgs) {
      $q.notify({
        color: 'negative',
        message: msg,
        actions: [
          { icon: 'close', color: 'white', round: true, handler: () => { /* ... */ } }],
        group: false,
        timeout: 0
      });
    }
  }

  return {
    book,
    books,
    author,
    authors,
    publisher,
    publishers,
    getBooks,
    getBookByISBN13,
    createBook,
    updateBook,
    deleteBook,
    getAuthors,
    getAuthorById,
    createAuthor,
    updateAuthor,
    deleteAuthor,
    getPublishers,
    getPublisherById,
    createPublisher,
    updatePublisher,
    deletePublisher
  }

});