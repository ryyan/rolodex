# rolodex

Simple address book, with a REST API built in Go

## Getting started

Build and run rolodex

```
go build
./rolodex
```

Use through the address route

```
http://localhost:8888/address
```

## Demo

Used [httpie](https://github.com/httpie/httpie) to run commands

```sh
# Create
$ http --form POST localhost:8888/address/ firstname=John lastname=Smith email=johnsmith@email.com phonenumber=123-456-7890
{
    "rolodex": [
        {
            "Email": "johnsmith@email.com",
            "FirstName": "John",
            "ID": "zYszZcecnwxL",
            "LastName": "Smith",
            "PhoneNumber": "123-456-7890"
        }
    ]
}

$ http --form POST localhost:8888/address/ firstname=Jane lastname=Smith email=janesmith@email.com phonenumber=321-654-0987
{
    "rolodex": [
        {
            "Email": "janesmith@email.com",
            "FirstName": "Jane",
            "ID": "qCMcgOPMtdbw",
            "LastName": "Smith",
            "PhoneNumber": "321-654-0987"
        }
    ]
}

# Update
$ http --form PUT localhost:8888/address/qCMcgOPMtdbw firstname=Jane lastname=Doe email=janedoe@email.com phonenumber=321-654-0987
{
    "rolodex": [
        {
            "Email": "janedoe@email.com",
            "FirstName": "Jane",
            "ID": "qCMcgOPMtdbw",
            "LastName": "Doe",
            "PhoneNumber": "321-654-0987"
        }
    ]
}

# Get
$ http GET localhost:8888/address/
{
    "rolodex": [
        {
            "Email": "johnsmith@email.com",
            "FirstName": "John",
            "ID": "zYszZcecnwxL",
            "LastName": "Smith",
            "PhoneNumber": "123-456-7890"
        },
        {
            "Email": "janedoe@email.com",
            "FirstName": "Jane",
            "ID": "qCMcgOPMtdbw",
            "LastName": "Doe",
            "PhoneNumber": "321-654-0987"
        }
    ]
}

# Delete
$ http --form DELETE localhost:8888/address/YhUkllfmMFFY
{
    "rolodex": null
}

$ http GET localhost:8888/address/
{
    "rolodex": [
        {
            "Email": "johnsmith@email.com",
            "FirstName": "John",
            "ID": "zYszZcecnwxL",
            "LastName": "Smith",
            "PhoneNumber": "123-456-7890"
        }
    ]
}
```
