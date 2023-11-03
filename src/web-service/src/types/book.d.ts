type Book = {
  id: number;
  name: string;
  description: string;
  authorid: number;
};

type UpdateBook = {
  name?: string;
  description?: string;
};
