## building

To build binaries just run Makefile:
```bash
make
```


Running the protocol:
```bash
./responder
./requester
```

## Example queries

Query for shipment:

```bash
query {
  personalDetails {
    name
    surname
    birth_date
    email
    BSN
  }
  passport {
    number
    expiration
    country
  }
  
	bankingDetails {
    IBAN
    bank
    nameOnCard
  }
}
```

Listing transactions:
```bash
query {
  permissionList {
    id
    transaction_id
    expiration
    title
    description
    requester_signature
    permissionNodes {
      node_id
    }
  }
}
```
