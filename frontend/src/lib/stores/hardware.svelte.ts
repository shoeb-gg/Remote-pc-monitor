import { fetchLatestMetrics } from '$lib/api/upstash';
import type { HardwareMetrics } from '$lib/types/hardware';

class HardwareStore {
	metrics = $state<HardwareMetrics | null>(null);
	loading = $state(true);
	error = $state<string | null>(null);
	lastUpdate = $state<Date | null>(null);

	async refresh() {
		try {
			this.loading = true;
			this.error = null;

			const data = await fetchLatestMetrics();

			if (data) {
				this.metrics = data;
				// Use backend's timestamp (convert Unix timestamp to Date)
				this.lastUpdate = new Date(data.timestamp * 1000);
			} else {
				this.error = 'No data available';
			}
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to fetch data';
		} finally {
			this.loading = false;
		}
	}
}

export const hardwareStore = new HardwareStore();
