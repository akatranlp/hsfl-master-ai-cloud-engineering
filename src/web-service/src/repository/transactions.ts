import { AxiosInstance } from "axios";

export class TransactionRepository {
  constructor(private apiClient: AxiosInstance) {}

  async createTransaction(chapterID: number, bookID: number) {
    const response = await this.apiClient.post<void>(`/transactions`, { chapterID, bookID });
    return response.data;
  }

  async getMyReceivedTransactions() {
    const response = await this.apiClient.get<Transaction[]>(`/transactions?receiving=True`);
    return response.data;
  }

  async getMyPaidTransactions() {
    const response = await this.apiClient.get<Transaction[]>(`/transactions`);
    return response.data;
  }

  async getBookFromTransaction(transaction: Transaction) {
    const response = await this.apiClient.get<Book>(`/books/${transaction.bookID}`);
    return response.data;
  }
}
