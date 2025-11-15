# CLAUDE.md

This file provides comprehensive guidance to Claude Code when working with this PC hardware monitoring system.

---

## Project Overview

This is a **real-time PC hardware monitoring system** that tracks CPU and GPU temperatures and power consumption with a clean, modern web interface.

### Architecture Components

```
┌─────────────────────────────────────────────────────────────┐
│  Hardware Monitor (Port 8085)                               │
│  ├─ LibreHardwareMonitor or similar                         │
│  └─ Exposes JSON data via HTTP                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (Go)                                                │
│  ├─ Polls hardware monitor every 10s                        │
│  ├─ Extracts metrics via config-based system                │
│  └─ Publishes to Upstash Redis                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  Upstash Redis                                               │
│  └─ Stores latest metrics (simple SET/GET)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  Frontend (SvelteKit PWA)                                    │
│  ├─ Fetches from Upstash REST API every 30s                 │
│  ├─ Displays metrics with customizable ranges               │
│  └─ Real-time status indicators                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Backend (Go)

**Location:** `backend/`

### Key Files

- **`main.go`**: Core application logic
- **`metrics-config.json`**: Configuration-based metric extraction
- **`.env`**: Environment variables (Redis credentials, hardware monitor URL)
- **`go.mod`**: Go dependencies

### How It Works

1. **Environment Setup**: Loads `.env` using `godotenv`
2. **Config Loading**: Reads `metrics-config.json` to determine which metrics to extract
3. **Data Collection**: Polls hardware monitor HTTP endpoint (default: `http://localhost:8085/data.json`)
4. **Metric Extraction**: Uses JSON path navigation to extract specific values:
   - `cpu_temp_tctl`: CPU Tctl/Tdie temperature (overall CPU temp)
   - `cpu_temp_ccd1`: CCD1 (Core Complex Die 1) temperature (actual core temp)
   - `cpu_power`: CPU package power consumption
   - `gpu_temp`: GPU core temperature
   - `gpu_power`: GPU package power consumption
5. **Redis Storage**: Stores simplified JSON in Upstash Redis using `SET` command (no expiry)
6. **Retry Logic**: 5 retry attempts over ~60 seconds on failures
7. **Loop**: Repeats every 10 seconds

### Configuration System

**`metrics-config.json`** format:
```json
{
  "metrics": [
    {
      "name": "metric_name",           // JSON key in output
      "description": "Human readable", // For documentation
      "path": ["Node1|Alt1", "Node2", "Node3"],  // Navigation path (| = OR)
      "unit": "°C"                     // Unit to strip from value
    }
  ]
}
```

**Path Navigation:**
- Searches hardware tree using pattern matching
- `|` separator allows multiple alternatives (e.g., `"AMD Ryzen|Ryzen"` matches either)
- Automatically extracts float values from strings like `"72.3 °C"` → `72.3`

### Output Format

Backend sends this simplified JSON to Redis:
```json
{
  "cpu_temp_tctl": 65.4,
  "cpu_temp_ccd1": 75.6,
  "cpu_power": 59.1,
  "gpu_temp": 53.0,
  "gpu_power": 173.8,
  "pc_name": "SHOEBRIG",
  "timestamp": 1731664320
}
```

### Dependencies

- `github.com/joho/godotenv`: Load `.env` files
- `github.com/redis/go-redis/v9`: Redis client (updated from deprecated v8)

### Running the Backend

```bash
# From project root
.\run.ps1

# Or manually from backend directory
cd backend
go run main.go

# Build executable
go build -o hardware-monitor.exe main.go
```

### Environment Variables

**Required:**
- `UPSTASH_REDIS_ADDR`: Redis host:port (e.g., `example.upstash.io:6379`)
- `UPSTASH_REDIS_PASSWORD`: Redis password/token

**Optional:**
- `HARDWARE_MONITOR_URL`: Defaults to `http://localhost:8085/data.json`

---

## Frontend (SvelteKit)

**Location:** `frontend/`

### Technology Stack

- **Svelte 5**: Using new runes API (`$state`, `$props`, `$effect`)
- **SvelteKit**: Static adapter for PWA deployment
- **Progressive Web App (PWA)**: Installable, offline-capable
- **Service Worker**: Network-first caching strategy
- **Tailwind CSS**: Utility-first styling
- **TypeScript**: Full type safety (no `any` types)
- **Upstash REST API**: Direct fetch to Redis via HTTP
- **Sharp**: SVG to PNG icon generation

