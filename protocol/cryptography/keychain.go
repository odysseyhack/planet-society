package cryptography

type Keychain struct {
	MainPublicKey       Key32
	MainPrivateKey      Key32
	StoragePublicKey    Key32
	StoragePrivateKey   Key32
	SignaturePublicKey  Key32
	SignaturePrivateKey Key64
}

func OneShotKeychain() (*Keychain, error) {
	signer, err := GenerateSigner()
	if err != nil {
		return nil, err
	}

	mainBox, err := NewOneShotBox()
	if err != nil {
		return nil, err
	}

	storageBox, err := NewOneShotBox()
	if err != nil {
		return nil, err
	}

	return &Keychain{MainPublicKey: mainBox.publicKey, MainPrivateKey: mainBox.privateKey,
		StoragePublicKey: storageBox.publicKey, StoragePrivateKey: storageBox.privateKey,
		SignaturePublicKey: signer.publicKey, SignaturePrivateKey: signer.privateKey}, nil
}
