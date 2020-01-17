package fib

func Run(n int) int {
  if n <= 1 { return n }
  return Run(n - 1) + Run(n - 2)
}
