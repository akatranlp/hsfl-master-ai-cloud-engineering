type Simplify<T> = {
  [R in keyof T]: T[R];
} & {};