### Project Structure

```
frontend/
├── src/
│   ├── routes/
│   │   ├── +page.svelte         # Main dashboard page
│   │   └── +layout.svelte       # Root layout + SW registration
│   ├── lib/
│   │   ├── api/
│   │   │   └── upstash.ts       # Redis API client (GET endpoint)
│   │   ├── stores/
│   │   │   └── hardware.svelte.ts  # Reactive store for metrics
│   │   ├── components/
│   │   │   ├── TemperatureCard.svelte
│   │   │   ├── PowerCard.svelte
│   │   │   ├── StatusIndicator.svelte
│   │   │   ├── RangeSettingsModal.svelte
│   │   │   └── RefreshSettingsModal.svelte
│   │   └── types/
│   │       └── hardware.ts      # TypeScript interfaces
│   └── app.html                 # PWA metadata
├── static/
│   ├── service-worker.js        # PWA service worker
│   ├── manifest.json            # PWA manifest
│   ├── icon.svg                 # Base icon (CPU design)
│   ├── icon-192.png             # PWA icon 192x192
│   ├── icon-512.png             # PWA icon 512x512
│   └── favicon.png              # Browser favicon 32x32
├── generate-icons.js            # Sharp script for icon generation
├── svelte.config.js
├── tailwind.config.js
├── package.json
└── .env
```

### Data Flow

1. **Initial Load**: `+page.svelte` calls `hardwareStore.refresh()` on mount
2. **API Fetch**: `upstash.ts` sends GET request to Upstash REST API
3. **Response Parsing**: `{"result": "{json}"}` → parsed to `HardwareMetrics`
4. **Store Update**: `hardwareStore` updates reactive state
5. **UI Update**: Components reactively display new values
6. **Auto-refresh**: `setTimeout` schedules next refresh (default 30s)

### Key Features

#### 1. Configuration-Based Metrics
- No hardcoded tree traversal
- Backend handles all parsing
- Frontend just displays pre-extracted values

#### 2. Customizable Ranges
- Each card has settings gear icon
- Min/max values for progress bars
- Stored in `localStorage` per metric type
- Enter key support for quick updates

#### 3. Smart Status Indicator
- **Green**: Connected, data fresh
- **Yellow**: Loading/updating
- **Orange**: Data stale (>2 minutes old)
- **Red**: Error occurred
- Live countdown: "Updated Xs ago" with optimized refresh intervals

#### 4. Adaptive Timer Optimization
- **0-60s**: Updates every 1 second
- **1-60min**: Updates every 60 seconds
- **>1 hour**: Updates every hour
- **Result**: 98% fewer re-renders after first minute

#### 5. Manual Refresh
- Button with spinning animation
- Restarts auto-refresh cycle
- Prevents spam with loading state

#### 6. Progressive Web App (PWA)
- **Installable**: Can be installed on desktop and mobile devices
- **Offline Support**: Service worker caches assets for offline access
- **Network-First Strategy**: Always tries network first, falls back to cache
- **Manifest**: Full PWA manifest with app metadata and icons
- **Icons**: Custom CPU-themed icons (192x192, 512x512, favicon)
- **Service Worker**: Auto-registers on app load via `+layout.svelte`

### TypeScript Interfaces

```typescript
// hardware.ts
export interface HardwareMetrics {
  cpu_temp_tctl: number;
  cpu_temp_ccd1: number;
  cpu_power: number;
  gpu_temp: number;
  gpu_power: number;
  pc_name: string;
  timestamp: number;
}
```

### State Management

**Svelte 5 Runes Pattern:**
```typescript
// Store (hardware.svelte.ts)
class HardwareStore {
  metrics = $state<HardwareMetrics | null>(null);
  loading = $state(true);
  error = $state<string | null>(null);
  lastUpdate = $state<Date | null>(null);
}

// Component usage
let { title, temperature }: Props = $props();
let isModalOpen = $state(false);

// Effects
$effect(() => {
  if (lastUpdate) {
    scheduleNextUpdate();
  }
});
```

### Styling System

