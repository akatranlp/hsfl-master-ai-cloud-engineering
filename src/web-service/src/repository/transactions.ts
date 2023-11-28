import axios from "axios";

export const createTransaction = async (chapterID: number, bookID: number) => {
  const response = await axios.post<void>(`/api/v1/transactions`, { chapterID, bookID });
  return response.data;
};
export const getMyReceivedTransactions = async () => {
  const response = await axios.get<Transaction[]>(`/api/v1/transactions?receiving=True`);
  return response.data;
};
export const getMyPaidTransactions = async () => {
  const response = await axios.get<Transaction[]>(`/api/v1/transactions`);
  return response.data;
};
export const getBookFromTransaction = async (transaction: Transaction) => {
  const response = await axios.get<Book>(`/api/v1/books/${transaction.bookID}`);
  return response.data;
};
