package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const (
	// Letras minúsculas
	charsetLowercase = "abcdefghijklmnopqrstuvwxyz"
	// Letras maiúsculas
	charsetUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Dígitos numéricos
	charsetNumbers = "0123456789"
	// Símbolos
	charsetSymbols = "!@#$%&*-+=?"
)

func generatePassword(length int, charSet string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("o comprimento da senha deve ser maior que zero")
	}

	if charSet == "" {
		return "", fmt.Errorf("nenhum tipo de caractere selecionado")
	}

	var password strings.Builder

	charSetLength := big.NewInt(int64(len(charSet)))

	for i := 0; i < length; i++ {
		// Gera um índice aleatório e criptograficamente seguro dentro do
		// intervalo de 0 até len(charSet) - 1.
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", fmt.Errorf("erro ao gerar índice aleatório: %w", err)
		}

		// Adiciona o caractere correspondente ao índice aleatório
		password.WriteByte(charSet[randomIndex.Int64()])
	}

	return password.String(), nil
}

func main() {
	// 1. Recebimento de Argumentos (Pacote flag)

	// Define o argumento para o comprimento da senha (default: 12)
	length := flag.Int("length", 12, "Comprimento da senha a ser gerada.")

	// Define as flags booleanas para os tipos de caracteres (default: todos true)
	useLower := flag.Bool("lower", true, "Incluir letras minúsculas (a-z).")
	useUpper := flag.Bool("upper", true, "Incluir letras maiúsculas (A-Z).")
	useNumber := flag.Bool("number", true, "Incluir números (0-9).")
	useSymbol := flag.Bool("symbol", true, "Incluir símbolos (!@#...).")

	// Analisa os argumentos da linha de comando
	flag.Parse()

	// 2. Construção do Conjunto de Caracteres

	var charSet string

	if *useLower {
		charSet += charsetLowercase
	}
	if *useUpper {
		charSet += charsetUppercase
	}
	if *useNumber {
		charSet += charsetNumbers
	}
	if *useSymbol {
		charSet += charsetSymbols
	}

	// Verifica se pelo menos um conjunto foi selecionado
	if charSet == "" {
		fmt.Fprintln(os.Stderr, "Erro: Pelo menos um tipo de caractere deve ser selecionado (e.g., --lower, --upper, --number, --symbol).")
		// Exibe a ajuda (opcional)
		flag.Usage()
		os.Exit(1)
	}

	// 3. Geração da Senha

	// Usa a função generatePassword com o conjunto de caracteres construído
	password, err := generatePassword(*length, charSet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gerar senha: %v\n", err)
		os.Exit(1)
	}

	// 4. Exibição do Resultado

	fmt.Println(password)
}
