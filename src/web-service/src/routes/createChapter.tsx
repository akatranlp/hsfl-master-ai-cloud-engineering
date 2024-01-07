import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormField, FormItem, Form, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import rehypeSanitize from "rehype-sanitize";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import MDEditor from "@uiw/react-md-editor";
import toast from "react-hot-toast";
import { useRepository } from "@/provider/repository-provider";

const createChapterSchema = z.object({
  name: z.string().min(1),
  price: z.coerce.number().min(0),
  content: z.string().min(1),
});

export const CreateChapter = () => {
  const { bookId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { bookRepo } = useRepository();

  const { mutate } = useMutation({
    mutationFn: (chapter: CreateChapter) => bookRepo.createChapter(chapter),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chapters"] });
      navigate(`/books/${bookId}`);
    },
    onError: (error) => {
      toast.error(error.message);
    },
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
        <Button variant="secondary" type="submit">
          Create Chapter
        </Button>
      </form>
    </Form>
  );
};
