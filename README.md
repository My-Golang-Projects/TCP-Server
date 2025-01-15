# TCP Server Documentation

## Overview

This Go program sets up a TCP server that listens for incoming client connections, processes the data sent by clients, and acknowledges the receipt of messages. The server uses goroutines and channels to handle multiple client connections concurrently.

> This is by no chance production ready and shouldn't be used in that use case

## Flow of Execution

1. **Server Creation:**
   - A `Server` object is created, which holds configuration information (e.g., listening address), a quit channel (`quitch`), and a message channel (`msgch`) to process incoming data.
2. **Starting the Server:**

   - The server listens on a specific address (e.g., `:3000`) using the `Start` method.
   - The `acceptLoop` method is invoked concurrently to accept incoming connections.

3. **Accept Loop:**

   - The `acceptLoop` is an infinite loop that continuously accepts incoming TCP connections.
   - For each new connection, the server spawns a new goroutine to handle communication with the client via the `readLoop` method.

4. **Read Loop:**

   - The `readLoop` listens for data from the client connection.
   - When data is received, it is read into a buffer and sent to the `msgch` channel for further processing.
   - The server acknowledges the client by sending a response: `"Thank You For Your Message!"`.

5. **Message Channel:**

   - The server listens for incoming messages through the `msgch` channel.
   - Each message includes the sender's address and the received payload.
   - The server prints out the details of each received message.

6. **Graceful Shutdown:**
   - The server waits for the `quitch` channel to close before shutting down.
   - When the server is done, it closes the message channel (`msgch`) to notify that no more messages will be received.

## Code Structure

### Types

- **Message:**

  - Represents a message received from a client.
  - Fields:
    - `from`: The remote address of the sender.
    - `payload`: The message content sent by the client.

- **Server:**
  - Contains all necessary information for the TCP server.
  - Fields:
    - `listenAddr`: The address the server listens on (e.g., `:3000`).
    - `ln`: The TCP listener instance.
    - `quitch`: A channel used for server shutdown signaling.
    - `msgch`: A channel used to handle incoming messages.

### Functions

- `**NewServer(listenAddr string) *Server:**`

  - Creates a new server with the provided listen address.
  - Initializes the `quitch` and `msgch` channels.

- `**(s *Server) Start() error:**`

  - Starts the TCP server by listening on the specified address.
  - Calls `acceptLoop` to accept incoming connections.
  - Waits for the `quitch` channel to close, signaling shutdown.

- `**(s *Server) acceptLoop():**`

  - Accepts incoming TCP connections in an infinite loop.
  - For each connection, it starts a new goroutine that reads data from the connection via `readLoop`.

- `**(s *Server) readLoop(conn net.Conn):**`
  - Handles the reading and processing of data from the client connection.
  - Each received message is sent to the `msgch` channel, and an acknowledgment is sent back to the client.

### Main Function

- The `main` function creates a new server, starts it, and listens for messages through the `msgch` channel.
- Each message is printed with the sender's address and the message content.

## Example Usage

1. Start the server:
   ```bash
   go run server.go
   ```
2. The server will now listen on port `3000`. You can connect to it using a TCP client (e.g., using `telnet`):

   ```bash
   telnet localhost 3000
   ```

3. After sending a message to the server, it will acknowledge the connection:

   ```bash
   Thank You For Your Message!
   ```

4. The server will print the message received from the client:
   ```bash
   Received message from connection (clientAddress):(messageContent):
   ```
