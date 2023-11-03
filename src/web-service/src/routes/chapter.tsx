import { useQuery } from "@tanstack/react-query";
import { getChapter } from "@/repository/books.ts";
import { useParams } from "react-router-dom";
import { useMemo } from "react";

export const Chapter = () => {
  const { bookId, chapterId } = useParams();

  const parsedBookId = useMemo(() => parseInt(bookId!), [bookId]);
  const parsedChapterId = useMemo(() => parseInt(chapterId!), [chapterId]);

  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["books", bookId, "chapters", chapterId],
    queryFn: () => getChapter(parsedBookId, parsedChapterId),
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
      <div>{data.name}</div>
      {data.content}
    </div>
  );
};
