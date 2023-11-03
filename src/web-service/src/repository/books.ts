import axios from "axios";

export const getAllBooks = async () => {
  const response = await axios.get<Book[]>("/api/v1/books");
  return response.data;
};

export const getBookById = async (bookId: number) => {
  const response = await axios.get<Book>(`/api/v1/books/${bookId}`);
  return response.data;
};

export const getChaptersByBookId = async (bookId: number) => {
  const response = await axios.get<Chapter[]>(`/api/v1/books/${bookId}/chapters`);
  return response.data;
};

export const getChapter = async (bookId: number, chapterId: number) => {
  const response = await axios.get<Chapter>(`/api/v1/books/${bookId}/chapters/${chapterId}`);
  return response.data;
};

export const getAllTransactions = async () => {
  const response = await axios.get<Transaction[]>(`/api/v1/transactions`);
  return response.data;
};

export const createBook = async (book: Book) => {
  const response = await axios.post<void>("/api/v1/books", book);
  return response.data;
};

export const createChapter = async (bookId: number, chapter: Chapter) => {
  const response = await axios.post<void>(`/api/v1/books/${bookId}`, chapter);
  return response.data;
};

export const createTransaction = async (chapterId: number) => {
  const response = await axios.post<void>(`/api/v1/transactions`, { chapterId });
  return response.data;
};
