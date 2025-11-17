# 🖥️ PC Hardware Monitor

> 🚀 Real-time hardware monitoring with a sleek, modern web interface

Monitor your PC's vital stats from anywhere! Track CPU/GPU temperatures and power consumption with beautiful, customizable cards and real-time updates.

![Status](https://img.shields.io/badge/status-active-success.svg)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)
![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?logo=svelte)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)

---

## ✨ Features

### 🌡️ **Dual CPU Temperature Monitoring**
- **Tctl/Tdie**: Overall CPU temperature (thermal control)
- **CCD1**: Core Complex Die temperature (actual hotspot)

### ⚡ **Power Tracking**
- Real-time CPU package power consumption
- GPU power draw monitoring

### 🎨 **Beautiful UI**
- 🌙 Dark theme with gradient borders
- 📱 Responsive design (mobile-friendly)
- 📲 **Progressive Web App (PWA)** - installable on desktop & mobile
- 🎯 Color-coded status indicators
- ⚙️ Customizable temperature/power ranges
- ⏱️ Live countdown timers
- 🔌 Offline support with service worker

### 🔧 **Powerful Configuration**
- 📝 JSON-based metric configuration
- 🔄 Easy to add new metrics (no code changes!)
- 💾 Persistent user preferences (localStorage)
- 🔁 Configurable refresh intervals

### ⚡ **Performance Optimized**
- 🚄 97% smaller payloads (150 bytes vs 5KB)
- 🔋 Battery-friendly (98% fewer updates after 1 min)
- 📊 Smart adaptive timers
- 💰 Upstash free tier friendly (~173k commands/month)

---

## 🏗️ Architecture

```
┌─────────────┐
│  Hardware   │  ← LibreHardwareMonitor (localhost:8085)
│   Monitor   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Backend    │  ← Go (extracts metrics every 10s)
│   (Go)      │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Upstash   │  ← Redis Cloud (stores latest data)
│    Redis    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Frontend   │  ← SvelteKit PWA (auto-refresh every 30s)
│ (SvelteKit) │
└─────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- 🔧 [Go 1.21+](https://go.dev/dl/)
- 📦 [Node.js 18+](https://nodejs.org/)
- 🔥 Hardware Monitor (e.g., [LibreHardwareMonitor](https://github.com/LibreHardwareMonitor/LibreHardwareMonitor))
- ☁️ [Upstash Redis](https://upstash.com/) account (free tier)

### 1️⃣ Setup Hardware Monitor

Download and run LibreHardwareMonitor with Remote Web Server enabled (port 8085).

### 2️⃣ Configure Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your Upstash credentials
```

**`.env` file:**
```bash
UPSTASH_REDIS_ADDR=your-instance.upstash.io:6379
UPSTASH_REDIS_PASSWORD=your_password_token
HARDWARE_MONITOR_URL=http://localhost:8085/data.json
```

### 3️⃣ Configure Frontend

```bash
cd frontend
cp .env.example .env
# Edit .env with your Upstash REST API credentials
```

**`.env` file:**
```bash
PUBLIC_UPSTASH_REDIS_URL=https://your-instance.upstash.io
PUBLIC_UPSTASH_REDIS_TOKEN=your_rest_api_token
```

### 4️⃣ Install Dependencies

```bash
# Backend (from backend/ directory)
cd backend
go mod download

# Frontend (from frontend/ directory)
cd frontend
npm install
```

### 5️⃣ Run Everything!

**Option A: Quick Start (Windows)**
```bash
# Start backend (hidden, no terminal window)
.\start-backend-hidden.vbs

# Start frontend
cd frontend
npm run dev
```

**Option B: PowerShell Script**
```bash
# From project root
.\run.ps1
```

**Option C: Management Tool**
```bash
# Interactive menu for start/stop/status
.\manage-backend.bat
```

**Open Browser:** 🌐 http://localhost:5173

### 🔧 Backend Management (Windows)

- **`start-backend-hidden.vbs`** - Start backend silently (no window)
- **`start-backend.bat`** - Start with visible terminal (debugging)
- **`stop-backend.bat`** - Stop the running backend
- **`manage-backend.bat`** - Interactive menu for all operations
- **`build-backend.bat`** - Rebuild executable after code changes

