import { useQuery } from "@tanstack/react-query";
import { getBoughtBooks } from "@/repository/books.ts";
import { Link } from "react-router-dom";
import { Separator } from "@/components/ui/separator.tsx";

const BookCard = ({ book }: { book: Book }) => {
  return (
    <>
      <Link to={`/books/${book.id}`}>
        <div className="px-6">
          <div className="text-2xl">{book.name}</div>
          <div>Description: {book.description}</div>
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

/*TODO: Get Bought Books*/
export const BoughtBooks = () => {
  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["myBooks"],
    queryFn: () => getBoughtBooks(),
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
      <div className="items-center pt-2.5">
        <BookList books={data} />
      </div>
    </div>
  );
};
