import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getBookById, getChapter, getChaptersByBookId } from "@/repository/books.ts";
import { createTransaction, getMyPaidTransactions } from "@/repository/transactions.ts";
import { useNavigate } from "react-router-dom";
import { useUserData } from "@/provider/user-provider.tsx";
import { useParams } from "react-router-dom";
import { useMemo } from "react";
import MDEditor from "@uiw/react-md-editor";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button.tsx";
import { toast } from "react-hot-toast";
//TODO: Publish chapter, draft view, edit chapter
export const Chapter = () => {
  const { bookId, chapterId } = useParams();
  const user = useUserData();
  const parsedBookId = useMemo(() => parseInt(bookId!), [bookId]);
  const parsedChapterId = useMemo(() => parseInt(chapterId!), [chapterId]);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  //Fetch data
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
    queryKey: ["books", bookId, "chapters", chapterId],
    queryFn: () => getChapter(parsedBookId, parsedChapterId),
  });
  const {
    data: allChaptersData,
    isError: isallChaptersError,
    isLoading: isallChaptersLoading,
    isSuccess: isallChaptersSuccess,
    error: allChaptersError,
  } = useQuery({
    queryKey: ["books", bookId, "chapters"],
    queryFn: () => getChaptersByBookId(parsedBookId),
  });
  const {
    data: transactionsData,
    isError: isTransactionsError,
    isLoading: isTransactionsLoading,
    isSuccess: isTransactionsSuccess,
    error: transactionsError,
  } = useQuery({
    queryKey: ["booksTransactions"],
    queryFn: () => getMyPaidTransactions(),
  });

  const { mutateAsync: buyChapter } = useMutation({
    mutationFn: ({ buyChapterId, buyBookId }: { buyChapterId: number; buyBookId: number }) => createTransaction(buyChapterId, buyBookId),
    onSuccess: (_, variables) => {
      const { buyChapterId, buyBookId } = variables;
      queryClient.invalidateQueries({ queryKey: ["booksTransactions"] });
      queryClient.invalidateQueries({ queryKey: ["books", bookId, "chapters"] });
      navigate(`/books/${buyBookId}/chapters/${buyChapterId}`);
    },
  });

  if (isChapterLoading || isallChaptersLoading || isTransactionsLoading || isBookLoading) {
    return <div>Loading...</div>;
  }
  if (isBookError) {
    return <div>Error {bookError.message}</div>;
  }

  if (isChapterError) {
    return <div>Error {chapterError.message}</div>;
  }
  if (isallChaptersError) {
    return <div>Error {allChaptersError.message}</div>;
  }
  if (isTransactionsError) {
    return <div>Error {transactionsError.message}</div>;
  }

  if (!isChapterSuccess || !isallChaptersSuccess || !isTransactionsSuccess || !isBookSuccess) {
    return <div>Something went wrong!</div>;
  }

  let isOwner = false;
  if (bookData) {
    isOwner = bookData.authorId === user.id;
  }
  //Does it have a previous chapter?
  const previousChapter = allChaptersData.find((chapter) => chapter.id === parsedChapterId - 1);
  const hasPreviousChapter = parsedChapterId > 1;
  //Does it have a published next chapter?
  let hasNextChapter = false;
  const nextChapter = allChaptersData.find((chapter) => chapter.id === parsedChapterId + 1);
  if (nextChapter) {
    hasNextChapter = nextChapter.status === 1;
  }
  //Is the previous chapter bought?
  let isPreviousChapterBought = false;
  if (previousChapter) {
    isPreviousChapterBought = transactionsData.some(
      (transaction) => transaction.chapterID === previousChapter!.id && transaction.bookID === previousChapter!.bookid
    );
  }
  let isPreviousChapterBuyable = false;
  if (previousChapter) {
    isPreviousChapterBuyable = previousChapter.price <= user.balance;
  }
  //Is the next chapter bought?
  let isNextChapterBought = false;
  if (nextChapter) {
    isNextChapterBought = transactionsData.some(
      (transaction) => transaction.chapterID === nextChapter!.id && transaction.bookID === nextChapter!.bookid
    );
  }
  let isNextChapterBuyable = false;
  if (nextChapter) {
    isNextChapterBuyable = nextChapter.price <= user.balance;
  }

  return (
    <div className="flex flex-col gap-10 justify-center">
      <div className="text-5xl font-bold underline text-center mt-5 break-words">{chapterData.name}</div>
      <ul className="flex flex-row gap-10 justify-center">
        <li>
          {hasPreviousChapter ? (
            isPreviousChapterBought || isOwner ? (
              <Link to={`/books/${parsedBookId}/chapters/${parsedChapterId - 1}`}>
                <Button variant="ghost">Previous Chapter</Button>
              </Link>
            ) : (
              <Button
                variant="ghost"
                onClick={() =>
                  isPreviousChapterBuyable
                    ? toast.promise(buyChapter({ buyChapterId: previousChapter!.id, buyBookId: previousChapter!.bookid }), {
                        loading: "isLoading",
                        error: (err) => err.message,
                        success: "Purchased successfully",
                      })
                    : toast.error("You don't have enough VV-Coins!")
                }
              >
                Buy Next Chapter
              </Button>
            )
          ) : null}
        </li>
        <li>
          <Link to={`/books/${parsedBookId}/`}>
            <Button variant="ghost">Book Overview</Button>
          </Link>
        </li>
        <li>
          {hasNextChapter ? (
            isNextChapterBought || isOwner ? (
              <Link to={`/books/${parsedBookId}/chapters/${parsedChapterId + 1}`}>
                <Button variant="ghost">Next Chapter</Button>
              </Link>
            ) : (
              <Button
                variant="ghost"
                onClick={() =>
                  isNextChapterBuyable
                    ? toast.promise(buyChapter({ buyChapterId: nextChapter!.id, buyBookId: nextChapter!.bookid }), {
                        loading: "isLoading",
                        error: (err) => err.message,
                        success: "Purchased successfully",
                      })
                    : toast.error("You don't have enough VV-Coins!")
                }
              >
                Buy Next Chapter
              </Button>
            )
          ) : null}
        </li>
      </ul>

      <MDEditor.Markdown className="p-10 border border-white" source={chapterData.content} style={{ whiteSpace: "pre-wrap" }} />
    </div>
  );
};
