# Devlog: Go


## 1: Structs

- Go is not OOP by design but it facilitates the creation of complex data types with the use of structs.
- Start the name of the struct with an uppercase letter if it has to be a public data member. For private data, have it lowercase.
- Methods are also supported (in a weird way) with the use of recievers. While writing recievers, we are dealing with pass-by-value by default.
- If we wish to do pass-by-reference to actually alter 'behaviours' of our object (struct data), we use pointers.
- While still a matter of contention, getters and setters are done in Go in the way below. Go doesn't mandate it and by its design, it seems unnecessary.
- In Go, it is good practice to write `Ram()` as a function name rather than `getRam()`.

```go
type Laptop struct {
    cpu string
    ram int
    manufacturer string
}

func (laptop *Laptop) setRam(newRam int) {
    laptop.ram = newRam
}

func (laptop *Laptop) getRam() int {
    return laptop.ram
}
```


## 2: Constants

- They work slightly weird in Go. The type system is a little complex and it shows here. For 99% of the use-cases, it isn't much of a big deal.
- Do constant declarations in one line like this: `const Pi = 3.14`. This is untyped constant declaration.
- Typed constant declarations exist too: `const Pi float = 3.14`. For all your regular use-cases both typed and untyped can be considered to be the same.
- The main difference occurs in this weird case. This is too niche of a case to occur in code, I feel:

```go
type NewBool bool
const isConsidered = true
const isNotConsidered bool = true
var isItReallyConsidered NewBool

// This will work because isConsidered is a type-inferred bool.
isItReallyConsidered = isConsidered

// This will not work because isNotConsidered is a strictly typed bool (and not NewBool).
isItReallyConsidered = isNotConsidered
```


## 3: Concurrency & Parallelism

- Concurrency is the handling of different jobs by a single thread. Thread A does a bit of Work 1, then a bit of work 2, then a bit of work 3. In this way, all Works get serviced. But it is one worker doing all the work. 
- Parallelism is handling of different jobs all at the same time by different workers. Multiple threads do different work.
- Goroutines are used for concurrency in Go. Main() itself is a goroutine that runs first. The important thing is that once the main goroutine ends, the whole program ends. This means that even if the other goroutines are not done with their jobs, the program will terminate. This is why, if we take an elementary example of goroutines, it comes with the time.sleep() statement to pause the main goroutine from ending before the others.

```go
func simpleConc() {
    for i := 0; i < 3; i++ {
        go func(index int) {
            fmt.Println(time.Now(), index)
        }(i)
    }

    time.Sleep(time.Second)
    fmt.Println("done")
}
```


## 4: Unbuffered Channels

- Channels are how data is exchanged in goroutines. They are unbuffered - data sent in a channel needs to be recieved right at that instance.
- Channels are strongly typed.
- This means sequential programming would not work in channels. The other way to make them work outside goroutines is to close them.
- A big use-case of channels is to block the main goroutine so that the other goroutines can execute completely before main exits. Such channels are generally created with `make(chan struct{})`.


## 5: Buffered Channels

- Here, data can be stored. No need to send and recieve at the exact moment.
- It cannot recieve data from an empty channel.
- Also, it cannot send more data to a channel than specified.
- This would work:

```go
func bufferedCh() {
    ch := make(chan int, 1)
    ch <- 1
    res := <-ch
    fmt.Println(res)
}
```


## 6: Deferring

- The defer keyword if used to execute a function after the current function returns (normally or due to an error).
- Its main uses are to close any connections, files or resources. This prevents memory leaks or any ugly situations. It is also used to release mutexes.


## 7: `go.mod`

- This file is used to track project dependencies and module information.
- It stores the go version being used as well as the main place where the module is available (be it a GitHub repo, a website, etc.).
- As dependencies are added in, their versions would get listed down. This is to keep builds reproducibale and consistent.
- It's kind of like `package.json`.


## 8: Walrus Operator

- It does type inference based on value of the right hand side.
- We can also do explicit type declarations. That would mean going the long way about it.

```go
var message string
message = "Hi!"
```
