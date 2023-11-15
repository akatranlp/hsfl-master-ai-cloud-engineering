import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { ThemeProvider } from "@/provider/theme-provider.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { Books } from "@/routes/books.tsx";
import { RootLayout } from "@/layouts/root-layout.tsx";
import { Book } from "@/routes/book.tsx";
import { Chapter } from "@/routes/chapter.tsx";
import { Transactions } from "@/routes/transactions.tsx";
import { MyBooks } from "@/routes/myBooks.tsx";
import { BoughtBooks } from "@/routes/boughtBooks.tsx";
import { CreateBook } from "@/routes/createBook.tsx";
import { CreateChapter } from "@/routes/createChapter.tsx";
import { Toaster } from "react-hot-toast";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      {
        index: true,
        Component: App,
      },
      {
        path: "books",
        Component: Books,
      },
      {
        path: "books/create",
        Component: Books,
      },
      {
        path: "books/:bookId",
        Component: Book,
      },
      {
        path: "books/:bookId/chapters/create",
        Component: Book,
      },
      {
        path: "books/:bookId/chapters/:chapterId",
        Component: Chapter,
      },
      {
        path: "transactions",
        Component: Transactions,
      },
      {
        path: "books/myBooks",
        Component: MyBooks,
      },
      {
        path: "books/boughtBooks",
        Component: BoughtBooks,
      },
      {
        path: "books/createBook",
        Component: CreateBook,
      },
      {
        path: "books/:bookId/chapters/createChapter",
        Component: CreateChapter,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <Toaster />
        <RouterProvider router={router} />
      </ThemeProvider>
    </QueryClientProvider>
  </React.StrictMode>
);
