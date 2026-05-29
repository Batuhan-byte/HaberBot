package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Topic struct {
	Name     string
	Slug     string
	Keywords string
	Sources  string
	Feeds    string
}

type Article struct {
	Title          string
	TurkishTitle   string
	OriginalURL    string
	SourceType     string
	OriginalContent string
	TurkishSummary string
	Score          int
	ImageURL       string
}

func main() {
	log.Println("Starting full database seeding process...")
	
	// Load environment variables
	godotenv.Load()
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer conn.Close(ctx)

	// 1. Define Topics
	topics := []Topic{
		{
			Name:     "Yapay Zeka",
			Slug:     "yapay-zeka",
			Keywords: `["AI", "LLM", "OpenAI", "ChatGPT", "Gemini", "Claude", "DeepMind", "Machine Learning"]`,
			Sources:  `["hackernews", "rss"]`,
			Feeds:    `["https://www.technologyreview.com/topic/artificial-intelligence/feed/", "https://www.wired.com/feed/tag/ai/latest/tgx", "https://www.artificialintelligence-news.com/feed/", "https://news.ycombinator.com/rss"]`,
		},
		{
			Name:     "Mobil",
			Slug:     "mobil",
			Keywords: `["iPhone", "Samsung", "Android", "iOS", "Xiaomi", "Smartphones", "Snapdragon", "Mobile App"]`,
			Sources:  `["hackernews", "rss"]`,
			Feeds:    `["https://www.phonearena.com/feed", "https://9to5mac.com/feed/", "https://9to5google.com/feed/", "https://news.ycombinator.com/rss"]`,
		},
		{
			Name:     "Uygulama/Yazılım",
			Slug:     "uygulama-yazilim",
			Keywords: `["React", "Go", "Docker", "DevOps", "Software", "Programming", "VS Code", "IDE", "Coding"]`,
			Sources:  `["hackernews", "rss"]`,
			Feeds:    `["https://dev.to/feed", "https://feed.infoq.com/", "https://www.smashingmagazine.com/feed/", "https://news.ycombinator.com/rss"]`,
		},
		{
			Name:     "Donanım",
			Slug:     "donanim",
			Keywords: `["RTX", "Nvidia", "Intel", "AMD", "Ryzen", "CPU", "GPU", "ASUS", "Hardware", "Overclock"]`,
			Sources:  `["hackernews", "rss"]`,
			Feeds:    `["https://www.tomshardware.com/feeds/all", "https://www.anandtech.com/rss/", "https://news.ycombinator.com/rss"]`,
		},
		{
			Name:     "Global Haberler",
			Slug:     "global-haberler",
			Keywords: `["SpaceX", "Tesla", "Starlink", "Apple", "Google", "Tech Industry", "Meta", "TSMC", "ASML"]`,
			Sources:  `["hackernews", "rss"]`,
			Feeds:    `["https://techcrunch.com/feed/", "https://www.theverge.com/rss/index.xml", "https://www.wired.com/feed/rss", "https://news.ycombinator.com/rss"]`,
		},
		{
			Name:     "Türkiye Haberleri",
			Slug:     "turkiye-haberleri",
			Keywords: `["Togg", "Baykar", "Teknofest", "Turkcell", "E-Devlet", "Yerli Yapay Zeka", "Hepsiburada", "Trendyol"]`,
			Sources:  `["rss"]`,
			Feeds:    `["https://webrazzi.com/feed/", "https://shiftdelete.net/feed", "https://www.webtekno.com/rss.xml", "https://www.donanimhaber.com/rss/tumhaberler.xml"]`,
		},
	}

	// Map topic slug -> topic ID
	topicIDs := make(map[string]string)

	for _, t := range topics {
		var topicID string
		// Check if topic exists
		err = conn.QueryRow(ctx, "SELECT id FROM topics WHERE slug = $1", t.Slug).Scan(&topicID)
		if err != nil {
			// Insert new topic
			topicID = uuid.New().String()
			_, err = conn.Exec(ctx, `
				INSERT INTO topics (id, name, slug, keywords, sources, rss_feeds, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				topicID, t.Name, t.Slug, t.Keywords, t.Sources, t.Feeds, true)
			if err != nil {
				log.Fatalf("Failed to insert topic %s: %v", t.Name, err)
			}
			log.Printf("Seeded new topic: %s", t.Name)
		} else {
			// Update topic info
			_, err = conn.Exec(ctx, `
				UPDATE topics 
				SET name = $1, keywords = $2, sources = $3, rss_feeds = $4 
				WHERE id = $5`,
				t.Name, t.Keywords, t.Sources, t.Feeds, topicID)
			if err != nil {
				log.Fatalf("Failed to update topic %s: %v", t.Name, err)
			}
			log.Printf("Updated topic: %s", t.Name)
		}
		topicIDs[t.Slug] = topicID
	}

	// 2. Define Sample Articles grouped by Topic
	articles := map[string][]Article{
		"yapay-zeka": {
			{
				Title:           "Google DeepMind introduces Gemini 1.5 Pro with 1M context",
				TurkishTitle:    "Google DeepMind, 1 Milyon Token Kapasiteli Gemini 1.5 Pro Modelini Duyurdu",
				OriginalURL:     "https://deepmind.google/technologies/gemini/1-5-pro-large",
				SourceType:      "rss",
				OriginalContent: "Google DeepMind has officially rolled out Gemini 1.5 Pro, featuring an unprecedented context window of up to 1 million tokens. This allows the model to process hours of video, massive audio files, and hundreds of thousands of lines of code in a single prompt.",
				TurkishSummary:  "Google DeepMind, yapay zekada bağlam sınırlarını ortadan kaldıran 1 milyon token kapasiteli Gemini 1.5 Pro'yu tanıttı. Yeni model, tek seferde 1 saatlik video, saatlerce ses kaydı veya 30.000 satırdan fazla kod bloğunu mükemmel bir şekilde analiz edebilmektedir.",
				Score:           2900,
				ImageURL:        "https://images.unsplash.com/photo-1677442136019-21780efad99a?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "OpenAI launches GPT-4o, an omni model for text, voice and vision",
				TurkishTitle:    "OpenAI, Ses ve Görüntü Destekli Yeni GPT-4o Modelini Tanıttı",
				OriginalURL:     "https://openai.com/blog/gpt-4o-realtime",
				SourceType:      "hackernews",
				OriginalContent: "OpenAI has announced GPT-4o, its newest flagship model that can reason across audio, vision, and text in real-time. GPT-4o responds to audio inputs in as little as 232 milliseconds, matching human response time in a conversation.",
				TurkishSummary:  "OpenAI, ses, metin ve görüntüyü aynı anda ve gerçek zamanlı olarak işleyebilen yeni amiral gemisi modeli GPT-4o'yu duyurdu. 232 milisaniye gibi rekor bir sürede sesli yanıt verebilen model, insan sohbet hızını yakalamış durumda.",
				Score:           2750,
				ImageURL:        "https://images.unsplash.com/photo-1620712943543-bcc4688e7485?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Anthropic releases Claude 3.5 Sonnet setting new industry benchmarks",
				TurkishTitle:    "Anthropic, Sektör Standartlarını Aşan Claude 3.5 Sonnet Modelini Yayınladı",
				OriginalURL:     "https://anthropic.com/claude/3-5-sonnet-release",
				SourceType:      "rss",
				OriginalContent: "Anthropic has released Claude 3.5 Sonnet, which sets new industry benchmarks for graduate-level reasoning, undergraduate-level knowledge, and coding proficiency. The model operates at twice the speed of Claude 3 Opus while costing significantly less.",
				TurkishSummary:  "Anthropic, akıl yürütme, lisans düzeyinde bilgi birikimi ve kodlama becerilerinde liderliği ele geçiren Claude 3.5 Sonnet modelini yayınladı. Model, Claude 3 Opus'a kıyasla iki kat daha hızlı çalışıyor.",
				Score:           2450,
				ImageURL:        "https://images.unsplash.com/photo-1546776310-eef45dd6d63c?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "AI chip maker NVIDIA hits $3 trillion market cap overtaking Apple",
				TurkishTitle:    "Yapay Zeka Devi NVIDIA, Apple'ı Geride Bırakarak 3 Trilyon Dolar Kulübüne Girdi",
				OriginalURL:     "https://nvidia.com/market-cap-3trillion",
				SourceType:      "rss",
				OriginalContent: "NVIDIA's stock surged once again, pushing the company's market capitalization past $3 trillion for the first time. The chipmaker has benefited immensely from the AI boom, supplying GPUs to major cloud service providers worldwide.",
				TurkishSummary:  "NVIDIA'nın hisseleri, yapay zeka çiplerine olan muazzam talep nedeniyle yükselişini sürdürerek şirketin değerini 3 trilyon doların üzerine çıkardı. NVIDIA bu hamleyle Apple'ı geçerek dünyanın en değerli ikinci şirketi oldu.",
				Score:           2100,
				ImageURL:        "https://images.unsplash.com/photo-1591453089816-0fbb971b454c?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Microsoft integrates advanced AI capabilities into Windows 11 Copilot",
				TurkishTitle:    "Microsoft, Windows 11 Copilot Sistemine Gelişmiş Yapay Zeka Özellikleri Ekledi",
				OriginalURL:     "https://microsoft.com/windows-11-copilot-update",
				SourceType:      "rss",
				OriginalContent: "Microsoft has announced a wave of new Copilot AI features built directly into Windows 11. This includes Recall, which helps users search their past actions on the PC, and real-time live translation for voice and system audio.",
				TurkishSummary:  "Microsoft, Windows 11 işletim sistemine derin entegrasyon sağlayan yeni yapay zeka özelliklerini duyurdu. Kullanıcının bilgisayardaki geçmişini aramasını sağlayan 'Recall' özelliği güncellemeler arasında öne çıkıyor.",
				Score:           1900,
				ImageURL:        "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?auto=format&fit=crop&w=800&q=80",
			},
		},
		"mobil": {
			{
				Title:           "Apple iPhone 16 Pro Max rumors point to super thin bezels",
				TurkishTitle:    "iPhone 16 Pro Max Sızıntıları: Dünyanın En İnce Ekran Çerçeveli Telefonu Geliyor",
				OriginalURL:     "https://macrumors.com/iphone-16-pro-bezels",
				SourceType:      "rss",
				OriginalContent: "New leaks surrounding the upcoming iPhone 16 Pro Max indicate Apple is using advanced Border Reduction Structure (BRS) technology to achieve the thinnest bezels ever seen on a smartphone, maximizing the display space.",
				TurkishSummary:  "Sızan yeni şemalara göre Apple, iPhone 16 Pro Max modelinde BRS adı verilen özel bir ekran teknolojisi kullanacak. Bu sayede cihaz, şimdiye kadar bir akıllı telefonda görülen en ince ekran çerçevelerine sahip olacak.",
				Score:           2800,
				ImageURL:        "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Samsung Galaxy S25 Ultra to feature complete design overhaul",
				TurkishTitle:    "Samsung Galaxy S25 Ultra Tasarımı Baştan Aşağı Değişiyor",
				OriginalURL:     "https://samsung.com/galaxy-s25-design-leaks",
				SourceType:      "rss",
				OriginalContent: "Samsung is planning a significant design overhaul for the Galaxy S25 Ultra, moving away from sharp rectangular corners towards slightly rounded edges. The company aims to improve ergonomics and hand comfort significantly.",
				TurkishSummary:  "Samsung, önümüzdeki yıl tanıtacağı Galaxy S25 Ultra modelinde köşeli tasarım dilini terk ederek hafif yuvarlatılmış kenarlara geçiş yapacak. Şirket bu hamleyle ergonomiyi ve tek elle kullanım konforunu artırmayı hedefliyor.",
				Score:           2400,
				ImageURL:        "https://images.unsplash.com/photo-1610945265064-0e34e5519bbf?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Google launches Android 15 Beta with enhanced privacy controls",
				TurkishTitle:    "Google, Gelişmiş Gizlilik Özelliklerine Sahip Android 15'i Duyurdu",
				OriginalURL:     "https://developer.android.com/android-15-beta",
				SourceType:      "rss",
				OriginalContent: "Google has officially released the first public beta of Android 15. The update introduces Private Space, which lets users hide sensitive apps in an isolated profile, and improved satellite connectivity support at OS level.",
				TurkishSummary:  "Google, mobil işletim sisteminin yeni sürümü Android 15'in halka açık ilk beta sürümünü yayınladı. Güncelleme, hassas uygulamaları gizlemeyi sağlayan 'Özel Alan' (Private Space) ve uydu bağlantısı desteği getiriyor.",
				Score:           2100,
				ImageURL:        "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Qualcomm Snapdragon 8 Gen 4 promises huge GPU performance leaps",
				TurkishTitle:    "Qualcomm Snapdragon 8 Gen 4 İşlemcisi Grafik Gücünde Sınırları Zorlayacak",
				OriginalURL:     "https://qualcomm.com/snapdragon-8-gen-4-leaks",
				SourceType:      "rss",
				OriginalContent: "Qualcomm's upcoming Snapdragon 8 Gen 4 processor, featuring custom Oryon CPU cores, is rumored to deliver a massive 40% performance jump in mobile GPU benchmarks, giving upcoming flagship phones unprecedented gaming power.",
				TurkishSummary:  "Qualcomm'un bu yılın sonlarında çıkaracağı Snapdragon 8 Gen 4 yonga seti, mobil grafik testlerinde selefine kıyasla %40 daha fazla performans vaat ediyor. Özel Oryon çekirdekleri, oyunlarda konsol kalitesinde grafik sunacak.",
				Score:           1800,
				ImageURL:        "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=800&q=80",
			},
		},
		"uygulama-yazilim": {
			{
				Title:           "React 19 RC is officially available: Compiler and Server Actions are here",
				TurkishTitle:    "React 19 Resmi Olarak Tanıtıldı: Otomatik Derleyici ve Sunucu Aksiyonları Geldi",
				OriginalURL:     "https://react.dev/blog/react-19-rc",
				SourceType:      "hackernews",
				OriginalContent: "The React team has officially released the Release Candidate for React 19. The highlights include the React Compiler, which eliminates the need for useMemo and useCallback, and built-in support for asynchronous Server Actions.",
				TurkishSummary:  "React ekibi, web geliştirme dünyasının en popüler kütüphanesinin yeni sürümü React 19'u resmi olarak duyurdu. Geliştiricileri 'useMemo' ve 'useCallback' kullanma yükünden kurtaran React Derleyici (React Compiler) en önemli yenilik.",
				Score:           2700,
				ImageURL:        "https://images.unsplash.com/photo-1633356122544-f134324a6cee?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Go 1.23 released with range over func and improved toolchains",
				TurkishTitle:    "Go 1.23 Sürümü Yayınlandı: Fonksiyonlarda Döngü Desteği ve Hızlı Derleme",
				OriginalURL:     "https://go.dev/blog/go1.23",
				SourceType:      "hackernews",
				OriginalContent: "The Go programming language has announced its version 1.23 release. The version introduces range over function iterators, allowing custom iteration patterns in loops, alongside standard library security enhancements.",
				TurkishSummary:  "Go programlama dili, 1.23 numaralı kararlı sürümünü yayınladı. Bu sürümle birlikte, döngülerde özel fonksiyon tabanlı yineleyicilerin (iterators) kullanılmasına imkan tanıyan dil güncellemeleri ve güvenlik yamaları eklendi.",
				Score:           2300,
				ImageURL:        "https://images.unsplash.com/photo-1607799279861-4dd421887fb3?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "VS Code introduces Native AI Copilot inside terminal and edit fields",
				TurkishTitle:    "VS Code Editörüne Entegre Terminal Copilot Yapay Zekası Eklendi",
				OriginalURL:     "https://code.visualstudio.com/updates/v1-90",
				SourceType:      "rss",
				OriginalContent: "Microsoft's latest Visual Studio Code update integrates GitHub Copilot directly into the terminal window. Users can now generate shell commands, troubleshoot build errors, and debug code natively from the console.",
				TurkishSummary:  "Popüler kod editörü VS Code, terminal alanına doğrudan yapay zeka entegrasyonu getirdi. Geliştiriciler artık konsol hatalarını çözmek veya komut oluşturmak için ayrı bir pencere açmadan Copilot'tan yararlanabilecek.",
				Score:           2100,
				ImageURL:        "https://images.unsplash.com/photo-1542831371-29b0f74f9713?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Docker Desktop update significantly boosts container launch speeds on Windows",
				TurkishTitle:    "Docker Windows Güncellemesi ile Konteyner Başlangıç Hızları 2 Kat Arttı",
				OriginalURL:     "https://docker.com/desktop-windows-performance",
				SourceType:      "rss",
				OriginalContent: "Docker has launched an update for Windows users using WSL2 backend. By optimizing memory reclamation and file sync operations, Docker Desktop now boots up and spins containers twice as fast as older versions.",
				TurkishSummary:  "Docker, Windows üzerinde WSL2 altyapısını kullanan geliştiriciler için büyük performans güncellemesi yayınladı. Dosya eşitleme ve bellek yönetimi iyileştirmeleriyle konteyner başlangıç süreleri %50 kısaltıldı.",
				Score:           1900,
				ImageURL:        "https://images.unsplash.com/photo-1587620962725-abab7fe55159?auto=format&fit=crop&w=800&q=80",
			},
		},
		"donanim": {
			{
				Title:           "NVIDIA GeForce RTX 5090 specs leaked: massive VRAM upgrade",
				TurkishTitle:    "NVIDIA GeForce RTX 5090 Sızıntıları: 32 GB GDDR7 VRAM ve Canavar Performans",
				OriginalURL:     "https://videocardz.com/rtx-5090-specs-leak",
				SourceType:      "rss",
				OriginalContent: "Leaked technical spec sheets for NVIDIA's upcoming flagship Blackwell architecture, GeForce RTX 5090, indicate a massive upgrade. The GPU will feature 32 GB GDDR7 VRAM on a 512-bit bus width, drawing up to 550W power.",
				TurkishSummary:  "NVIDIA'nın merakla beklenen yeni Blackwell mimarili amiral gemisi GeForce RTX 5090'ın teknik özellikleri sızdırıldı. Kartın 32 GB yeni nesil GDDR7 VRAM bellek ve 512-bit genişliğinde veri yolu barındıracağı belirtiliyor.",
				Score:           2850,
				ImageURL:        "https://images.unsplash.com/photo-1587202372775-e229f172b9d7?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Intel Core Ultra 9 285K benchmark scores reveal multi-threaded dominance",
				TurkishTitle:    "Intel Core Ultra 9 285K Test Sonuçları Sızdı: Çoklu Çekirdekte Büyük Fark",
				OriginalURL:     "https://wccftech.com/intel-core-ultra-9-285k-benchmarks",
				SourceType:      "rss",
				OriginalContent: "First engineering sample benchmark runs for Intel's upcoming Arrow Lake CPU, Core Ultra 9 285K, have surfaced. The results show impressive multi-threaded performance gains while significantly lowering operating temps.",
				TurkishSummary:  "Intel'in yeni Arrow Lake mimarisine dayanan Core Ultra 9 285K işlemcisinin ilk performans testleri ortaya çıktı. İşlemci, yüksek saat hızlarına ulaşırken güç tüketimini ve çalışma sıcaklıklarını düşürmeyi başarıyor.",
				Score:           2200,
				ImageURL:        "https://images.unsplash.com/photo-1591799264318-7e6ef8ddb7ea?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "AMD Ryzen 9000 series pricing details leaked online prior to retail launch",
				TurkishTitle:    "AMD Ryzen 9000 Serisi İşlemcilerin Fiyatları Perakende Satış Öncesi Sızdırıldı",
				OriginalURL:     "https://tomshardware.com/amd-ryzen-9000-pricing",
				SourceType:      "rss",
				OriginalContent: "Retailers have accidentally leaked the official MSRP for AMD's upcoming Zen 5 architecture, the Ryzen 9000 Series. The flagship Ryzen 9 9950X is priced competitively against Intel's current top tier desktop processors.",
				TurkishSummary:  "AMD'nin yeni Zen 5 işlemci ailesi Ryzen 9000 serisinin lansman öncesi perakende satış fiyatları sızdı. Tepe model Ryzen 9 9950X'in Intel'in fiyat politikasının altında kalarak agresif bir çıkış yapacağı görülüyor.",
				Score:           1950,
				ImageURL:        "https://images.unsplash.com/photo-1601524909162-be87252be298?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "ASUS launches massive 48-inch 4K OLED 240Hz gaming monitor",
				TurkishTitle:    "ASUS, Oyuncular İçin 48 İnç Boyutunda 240Hz 4K Dev OLED Monitör Çıkardı",
				OriginalURL:     "https://asus.com/rog-swift-oled-pg48uq",
				SourceType:      "rss",
				OriginalContent: "ASUS Republic of Gamers has rolled out the Swift OLED PG48UQ. Featuring a 48-inch 4K display, 240Hz refresh rate, 0.03ms response time, and a custom heatsink, this massive monitor aims at high-end console and PC gamers.",
				TurkishSummary:  "ASUS ROG, oyuncu monitörleri pazarında sınırları zorlayan Swift OLED modelini satışa sundu. 48 inçlik dev ekran, 4K çözünürlük, 240Hz tazeleme hızı ve 0.03ms inanılmaz tepki süresiyle profesyonel oyuncuları hedefliyor.",
				Score:           1750,
				ImageURL:        "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=800&q=80",
			},
		},
		"global-haberler": {
			{
				Title:           "SpaceX successfully launches Starship Flight 5 catching super heavy booster",
				TurkishTitle:    "SpaceX Starship 5. Uçuşunu Başarıyla Tamamladı: Dev Roket Kollarla Yakalandı!",
				OriginalURL:     "https://spacex.com/starship-flight-5-success",
				SourceType:      "rss",
				OriginalContent: "In a historical feat, SpaceX launched Starship Flight 5 and successfully caught the massive Super Heavy booster back at the launch tower using mechanical arms called Chopsticks, proving full rapid reusability is viable.",
				TurkishSummary:  "SpaceX, uzay havacılık tarihini baştan yazan bir başarıya imza atarak dev roket fırlatıcısı Super Heavy'i havada mekanik kollar vasıtasıyla yakalamayı başardı. Starship uzay aracı da hedeflenen rotaya sorunsuzca ulaştı.",
				Score:           2950,
				ImageURL:        "https://images.unsplash.com/photo-1541185933-ef5d8ed016c2?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "TSMC breaks ground on massive new 2nm chip fabrication facility",
				TurkishTitle:    "Çip Üretim Devi TSMC, 2 Nanometrelik Dev Fabrika İnşaatına Başladı",
				OriginalURL:     "https://tsmc.com/news-2nm-fab-construction",
				SourceType:      "rss",
				OriginalContent: "TSMC has broken ground on its highly anticipated new fabrication facility dedicated to advanced 2-nanometer nodes. The company plans to start commercial production by late 2026, supplying Apple, NVIDIA, and AMD.",
				TurkishSummary:  "Tayvan merkezli çip üreticisi TSMC, mobil ve masaüstü dünyasında yeni bir devir başlatacak 2 nanometrelik yarı iletken fabrikasının inşasına başladı. Tesisin 2026 sonunda Apple ve NVIDIA için üretime başlaması planlanıyor.",
				Score:           2150,
				ImageURL:        "https://images.unsplash.com/photo-1581092160607-ee22621dd758?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "ASML ships its first High-NA EUV lithography machine to Intel labs",
				TurkishTitle:    "Yarı İletken Lideri ASML, Yeni Nesil High-NA EUV Çip Makinesini Teslim Etti",
				OriginalURL:     "https://asml.com/intel-high-na-euv-delivery",
				SourceType:      "hackernews",
				OriginalContent: "ASML has successfully shipped its first revolutionary High-NA EUV lithography system, costing over $350 million, to Intel's research labs. The machine will print sub-2nm chip geometries for future generations of processors.",
				TurkishSummary:  "Hollandalı teknoloji devi ASML, dünyanın en gelişmiş çip üretim makinesi olan High-NA EUV sisteminin ilk sevkiyatını Intel tesislerine gerçekleştirdi. 350 milyon dolarlık bu makine, 2nm altındaki çipleri basabilecek.",
				Score:           2050,
				ImageURL:        "https://images.unsplash.com/photo-1562408590-e32931084e23?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "EU digital markets act forces Apple to open up iOS ecosystem completely",
				TurkishTitle:    "Avrupa Birliği Kuralları Apple'ı Dize Getirdi: iOS Ekosistemi Tamamen Açılıyor",
				OriginalURL:     "https://techcrunch.com/eu-dma-apple-ios-compliance",
				SourceType:      "rss",
				OriginalContent: "Under pressure from the EU's Digital Markets Act, Apple has announced historic changes to iOS, allowing European users to download alternative app marketplaces, choose default browsers, and access NFC chips natively.",
				TurkishSummary:  "Avrupa Birliği'nin Dijital Pazarlar Yasası (DMA) baskısıyla Apple, iOS işletim sisteminde tarihi değişikliklere gitti. Avrupa'daki kullanıcılar artık üçüncü parti uygulama mağazalarını kurabilecek ve NFC çipini serbestçe kullanabilecek.",
				Score:           1950,
				ImageURL:        "https://images.unsplash.com/photo-1563986768609-322da13575f3?auto=format&fit=crop&w=800&q=80",
			},
		},
		"turkiye-haberleri": {
			{
				Title:           "TOGG T10F sedan model public road tests started in Istanbul streets",
				TurkishTitle:    "Yerli Otomobil Togg'un Yeni Sedan Modeli T10F Yollara Çıktı: İstanbul'da Görüntülendi",
				OriginalURL:     "https://togg.com.tr/t10f-road-tests-istanbul",
				SourceType:      "rss",
				OriginalContent: "Turkey's national electric vehicle brand, Togg, has initiated public road and validation tests for its second model, the T10F fastback sedan. Istanbul commuters spotted the heavily camouflaged EV during range tests.",
				TurkishSummary:  "Türkiye'nin yerli elektrikli otomobil markası Togg'un yeni fastback sedan modeli T10F, yollara çıktı. Kamuflajlı tasarımıyla İstanbul sokaklarında menzil ve sürüş testleri yaparken görüntülenen araç büyük ilgi topladı.",
				Score:           2850,
				ImageURL:        "https://images.unsplash.com/photo-1563720223185-11003d516935?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Baykar completes successful test flights for Bayraktar TB3 carrier drone",
				TurkishTitle:    "Bayraktar TB3 SİHA, TCG Anadolu Gemisinden İlk Kalkış Testini Başarıyla Tamamladı",
				OriginalURL:     "https://baykartech.com/bayraktar-tb3-carrier-flight-success",
				SourceType:      "rss",
				OriginalContent: "Turkish defense industry leader Baykar announced that its naval-configured unmanned combat aerial vehicle, Bayraktar TB3, successfully completed take-off and landing tests simulation for short-runway aircraft carriers.",
				TurkishSummary:  "Savunma sanayiinde küresel bir marka olan Baykar, kısa pistli gemilere iniş-kalkış kabiliyetine sahip yeni SİHA'sı Bayraktar TB3'ün test uçuşlarını başarıyla sürdürüyor. TB3, TCG Anadolu gemisine konuşlandırılacak.",
				Score:           2600,
				ImageURL:        "https://images.unsplash.com/photo-1508614589041-895b88991e3e?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Teknofest 2026 application numbers break all previous global records",
				TurkishTitle:    "Teknofest 2026 Başvurularında Dünya Rekoru: 1.5 Milyondan Fazla Öğrenci Başvurdu",
				OriginalURL:     "https://teknofest.org/teknofest-2026-application-record",
				SourceType:      "rss",
				OriginalContent: "The world's largest aerospace and technology festival, Teknofest, has reported record-breaking application numbers for its 2026 competitions. Over 1.5 million young minds applied across high-tech drone, AI, and rocketry fields.",
				TurkishSummary:  "Dünyanın en büyük havacılık, uzay ve teknoloji festivali Teknofest 2026 için teknoloji yarışması başvuru sayıları açıklandı. 1.5 milyondan fazla genç araştırmacının başvurduğu festival, küresel alanda yeni bir rekor kırdı.",
				Score:           2250,
				ImageURL:        "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&w=800&q=80",
			},
			{
				Title:           "Turkey releases its first national large language model built for Turkish grammar",
				TurkishTitle:    "Türkiye'nin İlk Tamamen Yerli ve Türkçe Büyük Dil Modeli (LLM) Tanıtıldı",
				OriginalURL:     "https://tubitak.gov.tr/national-turkish-llm-launch",
				SourceType:      "rss",
				OriginalContent: "A collaboration between TÜBİTAK and local research institutes has led to the release of the first national large language model, optimized specifically for Turkish grammar, idioms, and cultural context.",
				TurkishSummary:  "TÜBİTAK öncülüğünde yerli mühendisler tarafından geliştirilen ve Türkçe'nin morfolojik yapısına, atasözlerine ve kültürel bağlamına göre optimize edilen ilk yerli Büyük Dil Modeli (LLM) resmi olarak kullanıma sunuldu.",
				Score:           2100,
				ImageURL:        "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=800&q=80",
			},
		},
	}

	// 3. Clear existing sample seeded articles for these topics to avoid duplication, then insert
	for slug, arts := range articles {
		topicID := topicIDs[slug]
		
		// Optional: delete existing articles for this topic to clean up
		_, err := conn.Exec(ctx, "DELETE FROM articles WHERE topic_id = $1", topicID)
		if err != nil {
			log.Printf("Warning: failed to clear old articles for %s: %v", slug, err)
		}

		for _, a := range arts {
			articleID := uuid.New().String()
			
			// Compute mock times
			createdTime := time.Now().Add(-time.Duration(3 + (a.Score % 48)) * time.Hour)
			
			_, err = conn.Exec(ctx, `
				INSERT INTO articles (
					id, title, turkish_title, original_url, source_type, 
					original_content, turkish_summary, score, topic_id, 
					image_url, processed_at, fetched_at, created_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
				)
				ON CONFLICT (original_url) DO NOTHING`,
				articleID,
				a.Title,
				a.TurkishTitle,
				a.OriginalURL,
				a.SourceType,
				a.OriginalContent,
				a.TurkishSummary,
				a.Score,
				topicID,
				a.ImageURL,
				time.Now(),
				createdTime,
				createdTime,
			)
			if err != nil {
				log.Fatalf("Failed to insert article %s: %v", a.Title, err)
			}
		}
		log.Printf("Seeded %d sample articles for category: %s", len(arts), slug)
	}

	log.Println("Database successfully seeded with all requested topics and mock articles!")
}