**Auto-start on Windows boot:**
1. Press `Win + R`, type `shell:startup`, press Enter
2. Create shortcut to `start-backend-hidden.vbs`
3. Backend will start silently on login

**Find in Task Manager:** Look for `pcmon.exe` in Details tab

---

## 📸 Screenshots

### 🎯 Main Dashboard
Beautiful, real-time hardware monitoring with gradient cards and status indicators.

### ⚙️ Customizable Ranges
Click the gear icon on any card to set custom min/max values for progress bars!

### 🔄 Refresh Settings
Configure auto-refresh intervals with quick presets (5s, 10s, 30s, 60s).

---

## 🛠️ Adding New Metrics

Want to track more metrics? It's super easy!

### Step 1: Find the Metric 🔍

Visit http://localhost:8085/data.json and find your desired metric in the JSON tree.

Example path for CCD1 temperature:
```
Children[0] → "AMD Ryzen 5 7600X" → "Temperatures" → "CCD1 (Tdie)"
```

### Step 2: Update Config 📝

Edit `backend/metrics-config.json`:

```json
{
  "metrics": [
    {
      "name": "cpu_voltage",
      "description": "CPU Core Voltage",
      "path": ["AMD Ryzen|Ryzen", "Voltages", "Core"],
      "unit": "V"
    }
  ]
}
```

### Step 3: Update Frontend Types 📦

Edit `frontend/src/lib/types/hardware.ts`:

```typescript
export interface HardwareMetrics {
  // ... existing fields
  cpu_voltage: number;
}
```

### Step 4: Add UI Card 🎨

Edit `frontend/src/routes/+page.svelte` and add your new card!

### Step 5: Restart ♻️

Restart the backend and you're done! The new metric will automatically appear.

---

## 🎨 Customization

### 🌡️ Temperature Ranges

Click the ⚙️ icon on any temperature card to set custom min/max ranges for the progress bar.

### 🔄 Refresh Intervals

Click the global ⚙️ icon in the header to configure auto-refresh timing:
- ⚡ 5 seconds (high frequency)
- 🔄 10 seconds
- ✅ 30 seconds (recommended)
- 💰 60 seconds (cost-effective)

### 🎨 Card Colors

Cards are color-coded by type:
- 🔵 CPU Temperature: Blue/Cyan gradient
- 🟢 GPU Temperature: Green/Emerald gradient
- 🔴 CPU Power: Pink gradient
- 🟠 GPU Power: Orange gradient

---

## 📊 Performance Stats

### Backend Optimizations
- ✅ Config-based extraction (no hardcoding)
- ✅ 5 retry attempts over 60 seconds
- ✅ Direct metric extraction (no double parsing)
- ✅ 97% payload reduction (5KB → 150 bytes)

### Frontend Optimizations
- ✅ Direct value access (no tree traversal)
- ✅ Smart adaptive timers (98% fewer updates)
- ✅ TypeScript type safety (zero `any` types)
- ✅ localStorage for user preferences

### Cost Efficiency
- 💰 ~173k Upstash commands/month
- ✅ Well within 500k free tier limit
- 📉 50% reduction vs original design

---

## 🐛 Troubleshooting

### ❌ Backend won't start

**Check:**
- ✅ `.env` file exists in `backend/` directory
- ✅ Upstash credentials are correct
- ✅ Hardware monitor is running on port 8085

**Test:** Visit http://localhost:8085/data.json in your browser

### ❌ Frontend shows 0.0

**Check:**
- ✅ Metric path in `metrics-config.json` matches exactly
- ✅ Backend console for extraction errors
- ✅ Use `|` separator for alternative names

### ⚠️ Stale Data Warning

**Check:**
- ✅ Backend is running and uploading
- ✅ Network connectivity to Upstash
- ✅ Redis credentials are correct

### 💸 High API Usage

**Solutions:**
- ⏱️ Increase refresh interval (30s → 60s)
- 🔍 Check for multiple backend instances
- ✅ Verify using `setTimeout` (not `setInterval`)

---

## 🛠️ Tech Stack

