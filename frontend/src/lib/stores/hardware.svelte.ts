import { fetchLatestMetrics, type HardwareMetrics } from '$lib/api/upstash';

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
				this.lastUpdate = new Date();
			} else {
				this.error = 'No data available';
			}
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to fetch data';
			console.error('Error refreshing hardware metrics:', err);
		} finally {
			this.loading = false;
		}
	}
}

export const hardwareStore = new HardwareStore();
