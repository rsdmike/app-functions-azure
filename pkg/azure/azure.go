package azure


type AzureAccountInfo struct{
	AccountName string
	AccountKey string
}
func NewAzureAccountInfo(accountName string, accountKey string) AzureAccountInfo {
	return AzureAccountInfo{
		AccountName: accountName,
		AccountKey: accountKey,
	}
}