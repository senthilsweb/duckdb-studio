package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	vault "github.com/mch1307/vaultlib"
)

// InitSecrects is
func InitSecrects(fiile string) (string, error) {
	token := ""
	data, err := ioutil.ReadFile(fiile)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	token = string(data)
	//fmt.Println(token)
	vcConf := vault.NewConfig()
	token = strings.Replace(token, "Environment=VAULT_TOKEN=", "", -1)
	token = strings.Replace(token, "[Service]", "", -1)
	token = strings.Replace(token, "\n", "", -1)
	token = strings.Replace(token, "\n", "", -1)
	vcConf.Token = token
	//fmt.Println(token)
	// Create new client
	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal("\n Cannot connect to vault ", err)
		//responses.JSON(w, http.StatusInternalServerError, secret)
		return "", err
	}
	// Get the Vault secret kv_v2/path/my-secret
	secret, err := vaultCli.GetSecret("secret/kv/postgres")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	/*fmt.Printf("My Secret is %v\n", secret.KV["vault.spring.datasource.password"])
	// We know our secret is a key/value so we use secret.KV
	// if our secret was of JSON type, we would have used secret.JSONSecret
	for k, v := range secret.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
		//log.Printf("Secret %v: %v\n", k, v)
	}*/
	DBURI := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", secret.KV["vault.spring.datasource.username"], "edm", secret.KV["vault.spring.datasource.password"])

	return DBURI, nil
}
