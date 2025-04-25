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
