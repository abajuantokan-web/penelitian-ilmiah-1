p align="center">
  <img src="https://img.shields.io/badge/Platform-OpenPeo-F5A623?style=for-the-badge&labelColor=1a1a40" alt="OpenPeo"/>
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/Vue.js-3-4FC08D?style=for-the-badge&logo=vuedotjs&logoColor=white" alt="Vue 3"/>
  <img src="https://img.shields.io/badge/MySQL-XAMPP-4479A1?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL"/>
</p>

# 🏝️ OpenPeo — Marketplace Produk Tradisional NTT

**OpenPeo** adalah platform e-commerce marketplace khusus yang menghubungkan pengrajin dan produsen lokal **Nusa Tenggara Timur (NTT)** dengan pembeli di seluruh Indonesia melalui sistem **pre-order** yang adil dan transparan.

Platform ini telah dikembangkan secara komprehensif untuk mendukung sistem **Authentication (Multi-role & Registrasi Akun)**, **Real-time Admin Sales Aggregator Dashboard**, dan **Isolated Real-time WebSocket Chat Rooms** berbasis user ID.

> Nama "OpenPeo" terinspirasi dari bahasa lokal NTT — sebuah wadah terbuka untuk masyarakat NTT berkolaborasi dan berkembang bersama.

---

## 📑 Daftar Isi

