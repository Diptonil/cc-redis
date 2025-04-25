# Go - JSON

We get this via `encoding/json`.


## 1: Approaches

- We use `Marshal()` or `NewEncoder().Encode()` for this.
- It is important to know that the benchmarked difference in performance are very slim between the two implementations. In all the cases, Encode() is better (by a pinch).
- It may make more sense to use Encode() if we are also rigging up our io. If we choose to use Marshal(), there is an extra step for conversion of the byte-array. It really is up to the implementation.
- The difference between the two implementations is that Marshal loads up the entire byte array into the memory and works with it. Encode() streams data in chunks, keeping memory overhead less.


## 2: Serialization With Marshal

```go
data := map[int]string {
    1: "One",
    2: "Two"
}

result, err := json.Marshal(&data)
if err != nil {
    return err
}

fmt.Print(string(result))
```


## 3: Deserialization With Unmarshal

```go
file, err := os.Open("file.json")
if err != nil {
    return err
}
defer file.Close()

data, err := io.ReadAll(file)
if err != nil {
    return err
}

var result map[int]string
json.Unmarshal(data, &result)

fmt.Print(string(result))
```