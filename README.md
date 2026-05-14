<p align="center">
  <img src="https://raw.githubusercontent.com/RayhanDitaAdam/NodeFlow-Adaptive-Orchestrator/main/Logo.png" width="200" alt="GoNode Logo">
</p>

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

## 📐 Flow Cara Kerja GoNode

Berikut adalah visualisasi alur kerja bagaimana GoNode mengelola infrastruktur aplikasi lu:

<p align="center">
  <img src="https://raw.githubusercontent.com/RayhanDitaAdam/NodeFlow-Adaptive-Orchestrator/main/WorkFlow%20GoNode.png" width="800" alt="GoNode Workflow">
</p>

### High-Level Architecture (Mermaid)
```mermaid
graph TD
    User((🌐 Internet/User)) -->|Port 80/443| GoNode[("🚀 GoNode Engine")]
    
    subgraph "Internal Infrastructure"
    GoNode -->|1. Profiling| Spec[Selection: Eco / Balanced / Power]
    GoNode -->|2. Process Mgmt| NodeApp["📦 Node.js Instance"]
    GoNode -->|3. Reverse Proxy| NodeApp
    end

    subgraph "Health Check & Recovery"
    NodeApp -.->|Exit Code| Monitor{Monitor}
    Monitor -->|Auto-Restart| NodeApp
    end

    style GoNode fill:#00ADD8,stroke:#333,stroke-width:2px,color:#fff
    style NodeApp fill:#339933,stroke:#333,stroke-width:2px,color:#fff
