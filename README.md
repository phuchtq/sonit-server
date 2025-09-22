# Sonit Server with Golang
## Introduction

## Folder hierarchy
```
└── sonit-server
    ├── api_route
    ├── cmd
    ├── constant
    │   ├── action_type
    │   ├── env
    │   ├── mail_const
    │   └── noti
    ├── data_access
    │   └── db // Handle database actions: connection, migration, ...
    ├── dto
    │   ├── request
    │   └── response
    ├── handler
    ├── html_template
    ├── interface
    │   ├── business_logic
    │   └── data_access
    ├── model 
    │   ├── dto
    │   │   ├── request
    │   │   ├── response
    │   ├── business_object
    ├── sql_script
    │   ├── migration
    │   ├── rollback
    │   └── seed
    ├── usecase
    │   ├── business_logic
    │   └── test
    ├── utils
    │   ├── cache // Redis
    │   ├── middleware
    └── main.go
```