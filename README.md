# Sonit Server with Golang
## Introduction
Sonit is a platform which supports users for purchasing, customizing billiard cues with multiple services seperated into wide range of prices. Sonit Server contains the core logic of platform to manage users' requests
## Folder hierarchy
```
└── sonit-server
    ├── api_route
    ├── cmd
    ├── constant
    ├── data_access
    │   └── db // Handle database actions: connection, migration, ...
    │   └── cache
    │   └── db_server
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
    │   ├── business_object
    ├── sql_script
    │   ├── migration
    │   ├── rollback
    │   └── seed
    ├── usecase
    │   ├── business_logic
    │   └── test
    ├── utils
    │   ├── file
    │   ├── middleware
    └── main.go
```


## Reference
Reference: [Original Repository Name](https://github.com/Sonit-Custom/sonit-server).  
