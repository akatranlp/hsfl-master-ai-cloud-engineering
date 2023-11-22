type UpdateChapter = {
  name?: string;
  price?: number;
  content?: string;
  status?: number;
};

type Chapter = {
  id: number;
  bookid: number;
  name: string;
  price: number;
  content: string;
  status: number;
};

type ChapterPreview = {
  id: number;
  bookid: number;
  name: string;
  price: number;
  status: number;
};

type CreateChapter = {
  name: string;
  bookid: number;
  price: number;
  content: string;
};
