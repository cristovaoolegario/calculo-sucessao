package model

import (
	"testing"
)

func TestAdicionarHerdeiro(t *testing.T) {
	heranca := NovaHeranca(1000000.00, "Falecido", "2024-04-01")
	herdeiro := heranca.AdicionarHerdeiro("Filho", "0000-00-00")
	if herdeiro == nil {
		t.Error("Erro ao adicionar herdeiro")
	}
	if heranca.Herdeiros["Filho"] == nil {
		t.Error("Herdeiro não adicionado corretamente ao mapa de herdeiros")
	}
}

func TestAdicionarConexao(t *testing.T) {
	heranca := NovaHeranca(1000000.00, "Falecido", "2024-04-01")
	heranca.AdicionarHerdeiro("Filho", "0000-00-00")
	heranca.AdicionarHerdeiro("Pai", "0000-00-00")

	heranca.AdicionarConexao("Falecido", "Filho", "filho")
	heranca.AdicionarConexao("Falecido", "Pai", "pai")

	if len(heranca.Herdeiros["Falecido"].Conexoes) != 2 {
		t.Error("Erro ao adicionar conexões")
	}
}

func TestCalcularHerancaPrimeiroGrau(t *testing.T) {
	heranca := NovaHeranca(1000000.00, "Falecido", "2024-04-01")
	heranca.AdicionarHerdeiro("Conjuge", "0000-00-00")
	heranca.AdicionarHerdeiro("Filho1", "0000-00-00")
	heranca.AdicionarHerdeiro("Filho2", "0000-00-00")

	heranca.AdicionarConexao("Falecido", "Conjuge", "conjuge")
	heranca.AdicionarConexao("Falecido", "Filho1", "filho")
	heranca.AdicionarConexao("Falecido", "Filho2", "filho")

	heranca.CalcularHeranca()

	conjugeHeranca := heranca.Herdeiros["Conjuge"].ValorHeranca
	filho1Heranca := heranca.Herdeiros["Filho1"].ValorHeranca
	filho2Heranca := heranca.Herdeiros["Filho2"].ValorHeranca

	if conjugeHeranca != 500000 {
		t.Errorf("Erro ao calcular herança do cônjuge. Esperado: 500000, Recebido: %f", conjugeHeranca)
	}
	if filho1Heranca != 250000 {
		t.Errorf("Erro ao calcular herança do Filho1. Esperado: 250000, Recebido: %f", filho1Heranca)
	}
	if filho2Heranca != 250000 {
		t.Errorf("Erro ao calcular herança do Filho2. Esperado: 250000, Recebido: %f", filho2Heranca)
	}
}

func TestCalcularHerancaSegundoGrau(t *testing.T) {
	heranca := NovaHeranca(1000000, "Falecido", "2024-04-01")
	heranca.AdicionarHerdeiro("Pai", "0000-00-00")
	heranca.AdicionarHerdeiro("Mãe", "0000-00-00")
	heranca.AdicionarHerdeiro("Irmão", "0000-00-00")

	heranca.AdicionarConexao("Falecido", "Pai", "pai")
	heranca.AdicionarConexao("Falecido", "Mãe", "mãe")
	heranca.AdicionarConexao("Falecido", "Irmão", "irmão")

	heranca.CalcularHeranca()

	paiHeranca := heranca.Herdeiros["Pai"].ValorHeranca
	maeHeranca := heranca.Herdeiros["Mãe"].ValorHeranca
	irmaoHeranca := heranca.Herdeiros["Irmão"].ValorHeranca

	if paiHeranca != 333333.3333333333 {
		t.Errorf("Erro ao calcular herança do Pai. Esperado: 333333.333333, Recebido: %f", paiHeranca)
	}
	if maeHeranca != 333333.3333333333 {
		t.Errorf("Erro ao calcular herança da Mãe. Esperado: 333333.333333, Recebido: %f", maeHeranca)
	}
	if irmaoHeranca != 333333.3333333333 {
		t.Errorf("Erro ao calcular herança do Irmão. Esperado: 333333.333333, Recebido: %f", irmaoHeranca)
	}
}
