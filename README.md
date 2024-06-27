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
- [ ] counters
  - [ ] total tasks
  - [ ] completed tasks


## Dependances
- [htmx](https://htmx.org) (already installed in static directory)
- [tailwindcss](https://tailwindcss.com)
- [templ](https://templ.guide)
- go

## Installation
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

## Run
```bash
air
```


## NOTES:
- hx-delete accept StatusOk to delete hx-target element
- do not use htmx put, delete requests with proxy without cors setup
- you may write the name of attributes wrong (i.e ht-post)
- it is better to use oob than js code with hx events
 
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
  
