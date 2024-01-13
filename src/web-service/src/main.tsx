import React from "react";
import ReactDOM from "react-dom/client";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "react-hot-toast";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { ThemeProvider } from "@/provider/theme-provider.tsx";
import { Books } from "@/routes/books.tsx";
import { RootLayout } from "@/layouts/root-layout.tsx";
import { MainLayout } from "@/layouts/main-layout.tsx";
import { Book } from "@/routes/book.tsx";
import { Chapter } from "@/routes/chapter.tsx";
import { Transactions } from "@/routes/transactions.tsx";
import { MyBooks } from "@/routes/myBooks.tsx";
import { BoughtBooks } from "@/routes/boughtBooks.tsx";
import { CreateBook } from "@/routes/createBook.tsx";
import { EditBook } from "@/routes/editBook.tsx";
import { CreateChapter } from "@/routes/createChapter.tsx";
import { EditChapter } from "@/routes/editChapter.tsx";
import { ManageCoins } from "@/routes/manageCoins.tsx";
import { Login } from "@/routes/login.tsx";
import { Register } from "@/routes/register.tsx";

import "./index.css";
import App from "./App.tsx";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      {
        path: "login",
        Component: Login,
      },
      {
        path: "register",
        Component: Register,
      },
      {
        path: "",
        Component: MainLayout,
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
            path: "books/:bookId/edit",
            Component: EditBook,
          },
          {
            path: "books/:bookId/chapters/createChapter",
            Component: CreateChapter,
          },
          {
            path: "books/:bookId/chapters/:chapterId/edit",
            Component: EditChapter,
          },
          {
            path: "user/manageCoins",
            Component: ManageCoins,
          },
        ],
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
      <ReactQueryDevtools initialIsOpen={false} buttonPosition="top-right" />
    </QueryClientProvider>
  </React.StrictMode>
);
