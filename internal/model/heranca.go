package model

import (
	"fmt"
	"os"
	"time"
)

/*
OBS: Relações de mesmo nível = irmãos e meio irmãos
Relações de mesmo nível = irmãos adotados ou de sangue tem os mesmos direitos
*/
type Herdeiro struct {
	Nome         string
	ValorHeranca float64
	DataObito    time.Time
	// Grau de parentesco - Descendente, ascendente, colateral em linha reta e colateral em linha transversal
	Conexoes []Conexao
}

type Conexao struct {
	Para string
	// Regime de união = universal, separação parcial, comunhão de bens, regime especial, união estável não declarada
	TipoRelacao string // Ex: "filho", "conjuge"
}

type Heranca struct {
	ValorTotal float64
	Falecido   string
	DataObito  time.Time
	Herdeiros  map[string]*Herdeiro
	/*
		Herdeiros necessários em ordem - descendentes, conjuge em caso de meação, ascendentes, colaterais até quarto grau
		OBS: Parentes mais próximos excluem os mais distantes
	*/
}

func NovaHeranca(valor float64, falecido, dataObito string) *Heranca {
	data, _ := time.Parse("2006-01-02", dataObito)
	return &Heranca{
		ValorTotal: valor,
		Falecido:   falecido,
		DataObito:  data,
		Herdeiros: map[string]*Herdeiro{
			falecido: {Nome: falecido, DataObito: data, Conexoes: []Conexao{}},
		},
	}
}

func (h *Heranca) AdicionarHerdeiro(nome string, dataObito string) *Herdeiro {
	data, _ := time.Parse("2006-01-02", dataObito)
	herdeiro := &Herdeiro{Nome: nome, DataObito: data, Conexoes: []Conexao{}}
	h.Herdeiros[nome] = herdeiro
	return herdeiro
}

func (h *Heranca) AdicionarConexao(de string, para string, tipoRelacao string) {
	if _, existe := h.Herdeiros[de]; !existe {
		fmt.Printf("Herdeiro '%s' não existe.\n", de)
		return
	}
	if _, existe := h.Herdeiros[para]; !existe {
		fmt.Printf("Herdeiro '%s' não existe.\n", para)
		return
	}

	h.Herdeiros[de].Conexoes = append(h.Herdeiros[de].Conexoes, Conexao{
		Para:        para,
		TipoRelacao: tipoRelacao,
	})
}

func (h *Heranca) GerarDot() {
	file, err := os.Create("arvore.dot")
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "node [shape=box];")

	// Nó invisível para união das linhas
	nomeInvisivel := "InvisibleUnion"
	fmt.Fprintf(file, "\"%s\" [label=\"\", style=invis, height=0.0, width=0.0];\n", nomeInvisivel)

	// Falecido e Conjuge no mesmo nível
	fmt.Fprintf(file, "{ rank=same; \"%s\"; \"%s\"; \"%s\";}\n", h.Falecido, "Conjuge", nomeInvisivel)
	// Falecido e Conjuge se conectam ao nó invisível
	fmt.Fprintf(file, "\"%s\" -> \"%s\" [dir=none];\n", h.Falecido, nomeInvisivel)
	fmt.Fprintf(file, "\"%s\" -> \"%s\" [dir=none];\n", nomeInvisivel, "Conjuge")
	// Nó invisível se conecta ao filho
	fmt.Fprintf(file, "\"%s\" -> \"Filho\";\n", nomeInvisivel)

	// Conexões dos herdeiros restantes
	for _, herdeiro := range h.Herdeiros {
		for _, conexao := range herdeiro.Conexoes {
			if conexao.Para != nomeInvisivel { // Ignorar conexões para o nó invisível
				fmt.Fprintf(file, "\"%s\" -> \"%s\" [label=\"%s\"];\n", herdeiro.Nome, conexao.Para, conexao.TipoRelacao)
			}
		}
	}

	fmt.Fprintln(file, "}")
}

func (h *Heranca) CalcularHeranca() {
	herdeirosDePrimeiroGrau := []string{}
	herdeirosDeSegundoGrau := []string{}
	conjuge := ""

	for _, conexao := range h.Herdeiros[h.Falecido].Conexoes {
		if conexao.TipoRelacao == "filho" {
			if herdeiro, existe := h.Herdeiros[conexao.Para]; existe && (herdeiro.DataObito.IsZero() || herdeiro.DataObito.After(h.DataObito)) {
				herdeirosDePrimeiroGrau = append(herdeirosDePrimeiroGrau, conexao.Para)
			}
		} else if conexao.TipoRelacao == "conjuge" {
			if herdeiro, existe := h.Herdeiros[conexao.Para]; existe && (herdeiro.DataObito.IsZero() || herdeiro.DataObito.After(h.DataObito)) {
				conjuge = conexao.Para
			}
		} else if conexao.TipoRelacao == "pai" || conexao.TipoRelacao == "mãe" || conexao.TipoRelacao == "irmão" {
			if herdeiro, existe := h.Herdeiros[conexao.Para]; existe && (herdeiro.DataObito.IsZero() || herdeiro.DataObito.After(h.DataObito)) {
				herdeirosDeSegundoGrau = append(herdeirosDeSegundoGrau, conexao.Para)
			}
		}
	}

	if conjuge != "" {
		parteConjuge := h.ValorTotal / 2
		h.Herdeiros[conjuge].ValorHeranca = parteConjuge

		if len(herdeirosDePrimeiroGrau) > 0 {
			parteFilhos := h.ValorTotal / 2 / float64(len(herdeirosDePrimeiroGrau))
			for _, nome := range herdeirosDePrimeiroGrau {
				h.Herdeiros[nome].ValorHeranca = parteFilhos
			}
		}
	} else if len(herdeirosDePrimeiroGrau) > 0 {
		parte := h.ValorTotal / float64(len(herdeirosDePrimeiroGrau))
		for _, nome := range herdeirosDePrimeiroGrau {
			h.Herdeiros[nome].ValorHeranca = parte
		}
	} else if len(herdeirosDeSegundoGrau) > 0 {
		parte := h.ValorTotal / float64(len(herdeirosDeSegundoGrau))
		for _, nome := range herdeirosDeSegundoGrau {
			h.Herdeiros[nome].ValorHeranca = parte
		}
	}
}

func (h *Heranca) MostrarValoresHeranca() {
	for nome, herdeiro := range h.Herdeiros {
		fmt.Printf("Herdeiro: %s, Valor da Herança: R$ %.2f\n", nome, herdeiro.ValorHeranca)
	}
}
