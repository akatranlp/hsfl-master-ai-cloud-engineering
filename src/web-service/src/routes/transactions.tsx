import { useQuery } from "@tanstack/react-query";
import { getAllTransactions } from "@/repository/books.ts";

export const Transactions = () => {
  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["transactions"],
    queryFn: getAllTransactions,
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
      <div>{JSON.stringify(data)}</div>
    </div>
  );
};
