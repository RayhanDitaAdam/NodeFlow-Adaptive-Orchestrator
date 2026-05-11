# GoNode Adaptive Infrastructure Engine 🚀

**GoNode** adalah *infrastructure-as-code orchestrator* berbasis Golang yang dirancang khusus untuk menjalankan aplikasi Node.js tanpa ketergantungan pada Nginx. Cukup satu binary untuk mengatur segalanya: mulai dari *reverse proxy*, manajemen proses, hingga optimasi sumber daya otomatis.

---

## 🧐 Kenapa Harus Pakai GoNode?

Biasanya, kalau kita mau *deploy* aplikasi Node.js (MERN/Next.js) ke VPS, kita harus melewati proses yang melelahkan:
1. **Install & Config Nginx** (Ribet urusan file di `/etc/nginx/sites-available`).
2. **Install PM2** (Menambah dependensi hanya untuk menjaga aplikasi tetap hidup).
3. **Setting SSL Manual** (Seringkali pusing dengan *symlink* certbot).
4. **Alokasi RAM & CPU Manual** (Risiko *crash* jika spek VPS tidak sesuai dengan beban kerja).

**GoNode memangkas semua itu.** GoNode bertindak sebagai "Bodyguard" sekaligus "Manajer" yang memastikan aplikasi Node.js lu tetap stabil, aman, dan kencang dengan *overhead* memori yang sangat rendah.

---

## 🛠️ Fitur Utama

- **Vite-Style Bootstrapping**: Menu interaktif di terminal untuk memilih spek server (Low, Mid, High) sebelum aplikasi berjalan.
- **Embedded Reverse Proxy**: Menggantikan Nginx dengan performa internal Go yang jauh lebih ringan RAM dan CPU.
- **Auto-Restart & Monitoring**: Menggantikan peran PM2. Jika Node.js mati karena *error*, GoNode akan otomatis menyalakannya kembali dalam hitungan milidetik.
- **Zero-Config Infrastructure**: Tidak ada lagi file `.conf` eksternal. Semua logika infrastruktur ada di dalam kode Go.
- **Smart Resource Allocation**: Otomatis menyesuaikan *worker pool*, *timeout*, dan *memory heap* Node.js sesuai profil spek yang dipilih.

---

## 📑 Studi Kasus: Masalah vs Solusi

### 1. Masalah: "VPS Murah (Shared RAM) Sering Crash"
* **Kondisi:** Lu punya VPS 512MB RAM. Pas lu jalanin Nginx + PM2 + Node.js, RAM sudah termakan banyak sebelum ada trafik. Begitu trafik naik, sistem langsung *OOM (Out Of Memory) Kill*.
* **Solusi GoNode:** Pilih profil **"Eco/Low Spec"**. GoNode akan mematikan fitur yang tidak perlu, membatasi *buffer size*, dan mengoptimalkan penggunaan memori agar Node.js punya ruang napas lebih lega di RAM yang terbatas.

### 2. Masalah: "Ribet Mindahin Project ke Server Baru"
* **Kondisi:** Harus pindah VPS tapi malas install ulang Nginx, mengatur konfigurasi yang sama, dan melakukan *symlink* ulang yang rawan *human error*.
* **Solusi GoNode:** Lu cukup *copy* folder proyek dan binary `gonode`. Jalankan di VPS mana pun (Ubuntu/Debian/CentOS), pilih spek, dan web langsung *live* di port 80/443. *Zero dependencies.*

### 3. Masalah: "Aplikasi Mati Pas Tengah Malem"
* **Kondisi:** Ada *logic error* di kode Node.js yang bikin aplikasi *exit* tiba-tiba saat lu lagi tidur.
* **Solusi GoNode:** Sebagai *parent process*, GoNode akan mendeteksi kematian proses Node.js secara *real-time*, melakukan *respawn* otomatis, dan menjaga layanan tetap bisa diakses tanpa campur tangan manual.

---

## 🚀 Cara Menjalankan

1. **Build Project**
   ```bash
   go build -o gonode
