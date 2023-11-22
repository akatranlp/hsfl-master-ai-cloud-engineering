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

export const createBook = async (book: CreateBook) => {
  const response = await axios.post<void>("/api/v1/books", book);
  return response.data;
};
export const editBook = async (book: UpdateBook, bookId: number) => {
  const response = await axios.patch<void>(`/api/v1/books/${bookId}`, book);
  return response.data;
};

export const createChapter = async (chapter: CreateChapter) => {
  const response = await axios.post<void>(`/api/v1/books/${chapter.bookid}/chapters`, chapter);
  return response.data;
};
export const editChapter = async (chapter: UpdateChapter, bookId: number, chapterId: number) => {
  const response = await axios.patch<void>(`/api/v1/books/${bookId}/chapters/${chapterId}`, chapter);
  return response.data;
};

export const getMyBooks = async (userId: number) => {
  const response = await axios.get<Book[]>(`/api/v1/books?userId=${userId}`);
  return response.data;
};

export const getBoughtBooks = async () => {
  const response = await axios.get<Book[]>("/api/v1/books");
  const response2 = await axios.get<Transaction[]>(`/api/v1/transactions`);
  const boughtBookIds = response2.data.map((transaction) => transaction.bookID);
  return response.data.filter((book) => boughtBookIds.includes(book.id));
};
