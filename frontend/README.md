# deweb Frontend

All functions that are in `/lib` of this repo are available inside of the frontend as a JS functions.

For example `/lib/getselfid.go` contain the following funcction:

```go
func GetSelfID() SelfID {
    [...]
    return SelfID{
		OK: true,
		ID: "aaaaa-aaaaa-aaaaaa-aaaaa-aaaaa[key=asdasdasdasdasdasdasd]",
	}
}
```

It means that in JavaScript to retrieve the following struct you need to write a function like this:
```javascript
getSelfID().then((v) => {
    console.log(v)    // [object Object]
    console.log(v.OK)    // true
    console.log(v.ID) // "aaaaa-aaaaa-aaaaaa-aaaaa-aaaaa[key=asdasdasdasdasdasdasd]"
})
```
No need to use .catch, golang does not return errors, instead usually there will be a `OK` field, together with an `Error`.

There is no JavaScript documentation, instead future frontend developer should follow the Golang documentation, optainable thru `godoc -v -http=localhost:6060` (to ensure that you run on the latest version).

