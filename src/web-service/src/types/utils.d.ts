type Simplify<T> = {
  [R in keyof T]: T[R];
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};
