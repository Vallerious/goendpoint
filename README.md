# goendpoint
Mock a JSON server with as little configuration as possible.
*This project is for educational purposes (me learning Go's syntax and std. library) and should not be used, at least in 
it's current form.*

## Installation
Clone the repo and run `go build` inside the main folder. The repo also contains a build for OSX
You can then add the executable to your /usr/bin folder, for Linux and Mac users. Make sure it has Execute rights
enabled for your user.
## Launch
`./goendpoint -f=test.json -p=3000 -u=user -s=pass`
where `test.json` is a file with a JSON object, defining the required properties of your schema.

```
{
  "firstName": "Valeri",
  "age": 29
}
```

The name of the file without the extension (`test`) will be the name of your resource.
So the rest call would be to `http://localhost:3000/test`.

## Supported methods
GET - all records
POST - create record
PUT - update record
