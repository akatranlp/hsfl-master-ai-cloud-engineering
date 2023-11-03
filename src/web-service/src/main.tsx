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
        Component: Books,
      },
      {
        path: "books/:bookId/chapters/:chapterId",
        Component: Chapter,
      },
      {
        path: "transactions",
        Component: Transactions,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <RouterProvider router={router} />
      </ThemeProvider>
    </QueryClientProvider>
  </React.StrictMode>
);
