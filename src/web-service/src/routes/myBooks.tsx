import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Separator } from "@/components/ui/separator.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useUserData } from "@/provider/user-provider.tsx";
import { useRepository } from "@/provider/repository-provider";

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

export const MyBooks = () => {
  const { bookRepo } = useRepository();
  const user = useUserData();

  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["myBooks"],
    queryFn: () => bookRepo.getMyBooks(user.id),
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
      <div className={"flex justify-end px-6 pt-2.5"}>
        <div className="m-4">
          <Link to="/books/createBook">
            <Button variant="secondary">Create a new Book</Button>
          </Link>
        </div>
      </div>
      <Separator />
      <div className="items-center pt-2.5">
        <BookList books={data} />
      </div>
    </div>
  );
};