### Backend
- 🐹 **Go 1.21+** - Fast, efficient, concurrent
- ☁️ **Upstash Redis** - Serverless Redis with REST API
- 🔐 **godotenv** - Environment configuration
- 🔄 **go-redis/v9** - Redis client library

### Frontend
- ⚡ **Svelte 5** - Modern reactive framework with runes
- 🎨 **SvelteKit** - Full-stack framework with static adapter
- 🎯 **TypeScript** - Full type safety
- 💅 **Tailwind CSS** - Utility-first styling
- 📲 **Progressive Web App (PWA)** - Installable, offline-capable
- 🔧 **Service Worker** - Network-first caching strategy
- 🎨 **Sharp** - Icon generation from SVG

---

## 📁 Project Structure

```
📦 PC-Hardware-Monitor/
├── 📂 backend/
│   ├── 🔧 main.go                  # Go backend application
│   ├── ⚙️ metrics-config.json      # Metric extraction config
│   ├── 🔐 .env                     # Environment variables
│   ├── 📝 go.mod                   # Go dependencies
│   └── 📝 go.sum                   # Dependency checksums
│
├── 📂 frontend/
│   ├── 📂 src/
│   │   ├── 📂 routes/
│   │   │   ├── 🏠 +page.svelte     # Main dashboard
│   │   │   └── +layout.svelte      # Service worker registration
│   │   ├── 📂 lib/
│   │   │   ├── 📂 api/
│   │   │   │   └── upstash.ts      # API client
│   │   │   ├── 📂 stores/
│   │   │   │   └── hardware.svelte.ts  # State store
│   │   │   ├── 📂 components/      # Reusable UI components
│   │   │   └── 📂 types/
│   │   │       └── hardware.ts     # TypeScript interfaces
│   │   └── app.html                # PWA metadata
│   ├── 📂 static/
│   │   ├── 🔧 service-worker.js    # PWA service worker
│   │   ├── 📋 manifest.json        # PWA manifest
│   │   ├── 🎨 icon.svg             # Base icon source
│   │   ├── 🖼️ icon-192.png         # PWA icon 192x192
│   │   ├── 🖼️ icon-512.png         # PWA icon 512x512
│   │   └── 🖼️ favicon.png          # Browser favicon
│   ├── 📝 generate-icons.js        # Icon generation script
│   ├── 📦 package.json
│   ├── ⚙️ svelte.config.js
│   └── 🎨 tailwind.config.js
│
├── 🚀 run.ps1                      # PowerShell launcher
├── 📖 README.md                    # This file!
├── 📋 CLAUDE.md                    # Comprehensive dev guide
└── 🙈 .gitignore
```

---

## 🎯 Roadmap

- [x] 📲 **Progressive Web App (PWA)** - Installable, works offline
- [ ] 📊 Historical data graphs
- [ ] 🔔 Alert thresholds with notifications
- [ ] 🖥️ Multiple PC monitoring
- [ ] 🌓 Dark/light theme toggle
- [ ] 📤 Export data to CSV
- [ ] 🧠 RAM usage monitoring
- [ ] 💾 Disk activity tracking
- [ ] 🌐 Network traffic monitoring
- [ ] 🔔 Push notifications for critical temperatures

---

## 🤝 Contributing

Contributions are welcome! Feel free to:
- 🐛 Report bugs
- 💡 Suggest features
- 🔧 Submit pull requests
- ⭐ Star the repository

---

## 📄 License

This project is open source and available under the MIT License.

---

## 💬 Support

Having issues? Check out:
- 📖 [CLAUDE.md](CLAUDE.md) - Comprehensive development guide
- 🐛 GitHub Issues
- 💡 Discussion board

---

## 🎉 Acknowledgments

- 🔥 [LibreHardwareMonitor](https://github.com/LibreHardwareMonitor/LibreHardwareMonitor) - Hardware monitoring
- ☁️ [Upstash](https://upstash.com/) - Serverless Redis
- ⚡ [Svelte](https://svelte.dev/) - Reactive framework
- 🐹 [Go](https://go.dev/) - Backend language

---

<div align="center">

### 🌟 Star this repo if you find it useful! 🌟

Made with ❤️ and ☕

**[⬆ Back to Top](#-pc-hardware-monitor)**

</div>
