CREATE DATABASE IF NOT EXISTS db_openpeo;

USE db_openpeo;

-- 1. TABEL USERS (Menyimpan data akun untuk Login)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) DEFAULT '',
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) DEFAULT '',
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'vendor', 'customer') DEFAULT 'customer',
    phone VARCHAR(20) DEFAULT '',
    address TEXT DEFAULT NULL,
    avatar VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- 2. TABEL PRODUCTS
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vendor_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(12, 2) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    po_duration INT NOT NULL,
    min_order INT DEFAULT 1,
    image_url VARCHAR(255) DEFAULT '',
    stock INT NOT NULL DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (vendor_id) REFERENCES users (id) ON DELETE CASCADE
);

-- 3. TABEL ORDERS (Untuk Dashboard Penjualan Realtime Admin)
CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL(12, 2) NOT NULL,
    status ENUM('pending', 'produced', 'shipped', 'completed', 'cancelled') DEFAULT 'pending',
    note TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

-- 4. TABEL MESSAGES (Riwayat Chat Berdasarkan Room User & Admin Spesifik)
CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    content TEXT NOT NULL,
    is_read TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

-- DUMMY DATA COCOK UNTUK TEST LOGIN & REALTIME DASHBOARD
INSERT INTO users (id, name, username, email, password, role, phone, address, avatar)
VALUES 
    (1, 'Admin Flores', 'admin_flores', 'admin@openpeo.id', 'password123', 'admin', '081234567890', 'Kupang, NTT', ''),
    (2, 'Admin Sumba', 'admin_sumba', 'sumba@openpeo.id', 'password123', 'admin', '081234567891', 'Waingapu, Sumba Timur, NTT', ''),
    (3, 'Pembeli Flores', 'pembeli_flores', 'flores@email.com', 'password123', 'customer', '081234567892', 'Ende, Flores, NTT', ''),
    (4, 'Pembeli Sumba', 'pembeli_sumba', 'budi@email.com', 'password123', 'customer', '081234567893', 'Jakarta, Indonesia', '');

-- Seeding all 8 NTT heritage products directly
INSERT INTO products (id, vendor_id, name, price, description, category, region, po_duration, min_order, image_url, stock, is_active)
VALUES
    (1, 1, 'Kain Tenun Ikat Sumba Pahikung', 1500000.00, 'Kain tenun ikat tradisional Sumba Timur bermotif Pahikung yang melambangkan kebangsawanan dan kekuatan. Dibuat dengan teknik tenun ikat pakan warisan turun-temurun menggunakan pewarna alami dari akar dan daun.', 'Tenun', 'Sumba', 30, 2, 'https://images.unsplash.com/photo-1558171813-4c088753af8f?w=400', 15, 1),
    (2, 1, 'Selimut Tenun Sumba Hinggi Kombu', 2500000.00, 'Selimut tenun Hinggi Kombu khas Sumba dengan motif hewan dan geometris. Simbol kehormatan dan status sosial dalam budaya Marapu.', 'Tenun', 'Sumba', 45, 1, 'https://images.unsplash.com/photo-1606722590583-6951b5ea92ad?w=400', 8, 1),
    (3, 1, 'Songke Manggarai — Kain Adat', 850000.00, 'Kain tenun songke khas Manggarai dengan motif wela mpuu (bunga penuh). Ditenun secara tradisional oleh para perempuan Manggarai untuk upacara adat.', 'Tenun', 'Manggarai', 21, 3, 'https://images.unsplash.com/photo-1621600411688-4be93cd68504?w=400', 12, 1),
    (4, 1, 'Kopi Flores Bajawa Single Origin', 120000.00, 'Kopi arabika single origin dari dataran tinggi Bajawa, Flores. Ditanam di ketinggian 1200-1600 mdpl dengan cita rasa coklat, karamel, dan sentuhan rempah khas.', 'Kuliner', 'Flores', 14, 5, 'https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=400', 50, 1),
    (5, 1, 'Manik-Manik Sumba Handmade Necklace', 350000.00, 'Kalung manik-manik handmade khas Sumba dengan warna-warna tanah. Setiap butir dibentuk dan diwarnai secara manual oleh pengrajin lokal.', 'Aksesoris', 'Sumba', 14, 3, 'https://images.unsplash.com/photo-1611085583191-a3b181a88401?w=400', 20, 1),
    (6, 1, 'Madu Hutan Timor Asli', 175000.00, 'Madu hutan asli dari pegunungan Timor, NTT. Dipanen langsung dari sarang lebah liar di hutan lindung dengan rasa yang kaya dan kental.', 'Kuliner', 'Timor', 10, 4, 'https://images.unsplash.com/photo-1587049352846-4a222e784d38?w=400', 35, 1),
    (7, 1, 'Patung Ukir Kayu Sandalwood Kupang', 950000.00, 'Patung ukiran kayu cendana (sandalwood) buatan tangan pengrajin Kupang. Motif tradisional NTT dengan aroma khas kayu cendana yang tahan lama.', 'Kerajinan', 'Kupang', 30, 1, 'https://images.unsplash.com/photo-1513519245088-0e12902e35ca?w=400', 5, 1),
    (8, 1, 'Syal Tenun Ende Lio', 450000.00, 'Syal tenun ikat Ende Lio dengan motif bunga dan geometris dalam warna-warna alam. Pewarna alami dari tumbuhan lokal daerah Ende.', 'Tenun', 'Flores', 21, 2, 'https://images.unsplash.com/photo-1601924921557-45e16393d8e1?w=400', 18, 1);

-- Transaksi dummy agar Dashboard Admin langsung terisi grafik/datanya
INSERT INTO orders (id, customer_id, product_id, quantity, total_price, status, note, created_at)
VALUES 
    (1, 3, 4, 5, 600000.00, 'pending', 'Kopi Bajawa tolong kirim biji utuh', NOW() - INTERVAL 2 DAY),
    (2, 4, 1, 2, 3000000.00, 'shipped', 'Kain tenun Sumba motif naga', NOW() - INTERVAL 1 DAY);