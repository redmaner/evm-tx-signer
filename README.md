# EVM transaction signer 
A simple binary to sign raw unsigned transactions for Ethereum Virtual Machine (EVM) based blockchains. 

## Installation 
Build binary, requires `go` and `make`
````
make build
````

## Usage 
```
evm-tx-signer [options] <unsigned_tx_string>

	evm-tx-signer signs an usinged transaction using supplied private key for EVM based blockchains like Ethereum and Polygon 

Options:
    --privatekey		Hex encoded private key (required)
    --chainid			The chain ID to use for signing transaction (required)
    --signer			The signer to use [london|legacy] (london by default)
```

Example for Polygon mainnet:
```
./bin/evm-tx-signer --privatekey <private_key_here> --chainid 137 <raw_unsigned_tx_here>
```