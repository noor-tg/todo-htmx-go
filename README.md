# TODO APP

## Features
- [x] add task
- [x] list tasks
- [x] edit task
- [x] remove task
- [x] toggle task status
- [x] filter
  - [x] by status
  - [x] by description
- [x] counters
  - [x] total tasks
  - [x] completed tasks
- [x] pwa with https


## Dependances
- [htmx](https://htmx.org) (already installed in static directory)
- [tailwindcss](https://tailwindcss.com)
- [templ](https://templ.guide)
- `go`
  - can use chi implementation
  - can use echo implementation
- `sqlite` as db engine
- [mkcert](https://github.com/FiloSottile/mkcert) setup ssl cert for development and test pwa

## Installation
- go dependances
```bash
go mod tidy
```

- tailwind
  - download [tailwind cli](https://github.com/tailwindlabs/tailwindcss/releases)
  - add to system path

- install air
```bash
go install github.com/cosmtrek/air@latest
```

- install `make`
- install `mkcert` from releases page
- run
```
  mkcert -install
  mkcert todo.local
```

## Run
```bash
air
```


## NOTES:
- hx-delete accept StatusOk to delete hx-target element
- do not use htmx put, delete requests with proxy without cors setup
- you may write the name of attributes wrong (i.e ht-post)
- it is better to use oob than js code with hx events
 
## TODO
- [x] refactor using echo framework
  - [x] install framework
  - [x] integrate with templ
  - [x] integrate static assets
  - [x] refactor handlers
    - [x] index
    - [x] create
    - [x] edit
    - [x] patch
    - [x] delete
    - [x] counters
  - [x] add integration tests
  - [x] validate user acceptance tests

## E2E Tests (with rod)
- [x] add task
- [x] remove task
- [x] edit task
- [x] list tasks
- [x] toggle task status
- [x] filter tasks
  - [x] by status
  - [x] by description
  - [x] by status and description
- [ ] counters
  
## License 
It is free and open source. use it however you want.
