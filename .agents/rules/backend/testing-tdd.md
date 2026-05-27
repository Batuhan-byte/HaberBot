---
trigger: on_change
glob: backend/**/*_test.go
description: HaberBot TDD, Birim Test ve Mocklama Kuralları. Go standard library ile table-driven test şablonları ve el yazımı port mock struct tasarımları.
---

# 🧪 TDD, Birim Test & Mocklama Kuralları (rules/backend/testing-tdd.md)

Bu kural dosyası, HaberBot backend katmanındaki kod doğruluğunu korumak, regresyonları engellemek ve yüksek birim test kapsamını (coverage) sürdürülebilir kılmak amacıyla tasarlanmıştır. Herhangi bir Go test dosyası (`*_test.go`) yazarken veya düzenlerken bu kurallara kesinlikle uymalısınız.

*Bu dosya `golang-pro` ve `test-driven-development` yetenekleri standartlarına göre yapılandırılmıştır.*

---

## 📋 1. Tablo Güdümlü Test Standartları (Table-Driven Tests)

Go dilinin en güçlü test pratiklerinden biri olan **Table-Driven Tests (Tablo Güdümlü Testler)** standart unit test yapımızdır. Tüm iş mantığı, doğrulama (validation) ve veri ayrıştırma (parsing) fonksiyonları bu kalıpla test edilmelidir.

### 🛠️ Table-Driven Test Şablonu:
```go
// internal/adapter/gateway/openai_processor_test.go
package gateway

import (
	"testing"
)

func TestExtractTranslation(t *testing.T) {
	// 1. Test senaryosu tablosunu tanımlıyoruz
	tests := []struct {
		name          string
		input         string
		expectedTitle string
		expectedBody  string
		expectError   bool
	}{
		{
			name:          "Başarılı Standart Çeviri Çıktısı",
			input:         "BASLIK: Harika Go Kitabı\nMETIN: Bu kitap Go dilini çok iyi anlatır.",
			expectedTitle: "Harika Go Kitabı",
			expectedBody:  "Bu kitap Go dilini çok iyi anlatır.",
			expectError:   false,
		},
		{
			name:          "Hatalı/Eksik Metin Durumu",
			input:         "BASLIK: Eksik Metin",
			expectedTitle: "",
			expectedBody:  "",
			expectError:   true,
		},
	}

	// 2. Tablodaki her senaryoyu döngüyle koşturuyoruz
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, body, err := ExtractTranslationHelper(tt.input) // Test edilecek fonksiyon

			// Hata beklentisi kontrolü
			if (err != nil) != tt.expectError {
				t.Errorf("ExtractTranslationHelper() hata = %v, beklenen hata beklentisi %v", err, tt.expectError)
				return
			}

			// Çıktı doğruluğu kontrolü
			if title != tt.expectedTitle {
				t.Errorf("ExtractTranslationHelper() başlık = %v, beklenen %v", title, tt.expectedTitle)
			}
			if body != tt.expectedBody {
				t.Errorf("ExtractTranslationHelper() metin = %v, beklenen %v", body, tt.expectedBody)
			}
		})
	}
}
```

---

## 🧱 2. El Yazımı Port Mock Standartları (Manual Mocks)

Clean Architecture bağımsızlığını ve test izolyasyonunu korumak amacıyla, UseCase testlerinde veritabanına veya harici API'ye gitmek yerine **domain port arayüzlerini taklit eden el yazımı mock yapılar** kullanırız. Bu sayede testler harici hiçbir ağ/dosya veya DB bağlantısına bağımlı kalmadan milisaniyeler içinde çalışır.

### 🛠️ El Yazımı Mock Arayüz ve Test Şablonu:
```go
// internal/usecase/summarize_article_test.go
package usecase

import (
	"context"
	"errors"
	"testing"
	"HaberBot/internal/domain/entity"
)

// 1. Arayüzü taklit eden mock yapıyı tanımlıyoruz (El yazımı Mock)
type MockArticleRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*entity.Article, error)
}

// Arayüz sözleşmesini implemente ediyoruz
func (m *MockArticleRepository) GetByID(ctx context.Context, id int64) (*entity.Article, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
// Portun diğer metodları da ihtiyaç durumunda buraya eklenir...

type MockAIProcessor struct {
	SummarizeFunc func(ctx context.Context, content string) (string, error)
}

func (m *MockAIProcessor) Summarize(ctx context.Context, content string) (string, error) {
	if m.SummarizeFunc != nil {
		return m.SummarizeFunc(ctx, content)
	}
	return "", nil
}

// 2. Mock yapılarla UseCase testi
func TestSummarizeUseCase_Execute(t *testing.T) {
	// Mock repo hazırlığı
	mockRepo := &MockArticleRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
			if id == 99 {
				return nil, errors.New("makale veritabanında yok")
			}
			return &entity.Article{
				ID:              id,
				OriginalContent: "Go is an open source programming language.",
			}, nil
		},
	}

	// Mock yapay zeka hazırlığı
	mockAI := &MockAIProcessor{
		SummarizeFunc: func(ctx context.Context, content string) (string, error) {
			return "Go açık kaynaklıdır.", nil
		},
	}

	// UseCase'e mock bağımlılıkları enjekte ediyoruz
	u := NewSummarizeUseCase(mockRepo, mockAI)

	t.Run("Başarılı Özet Oluşturma", func(t *testing.T) {
		res, err := u.Execute(context.Background(), 1)
		if err != nil {
			t.Fatalf("beklenmeyen hata: %v", err)
		}
		if res != "Go açık kaynaklıdır." {
			t.Errorf("beklenen özet 'Go açık kaynaklıdır.', alınan '%s'", res)
		}
	})

	t.Run("Veritabanı Bulunamadı Hatası", func(t *testing.T) {
		_, err := u.Execute(context.Background(), 99)
		if err == nil {
			t.Error("hata bekleniyordu ama alınmadı")
		}
	})
}
```

---

## 🚫 3. Test Anti-Patterns (Yapılmaması Gerekenler)

* **Birim Testlerde Dış API/Ağ Çağrısı Yapmayın:** Test koştururken (`go test ./...`) asla gerçek Gemini, Mistral veya internet üzerindeki herhangi bir web sitesine ağ isteği gönderilmemelidir. Bu tür harici entegrasyonlar tamamen mock yapılarıyla taklit edilmeli; gerçek çağrılar sadece bağımsız entegrasyon testlerinde yapılmalıdır.
* **Testlerde Global Durumları (State) Kirletmeyin:** Paralel koşan testlerin birbirini etkilemesini engellemek için test içinde ortak global değişkenleri, ortam değişkenlerini (`os.Setenv`) veya veritabanı bağlantılarını doğrudan modifiye etmeyin. Eğer modifiye etmeniz gerekiyorsa, her test çalıştıktan sonra `t.Cleanup(func() { ... })` ile eski durumunu geri yükleyin.
* **Assert İçin print() veya println() Kullanmayın:** Test hatalarını bildirmek için standart çıktı fonksiyonlarını kullanmayın. Her zaman `t.Errorf`, `t.Fatalf` veya `t.Fail()` metotlarını kullanarak hata durumlarını test koşum motoruna bildirin.
