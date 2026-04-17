//go:build localdev

package authn

func InsecureSkipVerify() bool {
	return true
}
