# cc-redis

An implementation of a lite Redis database (eighth in a series of coding challenges by John Crickett). This entails no external dependencies. The only built-ins used are:
- `encoding/json`: To serialize & deserialize data to and from files.
- `io`: For I/O operations on files.
- `bufio`: For buffered I/O. This is utilized in establishing new network connections.
- `net`: For TCP connection management (the way it was done in Redis).
- `log/slog`: For structured logging.
- `os`: For file operations.
- `strconv`: For interconversions between strings and integers.
- `strings`: For string operations (trimming, equality, etc.).
- `time`: To sleep the thread (or should I say goroutine...?).


## Files

- `Makefile`: Build automation tool.
- `main.go`: Application code.
- `go.mod`: Dependency tracking & versioning.
- `data.json`: File that gets written to when current data-store state is saved. Empty as such.


## Logic

The idea was to implement a 'lite' version of Redis without sweating over the critical implementations and optimizations that exist within the real app.

Here is what we have:
- A basic data store that supports all GET, SET & DELETE operations handled in Redis.
- A TCP-server capable of handling concurrent requests. This is made possible with the use of goroutines.
- Support for expiring data after a given time using EX & PX directives with SET.
- Support for SAVE directive. Saved data is loaded by default upon start-up.

Here is what it cannot do (because I did not have time to implement it and the lite version is a stripped down version anyway):
- Autosaves: This is a straightforward enhancement. Spin up a goroutine to call `Save()` for state persistence. We can either do this upon every request, or we may leave a goroutine running in a loop that periodically calls save and sleeps.
- Scale-Outs: I actually am not sure how to handle it. As of writing, I am not entirely sure (and haven't yet researched) how to run this application within a container and have replicas of the containers in sync.
- The niche Redis commands: There are a few of these. I would rather make my own data store than implement these here.

Shortcomings:
- Structure: This is my first Go project. It seems to be quite badly structured with just one main.go. We can definitely do better.
- Perhaps I can make use of channels to communicate when the expiry options are set wrongly instead of doing a server-aware return (which client would not know)


## How are Goroutines Implemented Here?

- One goroutine spins up for each connection made to the server.
- One goroutine spins up for handling each command request made to the server by a single connection.
- One goroutine spins up in case expiry of data is needed while `SET` is used.


## Usage

To simplify things, just run this to start the server:

```
make run-server
```

To run a client (running this on a Mac/Linux):

```
make run-client
```
