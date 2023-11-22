import { useQuery } from "@tanstack/react-query";
import { getChapter } from "@/repository/books.ts";
import { useParams } from "react-router-dom";
import { useMemo } from "react";
import MDEditor from "@uiw/react-md-editor";

//TODO: Publish chapter, draft view, edit chapter
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
    <div className="flex flex-col gap-10 justify-center">
      <div className="text-5xl font-bold underline text-center mt-5">{data.name}</div>
      <MDEditor.Markdown className="p-10 border border-white" source={data.content} style={{ whiteSpace: "pre-wrap" }} />
    </div>
  );
};
