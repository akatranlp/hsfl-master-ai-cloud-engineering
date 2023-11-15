import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getBookById, getChaptersByBookId } from "@/repository/books.ts";
import { createTransaction, getMyPaidTransactions } from "@/repository/transactions.ts";
import { Link, useParams } from "react-router-dom";
import { useMemo } from "react";
import { Button } from "@/components/ui/button.tsx";
import { toast } from "react-hot-toast";
import { Separator } from "@/components/ui/separator";
import { useUserData } from "@/provider/user-provider.tsx";

const ChapterCard = ({ transactions, authorId, chapter }: { transactions: Transaction[]; authorId: number; chapter: Chapter }) => {
  return (
    <div className="flex dark:bg-slate-700 bg-slate-100 border rounded-lg p-4 shadow-md mb-2 items-center">
      <div className="flex-grow">
        <p className="text-xl dark:text-white text-black font-semibold">{chapter.name}</p>
        <p className="text-lg dark:text-white text-black">Price: {chapter.price} VV-Coins</p>
      </div>
      <div className="flex-shrink-0 justify-end ">
        <ChapterButton transactions={transactions} authorId={authorId} chapter={chapter} />
      </div>
    </div>
  );
};

const CreateChapterButton = ({ bookid }: { bookid: number }) => {
  return (
    <>
      <div className={"justify-end px-6 pt-2.5 items-end"}>
        <div className="ml-auto m-4">
          <Link to={`/books/${bookid}/chapters/createChapter`}>
            <Button>Create a new Chapter</Button>
          </Link>
        </div>
      </div>
    </>
  );
};

const ChapterList = ({ transactions, authorId, chapters }: { transactions: Transaction[]; authorId: number; chapters: Chapter[] }) => {
  return (
    <div>
      {chapters.map((chapter) => (
        <ChapterCard key={chapter.id} chapter={chapter} authorId={authorId} transactions={transactions} />
      ))}
    </div>
  );
};

const ChapterButton = ({ transactions, authorId, chapter }: { transactions: Transaction[]; authorId: number; chapter: Chapter }) => {
  const user = useUserData();
  const queryClient = useQueryClient();

  const { mutateAsync: buyChapter } = useMutation({
    mutationFn: createTransaction,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["bookTransactions", "chapters"] }),
  });

  if (authorId === user.id) return;

  const isOwned = () => {
    let owned = false;

    transactions.forEach((transaction) => {
      if (transaction.chapterID === chapter.id) {
        owned = true;
      }
    });

    return owned;
  };

  const isBuyable = () => {
    if (user.balance >= chapter.price) return true;
  };

  return (
    <div>
      {isOwned() ? (
        <Link to={`/books/${chapter.bookid}/chapters/${chapter.id}`}>
          <Button variant="ghost">Read Chapter</Button>
        </Link>
      ) : (
        <Button
          variant="ghost"
          onClick={() =>
            isBuyable()
              ? toast.promise(buyChapter(chapter.id), { loading: "isLoading", error: (err) => err.message, success: "Purchased successfully" })
              : toast.error("You don't have enough VV-Coins!")
          }
        >
          Buy Chapter
        </Button>
      )}
    </div>
  );
};

//TODO: Edit book
export const Book = () => {
  const { bookId } = useParams();
  const parsedBookId = useMemo(() => parseInt(bookId!), [bookId]);

  const user = useUserData();

  const {
    data: bookData,
    isError: isBookError,
    isLoading: isBookLoading,
    isSuccess: isBookSuccess,
    error: bookError,
  } = useQuery({
    queryKey: ["book", bookId],
    queryFn: () => getBookById(parsedBookId),
  });

  const {
    data: chapterData,
    isError: isChapterError,
    isLoading: isChapterLoading,
    isSuccess: isChapterSuccess,
    error: chapterError,
  } = useQuery({
    queryKey: ["books", bookId, "chapters"],
    queryFn: () => getChaptersByBookId(parsedBookId),
  });

  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["booksTransactions"],
    queryFn: () => getMyPaidTransactions(),
  });

  if (isBookLoading || isChapterLoading || isLoading) {
    return <div>Loading...</div>;
  }

  if (isBookError) {
    return <div>Error {bookError.message}</div>;
  }

  if (isChapterError) {
    return <div>Error {chapterError.message}</div>;
  }

  if (isError) {
    return <div>Error {error.message}</div>;
  }

  if (!isBookSuccess || !isChapterSuccess || !isSuccess) {
    return <div>Something went wrong with loading the book data, please try again.</div>;
  }

  const isAuthor = bookData.authorId === user.id;
  return (
    <div>
      <ul className="flex">
        <li className={"m-6 text-2xl"}>{bookData.name}</li>
        <li className="ml-auto align-middle">{isAuthor && <CreateChapterButton bookid={bookData.id} />}</li>
      </ul>
      <Separator className={"my-2"} />
      <div>{bookData.description}</div>
      <Separator className={"my-2"} />
      <div className={"mt-4 mb-2"}>Chapters:</div>
      <div>
        <ChapterList transactions={data} authorId={bookData.authorId} chapters={chapterData} />
      </div>
    </div>
  );
};
