import { PUBLIC_UPSTASH_REDIS_URL, PUBLIC_UPSTASH_REDIS_TOKEN } from '$env/static/public';
import type { HardwareMetrics } from '$lib/types/hardware';

export async function fetchLatestMetrics(): Promise<HardwareMetrics | null> {
	try {
		const url = `${PUBLIC_UPSTASH_REDIS_URL}/get/hardware:metrics`;

		const response = await fetch(url, {
			headers: {
				'Authorization': `Bearer ${PUBLIC_UPSTASH_REDIS_TOKEN}`
			},
			// Abort a hung request so the polling loop can't deadlock with
			// loading stuck true (which also disables the manual refresh button).
			signal: AbortSignal.timeout(10000)
		});

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();

		// Upstash GET returns: {"result": "{json}"}
		if (data.result) {
			const metrics: HardwareMetrics = JSON.parse(data.result);
			return metrics;
		}

		return null;
	} catch (error) {
		console.error('Failed to fetch metrics:', error);
		return null;
	}
}
