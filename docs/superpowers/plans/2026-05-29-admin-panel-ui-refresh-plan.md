# Implementation Plan: Admin Panel UI Refresh

## Overview
AdminContent’teki görsel dil ve düzeni AdminDashboard ve AdminSources ekranlarına birebir taşımak için ortak AdminLayout ve paylaşılan stiller oluşturulacak. Veri akışı ve içerik yapısı değişmeyecek; yalnızca layout ve stil ortaklaştırması yapılacak.

## Architecture Decisions
- **Ortak AdminLayout**: Sidebar, header ve içerik iskeleti tek bir layout bileşeninde toplanacak; sayfalar `title`, `subtitle`, `actions`, `children` slotlarıyla içeriklerini sağlayacak.
- **Paylaşılan AdminStyles**: AdminContent’teki sınıflar/stil kalıpları ortaklaştırılacak; AdminDashboard ve AdminSources aynı tipografi, spacing ve renk sistemini kullanacak.

## Task List

### Phase 1: Foundation
#### Task 1: AdminLayout bileşenini oluştur
**Description:** AdminContent’teki sidebar + header + main iskeleti referans alınarak ortak bir AdminLayout bileşeni oluşturulacak. Aktif link vurgusu route bazlı olacak.

**Acceptance criteria:**
- [ ] AdminLayout `title`, `subtitle`, `actions`, `children` slotlarını destekler.
- [ ] Sidebar aktif route’a göre doğru vurgulanır.
- [ ] AdminLayout tek başına render edildiğinde AdminContent’in iskeletiyle aynı görsel dili sağlar.

**Verification:**
- [ ] Manual check: AdminLayout örnek sayfada render edilip sidebar/header görünümü AdminContent ile eşleşir.

**Dependencies:** None

**Files likely touched:**
- `frontend/src/components/Admin/AdminLayout.jsx` (new)
- `frontend/src/components/Admin/AdminSidebar.jsx` (new or refactor)
- `frontend/src/components/Admin/AdminHeader.jsx` (new or refactor)

**Estimated scope:** Medium (3-5 files)

#### Task 2: AdminStyles ortak stil katmanını çıkar
**Description:** AdminContent’te kullanılan ortak sınıf/stil kalıpları yeniden kullanılabilir hale getirilir.

**Acceptance criteria:**
- [ ] AdminContent’e özel sınıflar ortak bir yerde tanımlanır ve referanslanır.
- [ ] AdminDashboard ve AdminSources aynı stil tokenlarını kullanabilir.

**Verification:**
- [ ] Manual check: Yeni AdminStyles ile AdminContent görünümü bozulmaz.

**Dependencies:** Task 1

**Files likely touched:**
- `frontend/src/pages/AdminPage/AdminContent.jsx` (sınıf referansları)
- `frontend/src/pages/AdminPage/AdminContent.css` (varsa ortaklaştırma)
- `frontend/src/styles/admin.css` (new or shared)

**Estimated scope:** Medium (3-5 files)

### Checkpoint: Foundation
- [ ] AdminLayout render ediliyor ve AdminContent görünümü korunuyor.
- [ ] Stil ortaklaştırması sonrası AdminContent’te görsel regresyon yok.

### Phase 2: Core Features
#### Task 3: AdminDashboard’u AdminLayout’a taşı
**Description:** AdminDashboard içeriği AdminLayout içine alınacak ve AdminContent görsel dili birebir uygulanacak.

**Acceptance criteria:**
- [ ] AdminDashboard sayfası AdminLayout ile render edilir.
- [ ] Sidebar, header ve container görsel dili AdminContent ile eşleşir.
- [ ] AdminDashboard veri ve içerik yapısı değişmez.

**Verification:**
- [ ] Manual check: AdminDashboard’da sidebar/header/container ve buton stilleri AdminContent ile aynı görünür.

**Dependencies:** Task 1, Task 2

**Files likely touched:**
- `frontend/src/pages/AdminPage/AdminDashboard.jsx`
- `frontend/src/components/Admin/AdminLayout.jsx`

**Estimated scope:** Small (1-2 files)

#### Task 4: AdminSources’u AdminLayout’a taşı
**Description:** AdminSources içeriği AdminLayout içine alınacak ve AdminContent görsel dili birebir uygulanacak. “Add New Source” üst aksiyon alanına taşınacak.

**Acceptance criteria:**
- [ ] AdminSources sayfası AdminLayout ile render edilir.
- [ ] “Add New Source” aksiyonu layout `actions` slotunda görünür.
- [ ] AdminSources veri ve içerik yapısı değişmez.

**Verification:**
- [ ] Manual check: AdminSources’ta sidebar/header/container ve buton stilleri AdminContent ile aynı görünür.

**Dependencies:** Task 1, Task 2

**Files likely touched:**
- `frontend/src/pages/AdminPage/AdminSources.jsx`
- `frontend/src/components/Admin/AdminLayout.jsx`

**Estimated scope:** Small (1-2 files)

### Checkpoint: Core Features
- [ ] AdminDashboard ve AdminSources görsel olarak AdminContent ile tutarlı.
- [ ] Admin panelde tüm admin sayfaları aynı layoutu kullanıyor.

### Phase 3: Polish
#### Task 5: Son görsel hizalama ve küçük düzeltmeler
**Description:** Kenar boşlukları, hover/active state ve tipografi tutarlılığı son kez kontrol edilir.

**Acceptance criteria:**
- [ ] Sidebar, header, container, buton ve tipografi tutarlı.
- [ ] Görsel dil AdminContent ile birebir.

**Verification:**
- [ ] Manual check: AdminContent referansı ile yan yana kontrol.

**Dependencies:** Task 3, Task 4

**Files likely touched:**
- `frontend/src/pages/AdminPage/AdminDashboard.jsx`
- `frontend/src/pages/AdminPage/AdminSources.jsx`
- `frontend/src/styles/admin.css`

**Estimated scope:** Small (1-2 files)

### Checkpoint: Complete
- [ ] Tüm acceptance criteria karşılandı.
- [ ] Admin panel sayfaları AdminContent görsel diliyle tam uyumlu.

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| AdminContent sınıfları taşınırken görsel regresyon | Medium | Değişiklikleri küçük parçalara böl, her checkpoint sonrası manuel kontrol |
| Layout değişimi ile spacing kırılması | Low | AdminContent referansı ile birebir karşılaştırma |

## Open Questions
- None
