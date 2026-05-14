# 🧠 GoNode - Fundamental & Architecture Guide

File ini dibuat khusus buat ente biar paham "jeroan" GoNode dan dasar-dasar Golang yang kita pake di sini.

---

## 1. Ini Project Buat Apa Sih?
**GoNode** itu ibarat **Manager / Supervisor** buat aplikasi Node.js ente.
- **Masalah**: Biasanya kita ribet harus setup Nginx, atur RAM Node.js, monitor log, dan mastiin dia jalan di background.
- **Solusi**: GoNode otomatisasi itu semua. Ente tinggal pilih mau spek "Eco" atau "Power", GoNode yang bakal ngejalanin Node.js ente, bikinin config Nginx-nya, dan jagain log-nya biar nggak menuhin disk.

---

## 2. Struktur Folder & Fungsinya
Kita pake standar **Go Project Layout** biar rapi dan profesional:

- **`cmd/gonode/`**: 
  - **Isi**: `main.go`
  - **Fungsi**: Pintu masuk utama. Tugasnya cuma nerima perintah dari ente (start/stop/list) terus dilempar ke logic yang sesuai.
- **`pkg/engine/`** (Otak Utama):
  - `cli.go`: Ngurusin tampilan menu interaktif pas ente ngetik `./gonode start`.
  - `daemon.go`: Mesin yang jalan di background. Dia yang ngejalanin proses Node.js dan dengerin perintah lewat socket.
  - `detector.go`: "AI" sederhana (Heuristics) buat ngetebak aplikasi ente itu Next.js atau Node.js biasa.
  - `nginx.go`: Otomatisasi pembuatan file konfigurasi Nginx.
- **`pkg/logger/`** (Monitor):
  - Ngurusin penambahan waktu (timestamp) di log dan motong file log kalo udah kegedean (1MB).
- **`pkg/utils/`** (Assistant):
  - Isinya fungsi bantuan kayak nge-check apakah Go udah keinstall atau nampilin teks bantuan.
- **`docs/`**: Tempat naro dokumen kayak PRD dan file fundamental ini.
- **`examples/`**: Tempat naro contoh kode Node.js (`app.js`).

---

## 3. Konsep Golang yang Kita Pake
Biar ente makin jago Go, ini beberapa konsep kunci yang ada di project ini:

### A. Goroutines (`go func()`)
Golang itu jagonya *Concurrency*. Di `daemon.go`, kita pake `go` keyword buat jalanin logger di "jalur lain". Jadi si GoNode bisa dengerin socket sambil nulis log secara bersamaan tanpa nunggu satu sama lain.

### B. Unix Domain Sockets
Kita nggak pake HTTP buat komunikasi antar perintah (misal `./gonode list`). Kita pake file socket di `/tmp/gonode.sock`. Ini jauh lebih kenceng dan aman buat komunikasi antar proses di dalam satu server.

### C. Structs & Maps
Kita pake `struct` buat nentuin profil RAM (Eco/Balanced/Power). Ibaratnya kayak bikin "Template" data yang isinya udah fix.

### D. OS Exec
Fungsi utama GoNode adalah jalanin perintah sistem (kayak `node app.js` atau `npm start`). Kita pake library `os/exec` buat nge-spawn proses itu dan nangkep output-nya.

### E. Pointer & Interfaces
Pas kita nulis log, kita pake `io.ReadCloser`. Itu adalah *Interface* yang fleksibel, bisa nerima data dari mana aja (output Node.js atau error sistem).

---

## 4. Alur Kerja Aplikasi (Workflow)
1. **Setup**: Ente jalanin `./setup.sh` buat install Go/Node/Nginx.
2. **Build**: Ente jalanin `./install.sh` buat jadiin kode Go ente file binary `gonode`.
3. **Start**: Ente jalanin `./gonode start`.
4. **Daemon**: GoNode bakal "melepas diri" ke background, ngejalanin Node.js, dan nulis log ke `gonode.log`.
5. **Nginx**: GoNode bikinin file `.conf` buat Nginx ente biar website ente bisa diakses lewat domain.

---

Semoga ini ngebantu ente buat makin paham sama project yang ente bangun ini, bre! Semangat belajarnya! 🚀
