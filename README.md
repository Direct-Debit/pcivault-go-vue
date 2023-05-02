# PCI Vault PCD Form Golang-Vue example

This is a working example on how to use the [PCI Vault API](https://api.pcivault.io) with a Vue frontend and Golang backend.

This project has code that:
- Incorporates the [PCI Vault PCD Form](https://api.pcivault.io/pcd/how-to-capture-and-tokenize-payment-card-data.html)
- Sends data to Stripe with a [proxy endpoint](https://api.pcivault.io/#/Proxy/post-proxy-post).


This example is for reference purposes only. Use it to copy code into your
project, but make sure to read and understand the code first. You will also
have to add your own error handling, this code handles the happy case only.

## Running
In order to run the example, set the following environment variables:
```sh
PCI_KEY=key-user
PCI_PASSPHRASE=key-passphrase
PCI_BASIC_AUTH=username:password
```
then `cd` to the backend directory and run `go run main.go`.

In a separate terminal `cd` to the frontend directory 
and run `npm install` followed by `npm run dev`.

### Configuration
These environment variable are also significant:
```sh
DEBUG_MODE=true # send proxy requests with debug flag enabled (default is false)
STRIPE_KEY=sk_your_stripe_key # Stripe key for testing Stripe integration
```

## More Information
- [PCI Vault](https://pcivault.io)
- [PCI Vault API](https://api.pcivault.io)
- [PCI Vault PCD Form](https://api.pcivault.io/pcd/how-to-capture-and-tokenize-payment-card-data.html)
