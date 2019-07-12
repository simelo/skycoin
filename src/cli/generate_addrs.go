package cli

import (
	"os"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func walletAddAddressesCmd() *gcli.Command {
	walletAddAddressesCmd := &gcli.Command{
		Use:   "walletAddAddresses",
		Short: "Generate additional addresses for a wallet",
		Long: fmt.Sprintf(`The default wallet (%s) will be used if no wallet was specified.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`, cliConfig.FullWalletPath()),
		SilenceUsage: true,
		RunE:         generateAddrs,
	}

	walletAddAddressesCmd.Flags().Uint64P("num", "n", 1, "Number of addresses to generate")
	walletAddAddressesCmd.Flags().StringP("wallet-file", "f", cliConfig.FullWalletPath(), "Generate addresses in the wallet")
	walletAddAddressesCmd.Flags().StringP("password", "p", "", "wallet password")
	walletAddAddressesCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format")

	return walletAddAddressesCmd
}

func generateAddrs(c *gcli.Command, _ []string) error {

	f, err2:= os.Create("/tmp/my_output")
	if err2 != nil{
		defer f.Close()
	}

	// get number of address that are need to be generated.
	f.WriteString("Checkpoint 1:\n")
	num, err := c.Flags().GetUint64("num")
	if err != nil {
		f.WriteString(err.Error())
		f.Close()
		return err
	}
	
	if num == 0 {
		return errors.New("-n must > 0")
	}
	
	f.WriteString("Checkpoint 2:\n")
	jsonFmt, err := c.Flags().GetBool("json")
	if err != nil {
		f.WriteString(err.Error())
		f.Close()
		return err
	}

	f.WriteString(cliConfig.FullWalletPath())
	f.WriteString("\n")
	f.WriteString(cliConfig.WalletDir)
	f.WriteString("\n")

	f.WriteString("Checkpoint 3:\n")
	w, err := resolveWalletPath(cliConfig, c.Flag("wallet-file").Value.String())
	f.WriteString(w)
	f.WriteString("\n")
	if err != nil {
		f.WriteString(err.Error())
		f.Close()
		return err
	}
	defer f.Close()
	
	f.WriteString("Checkpoint 4:\n")
	pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))
	
	addrs, err := GenerateAddressesInFile(w, num, pr)
	f.WriteString("asdfe")
	f.WriteString(err.Error())
	f.WriteString("\n")
	switch err.(type) {
	case nil:
		f.WriteString("1\n")
	case WalletLoadError:
		f.WriteString("2\n")
		printHelp(c)
		return err
	default:
		f.WriteString("3\n")
		return err
	}
	
	f.WriteString("Checkpoint 5:\n")
	if jsonFmt {
		s, err := FormatAddressesAsJSON(addrs)
		if err != nil {
			f.WriteString("Checkpoint 5:\n")
			f.WriteString(err.Error())
			f.WriteString("\n")
			return err
		}
		fmt.Println(s)
	} else {
		fmt.Println(FormatAddressesAsJoinedArray(addrs))
	}
	f.Close()
	return nil
}

// GenerateAddressesInFile generates addresses in given wallet file
func GenerateAddressesInFile(walletFile string, num uint64, pr PasswordReader) ([]cipher.Addresser, error) {
	
	f, err2:= os.OpenFile("/tmp/my_output", os.O_APPEND | os.O_WRONLY, 0600)
	if err2 != nil{
		defer f.Close()
	}
	defer f.Close()
	f.WriteString("Checkpoint 2.0\n")
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}
	
	f.WriteString("Checkpoint 2.1\n")
	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}
	f.WriteString("Checkpoint 2.2\n")
	genAddrsInWallet := func(w *wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
		return w.GenerateAddresses(n)
	}
	f.WriteString("Checkpoint 2.3\n")
	if wlt.IsEncrypted() {
		genAddrsInWallet = func(w *wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
			password, err := pr.Password()
			if err != nil {
				return nil, err
			}

			var addrs []cipher.Addresser
			if err := w.GuardUpdate(password, func(wlt *wallet.Wallet) error {
				var err error
				addrs, err = wlt.GenerateAddresses(n)
				return err
			}); err != nil {
				return nil, err
			}

			return addrs, nil
		}
	}
	f.WriteString("Checkpoint 2.4\n")
	addrs, err := genAddrsInWallet(wlt, num)
	f.WriteString(addrs)
	f.WriteString("\n")
	if err != nil {
		return nil, err
	}
	f.WriteString("Checkpoint 2.5\n")
	dir, err := filepath.Abs(filepath.Dir(walletFile))
	f.WriteString(dir)
	f.WriteString("\n")
	if err != nil {
		return nil, err
	}
	f.WriteString("Checkpoint 2.6\n")
	if err := wlt.Save(dir); err != nil {
		return nil, WalletSaveError{err}
	}
	f.WriteString("Checkpoint 2.7\n")
	return addrs, nil
}

// FormatAddressesAsJSON converts []cipher.Address to strings and formats the array into a standard JSON object wrapper
func FormatAddressesAsJSON(addrs []cipher.Addresser) (string, error) {
	d, err := formatJSON(struct {
		Addresses []string `json:"addresses"`
	}{
		Addresses: AddressesToStrings(addrs),
	})

	if err != nil {
		return "", err
	}

	return string(d), nil
}

// FormatAddressesAsJoinedArray converts []cipher.Address to strings and concatenates them with a comma
func FormatAddressesAsJoinedArray(addrs []cipher.Addresser) string {
	return strings.Join(AddressesToStrings(addrs), ",")
}

// AddressesToStrings converts []cipher.Address to []string
func AddressesToStrings(addrs []cipher.Addresser) []string {
	if addrs == nil {
		return nil
	}

	addrsStr := make([]string, len(addrs))
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}

	return addrsStr
}
