# HaberBot

Yapay zeka destekli Türkçe teknoloji haber kürasyon platformu.

HaberBot, teknoloji haberlerini HackerNews ve çeşitli RSS kaynaklarından (Yapay Zeka ve Yazılım Geliştirme odaklı) otomatik olarak derler. İngilizce makaleleri çeker, Google Gemini yapay zekasını kullanarak bu makaleleri çevirir, Türkçe "hap bilgi" (kısa ve öz özetler) haline getirir ve modern, premium bir React web arayüzünde sunar.

## Hızlı Başlangıç

### Ön Koşullar
- Go 1.22+
- Node.js 20+
- PostgreSQL veritabanı
- Google Gemini API Anahtarı

### 1. Klonlama ve Kurulum
```bash
git clone <repository-url>
cd haberbot
```

### 2. Backend Kurulumu
```bash
cd backend
cp .env.example .env
# .env dosyasını DATABASE_URL ve GEMINI_API_KEY ile düzenleyin
go run run_migrations.go  # Veritabanı migrasyonlarını çalıştır
go run cmd/server/main.go # API sunucusunu :8080 portunda başlat
```

### 3. Frontend Kurulumu
```bash
cd frontend
npm install
npm run dev # Vite geliştirme sunucusunu :5173 portunda başlat
```

## Komutlar

| Komut | Açıklama | Dizin |
|---------|-------------|-----------|
| `go run cmd/server/main.go` | Backend sunucusunu başlatır | `/backend` |
| `go run run_migrations.go` | Veritabanı migrasyonlarını çalıştırır | `/backend` |
| `npm run dev` | Frontend geliştirme sunucusunu başlatır | `/frontend` |
| `npm run build` | Frontend production build'ini alır | `/frontend` |

## Mimari

HaberBot, iş mantığını altyapı detaylarından ayırmak için **Clean Architecture** prensiplerini izler.
- **Domain:** Temel varlıklar (`Article`, `Topic`) ve arayüzler (interfaces).
- **UseCase:** Uygulamaya özgü iş kuralları (Veri çekme, Yapay zeka işleme).
- **Adapter:** Repository'ler, Ağ Geçitleri (Gemini) ve HTTP İşleyicileri (Handlers) için somut uygulamalar.
- **Infrastructure:** Framework'ler ve sürücüler (Fiber, PostgreSQL, Cron).

Detaylı teknik kararlar için [Mimari Karar Kayıtları (ADRs)](docs/decisions/) bölümüne bakın.
- [ADR-001: Go ve Fiber](docs/decisions/ADR-001.md)
- [ADR-002: PostgreSQL](docs/decisions/ADR-002.md)
- [ADR-003: Google Gemini Yapay Zekası](docs/decisions/ADR-003.md)
- [ADR-004: Clean Architecture](docs/decisions/ADR-004.md)

## Teknoloji Yığını
- **Backend:** Go, Fiber, pgx (PostgreSQL sürücüsü)
- **Frontend:** React, Vite, TailwindCSS (utility olarak), React Query
- **Veritabanı:** PostgreSQL (Neon)
- **Yapay Zeka:** Google Gemini (gemini-flash-latest)

## API Dokümantasyonu
Tam REST API referansı için [API.md](docs/api/API.md) dosyasına bakın.
