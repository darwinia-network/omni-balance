# Omni Balance
Omni Balance is an intelligent liquidity management tool designed to achieve seamless interoperability between decentralized exchanges (DEX), centralized exchanges (CEX), and bridging services. It integrates automated cross-chain operations to help users effortlessly balance their liquidity pools while ensuring token balances across different addresses are replenished.

## Table of Contents
* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Configuration](#configuration)

### Features

* **Multi-address Monitoring**: Real-time monitoring of specific token balances in different addresses, ensuring continuous supervision.
* **Automatic Token Balance Adjustment**: Allows setting custom balance thresholds, and once below the threshold, it automatically adjusts through the best route, including transfers from DEX, CEX, and bridges.
* **Cross-chain Support**: Enables seamless cross-chain transactions, expanding liquidity options.
* **Address Isolation**: Allows separation of monitoring addresses and operational addresses.
* **Supports Most EVM-compatible Chains** (depending on the provider you choose).

#### Supported Providers

| name                                                                             | swap | cross-chain |
|----------------------------------------------------------------------------------|------|-------------|
| [UniSwap](https://uniswap.org/)                                                  | ✅    |             |
| [Gate.io](https://gate.io)                                                       | ✅    | ✅           |
| [Helix Bridge](https://helixbridge.app/)                                         |      | ✅           |
| [Darwinia Bridge](https://bridge.darwinia.network/)                              |      | ✅           |
| [OKX web3](https://www.okx.com/zh-hans/web3/build/docs/waas/dex-crosschain-swap) | ✅    | ✅           |
| [Router Nitro](https://app.routernitro.com/swap)                                 | ✅    | ✅           |
| [Li.FI](https://li.fi/)                                                          | ✅    | ✅           |


### Installation

#### Download Binary Files from GitHub Releases

1. Visit the [Releases](https://github.com/darwinia-network/omni-balance/releases) page.
2. Download the latest version of the binary file for your operating system.

#### Using Docker

1. Pull the Docker image: `docker pull ghcr.io/darwinia-network/omni-balance:latest`

#### Building from Source

1. Ensure Golang 1.22 or later is installed.
2. Clone the repository: `git clone https://github.com/darwinia-network/omni-balance && cd omni-balance`
3. Build the project: `go build ./cmd`

#### Using Golang Install

1. Run the command: `go install github.com/darwinia-network/omni-balance@latest`

### Usage

#### Supported Commands

```sh
USAGE:
   omni-balance [global options] command [command options] 

COMMANDS:
   gate_liquidity  Create an order for the liquidity of Gate.
   del_order       delete order by id
   version, v      show version
   list            list supported providers and docs
   tasks           list supported tasks
   example         create a example config file
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   
    --conf value, -c value              (default: "./config.yaml")        
          config file path
   
    --placeholder, -p                   (default: false)                  
          enable placeholder, you can use placeholder to replace private key, Example:
          Fill '{{privateKey}}' in config.yaml.Run with -p to enable placeholder, Example:
          SERVER_PORT=:8080
          /var/folders/b9/kgdbsh096b76234g3nhxg2qh0000gn/T/go-build3443181796/b001/exe/cmd
          -c ./config.yaml -pWaiting for 'waiting for placeholder...' log, send
          placeholder data according to the prompt.
   
    --port value                        (default: ":8080")                
          When the placeholder parameter is set to true, you can specify and set the
          listening address of the HTTP server that receives the placeholder.

   MISC

   
    --help, -h                          (default: false)                  
          show help
```
#### Configuration

##### Example

* [global_config.yaml](./example/configs/global_config.yaml): A generic example configuration for global settings.
* [gate_liquidity_config.yaml](./example/configs/gate_liquidity_config.yaml): An example configuration for providing liquidity to the Gate.io exchange.
* [general_config.yaml](./example/configs/general_config.yaml): An example configuration for standard addresses swapping liquidity between two chains.
* [helix_liquidity_config.yaml](./example/configs/helix_liquidity_config.yaml): An example configuration that, when checking address balances, includes both on-chain balances and unclaimed HelixBridge balances for evaluation.
* [operator_safe_config.yaml](./example/configs/operator_safe_config.yaml): An example configuration for scenarios where the operator address utilizes a multi-signature setup.

##### Configuration Reference
```
debug: true
apiKey: "999999999999"
# All utilized blockchain networks along with their corresponding tokens must be defined here.
chains:
  - id: 42161 # Chain ID
    name: arbitrum # Chain name
    nativeToken: "ETH" # Native token of the chain; if not provided, defaults to the token with contract address set to zero.
    rpcEndpoints: # List of RPC endpoints; multiple entries can be specified, and requests will cycle through these.
      - https://arb1.arbitrum.io/rpc
      - https://arbitrum.llamarpc.com
      - https://arb-mainnet-public.unifra.io
      - https://arbitrum.rpc.subquery.network/public
      - https://1rpc.io/arb
    tokens: # All required tokens must be configured.
      - name: ETH # Token name
        contractAddress: "0x0000000000000000000000000000000000000000" # Token contract address
        decimals: 18 # Token decimal places; must be specified.
      - name: USDT
        contractAddress: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"
        decimals: 6
      - name: RING
        contractAddress: "0x9e523234D36973f9e38642886197D023C88e307e"
        decimals: 18
  - id: 137
    name: polygon
    nativeToken: "MATIC"
    rpcEndpoints:
      - https://polygon-rpc.com
    tokens:
      - name: MATIC
        contractAddress: "0x0000000000000000000000000000000000000000"
        decimals: 18
      - name: USDT
        contractAddress: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
        decimals: 6
  - id: 1
    name: ethereum
    nativeToken: "ETH"
    rpcEndpoints:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth
      - https://1rpc.io/eth
      - https://ethereum-rpc.publicnode.com
    tokens:
      - name: ETH
        contractAddress: "0x0000000000000000000000000000000000000000"
        decimals: 18
      - name: RING
        contractAddress: 0x9469D013805bFfB7D3DEBe5E7839237e535ec483
        decimals: 18
  - id: 534352
    name: scroll
    nativeToken: "ETH"
    rpcEndpoints:
      - https://rpc.scroll.io
      - https://scroll.drpc.org
      - https://1rpc.io/scroll
    tokens:
      - name: ETH
        contractAddress: "0x0000000000000000000000000000000000000000"
        decimals: 18
      - name: USDT
        contractAddress: 0xf55BEC9cafDbE8730f096Aa55dad6D22d44099Df
        decimals: 6

# During rebalancing, tokens from the following configuration will be selected to determine the optimal route for swaps.
sourceTokens:
  - name: ETH # Token name, must match the token name in the 'chains' section.
    chains: # On which chains this token can be used as a source token.
      - arbitrum
      - scroll
      - ethereum
  - name: USDT
    chains:
      - arbitrum
      - scroll
      - ethereum
  - name: RING
    chains:
      - arbitrum
      - ethereum

# Providers that need to be enabled.
providers:
  - type: CEX
    name: gate.io
    config:
      key: <gate_api_key>
      secret: <gate_api_secret>
  - type: DEX
    name: uniswap
  - type: Bridge
    name: helixbridge
    config: {}
  - type: Bridge
    name: darwinia-bridge
  - type: Bridge
    liquidityName: router_nitro
    config: {}
  - type: Bridge
    liquidityName: okx
    config:
      apiKey: "xxxx"
      secretKey: "xxxx"
      passphrase: "xxxxx"
      project: "xxxxx"

# Database storage settings.
db:
  type: SQLite # Database type; options include: MySQL, PostgreSQL, SQLite.
  SQLITE: # SQLite configuration.
    path: ./omni-balance.db

# Task execution intervals.
taskInterval:
  crossChain: 30m # Interval for processing cross-chain order tasks.
  getTokenPriceInUsdt: 30m # Interval for fetching token price tasks.
  bot: 2m # Interval for running bot tasks, primarily for creating orders.
  market: 2m # Interval for executing order tasks.


# Wallets configuration for monitoring purposes where direct operations on these addresses are not allowed.
# All token operations like withdrawals or deposits will be executed through the associated operator address.
wallets:
  - address: 0x888888c278443CFC4D6E60b1964c2F1BAc3Ef257
    # Private key for the wallet associated with the operator address for transaction signing.
    privateKey: <0x9273283412f0A26C2cB99BBD874D54AD18540101_privateKey>
    tokens:
      # Configuration for triggering the gate_liquidity bot when the balance of a specific token at this address meets certain conditions.
      - name: USDT
        # Current amount of the token at the address.
        amount: "1.5"
        # Threshold amount; when the balance is less than or equal to this value, the gate_liquidity bot is activated.
        threshold: "1"
        # List of chains where the token balance is to be monitored.
        chains:
          - arbitrum
          - polygon
```

#### Run
##### Using Docker
1. Run the Docker image:
```
docker run -d --name omni-balance -p 8080:8080 \
    -v <you-host-path>:/data \
    --restart=always \
    --name omni-balance \
     omni-balance:latest -c /data/config.yaml -p
```
2. Check logs: `docker logs omni-balance`

##### Using command
1. Run the command: `omni-balance -c <your-config-path> -p`

##### Replacing Placeholders (Optional)
**For security reasons, if you do not want to write private keys in plain text in the configuration file, you can replace placeholders via a POST request**

1. Run the command:
```shell
curl -X POST http://localhost:8080 \
  -d '{"<wallet1_private_key>":"xxxx", "<wallet2_private_key>":"xxxx"}'
```
2. Clear history
```shell
history -c
```
