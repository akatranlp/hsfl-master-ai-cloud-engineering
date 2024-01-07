import { AxiosInstance } from "axios";

export class BookRepository {
  constructor(private apiClient: AxiosInstance) {}

  async getAllBooks() {
    const response = await this.apiClient.get<Book[]>("/books");
    return response.data;
  }

  async getBookById(bookId: number) {
    const response = await this.apiClient.get<Book>(`/books/${bookId}`);
    return response.data;
  }

  async getChaptersByBookId(bookId: number) {
    const response = await this.apiClient.get<ChapterPreview[]>(`/books/${bookId}/chapters`);
    return response.data;
  }

  async getChapter(bookId: number, chapterId: number) {
    const response = await this.apiClient.get<Chapter>(`/books/${bookId}/chapters/${chapterId}`);
    return response.data;
  }

  async createBook(book: CreateBook) {
    const response = await this.apiClient.post<void>("/books", book);
    return response.data;
  }
  async editBook(book: UpdateBook, bookId: number) {
    const response = await this.apiClient.patch<void>(`/books/${bookId}`, book);
    return response.data;
  }

  async createChapter(chapter: CreateChapter) {
    const response = await this.apiClient.post<void>(`/books/${chapter.bookid}/chapters`, chapter);
    return response.data;
  }
  async editChapter(chapter: UpdateChapter, bookId: number, chapterId: number) {
    const response = await this.apiClient.patch<void>(`/books/${bookId}/chapters/${chapterId}`, chapter);
    return response.data;
  }

  async getMyBooks(userId: number) {
    const response = await this.apiClient.get<Book[]>(`/books?userId=${userId}`);
    return response.data;
  }

  async getBoughtBooks() {
    const response = await this.apiClient.get<Book[]>("/books");
    const response2 = await this.apiClient.get<Transaction[]>(`/transactions`);
    const boughtBookIds = response2.data.map((transaction) => transaction.bookID);
    return response.data.filter((book) => boughtBookIds.includes(book.id));
  }
}
