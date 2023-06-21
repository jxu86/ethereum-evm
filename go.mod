module ethereum-evm

go 1.18

require (
	github.com/holiman/uint256 v1.2.2
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/tsdb v0.10.0
	golang.org/x/crypto v0.9.0
)

require github.com/go-logfmt/logfmt v0.5.1 // indirect

require (
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
)

require (
	github.com/cloudflare/cfssl v1.6.4
	github.com/ethereum/go-ethereum v1.12.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	golang.org/x/sys v0.8.0 // indirect
)

replace ethereum-evm/core v0.0.0 => ./core
