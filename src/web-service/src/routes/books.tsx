import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getAllBooks } from "@/repository/books.ts";

const BookCard = ({ book }: { book: Book }) => {
  return (
    <div>
      <Link to={`/books/${book.id}`}>{book.name} </Link>
    </div>
  );
};

const BookList = ({ books }: { books: Book[] }) => {
  return (
    <div>
      {books.map((book) => (
        <BookCard key={book.id} book={book} />
      ))}
    </div>
  );
};

export const Books = () => {
  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["books"],
    queryFn: getAllBooks,
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Error {error.message}</div>;
  }

  if (!isSuccess) {
    return <div>Something went wrong!</div>;
  }

  return (
    <div>
      <div className="flex items-center">
        <BookList books={data} />
      </div>
    </div>
  );
};
