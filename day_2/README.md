Write a function that parses given input data into 2 list of structs and sorts 


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
        },
 		{
            "name": "Albert",
            "age": 3
        },
		{
            "city": "Dnipro",
            "country": "Ukraine"
        },
		{
            "name": "Roman",
            "age": 32
        },
		{
            "city": "New York City",
            "country": "US"
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
4. Pritn slices len
5. Write sort functions for each slice
6. Call sort on resulted slices
7. Print the slices
8. Output should be equal to:
```
3 4
[{Bob 36} {Alice 37} {Albert 3}] [{Ipoh Malaysia} {Dnipro Ukraine} {Northampton England} {New York City US}]

```