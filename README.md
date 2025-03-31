# Go Auth Backend

This is the backend for the Go Auth application. It is a RESTful API that is built using the Echo framework. The API is used to manage users and authentication.

## Get Started

```bash
$ air -c .air.toml
```

### Generate JWT key

Generate a private key
```bash
$ openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```

Generate a public key from the private key
```bash
$ openssl rsa -pubout -in private_key.pem -out public_key.pem
```
