




(function () {
  'use strict';

  
  
  
  const header = document.getElementById('site-header');
  const SCROLL_THRESHOLD = 60;

  function handleHeaderScroll() {
    if (window.scrollY > SCROLL_THRESHOLD) {
      header.classList.add('scrolled');
    } else {
      header.classList.remove('scrolled');
    }
  }

  window.addEventListener('scroll', handleHeaderScroll, { passive: true });
  handleHeaderScroll(); 


  
  
  
  const menuToggle = document.getElementById('menu-toggle');
  const mobileNav = document.getElementById('mobile-nav');

  if (menuToggle && mobileNav) {
    menuToggle.addEventListener('click', function () {
      const isOpen = mobileNav.classList.toggle('open');
      menuToggle.classList.toggle('active');
      menuToggle.setAttribute('aria-expanded', isOpen);

      
      if (isOpen) {
        header.classList.add('scrolled');
        document.body.style.overflow = 'hidden';
      } else {
        document.body.style.overflow = '';
        handleHeaderScroll();
      }
    });

    
    mobileNav.querySelectorAll('a').forEach(function (link) {
      link.addEventListener('click', function () {
        mobileNav.classList.remove('open');
        menuToggle.classList.remove('active');
        menuToggle.setAttribute('aria-expanded', 'false');
        document.body.style.overflow = '';
        handleHeaderScroll();
      });
    });
  }


  
  
  
  const revealElements = document.querySelectorAll('.reveal');

  if ('IntersectionObserver' in window && revealElements.length > 0) {
    const revealObserver = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            entry.target.classList.add('visible');
            revealObserver.unobserve(entry.target);
          }
        });
      },
      {
        threshold: 0.15,
        rootMargin: '0px 0px -40px 0px',
      }
    );

    revealElements.forEach(function (el) {
      revealObserver.observe(el);
    });
  } else {
    
    revealElements.forEach(function (el) {
      el.classList.add('visible');
    });
  }


  
  
  
  document.querySelectorAll('a[href^="#"]').forEach(function (anchor) {
    anchor.addEventListener('click', function (e) {
      const targetId = this.getAttribute('href');
      if (targetId === '#') return;

      const targetEl = document.querySelector(targetId);
      if (targetEl) {
        e.preventDefault();
        const headerHeight = header.offsetHeight;
        const targetPosition = targetEl.getBoundingClientRect().top + window.scrollY - headerHeight;

        window.scrollTo({
          top: targetPosition,
          behavior: 'smooth',
        });
      }
    });
  });

  
  
  
  const pdpModal = document.getElementById('pdp-modal');
  const pdpClose = document.getElementById('pdp-close');
  const pdpOverlay = document.getElementById('pdp-overlay');
  
  const pdpCategory = document.getElementById('pdp-category');
  const pdpTitle = document.getElementById('pdp-title');
  const pdpPrice = document.getElementById('pdp-price');
  const pdpPreorderText = document.getElementById('pdp-preorder-text');
  const pdpDescText = document.getElementById('pdp-desc-text');
  const pdpDescNote = document.getElementById('pdp-desc-note');
  const pdpMainImg = document.getElementById('pdp-main-img');
  const pdpThumbnails = document.getElementById('pdp-thumbnails');
  
  const qtyMinus = document.getElementById('pdp-qty-minus');
  const qtyPlus = document.getElementById('pdp-qty-plus');
  const qtyValue = document.getElementById('pdp-qty-value');
  let currentQty = 1;

  function formatRupiah(number) {
    return 'Rp ' + number.toLocaleString('id-ID');
  }

  function openPDP(productId) {
    
    const product = (typeof productsData !== 'undefined') ? productsData.find(p => p.id === productId) : null;
    if (!product) return;

    
    pdpCategory.textContent = product.category;
    pdpTitle.textContent = product.name;
    pdpPrice.textContent = formatRupiah(product.price);
    pdpPreorderText.textContent = `Estimasi Pre-order: ${product.preOrderDays} hari kerja`;
    pdpDescText.textContent = product.description;
    
    
    if (product.category === "Koleksi Tenun NTT") {
      pdpDescNote.textContent = "Tersedia layanan custom size (request ukuran sesuai keinginan Anda).";
      pdpDescNote.style.display = "block";
    } else {
      pdpDescNote.style.display = "none";
    }

    
    pdpMainImg.src = product.images[0];
    pdpMainImg.alt = product.name;
    
    pdpThumbnails.innerHTML = '';
    product.images.forEach((imgSrc, index) => {
      const img = document.createElement('img');
      img.src = imgSrc;
      img.className = 'pdp-thumbnail' + (index === 0 ? ' active' : '');
      img.addEventListener('click', function() {
        pdpMainImg.src = imgSrc;
        document.querySelectorAll('.pdp-thumbnail').forEach(el => el.classList.remove('active'));
        img.classList.add('active');
      });
      pdpThumbnails.appendChild(img);
    });

    
    currentQty = 1;
    qtyValue.textContent = currentQty;

    
    if (pdpModal) pdpModal.classList.add('open');
    document.body.style.overflow = 'hidden'; 
  }

  function closePDP() {
    if (pdpModal) pdpModal.classList.remove('open');
    document.body.style.overflow = ''; 
  }

  
  document.querySelectorAll('.product-card').forEach(card => {
    card.addEventListener('click', function(e) {
      e.preventDefault();
      const productId = this.getAttribute('id');
      openPDP(productId);
    });
  });

  
  if(pdpClose) pdpClose.addEventListener('click', closePDP);
  if(pdpOverlay) pdpOverlay.addEventListener('click', closePDP);

  
  if(qtyMinus && qtyPlus) {
    qtyMinus.addEventListener('click', () => {
      if (currentQty > 1) {
        currentQty--;
        qtyValue.textContent = currentQty;
      }
    });
    qtyPlus.addEventListener('click', () => {
      currentQty++;
      qtyValue.textContent = currentQty;
    });
  }

})();
