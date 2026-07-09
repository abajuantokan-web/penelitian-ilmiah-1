const productsData = [
  // SECTION 1: Koleksi Tenun NTT
  {
    id: "product-1",
    name: "Kaftan Indigo Sumba",
    price: 2850000,
    category: "Koleksi Tenun NTT",
    images: ["images/product-1.png", "images/product-1.png", "images/product-1.png"],
    description: "Kaftan tenun ikat berwarna indigo biru tua dengan motif geometris putih dari Sumba. Ditenun secara tradisional memakan waktu berbulan-bulan, menghasilkan karya seni yang elegan dan otentik.",
    preOrderDays: 14,
    stockStatus: "Pre-order"
  },
  {
    id: "product-2",
    name: "Selendang Merah Sabu",
    price: 1950000,
    category: "Koleksi Tenun NTT",
    images: ["images/product-2.png", "images/product-2.png", "images/product-2.png"],
    description: "Selendang tenun berwarna merah coral dengan garis-garis halus dan motif tradisional Sabu. Sangat cocok dipadukan dengan pakaian modern maupun tradisional.",
    preOrderDays: 14,
    stockStatus: "Pre-order"
  },
  {
    id: "product-3",
    name: "Syal Tenun Amarasi",
    price: 1250000,
    category: "Koleksi Tenun NTT",
    images: ["images/product-3.png", "images/product-3.png", "images/product-3.png"],
    description: "Syal tenun ikat berwarna coklat dan krem dengan motif tradisional Amarasi. Lembut dan hangat, memberikan sentuhan etnik pada setiap penampilan.",
    preOrderDays: 7,
    stockStatus: "Pre-order"
  },
  {
    id: "product-4",
    name: "Tunik Hitam Rote",
    price: 3150000,
    category: "Koleksi Tenun NTT",
    images: ["images/product-4.png", "images/product-4.png", "images/product-4.png"],
    description: "Tunik tenun ikat berwarna hitam dengan motif emas dan krem dari Pulau Rote. Desain timeless dengan kualitas benang premium terbaik.",
    preOrderDays: 21,
    stockStatus: "Pre-order"
  },
  {
    id: "product-5",
    name: "Outer Tenun Ikat Sumba",
    price: 1450000,
    category: "Koleksi Tenun NTT",
    images: ["images/tenun-outer.png", "images/tenun-outer.png", "images/tenun-outer.png"],
    description: "Mannequin/model wearing a chic, modern-cut outer featuring bold Sumba horse motifs. Elegan untuk acara formal maupun semi-formal.",
    preOrderDays: 14,
    stockStatus: "Pre-order"
  },
  {
    id: "product-6",
    name: "Kemeja Tenun Ende Pria",
    price: 950000,
    category: "Koleksi Tenun NTT",
    images: ["images/tenun-kemeja.png", "images/tenun-kemeja.png", "images/tenun-kemeja.png"],
    description: "Professional man's shirt with classic, geometric Ende patterns. Memadukan kesan profesional dan tradisi, pas untuk gaya kantor modern.",
    preOrderDays: 10,
    stockStatus: "Pre-order"
  },
  {
    id: "product-7",
    name: "Dress Tenun Manggarai",
    price: 2100000,
    category: "Koleksi Tenun NTT",
    images: ["images/tenun-dress.png", "images/tenun-dress.png", "images/tenun-dress.png"],
    description: "Elegant women's dress with premium, detailed Manggarai weave patterns. Siluet yang menawan untuk acara spesial.",
    preOrderDays: 21,
    stockStatus: "Pre-order"
  },
  {
    id: "product-8",
    name: "Blouse Tenun Alor",
    price: 1150000,
    category: "Koleksi Tenun NTT",
    images: ["images/tenun-blouse.png", "images/tenun-blouse.png", "images/tenun-blouse.png"],
    description: "Stylish blouse featuring the distinct, colorful Alor textile designs. Nyaman dikenakan sepanjang hari dengan potongan yang modern.",
    preOrderDays: 14,
    stockStatus: "Pre-order"
  },

  // SECTION 2: Cita Rasa Lokal
  {
    id: "product-sei-babi",
    name: "Sei Babi Asap Kupang",
    price: 185000,
    category: "Cita Rasa Lokal",
    images: ["images/sei-babi.png", "images/sei-babi.png", "images/sei-babi.png"],
    description: "Sei babi asap tradisional dari Kupang dengan irisan daging berwarna amber keemasan. Diproses dengan kayu kosambi yang memberikan aroma khas.",
    preOrderDays: 3,
    stockStatus: "Pre-order"
  },
  {
    id: "product-madu-hutan",
    name: "Madu Hutan Timor Murni",
    price: 145000,
    category: "Cita Rasa Lokal",
    images: ["images/madu-hutan.png", "images/madu-hutan.png", "images/madu-hutan.png"],
    description: "Toples madu hutan murni dari Timor berwarna emas tua. Diperoleh langsung dari hutan belantara, kaya akan nutrisi dan antioksidan.",
    preOrderDays: 3,
    stockStatus: "Pre-order"
  },
  {
    id: "product-kopi-flores",
    name: "Kopi Arabika Flores Bajawa",
    price: 95000,
    category: "Cita Rasa Lokal",
    images: ["images/kopi-flores.png", "images/kopi-flores.png", "images/kopi-flores.png"],
    description: "Kemasan kopi arabika Flores Bajawa dengan biji kopi roasting gelap. Memiliki hints cokelat dan rempah yang kuat dengan acidity medium.",
    preOrderDays: 5,
    stockStatus: "Pre-order"
  },
  {
    id: "product-garam-gunung",
    name: "Garam Gunung Organik Khas NTT",
    price: 65000,
    category: "Cita Rasa Lokal",
    images: ["images/garam-gunung.png", "images/garam-gunung.png", "images/garam-gunung.png"],
    description: "Garam gunung organik NTT berwarna merah muda keabu-abuan. Tinggi mineral alami dan sangat cocok untuk fine dining.",
    preOrderDays: 5,
    stockStatus: "Pre-order"
  },
  {
    id: "product-jagung-titi",
    name: "Jagung Titi Khas NTT",
    price: 35000,
    category: "Cita Rasa Lokal",
    images: ["images/jagung-titi.png", "images/jagung-titi.png", "images/jagung-titi.png"],
    description: "Close-up of crispy, flat traditional corn snack in a premium rustic wooden bowl. Camilan ringan dan sehat untuk teman minum kopi.",
    preOrderDays: 2,
    stockStatus: "Pre-order"
  },
  {
    id: "product-sambal-luat",
    name: "Sambal Luat Kupang",
    price: 45000,
    category: "Cita Rasa Lokal",
    images: ["images/sambal-luat.png", "images/sambal-luat.png", "images/sambal-luat.png"],
    description: "Macro shot of authentic, spicy red Sambal Luat in a small glass jar with fresh chili garnishes. Pedas yang menyegarkan dengan sentuhan jeruk nipis.",
    preOrderDays: 2,
    stockStatus: "Pre-order"
  },
  {
    id: "product-kacang-sembunyi",
    name: "Kacang Sembunyi NTT",
    price: 30000,
    category: "Cita Rasa Lokal",
    images: ["images/kacang-sembunyi.png", "images/kacang-sembunyi.png", "images/kacang-sembunyi.png"],
    description: "Pile of crunchy, golden-brown caramelized nuts, studio lighting, clean background. Gurih manis yang adiktif.",
    preOrderDays: 2,
    stockStatus: "Pre-order"
  },
  {
    id: "product-gula-air",
    name: "Gula Air Sabu",
    price: 55000,
    category: "Cita Rasa Lokal",
    images: ["images/gula-air.png", "images/gula-air.png", "images/gula-air.png"],
    description: "Close-up of authentic, thick palm nectar in a clear bottle, honey-colored glow. Pemanis alami rendah indeks glikemik.",
    preOrderDays: 3,
    stockStatus: "Pre-order"
  },

  // SECTION 3: Koleksi Aksesoris
  {
    id: "product-headband",
    name: "Headband Kain Tenun Sumba",
    price: 275000,
    category: "Koleksi Aksesoris",
    images: ["images/headband-tenun.png", "images/headband-tenun.png", "images/headband-tenun.png"],
    description: "Headband kain tenun ikat Sumba berwarna indigo biru dengan motif geometris krem. Aksesori chic yang mudah dipadupadankan.",
    preOrderDays: 7,
    stockStatus: "Pre-order"
  },
  {
    id: "product-gelang-kalung",
    name: "Gelang & Kalung Serbuk Gading Maumere",
    price: 450000,
    category: "Koleksi Aksesoris",
    images: ["images/gelang-kalung.png", "images/gelang-kalung.png", "images/gelang-kalung.png"],
    description: "Set gelang dan kalung artisan dari Maumere dengan manik-manik serbuk gading dan kuningan. Statement piece untuk berbagai gaya.",
    preOrderDays: 14,
    stockStatus: "Pre-order"
  },
  {
    id: "product-mahkota",
    name: "Mahkota Ti'i Langga Rote",
    price: 850000,
    category: "Koleksi Aksesoris",
    images: ["images/mahkota-tiilangga.png", "images/mahkota-tiilangga.png", "images/mahkota-tiilangga.png"],
    description: "Mahkota Ti'i Langga tradisional dari Rote yang terbuat dari anyaman daun lontar. Sebuah karya seni ikonis yang bersejarah.",
    preOrderDays: 21,
    stockStatus: "Pre-order"
  },
  {
    id: "product-cincin-penyu",
    name: "Cincin & Aksesoris Ornamen Penyu Tradisional",
    price: 320000,
    category: "Koleksi Aksesoris",
    images: ["images/cincin-penyu.png", "images/cincin-penyu.png", "images/cincin-penyu.png"],
    description: "Koleksi cincin dan aksesoris ornamen penyu tradisional NTT dari bahan alami yang lestari.",
    preOrderDays: 7,
    stockStatus: "Pre-order"
  },
  {
    id: "product-kalung-timor",
    name: "Kalung Khas Timor",
    price: 150000,
    category: "Koleksi Aksesoris",
    images: ["images/kalung-timor.png", "images/kalung-timor.png", "images/kalung-timor.png"],
    description: "Macro shot of artisan shell and beadwork necklace on a dark velvet display. Detail karya tangan perajin lokal yang teliti.",
    preOrderDays: 7,
    stockStatus: "Pre-order"
  },
  {
    id: "product-sisir-sumba",
    name: "Sisir Adat Sumba",
    price: 225000,
    category: "Koleksi Aksesoris",
    images: ["images/sisir-sumba.png", "images/sisir-sumba.png", "images/sisir-sumba.png"],
    description: "Intricate, hand-carved traditional wooden hair comb with traditional motifs. Aksesori rambut eksotis khas pulau Sumba.",
    preOrderDays: 10,
    stockStatus: "Pre-order"
  },
  {
    id: "product-gelang-kerbau",
    name: "Gelang Kulit Kerbau Rote",
    price: 85000,
    category: "Koleksi Aksesoris",
    images: ["images/gelang-kerbau.png", "images/gelang-kerbau.png", "images/gelang-kerbau.png"],
    description: "Detailed shot of handcrafted leather and bead bracelet. Tahan lama dan bertekstur klasik.",
    preOrderDays: 5,
    stockStatus: "Pre-order"
  },
  {
    id: "product-anting-penyu",
    name: "Anting Motif Penyu Flores",
    price: 120000,
    category: "Koleksi Aksesoris",
    images: ["images/anting-penyu.png", "images/anting-penyu.png", "images/anting-penyu.png"],
    description: "Elegant metalwork earrings with traditional sea-turtle-inspired motifs. Menambah kesan ayu pada penampilan elegan Anda.",
    preOrderDays: 7,
    stockStatus: "Pre-order"
  }
];
