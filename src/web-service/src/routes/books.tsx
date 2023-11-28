import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getAllBooks } from "@/repository/books.ts";
import { Separator } from "@/components/ui/separator";

const BookCard = ({ book }: { book: Book }) => {
  return (
    <>
      <Link to={`/books/${book.id}`}>
        <div className="px-6">
          <div className="text-2xl overflow-hidden whitespace-nowrap overflow-ellipsis">{book.name}</div>
          <div className="overflow-hidden overflow-ellipsis line-clamp-2 break-words">Description: {book.description}</div>
        </div>
      </Link>
      <Separator className="my-4" />
    </>
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
      <div className="text-2xl m-5 ">Recently edited books</div>
      <Separator className="my-4" />
      <div className="items-center pt-2.5">
        <BookList books={data} />
      </div>
    </div>
  );
};
