How to Use

Start the Server:
Build and run the server program. It will print out the session seed (session ID and secret) and then wait for QR code input.
go build -o server main.go
./server

Generate QR Codes:
In another terminal, build and run the client program. Supply the session ID and secret (as printed by the server) via flags.
go build -o client main.go
./client -session=<SESSION_ID> -secret=<SECRET>

The client will output 10 iterations (one per second). Choose one of these (copy the entire CSV string).

Submit the QR Code:
Paste the chosen QR code content into the server’s prompt and hit Enter. The server will verify the token and print whether it is valid or not. Regardless of the result, the session is invalidated after the first submission.

This implementation uses a single initial seed call, generates animated QR code tokens that change every second and enforces a 10‑second validity window.