- [Arsitektur Sistem](#-arsitektur-sistem)
- [Tech Stack — Backend](#-tech-stack--backend-engine)
- [Tech Stack — Frontend](#-tech-stack--frontend-architecture)
- [Tech Stack — Database](#-tech-stack--database)
- [Struktur Project](#-struktur-project)
- [API Endpoints](#-api-endpoints)
- [Prasyarat Instalasi](#-prasyarat-instalasi)
- [Langkah-Langkah Menjalankan Project](#-langkah-langkah-menjalankan-project)
- [Panduan Pengujian Fitur](#-panduan-pengujian-fitur)
- [Catatan Pengembangan](#-catatan-pengembangan)

---

## 🏗️ Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────────────┐
│                        BROWSER (Client)                         │
│                     http://localhost:5173                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │             FRONTEND (Vue.js 3 + Vue Router)            │   │
│   │                                                         │   │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │   │
│   │  │  LoginPage   │  │   ShopPage   │  │AdminDashboard│  │   │
│   │  │(Role Selector│  │ (Katalog &   │  │ (KPIs, SVG   │  │   │
│   │  │  Form Auth)  │  │ Pre-Order)   │  │Chart, Table) │  │   │
│   │  └──────────────┘  └──────────────┘  └──────────────┘  │   │
│   │                                                         │   │
│   │  ┌──────────────────────────────────────────────────┐   │   │
│   │  │       ChatBox (Dynamic WS Chat Component)        │   │   │
│   │  └──────────────────────────────────────────────────┘   │   │
│   └───────────┬──────────────────────────┬──────────────────┘   │
│               │ HTTP REST (fetch API)    │ WebSocket (ws://)    │
│               ▼                          ▼                      │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │             BACKEND (Golang + Gin Framework)            │   │
│   │                  http://localhost:8080                   │   │
│   │                                                         │   │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │   │
│   │  │ Auth/Admin   │  │   Product    │  │     Chat     │  │   │
│   │  │   Handler    │  │   Handler    │  │   Handler    │  │   │
│   │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  │   │
│   │         │                  │                  │          │   │
│   │         ▼                  ▼                  ▼          │   │
│   │  ┌──────────────────────────────────────────────────┐   │   │
│   │  │              GORM (ORM Layer)                     │   │   │
│   │  └──────────────────────┬───────────────────────────┘   │   │
│   └─────────────────────────┼───────────────────────────────┘   │
│                             │                                    │
│                             ▼                                    │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │              MySQL Database (XAMPP)                       │   │
│   │              Database: db_openpeo                         │   │
│   │                                                          │   │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │   │
│   │  │  users   │ │ products │ │  orders  │ │ messages  │   │   │
│   │  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │   │
│   └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔧 Tech Stack — Backend Engine

### 1. Go (Golang) v1.22+
* **Fungsi**: Runtime utama untuk seluruh logika server-side, concurrency websocket, dan integrasi data MySQL.
* **Mengapa dipilih?**: Performa yang andal mendekati C/C++, hemat konsumsi memori, serta native goroutines yang mempermudah koneksi WebSocket persisten secara real-time tanpa overhead resource yang besar.

### 2. Gin Framework (github.com/gin-gonic/gin)
* **Fungsi**: Framework web untuk menangani routing HTTP, HTTP status codes, parsing query parameters, body data binding, dan middleware CORS.
* **Kegunaan Spesifik**:
  * `POST /api/login` - Validasi kredensial pengguna ke database.
  * `POST /api/register` - Mendaftarkan akun customer baru secara dinamis ke database.
  * `GET /api/admin/sales-data` - Menyajikan data agregat dan transaksi riil untuk Admin Dashboard.
  * `GET /api/chat/contacts` - Membaca ID & Role pengguna yang aktif untuk memuat kontak chat dengan data jumlah pesan belum terbaca (unread badges).

### 3. GORM (gorm.io/gorm)
* **Fungsi**: Object-Relational Mapper untuk berinteraksi dengan database MySQL menggunakan struct Go secara aman (mencegah SQL Injection).
* **Kegunaan Spesifik**:
  * `AutoMigrate()` - Sinkronisasi otomatis struktur tabel berdasarkan pembaruan models.
  * `Preload()` - Memuat data relasional (misalnya data Customer dan Product ketika query log Order secara real-time).
  * `Aggregate Queries`: Menggunakan query agregasi GORM `.Select("COALESCE(SUM(total_price), 0)").Scan(...)` untuk kalkulasi data riil finansial dari database.

### 4. Gorilla WebSocket (github.com/gorilla/websocket)
* **Fungsi**: Protokol komunikasi full-duplex dua arah untuk live chat.
* **Upgrade Isolasi Hub**: Hub diperbarui menggunakan `sync.RWMutex` thread-safe dengan registry map `connections: make(map[int32]*websocket.Conn)`. Server memisahkan jalur pengiriman pesan dan mendistribusikannya HANYA ke target penerima yang sesuai di dalam registry (tidak lagi melakukan broadcast global).
* **Koneksi Stabil (Tanpa Re-koneksi)**: Front-end mempertahankan satu koneksi WebSocket persisten per sesi login. Mengubah kontak tujuan chat dilakukan tanpa memutus koneksi websocket utama, mencegah delay atau hilangnya pesan.

---

## 🎨 Tech Stack — Frontend Architecture

### 1. Vue.js 3 (Composition API)
* **Fungsi**: Framework reaktif utama untuk antarmuka pengguna berbasis komponen modular.
* **Kegunaan Spesifik**:
  * Menggunakan Composition API `<script setup>` untuk reaktivitas state.
  * Menyimpan dan memvalidasi kredensial login aktif di `localStorage`.
  * Dynamic conditional rendering untuk komponen `ChatBox` dan `OrderModal`.

### 2. Vue Router (vue-router v4)
* **Fungsi**: Penanganan routing halaman client-side dengan **Navigation Guards**.
* **Kegunaan Spesifik**:
  * Melindungi route `/shop` khusus untuk role `customer`.
  * Melindungi route `/admin/dashboard` khusus untuk role `admin`.
  * Mengalihkan pengunjung yang tidak memiliki sesi login aktif kembali ke halaman `/login`.

### 3. Anime.js v4 (Animation Library)
* **Fungsi**: Library animasi performa tinggi untuk mempercantik UI/UX.
* **Kegunaan Spesifik**:
  * Stagger entrance fade-in untuk produk katalog.
  * Rotasi 3D interaktif pada hover kartu produk (`ProductCard.vue`).
  * Efek elastis bounce saat reset kartu produk dan pop-up modal.

### 4. Custom CSS Design System (NTT Heritage Style)
* **Fungsi**: Memberikan visualisasi bercorak kebudayaan tenun & bumi NTT dengan Glassmorphism efek kaca buram dan palette warna eksklusif (`#1a1a40`, `#f5a623`, `#d35400`, dll.).

---

## 🗄️ Tech Stack — Database

### MySQL (via XAMPP)
* **Fungsi**: Penyimpanan data persisten.
* **Tabel Terkait**:
  1. `users`: Akun pengguna dengan kolom `id`, `name`, `username`, `email`, `password`, `role` (`admin`, `vendor`, `customer`), `phone`, `address`, dan `avatar`.
  2. `products`: Katalog produk pre-order dengan kolom `id`, `vendor_id`, `name`, `price`, `description`, `category`, `region`, `po_duration`, `min_order`, `image_url`, `stock` (stok yang ditetapkan admin), dan `is_active` (status tayang).
  3. `orders`: Transaksi pre-order dengan kolom `id`, `customer_id`, `product_id`, `quantity`, `total_price`, `status` (`pending`, `produced`, `shipped`, `completed`, `cancelled`), dan `note`.
  4. `messages`: Percakapan chat dengan kolom `id`, `sender_id`, `receiver_id`, `content` (isi pesan), dan `is_read` (status terbaca).

---

## 📁 Struktur Project

```
PenelitianIlmiahJuan/
│
├── 📂 backend/                          # 🔧 Backend Engine (Golang)
│   ├── main.go                          # Entry point: db sync, route, CORS middleware
│   ├── go.mod / go.sum                  # Manajemen dependensi backend
│   │
│   ├── 📂 config/
│   │   └── database.go                  # Koneksi GORM ke MySQL port 3306
│   │
│   ├── 📂 models/                       # GORM Struct Models
│   │   ├── user.go                      # Model Akun Pengguna
│   │   ├── product.go                   # Model Produk
│   │   ├── order.go                     # Model Transaksi
│   │   └── message.go                   # Model Pesan & WS Payload
│   │
│   └── 📂 handlers/                     # API Controllers
│       ├── auth_handler.go              # Login, Register, GetCurrentUser, GetChatContacts
│       ├── admin_handler.go             # Sales Aggregator Engine (KPI & charts)
│       ├── product_handler.go           # CRUD produk
│       ├── order_handler.go             # CRUD transaksi & Pengurangan Stok
│       └── chat_handler.go              # WebSocket Hub & Chat Logs history
│
└── 📂 frontend/                         # 🎨 Frontend Architecture (Vue 3)
    ├── index.html                       # HTML Entry + SEO Meta Tags
    ├── package.json                     # Scripts & dep (animejs, router, vue)
    ├── vite.config.js                   # Proxy API port 8080 & build settings
    │
    └── 📂 src/
        ├── main.js                      # Mount app & inject Router
        ├── App.vue                      # Router-view viewport shell
        │
        ├── 📂 assets/
        │   └── main.css                 # NTT design system & Glassmorphism
        │
        ├── 📂 router/
        │   └── index.js                 # Routing rules & auth guards
        │
        ├── 📂 views/                    # Routed Screens (Views)
        │   ├── HomePage.vue             # Beranda Utama dengan tombol dinamis
        │   ├── LoginPage.vue            # Login Screen & Registrasi Akun baru
        │   ├── ShopPage.vue             # Halaman belanja & chat seller
        │   └── AdminDashboard.vue       # Panel admin + SVG Chart + Chat panel
        │
        └── 📂 components/
            ├── ProductCard.vue          # Kartu produk 3D
            ├── ProductList.vue          # Grid produk NTT
            ├── OrderModal.vue           # Modal checkout pre-order
            └── ChatBox.vue              # Live Chat dengan dynamic contact list
```

---

## 🌐 API Endpoints

### REST API (`http://localhost:8080/api`)

| Method | Endpoint | Deskripsi | Parameter/Payload |
|--------|----------|-----------|-------------------|
| `POST` | `/api/login` | Otentikasi login pengguna | `{ "username": "...", "password": "..." }` |
| `POST` | `/api/register` | Pendaftaran akun customer baru | `{ "name": "...", "username": "...", "email": "...", "password": "...", "phone": "...", "address": "..." }` |
| `GET` | `/api/user/:id` | Mengambil data profile user | — |
| `GET` | `/api/chat/contacts` | Daftar kontak chat beserta jumlah unread counts | `?user_id=X&role=Y` |
| `GET` | `/api/admin/sales-data` | Agregasi statistik & log transaksi (Admin) | — |
| `POST` | `/api/products` | Unggah produk baru (Hanya Admin) | `{ name, price, region, category, min_order, po_duration, stock, ... }` |
| `GET` | `/api/products` | Ambil katalog produk dengan filter | `?region=Sumba&category=Tenun` |
| `PUT` | `/api/products/:id` | Edit info produk & kelola stok (Hanya Admin) | `{ name, price, stock, is_active, ... }` |
| `DELETE` | `/api/products/:id` | Hapus produk dari database (Hanya Admin) | `?user_id=ADMIN_ID` |
| `POST` | `/api/orders` | Checkout pre-order produk (Deduct Stock) | `{ customer_id, product_id, quantity, note }` |
| `GET` | `/api/orders` | Ambil daftar transaksi pre-order | `?customer_id=X&product_id=Y&status=Z` |
| `DELETE` | `/api/orders/:id` | Hapus transaksi dari log (Hanya Admin) | `?user_id=ADMIN_ID` |
| `GET` | `/api/messages` | Riwayat obrolan dua pengguna (Mark Read) | `?sender_id=X&receiver_id=Y` |
| `GET` | `/health` | Pemeriksaan kesehatan server | — |

### WebSocket (`ws://localhost:8080/ws/chat`)

| Parameter | Deskripsi |
|-----------|-----------|
| `sender_id` | ID user pengirim WebSocket (Wajib) |
| `receiver_id` | ID target user penerima (Opsional) |

---

## 📋 Prasyarat Instalasi

Pastikan program berikut telah terpasang di komputer Anda:
* **XAMPP** (dengan module MySQL aktif)
* **Go (Golang)** v1.22+
* **Node.js** v18+ dan **npm**

---

## 🚀 Langkah-Langkah Menjalankan Project

### Langkah 1: Jalankan MySQL Server
1. Buka **XAMPP Control Panel** dan klik **Start** pada module **MySQL**.
2. Masuk ke phpMyAdmin di: [http://localhost/phpmyadmin](http://localhost/phpmyadmin).
3. Buat database baru bernama: `db_openpeo` (collation: `utf8mb4_general_ci`).
4. *(Opsional)* Jika ingin mengimpor schema awal secara manual, jalankan skrip di file [db.sql](file:///c:/Users/lenovo/Projects/PenelitianIlmiahJuan/db.sql). Jika tidak, GORM akan meng-automigrate skema tabel secara otomatis saat backend dijalankan.

### Langkah 2: Jalankan Backend Server
1. Buka terminal baru dan jalankan perintah:
   ```bash
   cd c:\Users\lenovo\Projects\PenelitianIlmiahJuan\backend
   go mod tidy
   go run main.go
   ```
2. Server backend akan menyala di port `:8080`. Backend juga akan mendeteksi database kosong dan melakukan auto-seed **4 user akun demo** beserta **8 produk NTT**.

### Langkah 3: Jalankan Frontend Server
1. Buka terminal baru yang terpisah dan jalankan perintah:
   ```bash
   cd c:\Users\lenovo\Projects\PenelitianIlmiahJuan\frontend
   npm install
   npm run dev
   ```
2. Server frontend web akan berjalan pada url: **[http://localhost:5173](http://localhost:5173)**.

---

## 🔑 Kredensial Uji Coba

Gunakan kredensial berikut pada halaman login untuk masuk ke halaman role masing-masing:

### 1. Akun Admin
* **Username**: `admin_flores`
* **Password**: `password123`
* **Role**: Admin
* **Tujuan Pengalihan**: `/admin/dashboard`

### 2. Akun Customer (Seeded)
* **Username**: `pembeli_flores`
* **Password**: `password123`
* **Role**: Customer
* **Tujuan Pengalihan**: `/shop`

---

## 🧪 Panduan Pengujian Fitur

### 1. Registrasi & Login Customer Baru
* Buka browser ke halaman utama [http://localhost:5173](http://localhost:5173) dan klik **Masuk**.
* Klik tombol **Daftar sekarang** di bawah form.
* Masukkan data pendaftaran Anda (Nama Lengkap, Username, Email, Telepon, Alamat, Password minimal 6 karakter).
* Klik **Daftar**. Sistem akan memproses dan mengarahkan kembali ke form login dengan kredensial terisi otomatis.
* Klik **Masuk** -> Anda akan langsung dialihkan ke katalog belanja `/shop`.

### 2. Pengujian Stok Produk & Pre-Order
* Login sebagai Customer dan buka halaman katalog belanja `/shop`.
* Perhatikan nominal stok di masing-masing produk (misal: Sisa Stok `15 pcs`).
* Klik **Pre-Order** pada salah satu produk untuk membuka modal pesanan.
* Coba klik tombol tambah (`+`) hingga melebihi stok yang tersedia. Sistem akan membatasi kuantitas agar tidak melampaui stok dan menampilkan teks peringatan merah *"Stok tidak mencukupi"*.
* Lakukan order sukses dalam batas stok. Anda akan melihat stok di kartu produk otomatis berkurang secara langsung.
* Jika stok produk diubah admin menjadi `0`, tombol pre-order produk tersebut akan nonaktif dan bertuliskan *"Habis"*.

### 3. Pengujian Dashboard Admin & Log Transaksi
* Login sebagai Admin (`admin_flores` / `password123`) -> Masuk ke `/admin/dashboard`.
* Pada tab **Ringkasan**, perhatikan 4 kartu KPI statistik teratas (**Total Pendapatan, Jumlah Pre-Order, Produk Aktif, dan Total Pelanggan**) yang sekarang terisi data riil dari database.
* Di bawah log transaksi, klik tombol **Hapus** pada salah satu baris transaksi. Seluruh KPI (seperti Total Pendapatan dan grafik harian) akan dihitung ulang secara real-time.
* Jika semua transaksi dihapus, visual grafik wilayah akan menampilkan tulisan *"Belum ada catatan pendapatan wilayah"* secara bersih.
* Buka tab **Kelola Produk** untuk mencoba fitur CRUD (menambah produk baru, mengedit info stok/harga, atau menghapus produk) yang diotorisasi penuh khusus untuk Admin.

### 4. Pengujian WebSocket Live Chat (Tanpa Reconnect & Badge Unread)
* Pengujian chat disarankan menggunakan **dua browser berbeda** (misal Google Chrome dan jendela Incognito):
  * **Browser 1 (Customer)**: Login menggunakan akun baru Anda. Buka widget chat di pojok kanan bawah. Kirim pesan ke **Admin Flores**.
  * **Browser 2 (Admin)**: Login sebagai Admin. Buka widget chat. 
  * Di browser Admin, nama Customer baru tersebut akan langsung muncul di dropdown chat admin beserta notifikasi unread badge: **`Username (CUSTOMER) (💬 1 baru)`** tanpa perlu memuat ulang halaman.
  * Di browser Admin, coba ganti penerima chat di select dropdown -> Koneksi WebSocket akan tetap tersambung secara stabil tanpa terjadi putus nyambung (reconnect).
  * Balas pesan dari sisi Admin -> Pesan terkirim dan diterima secara real-time oleh Customer.

---

## 📝 Catatan Pengembangan

### Penanganan CORS
CORS telah dikonfigurasi di `backend/main.go` untuk mengizinkan request dari origin `http://localhost:5173` demi mempermudah debugging local development server.

### Build Production
Apabila Anda ingin mengevaluasi bundling statis frontend untuk di-deploy ke hosting server, jalankan command:
```bash
cd c:\Users\lenovo\Projects\PenelitianIlmiahJuan\frontend
npm run build
```
Vite akan mem-bundle komponen Single File Vue menjadi static files ter-kompresi di folder `/dist` dalam waktu singkat.

---

<p align="center">
  <strong>Dibuat dengan ❤️ untuk Nusa Tenggara Timur</strong><br/>
  <sub>OpenPeo — Marketplace Produk Tradisional NTT</sub>
</p>
