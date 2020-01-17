package worker

func Spawn(jobs <-chan int, results chan<- int, work func(int) int)  {
  for n := range jobs { results <- work(n) }
}
