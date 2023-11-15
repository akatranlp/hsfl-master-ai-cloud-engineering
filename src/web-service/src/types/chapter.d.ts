type UpdateChapter = {
  name: string;
  price: number;
  content: string;
};

type Chapter = {
  id: number;
  bookid: number;
  name: string;
  price: number;
  content: string;
};

type ChapterPreview = {
  id: number;
  bookid: number;
  name: string;
  price: number;
};

type CreateChapter = {
  name: string;
  bookid: number;
  price: number;
  content: string;
};
