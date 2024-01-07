import { useQuery } from "@tanstack/react-query";
import { useRepository } from "@/provider/repository-provider";
import { Separator } from "@/components/ui/separator.tsx";
import { useUserData } from "@/provider/user-provider.tsx";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

const TransactionCard = ({ transaction, bookChapter }: { transaction: Transaction; bookChapter: { book: Book; chapter: ChapterPreview } }) => {
  const user = useUserData();
  return (
    <>
      <div className="px-6">
        {transaction.payingUserID === user.id ? (
          <div className="text-2xl text-red-600">Amount: -{transaction.amount} VV-Coins</div>
        ) : (
          <div className="text-2xl text-green-600">Amount: +{transaction.amount} VV-Coins</div>
        )}
        <div>Book: {bookChapter.book.name}</div>
        <div>Chapter: {bookChapter.chapter.name}</div>
      </div>
      <Separator className="my-4" />
    </>
  );
};

const TransactionListWrapper = ({ transactionLists }: { transactionLists: [Transaction[], Transaction[]] }) => {
  const { transactionRepo, bookRepo } = useRepository();
  const {
    data: transactionsBookChapterData,
    isSuccess,
    isError,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["GetAllBooksAndChapterForAllTransactions"],
    queryFn: async () => {
      const paid = transactionLists[0];
      const received = transactionLists[1];

      const cb = (transaction: Transaction) => {
        return Promise.all([
          transactionRepo.getBookFromTransaction(transaction),
          bookRepo.getChapter(transaction.bookID, transaction.chapterID),
        ]).then(([book, chapter]) => ({
          book,
          chapter,
        }));
      };

      const paidList = Promise.all(paid.map(cb));
      const receivedList = Promise.all(received.map(cb));

      return Promise.all([paidList, receivedList]);
    },
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
      <div className="items-center pt-2.5">
        <Tabs defaultValue="paid">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="paid">Bought Chapters</TabsTrigger>
            <TabsTrigger value="received">Earnings</TabsTrigger>
          </TabsList>
          <TabsContent value="paid">
            <TransactionList transactions={transactionLists[0]} bookChapterList={transactionsBookChapterData[0]} />
          </TabsContent>
          <TabsContent value="received">
            <TransactionList transactions={transactionLists[1]} bookChapterList={transactionsBookChapterData[1]} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

const TransactionList = ({
  transactions,
  bookChapterList,
}: {
  transactions: Transaction[];
  bookChapterList: { book: Book; chapter: ChapterPreview }[];
}) => {
  if (transactions.length === 0) {
    return (
      <>
        <div className="px-6 text-2xl">No entries</div>
        <Separator className="my-4" />
      </>
    );
  }
  return (
    <div>
      {transactions.map((transaction, index) => (
        <TransactionCard key={transaction.id} transaction={transaction} bookChapter={bookChapterList[index]} />
      ))}
    </div>
  );
};

export const Transactions = () => {
  const { transactionRepo } = useRepository();
  //Get Data
  const { data, isError, isLoading, isSuccess, error } = useQuery({
    queryKey: ["transactions"],
    queryFn: async () => {
      const paidTransactions = await transactionRepo.getMyPaidTransactions();
      const receivedTransactions = await transactionRepo.getMyReceivedTransactions();
      return [paidTransactions, receivedTransactions] as [Transaction[], Transaction[]];
    },
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
  //Process Data
  const paidTransactions = data[0].filter((transaction) => transaction.amount !== 0);
  const receivedTransactions = data[1].filter((transaction) => transaction.amount !== 0);
  //Output Data
  return <TransactionListWrapper transactionLists={[paidTransactions, receivedTransactions]} />;
};
