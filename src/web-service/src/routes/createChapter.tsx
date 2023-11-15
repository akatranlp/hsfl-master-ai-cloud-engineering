import { useMutation } from "@tanstack/react-query";
import { createChapter } from "@/repository/books.ts";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField, FormItem, Form, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";

const createChapterSchema = z.object({
  name: z.string().min(1),
  price: z.coerce.number().min(0),
  content: z.string().min(1),
});

export const CreateChapter = () => {
  const { bookId } = useParams();
  const navigate = useNavigate();

  const { mutate } = useMutation({
    mutationFn: (chapter: CreateChapter) => createChapter(chapter),
  });

  const form = useForm<z.infer<typeof createChapterSchema>>({
    resolver: zodResolver(createChapterSchema),
    defaultValues: {
      name: "",
      price: 0,
      content: "",
    },
  });

  const onSubmit = (values: z.infer<typeof createChapterSchema>) => {
    const bookid = parseInt(bookId!, 10);
    mutate({ ...values, bookid });
    navigate(`/books/${bookid}`); //TODO
  };

  //TODO if you're not book-author, redirect to main page
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
                <Textarea placeholder="Content" {...field} />
              </FormControl>
              <FormDescription></FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Create Chapter</Button>
      </form>
    </Form>
  );
};
