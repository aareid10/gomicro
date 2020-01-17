package main

import (
  "gomicro/microstr"
  "gomicro/micronum"
  "net/http"
  "log"
  "os"
)

func main()  {
  http.Handle("/uppercase", microstr.InitializeUppercase())
  http.Handle("/count", microstr.InitializeCount())
	http.Handle("/fibc", micronum.InitializeFibc())
  port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Fatal(http.ListenAndServe(":"+ port, nil))
}
