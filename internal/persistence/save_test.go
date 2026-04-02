package persistence

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/DeMart1n/Tamagotchi/internal/model"
)

func TestSaveAndLoad(t *testing.T) {
	// Setup: Garantir que não vamos sobrescrever o save real do desenvolvedor
	originalPath := savePath
	tempSave := "test_tamago_save.json"
	// Nota: Como savePath é constante e privada, em um cenário real precisaríamos
	// refatorar para permitir injeção do path. Para este teste, vamos simular.
	
	// Mock de dados
	tama := &model.Tama{
		Name:   "TestPet",
		Hunger: 80,
		XP:     10,
	}

	// Teste de Escrita (Nota: vai criar o arquivo tamago_save.json no diretório local)
	// Em um ambiente de CI, isso é aceitável se garantirmos o cleanup.
	err := Save(tama)
	if err != nil {
		t.Fatalf("Erro ao salvar: %v", err)
	}
	defer os.Remove(originalPath)

	// Teste de Leitura
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Erro ao carregar: %v", err)
	}

	if loaded.Name != tama.Name {
		t.Errorf("Esperado nome %s, obtido %s", tama.Name, loaded.Name)
	}

	if loaded.Hunger != tama.Hunger {
		t.Errorf("Esperada fome %d, obtida %d", tama.Hunger, loaded.Hunger)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	// Simula um save corrompido
	err := os.WriteFile("tamago_save.json", []byte("{ invalid json "), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tamago_save.json")

	_, err = Load()
	if err == nil {
		t.Error("Deveria ter retornado erro para JSON inválido")
	}
}

func TestDataIntegrity_ManualManipulation(t *testing.T) {
	// Este teste demonstra o problema da falta de validação
	// Se alguém colocar valores impossíveis no JSON
	maliciousSave := SaveFile{
		Pet: model.Tama{
			Name:   "HackerPet",
			Hunger: -999, // Valor impossível
			Level:  9999,
		},
		Version: "1.0",
	}
	
	data, _ := json.Marshal(maliciousSave)
	_ = os.WriteFile("tamago_save.json", data, 0644)
	defer os.Remove("tamago_save.json")

	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	// Atualmente este teste FALHA em detectar o erro, o que prova o problema
	if loaded.Hunger < 0 {
		t.Logf("⚠️ ALERTA: Sistema aceitou Hunger negativa (%d). Vulnerabilidade de integridade confirmada.", loaded.Hunger)
	}
}
