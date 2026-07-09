<template>
  <div class="tentang-page">
    <div class="tentang-header">
      <div class="container">
        <h1>Tentang OpenPeo</h1>
        <p>Lebih dari sekadar platform, kami adalah jembatan yang menghubungkan karya tangan autentik dengan apresiasi global.</p>
      </div>
    </div>
    
    <div class="container tentang-content">
      <section class="tentang-section">
        <h2>Cerita Kami</h2>
        <p>OpenPeo lahir dari kecintaan kami terhadap warisan budaya Nusa Tenggara Timur. Kami menyadari bahwa banyak pengrajin lokal berbakat yang belum memiliki akses ke pasar yang lebih luas. Melalui platform ini, kami berkomitmen untuk melestarikan tradisi sekaligus memberdayakan komunitas pengrajin lokal dengan sistem pre-order yang adil dan berkelanjutan.</p>
      </section>

      <section class="tentang-section">
        <h2>Pertanyaan yang Sering Diajukan (FAQ)</h2>
        <div class="faq-list">
          <div class="faq-item">
            <h3>Bagaimana sistem pre-order bekerja?</h3>
            <p>Anda dapat memesan produk yang sedang diproduksi oleh pengrajin kami. Waktu pengerjaan bervariasi tergantung kerumitan produk, biasanya berkisar antara 2 hingga 4 minggu.</p>
          </div>
          <div class="faq-item">
            <h3>Apakah produk dijamin autentik?</h3>
            <p>Ya, semua produk kami dikurasi secara ketat dan langsung dikirim dari pengrajin lokal di Nusa Tenggara Timur untuk menjamin keaslian dan kualitas terbaik.</p>
          </div>
          <div class="faq-item">
            <h3>Metode pembayaran apa yang tersedia?</h3>
            <p>Kami menerima transfer bank, kartu kredit, dan berbagai e-wallet populer di Indonesia untuk memudahkan transaksi Anda.</p>
          </div>
        </div>
      </section>

      <section class="tentang-section contact-section">
        <h2>Hubungi Kami</h2>
        <p>Memiliki pertanyaan lebih lanjut? Tim kami siap membantu Anda.</p>
        <div class="contact-info">
          <p><strong>Email:</strong> halo@openpeo.com</p>
          <p><strong>Telepon / WhatsApp:</strong> +62 812 3456 7890</p>
          <p><strong>Alamat:</strong> Jl. El Tari No. 123, Kupang, Nusa Tenggara Timur, Indonesia</p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.tentang-page {
  background-color: #fafafa;
  min-height: 100vh;
  padding-bottom: 80px;
}

.tentang-header {
  padding: 80px 20px;
  text-align: center;
  background-color: #ffffff;
  border-bottom: 1px solid #eaeaea;
}

.tentang-header h1 {
  font-family: 'Playfair Display', serif;
  font-size: 2.75rem;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 16px;
}

.tentang-header p {
  font-family: 'Montserrat', sans-serif;
  color: #666;
  font-size: 1.1rem;
  line-height: 1.6;
  max-width: 600px;
  margin: 0 auto;
}

.tentang-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 60px 20px;
}

.tentang-section {
  margin-bottom: 60px;
}

.tentang-section h2 {
  font-family: 'Playfair Display', serif;
  font-size: 2rem;
  color: #1a1a1a;
  margin-bottom: 24px;
  position: relative;
  display: inline-block;
}

.tentang-section h2::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  width: 40px;
  height: 2px;
  background-color: #d4af37;
}

.tentang-section p {
  font-family: 'Montserrat', sans-serif;
  color: #555;
  font-size: 1.05rem;
  line-height: 1.8;
}

.faq-item {
  margin-bottom: 32px;
}

.faq-item h3 {
  font-size: 1.2rem;
  color: #222;
  margin-bottom: 12px;
  font-weight: 600;
}

.contact-info {
  margin-top: 24px;
  background-color: #fff;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.03);
}

.contact-info p {
  margin-bottom: 12px;
  font-size: 1rem;
}

.contact-info p:last-child {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .tentang-header h1 {
    font-size: 2.25rem;
  }
  
  .tentang-section h2 {
    font-size: 1.75rem;
  }
}
</style>
