type Book = {
  id: number;
  name: string;
  description: string;
  authorId: number;
};

type CreateBook = {
  name: string;
  description: string;
};

type UpdateBook = {
  name?: string;
  description?: string;
};
