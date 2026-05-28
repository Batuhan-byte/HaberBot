# Admin Panel UI Refresh Design

## Context
AdminContent ekranı yeni görsel dile geçirildi. AdminDashboard ve AdminSources ekranları aynı görsel sistemi kullanmıyor. İstenen: AdminContent’in görsel dili ve düzeni tüm admin paneline birebir taşınsın; içerik ve veri akışı değişmesin.

## Goals
- AdminDashboard ve AdminSources, AdminContent ile aynı görsel sistemi kullanır.
- Ortak sidebar, header ve içerik iskeleti tek bir layout üzerinden paylaşılır.
- Route ve mevcut veri akışları değişmeden korunur.

## Non-goals
- Yeni özellik eklemek, içerik yapısını değiştirmek.
- Yeni global hata/boş durum sistemi tasarlamak.
- API, veri modeli veya React Query akışlarını değiştirmek.

## Proposed Approach
### AdminLayout (Ortak İskelet)
- Sabit sidebar, header alanı, ana içerik alanı ve sayfa içi aksiyon slotlarını sağlar.
- `title`, `subtitle`, `actions`, `children` slotları ile sayfalar kendi içeriğini taşır.
- Sidebar aktif state route’a göre otomatik vurgulanır.

### Shared AdminStyles
- AdminContent’te kullanılan renk paleti, tipografi, spacing, border, shadow ve buton stilleri tek bir yerde yeniden kullanılabilir hale getirilir.
- AdminDashboard ve AdminSources bu stilleri doğrudan kullanır.

## Visual Consistency Requirements
- Sidebar: AdminContent ile aynı arka plan, border, padding, link hover/active stili.
- Header: aynı başlık fontu, alt açıklama, sağ üst aksiyon butonu stili.
- Container: aynı background, border ve radius.
- Butonlar: AdminContent’teki primary/secondary hover ve focus davranışı.
- Tablo/kart stil dili: AdminContent’e yakın gradient, divider ve text stili.

## Data Flow & Interaction
- AdminDashboard ve AdminSources mevcut state yönetimini korur.
- Sadece layout ve sınıf düzeyi değişiklik yapılır.
- “Hemen Haber Çek & İşle” AdminContent’teki `actions` slotunda kalır.
- “Add New Source” AdminSources’ta üst aksiyon alanına taşınır (layout slotu üzerinden).

## Error Handling
Mevcut sayfa içi hata mesajları/uyarılar korunur; yeni global error UI eklenmez.

## Testing & Verification
- AdminContent referans alınarak AdminDashboard ve AdminSources görsel hizası manuel doğrulanır.
- Sidebar, header, container, buton, tipografi ve spacing birebir kontrol edilir.

## Rollout
Tek aşamalı değişiklik: layout ve stil ortaklaştırma ile tamamlanır.
