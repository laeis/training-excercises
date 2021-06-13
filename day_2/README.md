Write a function that parses given input data into 2 list of structs 


```
var jsonStr = []byte(`
{
    "things": [
        {
            "name": "Alice",
            "age": 37
        },
        {
            "city": "Ipoh",
            "country": "Malaysia"
        },
        {
            "name": "Bob",
            "age": 36
        },
        {
            "city": "Northampton",
            "country": "England"
        }
    ]
}`)

```

With the following structure

```
func solution(jsonStr []byte) (people []Person, places []Place) {
    // your code
}
```

1. Define structs 
2. Parse json
3. Put data into corresonding slices
4. Print the slices
5. Output should be equal to:
```
[{Alice 37} {Bob 36}] [{Ipoh Malaysia} {Northampton England}]

```