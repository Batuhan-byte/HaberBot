# Tasarım Sistemi: HaberBot
**Proje ID:** HaberBot Web Arayüzü

## 1. Görsel Tema ve Atmosfer
HaberBot, premium, son teknoloji ürünü bir "Karanlık Mod (Dark Mode)" estetiğine sahiptir. Atmosfer şık, teknik ve dijital önceliklidir; modern mühendislik ve editoryal kalite hissi uyandırmak üzere tasarlanmıştır. Düz ve jenerik UI öğelerine dayanmadan derinlik yaratmak için ince vuruşlar (strokes) ve derin kontrast kullanır. Genel havası "Gece Yarısı Neon (Midnight Neon)" tarzıdır — profesyonel ancak görsel olarak çarpıcı, modern geliştirici odaklı platformlardan büyük ölçüde ilham alınmıştır.

## 2. Renk Paleti ve Rolleri

* **Zifiri Siyah Arka Plan** (`#000000`): Uygulamanın mutlak karanlık temeli, maksimum kontrast sağlar.
* **Derin Yüzey Grisi** (`#0d0d0d`): Kod blokları gibi ikincil yüzeylerde kullanılarak arka plandan hafifçe ayrılmalarını sağlar.
* **İnce Kenarlık Grisi** (`#222222`): Görsel gürültü yaratmadan yapısal sınırları tanımlamak için ince kenarlık vuruşlarında kullanılır.
* **Vurgu (Hover) Durumu Grisi** (`#444444`): Kullanıcı etkileşimli kartların veya öğelerin üzerine geldiğinde etkileşimi belirtmek için kullanılır.
* **Saf Beyaz Metin** (`#ffffff`): Birincil tipografi, başlıklar ve yüksek vurgulu veriler için kullanılır.
* **Sönük Metin Grisi** (`#e5e7eb`): Rahat okuma kontrastı için gövde kopyası ve paragraf metinlerinde kullanılır.
* **Soluk Bağlam Grisi** (`#9ca3af`): İkincil metinler, zaman damgaları (timestamps) ve alıntı (blockquote) stilleri için kullanılır.
* **Elektrik Mavisi Vurgu** (`#3b82f6`): Birincil bağlantılar, etkileşimli göstergeler ve alıntı (blockquote) vurguları için kullanılır.
* **Hover Vurgu Mavisi** (`#60a5fa`): Bağlantıların vurgu (hover) durumları için kullanılır.
* **Neon Pembe Kod** (`#ff79c6`): Karanlık arka planda öne çıkması için özellikle satır içi (inline) kod parçacıklarında kullanılır.

## 3. Tipografi Kuralları
* **Başlıklar:** Kalın ve buyurgan. `h2`, kompakt ve modern bir his yaratmak için 1.75rem (28px) boyutunda, 700 font ağırlığına ve sıkı bir harf aralığına (`-0.025em`) sahiptir. `h3` ise 1.375rem (22px) ile biraz daha küçüktür.
* **Gövde Metni (Body Text):** Editoryal okuma akışı için tasarlanmıştır. Standart paragraflar 1.125rem (18px) boyutunu ve rahat bir okuma için cömert bir satır yüksekliğini (`1.8`) kullanır.
* **Kod Metni:** Teknik içerikler için eş aralıklı (Monospaced - `'JetBrains Mono', 'Fira Code', monospace`) kullanılır, kod bloklarında netlik sağlar.

## 4. Bileşen (Component) Stilleri
* **Kartlar/Konteynerler:** İnce, belli belirsiz kenarlıklara (`1px solid #222222`) sahip minimalist yapı. Kullanıcı etkileşime dokunsal geri bildirim sağlamak amacıyla belirgin bir kenarlık rengi değişimi (`#444444`) ile tepki verirler (`card-border-hover`).
* **Zengin Metin (Rich Text) İçeriği:** Son derece özenle derlenmiş HTML oluşturma. Alıntılar (Blockquotes), hafif mavi tonlu bir arka plana ve belirgin Elektrik Mavisi sol kenarlığa sahiptir. Tablolar, koyu başlıklar (`#111111`) ve satırlardaki ince hover efektleriyle yoğun bir şekilde stillendirilmiştir.
* **Kaydırma Çubukları (Scrollbars):** Karanlık estetikle kusursuz bir şekilde bütünleşen, koyu renkli bir tutamağa (`#222222`) sahip özel, ultra ince (4px) kaydırma çubukları.

## 5. Düzen (Layout) Prensipleri
* **Beyaz Boşluk (Whitespace):** Cömert ve kasıtlı kullanım. Öğelerin nefes almasına izin vermek için öğeler önemli bir dikey ritme (vertical rhythm) sahiptir (örn. paragraflar için `margin-bottom: 1.5rem`, başlıklar için `margin-top: 2.5rem`).
* **Yapısal Çizgiler:** Ayırma işlemi, ağır arka planlar veya alt gölgeler (drop shadows) yerine tipografi ölçeği ve ince kenarlıklarla sağlanır.
