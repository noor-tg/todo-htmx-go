# TODO APP


## dependances
- [htmx](https://htmx.org) (already installed in static directory)
- [tailwindcss](https://tailwindcss.com)
- [templ](https://templ.guide)
- go

## install
- go dependances
```bash
go mod tidy
```

- tailwind
  - download [tailwind cli](https://github.com/tailwindlabs/tailwindcss/releases)
  - add to path

- install air
```bash
go install github.com/cosmtrek/air@latest
```

- install `make`

## run
```bash
air
```


## NOTES:
- hx-delete accept StatusOk to delete hx-target element
- do not use htmx put, delete requests with proxy without cors setup
- you may write the name of attributes wrong (i.e ht-post)
- it is better to use oob than js code with hx events
 
