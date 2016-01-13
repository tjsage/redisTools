package redis

type ScanResults struct {
	Iterator int
	Keys []string
}
