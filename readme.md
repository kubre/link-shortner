# Link Shortener

I wrote this as just try out GO language

### To run

```sh
go run main.go domain.com
```

Shorten a link

```sh
curl -X POST -H "Content-Type: application/json" -d '{"link":"https://google.com"}' https://domain.com

{"success": true, "link": "domain.com/abcdef"}

```

Get Full link

```sh
curl https://domain.com/abcdef

{"success": true, "link":"https://google.com"}
```
