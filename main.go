package main

import (
  "gomicro/microstr"
  "gomicro/micronum"
  "gomicro/microga"
  "net/http"
  "log"
  "os"
)

func main()  {
  http.Handle("/uppercase", microstr.InitializeUppercase())
  http.Handle("/count", microstr.InitializeCount())
  http.Handle("/fibc", micronum.InitializeFibc())
	http.Handle("/gac", microga.InitializeGAc())
  port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Fatal(http.ListenAndServe(":"+ port, nil))
}
