import { useQuery } from "@tanstack/react-query";
import { getBookById, getChaptersByBookId } from "@/repository/books.ts";
import { Link, useParams } from "react-router-dom";
import { useMemo } from "react";

const ChapterCard = ({ bookId, chapter }: { bookId: number; chapter: Chapter }) => {
  return (
    <div>
      <Link to={`/books/${bookId}/chapters/${chapter.id}`}>{chapter.name}</Link>
    </div>
  );
};

const ChapterList = ({ bookId, chapters }: { bookId: number; chapters: Chapter[] }) => {
  return (
    <div>
      {chapters.map((chapter) => (
        <ChapterCard key={chapter.id} bookId={bookId} chapter={chapter} />
      ))}
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

  if (isBookLoading || isChapterLoading) {
    return <div>Loading...</div>;
  }

  if (isBookError) {
    return <div>Error {bookError.message}</div>;
  }

  if (isChapterError) {
    return <div>Error {chapterError.message}</div>;
  }

  if (!isBookSuccess || !isChapterSuccess) {
    return <div>Something went wrong with loading the book data, please try again.</div>;
  }

  return (
    <div>
      <div>Book: {bookData.name}</div>
      <div>Description: {bookData.description}</div>
      <div>Chapters:</div>
      <div>
        <ChapterList bookId={bookData.id} chapters={chapterData} />
      </div>
    </div>
  );
};
