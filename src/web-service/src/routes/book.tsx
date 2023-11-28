import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getBookById, getChaptersByBookId, getUserById } from "@/repository/books.ts";
import { createTransaction, getMyPaidTransactions } from "@/repository/transactions.ts";
import { Link, useParams } from "react-router-dom";
import { useMemo } from "react";
import { Button } from "@/components/ui/button.tsx";
import { toast } from "react-hot-toast";
import { Separator } from "@/components/ui/separator";
import { useUserData } from "@/provider/user-provider.tsx";
import { useNavigate } from "react-router-dom";
import { editChapter } from "@/repository/books.ts";
const ChapterCard = ({ transactions, authorId, chapter }: { transactions: Transaction[]; authorId: number; chapter: ChapterPreview }) => {
  const user = useUserData();
  const isOwned = transactions.some((transaction) => transaction.chapterID === chapter.id);
  return (
    //limit flex to not overflow elements

    <div className="flex dark:bg-slate-700 bg-slate-100 border rounded-lg p-4 shadow-md mb-2 items-center">
      <div className="flex-1 w-32">
        <p className="text-xl dark:text-white text-black font-semibold overflow-hidden overflow-ellipsis break-words line-clamp-2">{chapter.name}</p>
        <p className="text-lg dark:text-white text-black">
          {authorId === user.id &&
            (chapter.status === 1 ? (
              <span className="bg-green-100 text-green-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded dark:bg-green-900 dark:text-green-300">
                Published
              </span>
            ) : (
              <span className="bg-red-100 text-red-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded dark:bg-red-900 dark:text-red-300">Draft</span>
            ))}
          {!isOwned || authorId === user.id ? (
            <span className="bg-yellow-100 text-yellow-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded dark:bg-yellow-900 dark:text-yellow-300">
              {chapter.price} VV-Coins
            </span>
          ) : (
            <span className="bg-green-100 text-green-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded dark:bg-green-900 dark:text-green-300">
              Owned
            </span>
          )}
        </p>
      </div>
      <div className="flex-none justify-end">
        <ChapterButton isOwned={isOwned} authorId={authorId} chapter={chapter} />
      </div>
    </div>
  );
};

const CreateChapterButton = ({ bookid }: { bookid: number }) => {
  return (
    <>
      <div className={"ml-4"}>
        <Link to={`/books/${bookid}/chapters/createChapter`}>
          <Button variant="secondary">Create a new Chapter</Button>
        </Link>
      </div>
    </>
  );
};
const EditBookButton = ({ bookid }: { bookid: number }) => {
  return (
    <>
      <div className={"ml-4"}>
        <Link to={`/books/${bookid}/edit`}>
          <Button variant="secondary">Edit Book Information</Button>
        </Link>
      </div>
    </>
  );
};
const ChapterList = ({ transactions, authorId, chapters }: { transactions: Transaction[]; authorId: number; chapters: ChapterPreview[] }) => {
  return (
    <div>
      {chapters.map((chapter) => (
        <ChapterCard key={chapter.id} chapter={chapter} authorId={authorId} transactions={transactions} />
      ))}
    </div>
  );
};

const ChapterButton = ({ isOwned, authorId, chapter }: { isOwned: boolean; authorId: number; chapter: ChapterPreview }) => {
  const user = useUserData();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { mutateAsync: buyChapter } = useMutation({
    mutationFn: ({ chapterId, bookId }: { chapterId: number; bookId: number }) => createTransaction(chapterId, bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bookTransactions"] });
      navigate(`/books/${chapter.bookid}/chapters/${chapter.id}`);
    },
  });
  const { mutate } = useMutation({
    mutationFn: (updateChapter: UpdateChapter) => editChapter(updateChapter, chapter.bookid, chapter.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books", "bookTransactions", "chapters"] });
      navigate(0); // navigate to the same page to refresh the data, this is a hacky solution, "`/books/${chapter.bookid}`" would be better but it doesn't work for some reason
    },
  });

  const isOwner = () => {
    return authorId === user.id;
  };
  const isPublished = () => {
    return chapter.status === 1;
  };

  const isBuyable = () => {
    if (user.balance >= chapter.price) return true;
  };

  return (
    <div>
      {isOwner() ? (
        <li className="list-none">
          {!isPublished() && (
            <ul>
              <Button
                variant="ghost"
                onClick={() => {
                  mutate({ status: 1 });
                }}
              >
                Publish Chapter
              </Button>
            </ul>
          )}
          <ul>
            <Link to={`/books/${chapter.bookid}/chapters/${chapter.id}`}>
              <Button variant="ghost">Read Chapter</Button>
            </Link>
          </ul>
          <ul className={"align-middle"}>
            <Link to={`/books/${chapter.bookid}/chapters/${chapter.id}/edit`}>
              <Button variant="ghost">Edit Chapter</Button>
            </Link>
          </ul>
        </li>
      ) : isOwned ? (
        <Link to={`/books/${chapter.bookid}/chapters/${chapter.id}`}>
          <Button variant="ghost">Read Chapter</Button>
        </Link>
      ) : (
        <Button
          variant="ghost"
          onClick={() =>
            isBuyable()
              ? toast.promise(buyChapter({ chapterId: chapter.id, bookId: chapter.bookid }), {
                  loading: "isLoading",
                  error: (err) => err.message,
                  success: "Purchased successfully",
                })
              : toast.error("You don't have enough VV-Coins!")
          }
        >
          Buy Chapter
        </Button>
      )}
    </div>
  );
};
const FetchedDataBook = ({ transactions, bookData, chapters }: { transactions: Transaction[]; bookData: Book; chapters: ChapterPreview[] }) => {
  const user = useUserData();
  const {
    data: authorData,
    isError: isAuthorError,
    isLoading: isAuthorLoading,
    isSuccess: isAuthorSuccess,
    error: authorError,
  } = useQuery({
    queryKey: ["users", "books"],
    queryFn: () => getUserById(bookData.authorId),
  });
  if (isAuthorLoading) {
    return <div>Loading...</div>;
  }
  if (isAuthorError) {
    return <div>Error {authorError.message}</div>;
  }
  if (!isAuthorSuccess) {
    return <div>Something went wrong with loading the book author, please try again.</div>;
  }

  const isAuthor = bookData.authorId === user.id;
  return (
    <div>
      <div className={"mt-6 mb-4 text-2xl"}>{bookData.name}</div>
      <div className={"mb-2 text-xl"}>Written by {authorData.profileName}</div>
      <div className="flex flex-row-reverse">
        {isAuthor && <CreateChapterButton bookid={bookData.id} />}
        {isAuthor && <EditBookButton bookid={bookData.id} />}
      </div>

      <Separator className={"my-2"} />
      <div className="break-words">{bookData.description}</div>
      <Separator className={"my-2"} />
      <div className={"mt-4 mb-2 text-xl"}>Chapters</div>
      <div>
        <ChapterList transactions={transactions} authorId={bookData.authorId} chapters={chapters} />
      </div>
    </div>
  );
};

export const Book = () => {
  const { bookId } = useParams();
  const parsedBookId = useMemo(() => parseInt(bookId!), [bookId]);

  const {
    data: bookData,
    isError: isBookError,
    isLoading: isBookLoading,
    isSuccess: isBookSuccess,
    error: bookError,
  } = useQuery({
    queryKey: ["books", bookId],
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

  return <FetchedDataBook transactions={data} bookData={bookData} chapters={chapterData} />;
};
