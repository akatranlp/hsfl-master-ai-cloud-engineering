import { z } from "zod";
import { FormField, FormItem, Form, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form.tsx";
import { login } from "@/repository/user.ts";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Link, useNavigate, Navigate } from "react-router-dom";
import toast from "react-hot-toast";
import { useUserData } from "@/provider/user-provider.tsx";

const loginSchema = z.object({
  email: z.string().min(1),
  password: z.string().min(1),
});

export const Login = () => {
  const user = useUserData();
  const navigate = useNavigate();

  const { mutate } = useMutation<void, unknown, z.infer<typeof loginSchema>>({
    mutationFn: (data) => login(data.email, data.password),
    onSuccess: () => {
      navigate("/books");
    },
    onError: () => {
      toast.error("An error occurred. Please try again.");
    },
  });

  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = (values: z.infer<typeof loginSchema>) => {
    mutate(values);
  };

  if (user.id !== 0) {
    <Navigate to="/books" />;
  }

  return (
    <div className="flex justify-center pt-5">
      <Card className="w-1/2">
        <CardHeader>
          <CardTitle>Login</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder="Email" {...field} />
                    </FormControl>
                    <FormDescription></FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder="Password" {...field} />
                    </FormControl>
                    <FormDescription></FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className={"flex justify-end px-6"}>
                <div className="m-4">
                  <Button variant="secondary" type="submit">
                    Login
                  </Button>
                </div>
              </div>
            </form>
          </Form>
        </CardContent>
        <CardFooter>
          <Button variant="secondary">
            <Link to="/register">Register instead</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};