- **Dark Theme**: `bg-gray-900` base
- **Gradient Borders**: Nested divs with `bg-gradient-to-br` + inner `bg-gray-900`
- **Color Coding**:
  - CPU cards: Blue/Cyan gradients
  - GPU cards: Green/Emerald gradients
  - Power cards: Pink/Orange gradients
- **Responsive**: 1 column mobile, 3 columns desktop (`md:grid-cols-3`)

### Running the Frontend

```bash
cd frontend

# Install dependencies
npm install

# Development server (http://localhost:5173)
npm run dev

# Type checking
npm run check

# Build for production
npm run build

# Preview production build
npm run preview

# Generate PWA icons from SVG (if icon.svg is modified)
npm run generate-icons
```

### Environment Variables

```bash
PUBLIC_UPSTASH_REDIS_URL=https://your-instance.upstash.io
PUBLIC_UPSTASH_REDIS_TOKEN=your_rest_api_token
```

**Note:** Uses Upstash REST API (not direct Redis connection)

---

## Performance Optimizations

### Backend
- ✅ **No double JSON parsing**: Direct metric extraction
- ✅ **97% smaller payloads**: 150 bytes vs 5KB
- ✅ **Retry logic**: 5 attempts over 1 minute
- ✅ **Config-based extraction**: No hardcoded paths

### Frontend
- ✅ **Direct value access**: No tree traversal
- ✅ **Smart timer**: 98% fewer updates after 1 minute
- ✅ **setTimeout vs setInterval**: Waits for response before next refresh
- ✅ **Type safety**: All TypeScript, no `any` types
- ✅ **localStorage caching**: Persistent user preferences
- ✅ **PWA optimizations**: Service worker caching, offline support
- ✅ **Static adapter**: Pre-rendered SPA for maximum performance

### Redis Cost Optimization
- **Backend**: Simple `SET` (no streams overhead)
- **Frontend**: Simple `GET` (no complex queries)
- **Monthly Usage**: ~173k commands (within 500k free tier)
- **Payload Size**: 97% reduction saves bandwidth

---

## Adding New Metrics

### Step 1: Find the Metric Path

1. Check hardware monitor at `http://localhost:8085/data.json`
2. Navigate JSON tree to find desired metric
3. Note the exact text values along the path

Example: To find CCD1 temp
```
Children[0] (Computer) →
  Children[] → Find "AMD Ryzen 5 7600X" →
    Children[] → Find "Temperatures" →
      Children[] → Find "CCD1 (Tdie)" →
        Value: "75.6 °C"
```

### Step 2: Update `metrics-config.json`

```json
{
  "name": "new_metric_name",
  "description": "Human readable description",
  "path": ["Hardware|AltName", "Section", "Metric Name"],
  "unit": "°C"
}
```

### Step 3: Update Frontend Types

```typescript
// frontend/src/lib/types/hardware.ts
export interface HardwareMetrics {
  // ... existing fields
  new_metric_name: number;
}
```

### Step 4: Update Frontend UI

```svelte
<!-- frontend/src/routes/+page.svelte -->
<TemperatureCard
  title="New Metric"
  temperature={hardwareStore.metrics?.new_metric_name}
  type="cpu"
  loading={hardwareStore.loading}
  error={!!hardwareStore.error}
/>
```

### Step 5: Restart Backend

The backend will automatically load the new config and start extracting the metric!

---

## Common Issues & Solutions

### Backend won't start
- **Check `.env` file exists** in `backend/` directory
- **Verify credentials** are set correctly
- **Test hardware monitor**: Visit `http://localhost:8085/data.json` in browser

### Frontend shows 0.0 for a metric
- **Check metric path** in `metrics-config.json`
- **Verify exact text match** (case-sensitive, include spaces)
- **Look at backend logs** for extraction errors
- **Use pattern alternatives** with `|` separator

### Stale data warning
- **Check backend is running** and uploading
- **Verify Redis credentials** are correct
- **Check network connectivity** to Upstash
- **Look for errors** in backend console

### High API usage
- **Increase refresh interval** (30s → 60s in frontend)
- **Check for multiple instances** running
- **Verify auto-refresh** is using setTimeout (not setInterval)

---

## Development Workflow

### Making Changes

