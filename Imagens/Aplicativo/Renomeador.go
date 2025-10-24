package main

// Imports usados no codigo
import (
	"flag"       	 // Lida com argumentos (flags) de linha de comando, como -d e -p
	"fmt"        	 // Implementa formatação de entrada e saída (como fmt.Println e fmt.Printf)
	"os"         	 // Fornece funções para interagir com o sistema operacional (ler diretórios, checar arquivos, renomear)
	"path/filepath"  // Lida com caminhos de arquivos de forma segura, independente do sistema operacional
	"strings"        // Fornece funções para manipular strings (Como por exemplo, verificar e remover os prefixos dos arquivos renomeados)
	"time"       	 // Fornece funções para medir e exibir o tempo (usado para obter a data atual)
)

// Variáveis para armazenar os argumentos de linha de comando
var (
	dirPath string
	prefix  string
)

func init() {
	// 1. Configura os argumentos de linha de comando usando o pacote 'flag'
	flag.StringVar(&dirPath, "dir", "", "O caminho do diretório onde os arquivos serão renomeados.")
	flag.StringVar(&dirPath, "d", "", "O caminho do diretório (atalho para --dir).")

	flag.StringVar(&prefix, "prefix", "", "O prefixo a ser adicionado ao nome de cada arquivo.")
	flag.StringVar(&prefix, "p", "", "O prefixo a ser adicionado (atalho para --prefix).")
}

func main() {
	flag.Parse()

	// O VALOR DE 'prefix' FOI CARREGADO AQUI PELO flag.Parse()

	if dirPath == "" || prefix == "" {
		fmt.Println("Erro: O diretório (--dir) e o prefixo (--prefix) são obrigatórios.")
		flag.Usage()
		os.Exit(1)
	}

	// === Adiciona a data ao PREFIXO ===

	// 1. Obtém a data e hora atuais
	now := time.Now()

	// 2. Formata a data como AAAA-MM-DD
	dateString := now.Format("2006-01-02")

	// 3. Constrói o novo prefixo, usando a variável global 'prefix' lida pelo flag.Parse()
	finalPrefix := dateString + "-" + prefix // AQUI VOCÊ DEVE TER CERTEZA QUE 'prefix' É RECONHECIDA

	// =======================================

	fmt.Printf("Diretório alvo: %s\n", dirPath)
	fmt.Printf("Prefixo original (do usuário): %s\n", prefix)
	fmt.Printf("Prefixo final a aplicar: %s\n", finalPrefix)

	// Verifica se o diretório existe
	err := checkDirectory(dirPath)
	if err != nil {
		fmt.Printf("Erro no diretório: %v\n", err)
		os.Exit(1)
	}

	// Renomeia os arquivos, passando o NOVO prefixo
	err = renameFiles(dirPath, finalPrefix)
	if err != nil {
		fmt.Printf("Erro ao renomear arquivos: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nProcesso concluído com sucesso!")
}

// checkDirectory verifica se o caminho é um diretório e se ele existe.
func checkDirectory(path string) error {
	// A função os.Stat retorna informações sobre o arquivo/diretório
	fileInfo, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("o diretório '%s' não existe", path)
		}
		// Outro erro (ex: permissão negada)
		return fmt.Errorf("não foi possível acessar o diretório '%s': %w", path, err)
	}

	// Verifica se é realmente um diretório
	if !fileInfo.IsDir() {
		return fmt.Errorf("o caminho '%s' existe, mas não é um diretório", path)
	}

	return nil // Retorna nil se estiver tudo OK
}

func renameFiles(dir string, prefix string) error {

	// Nota: entries é uma fatia de os.DirEntry, que não é diretamente FileInfo, mas funciona
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("falha ao ler o diretório: %w", err)
	}

	renamedCount := 0

	// A Lógica de data deve ser calculada na main() e passada como 'prefix' aqui.
	// Por exemplo: "2025-10-23-NOVO-"

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		originalName := entry.Name()
		currentName := originalName

		// 1. Lógica de Limpeza: Verifica se o arquivo já tem o prefixo.

		// Se o nome atual começar com o prefixo que estamos prestes a aplicar,
		// precisamos remover a parte do prefixo.
		if strings.HasPrefix(currentName, prefix) {
			// Remove o prefixo antigo (que é o mesmo que o novo que será aplicado)
			// Exemplo: Se 'prefix' é "2025-10-23-NOVO-" e 'currentName' é "2025-10-23-NOVO-foto.jpg",
			// o resultado será "foto.jpg"
			currentName = strings.TrimPrefix(currentName, prefix)

			fmt.Printf("Aviso: Limpando prefixo antigo de '%s' -> '%s'\n", originalName, currentName)
		}

		// 2. Lógica de Aplicação: Adiciona o prefixo LIMPO ao nome base

		newName := prefix + currentName

		// Constrói os caminhos
		oldPath := filepath.Join(dir, originalName)
		newPath := filepath.Join(dir, newName)

		// Verifica se o nome novo é diferente do nome original antes de tentar renomear
		if originalName == newName {
			fmt.Printf("Ignorado: O arquivo '%s' já possui o prefixo.\n", originalName)
			continue
		}

		// === RENOMEAÇÃO REAL ===
		err = os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("ERRO: Falha ao renomear '%s': %v\n", originalName, err)
			continue
		}

		fmt.Printf("Sucesso: '%s' -> '%s'\n", originalName, newName)
		renamedCount++
	}

	if renamedCount == 0 {
		fmt.Println("\nNenhum arquivo novo foi renomeado (arquivos existentes ignorados).")
	} else {
		fmt.Printf("\nProcessamento concluído. %d arquivo(s) renomeado(s).\n", renamedCount)
	}

	return nil
}
