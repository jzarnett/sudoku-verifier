# sudoku-verifier
Simple golang sudoku solution verifier. This is by no 
means a robust and optimal system, but that's not the 
point -- it just needs to work for verifying the 
solution to the sudoku.

What's a sudoku? Allow Wikipedia to tell you: [Sudoku](https://en.wikipedia.org/wiki/Sudoku)

#Endpoints

It listens on the port 4590 for incoming HTTP POST
operations to the `/verify` endpoint. This is the only implemented endpoint and the only valid verb. 

The endpoint will return HTTP 400 if anything is wrong
with the provided data (including if the body is missing). It will return HTTP 200 if it is capable of parsing and understanding the data; if it is a valid Sudoku solution then a body of `1` is sent back; otherwise `0` is returned.

The implementation has a 50 ms delay built into it to simulate sending data over the actual internet. It is likely that for the assignment all the servers and clients will be in the same network (or at least geographically nearby) and the delay makes the test scenario a bit more plausible.

# Deployment

This requires golang to be installed. Use `go build` to compile; `go run` to execute. It is suggested to run this in a `screen` session so it remains resident and so the output doesn't spam your console.

# Testing
There are a few simple test files provided in the testfiles directory. You can test at the commandline with the curl utility using:


```
curl -X POST -H "Content-Type: application/json" --data '@filename'  -v server:4590/verify
```

Just replace the filename with the one you want and the server with the correct servername (e.g., localhost).