1. **Backend changes**: Edit `main.go` or `metrics-config.json`, restart backend
2. **Frontend changes**: SvelteKit auto-reloads in dev mode
3. **Type changes**: Update both `backend/main.go` structs and `frontend/src/lib/types/hardware.ts`

### Testing

```bash
# Backend
cd backend
go run main.go  # Watch console output for errors

# Frontend
cd frontend
npm run dev     # Check browser console
npm run check   # TypeScript validation
```

### Production Deployment

```bash
# Backend
cd backend
go build -o hardware-monitor.exe main.go

# Frontend
cd frontend
npm run build
# Deploy `build/` folder to static hosting (Vercel, Netlify, etc.)
```

---

## Progressive Web App (PWA) Setup

### Overview

The frontend is a fully functional PWA that can be installed on desktop and mobile devices, with offline support via service worker.

### PWA Components

**1. Web App Manifest** (`static/manifest.json`)
```json
{
  "name": "PC Hardware Monitor",
  "short_name": "PC Monitor",
  "display": "standalone",
  "theme_color": "#111827",
  "icons": [...]
}
```

**2. Service Worker** (`static/service-worker.js`)
- **Strategy**: Network-first with cache fallback
- **Cache Name**: `pc-monitor-v1`
- **Cached Assets**: Root page, manifest
- **Behavior**:
  - Tries network first for fresh data
  - Falls back to cache if network fails
  - Skips cross-origin requests (CORS-safe)
  - Auto-updates cache with successful responses

**3. Service Worker Registration** (`src/routes/+layout.svelte`)
```typescript
onMount(() => {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/service-worker.js');
  }
});
```

**4. App Icons**
- `icon.svg` - Base SVG design (blue CPU chip with red temperature indicator)
- `icon-192.png` - 192x192 PNG for PWA
- `icon-512.png` - 512x512 PNG for PWA
- `favicon.png` - 32x32 PNG for browser tab

**5. PWA Metadata** (`src/app.html`)
- Manifest link
- Theme color meta tag
- Apple touch icon
- Favicon link

### Generating Custom Icons

If you want to customize the app icon:

1. Edit `frontend/static/icon.svg` with your design
2. Run icon generation script:
```bash
cd frontend
npm run generate-icons
```
3. This creates all required PNG sizes using Sharp

### Testing PWA Functionality

**Chrome/Edge DevTools:**
1. Open DevTools → Application tab
2. Check Manifest section (should show app info)
3. Check Service Workers (should show registered worker)
4. Install button should appear in address bar

**Installation:**
- Desktop: Click install button in browser
- Mobile: Add to Home Screen from browser menu

**Offline Testing:**
1. Load the app once while online
2. Open DevTools → Network → Enable "Offline"
3. Refresh page - should still load from cache

---

## Security Notes

- ✅ **No hardcoded credentials**: All in `.env` files
- ✅ **Environment validation**: Backend fails fast if vars missing
- ✅ **TLS for Redis**: Encrypted connection to Upstash
- ✅ **No secrets in frontend**: Only public REST API tokens
- ⚠️ **Local hardware monitor**: Not exposed to internet (localhost only)

---

## Code Style & Conventions

### Go
- Config-based over hardcoded logic
- Clear error messages with context
- Use structs for type safety
- Pattern matching for flexibility

### TypeScript/Svelte
- **No `any` types** - always use proper interfaces
- Svelte 5 runes (`$state`, `$props`, `$effect`)
- Component composition over monolithic components
- Tailwind utility classes over custom CSS

---

## Future Enhancement Ideas

- [x] **Progressive Web App** - Fully implemented!
- [ ] Add more metrics (RAM, disk, network)
- [ ] Historical data graphs
- [ ] Alert thresholds with notifications
- [ ] Multiple PC monitoring
- [ ] Push notifications for critical temperatures
- [ ] Dark/light theme toggle
- [ ] Export data to CSV

---

## Quick Reference

**Start Everything:**
```bash
# Terminal 1 - Backend
.\run.ps1

# Terminal 2 - Frontend
cd frontend && npm run dev
```

**Update Metrics:**
1. Edit `backend/metrics-config.json`
2. Restart backend
3. Done! Frontend auto-updates

**Troubleshoot:**
- Backend logs: Console output
- Frontend logs: Browser DevTools
- Redis data: Upstash dashboard
- Hardware data: http://localhost:8085/data.json
