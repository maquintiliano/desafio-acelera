package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Answer struct {
	NumeroCasas int    `json:"numero_casas"`
	Token       string `json:"token"`
	Cifrado     string `json:"cifrado"`
	Decifrado   string `json:"decifrado"`
	Resumo      string `json:"resumo_criptografico"`
}

func main() {

	//declarando valor da url
	url := "https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token=92693f8e09cf0f268c1334ca86e0e4e3af8b155c"

	//fazendo a requisição http get
	resp, err := http.Get(url)

	if err != nil {

		fmt.Println(err.Error())

	}

	//ler dados da api
	lerDados, _ := ioutil.ReadAll(resp.Body)

	//criando o arquivo answer.json
	arquivo, err := os.OpenFile("answer", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}

	//abrindo arquivo
	jsonFile, err := os.Open("answer")

	if err != nil {
		fmt.Println(err.Error())
	}

	//escrevendo no arquivo answer.json
	arquivo.Write(lerDados)

	//fechando arquivo
	defer jsonFile.Close()

	//Aqui o arquivo é convertido para uma variável array de bytes, através do pacote "io/ioutil"
	byteValueJSON, _ := ioutil.ReadAll(jsonFile)

	//Declaração abreviada de um objeto do tipo Answer
	objAnswer := Answer{}

	//Conversão da variável byte em um objeto do tipo struct Answer
	json.Unmarshal(byteValueJSON, &objAnswer)

	//mudando dado do campo decifrado
	objAnswer.Decifrado = "prolific programmers contribute to certain disaster. niklaus wirth"

	//criando resumo criptografico
	data := []byte(objAnswer.Decifrado)
	resumo := sha1.Sum(data)

	objAnswer.Resumo = hex.EncodeToString(resumo[:])

	//Aqui vamos converter nosso objAnswer com o nome alterado em bytes
	byteValueJSON, err = json.Marshal(objAnswer)
	if err != nil {
		fmt.Println(err)
	}

	//Por fim vamos salvar em um arquivo JSON alterado.
	err = ioutil.WriteFile("answer", byteValueJSON, 0600)
	if err != nil {
		fmt.Println(err)
	}

	//definindo o valor da url de submit
	urlSubmit := "https://api.codenation.dev/v1/challenge/dev-ps/submit-solution?token=92693f8e09cf0f268c1334ca86e0e4e3af8b155c"

	//transformando a file em multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(jsonFile.Name()))
	io.Copy(part, jsonFile)
	writer.Close()

	//fazendo requisição post
	r, _ := http.NewRequest("POST", urlSubmit, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	client.Do(r)

}
