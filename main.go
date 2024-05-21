package main

import "github.com/cristovaoolegario/calculo-sucessao/internal/model"

func main() {
	heranca := model.NovaHeranca(100000.00, "Falecido", "2024-04-01") // Suponha uma herança de 1 milhão
	//heranca.AdicionarHerdeiro("Falecido", "2024-04-01")
	heranca.AdicionarHerdeiro("Conjuge", "0000-00-00")
	heranca.AdicionarHerdeiro("Filho", "0000-00-00")
	heranca.AdicionarHerdeiro("Filha", "0000-00-00")
	heranca.AdicionarHerdeiro("Neto", "0000-00-00")
	heranca.AdicionarHerdeiro("Neta", "0000-00-00")
	heranca.AdicionarHerdeiro("Bisneto", "0000-00-00")
	heranca.AdicionarHerdeiro("Tio", "0000-00-00")
	heranca.AdicionarHerdeiro("Primo", "0000-00-00")
	heranca.AdicionarHerdeiro("PrimoSegundoGrau", "0000-00-00")

	// Conexões do falecido para o cônjuge e filhos
	heranca.AdicionarConexao("Falecido", "Conjuge", "conjuge")
	heranca.AdicionarConexao("Falecido", "Filho", "filho")
	heranca.AdicionarConexao("Falecido", "Filha", "filho")

	// Conexões dos filhos para seus próprios filhos (netos do falecido)
	heranca.AdicionarConexao("Filho", "Neto", "neto")
	heranca.AdicionarConexao("Filha", "Neta", "neto")

	// Conexões para representar descendentes mais distantes
	heranca.AdicionarConexao("Neto", "Bisneto", "bisneto")

	// Conexões para representar parentes colaterais
	heranca.AdicionarConexao("Falecido", "Tio", "irmão")
	heranca.AdicionarConexao("Tio", "Primo", "filho")
	heranca.AdicionarConexao("Primo", "PrimoSegundoGrau", "filho")

	// Calcular herança
	heranca.CalcularHeranca()

	heranca.MostrarValoresHeranca()

	// Gerar o arquivo .dot
	heranca.GerarDot()
}
