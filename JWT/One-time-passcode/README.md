Sample code based on [Workvivo's APIs unwired users one time passcodes](https://developer.workvivo.com/#aa34c835-aefb-4ff4-b1ad-232d00d37a9a)

To run the sample code you do not need to create a new key pair, they have already been created. However you will need to create your own key pair for your Workvivo environments. 

*Never share your private keys, and always create new key pairs for each of your Workvivo organisations.*

Create a public / private key pair

* Create a key pair, this creates a file ``private.pem``
```
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096 > /dev/null 2>&1
```

* Extract the public key, this creates a file: ``public.pem``

```
openssl rsa -in private.pem -pubout -out public.pem
```

* Create JWKS from public.pem

```
python3 create_jwks.py >> jwks.json
```