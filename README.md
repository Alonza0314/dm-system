# Device Management System

This is a useful device management for IT.

## Develop Environment

| DevOpts | Version |
| - | - |
| OS | Ubuntu 25.04 |
| go | 1.26.2 |
| nodejs | v20.20.0 |
| yarn | 1.22.22 |

## Make

| Type | Command |
| - | - |
| Make all | `make` |
| Backend | `make backend` |
| Frontend | `make frontend` |
| Run | `make run` |
| Tidy | `make tidy` |
| Lint | `make lint` |
| Test | `make test` |
| Docker Image | `make docker` |

## API Level

```bash
/api
    └─/login(POST)
    └─/logout(POST)
    └─/category(GET, POST)
    │   └─/:cate(GET, Delete)
    └─/device(POST)
        └─/:cate(GET)
        └─/:cate/:dev(GET, Delete)
```
