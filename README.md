# PCI Vault PCD Form Golang-Vue example

This is a working example on how to use the PCI Vault PCD Form
with a Vue frontend and Golang backend.

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
the `cd` to the backend directory and run `go run main.go`.

In a separate terminal `cd` to the frontend directory 
and run `npm install` followed by `npm run dev`.
