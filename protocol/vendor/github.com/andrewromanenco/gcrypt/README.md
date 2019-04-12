# gcrypt
Utility for GoLang to simplify data encryption.

## Supported features
 - Key derivation from a user provided password (256 bit)
 - AES 256 Data encryption with HMAC (Encrypt-and-Mac)
 - Data decryption with HMAC validation

## Key generation and salt
Key has to have a size of 256 bits. Derivation function creates a key based on user input and random salt (either generated as well or provided).
If this is a first time user, key is generated with salt. Salt has to be saved and used in key derivation for next session. Salt is not a secured piece of information; but it's is crucial to keep it persistent.  If salt is lot, the key is lost as well.

## Encryption/Decryption
Encrypt function uses recommended process of encrypting a message first and applying mac calculation after. HMAC has 256 bit and is appended to encrypted data.
Decrypt checks HMAC first and decrypts the data if the mac is valid.

## Example

```
    // first session
    key, salt, _ := gcrypt.DerivateKey256("password")
    // salt must be stored
    // returning session, salt is read from a storage
    // key, salt, _ := gcrypt.DerivateKey256WithSalt("password, salt")
    data := []byte("data")
    fmt.Println(data)
    ct, _ := gcrypt.Encrypt(key, data)
    fmt.Println(ct)
    plain, _ := gcrypt.Decrypt(key, ct)
    fmt.Println(plain)
```