import { useMutation, useQueryClient } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField, FormItem, Form, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useRepository } from "@/provider/repository-provider";

export const editBookSchema = z.object({
  name: z.string().min(1),
  description: z.string().min(1),
});

const EditBookForm = ({ book }: { book: Book }) => {
  const form = useForm<z.infer<typeof editBookSchema>>({
    resolver: zodResolver(editBookSchema),
    defaultValues: {
      name: book.name,
      description: book.description,
    },
  });
  const queryClient = useQueryClient();
  const { bookRepo } = useRepository();
  const navigate = useNavigate();
  const { mutate } = useMutation({
    mutationFn: (updateBook: UpdateBook) => bookRepo.editBook(updateBook, book.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books", book.id] });
      navigate(`/books/${book.id}`);
    },
  });

  const onSubmit = (values: z.infer<typeof editBookSchema>) => {
    mutate(values);
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
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea placeholder="Description" {...field} />
              </FormControl>
              <FormDescription></FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button variant="secondary" type="submit">
          Update Book
        </Button>
      </form>
    </Form>
  );
};

//Takes the bookId and chapterId from the url and requests the chapter from the server to fill the form, then the auther can submit the form to the server to update the chapter
export const EditBook = () => {
  const { bookId } = useParams();
  const parsedBookId = parseInt(bookId!);
  const { bookRepo } = useRepository();

  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["books", bookId],
    queryFn: () => bookRepo.getBookById(parsedBookId),
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
  return <EditBookForm book={data}></EditBookForm>;
};
