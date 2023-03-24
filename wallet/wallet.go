package wallet

// Wallet 钱包接口
type Wallet interface {
	ChainId() int

	// 符号
	Symbol() string

	DeriveAddress() string
	// 派生公钥
	DerivePublicKey() string

	// 派生私钥
	DerivePrivateKey() string
}
