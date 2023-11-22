import { useMutation, useQueryClient } from "@tanstack/react-query";
import { editChapter } from "@/repository/books.ts";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField, FormItem, Form, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { getChapter } from "@/repository/books.ts";
import { useQuery } from "@tanstack/react-query";
import rehypeSanitize from "rehype-sanitize";
import MDEditor from "@uiw/react-md-editor";

const editChapterSchema = z.object({
  name: z.string().min(1),
  price: z.coerce.number().min(0),
  content: z.string().min(1),
  status: z.coerce.number().min(0),
});

const EditChapterForm = ({ chapter }: { chapter: Chapter }) => {
  const form = useForm<z.infer<typeof editChapterSchema>>({
    resolver: zodResolver(editChapterSchema),
    defaultValues: {
      name: chapter.name,
      price: chapter.price,
      content: chapter.content,
      status: chapter.status,
    },
  });
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { mutate } = useMutation({
    mutationFn: (updateChapter: UpdateChapter) => editChapter(updateChapter, chapter.bookid, chapter.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chapters", chapter.id] });
      navigate(`/books/${chapter.bookid}/chapters/${chapter.id}`);
    },
  });

  const onSubmit = (values: z.infer<typeof editChapterSchema>) => {
    mutate(values);
  };

  const isPublished = () => {
    return chapter.status === 1;
  };
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Title" {...field} />
              </FormControl>
              <FormDescription></FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="price"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Price</FormLabel>
              <FormControl>
                <Input type="number" placeholder="Price" {...field} />
              </FormControl>
              <FormDescription></FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="content"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Content</FormLabel>
              <FormControl>
                <MDEditor
                  {...field}
                  previewOptions={{
                    rehypePlugins: [[rehypeSanitize]],
                  }}
                />
              </FormControl>
              <FormDescription></FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Update Chapter</Button>
        {!isPublished() && (
          <Button
            type="submit"
            onClick={() => {
              form.setValue("status", 1);
            }}
          >
            Update & Publish Chapter
          </Button>
        )}
      </form>
    </Form>
  );
};

//Takes the bookId and chapterId from the url and requests the chapter from the server to fill the form, then the auther can submit the form to the server to update the chapter
export const EditChapter = () => {
  const { bookId, chapterId } = useParams();
  const parsedBookId = parseInt(bookId!);
  const parsedChapterId = parseInt(chapterId!);

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

  //TODO if you're not book-author, redirect to main page
  return <EditChapterForm chapter={data}></EditChapterForm>;
};
