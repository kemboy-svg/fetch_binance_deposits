package models

type Deposit struct {
    ID            string `json:"id"`
    Amount        string `json:"amount"`
    Coin          string `json:"coin"`
    Network       string `json:"network"`
    Status        int    `json:"status"`
    Address       string `json:"address"`
    AddressTag    string `json:"addressTag"`
    TxID          string `json:"txId"`
    InsertTime    int64  `json:"insertTime"`
    TransferType  int    `json:"transferType"`
    ConfirmTimes  string `json:"confirmTimes"`
    UnlockConfirm int    `json:"unlockConfirm"`
    WalletType    int    `json:"walletType"`
}

func (Deposit) TableName() string {
	return "deposits"
}
