import { PUBLIC_UPSTASH_REDIS_URL, PUBLIC_UPSTASH_REDIS_TOKEN } from '$env/static/public';

export interface HardwareMetrics {
	cpuTemp?: number;
	gpuTemp?: number;
	timestamp?: number;
	[key: string]: any;
}

export async function fetchLatestMetrics(): Promise<HardwareMetrics | null> {
	try {
		const url = `${PUBLIC_UPSTASH_REDIS_URL}/xrevrange/hardware:metrics/+/-/COUNT/1`;
		console.log('Fetching from:', url);

		const response = await fetch(url, {
			headers: {
				'Authorization': `Bearer ${PUBLIC_UPSTASH_REDIS_TOKEN}`
			}
		});

		console.log('Response status:', response.status);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		console.log('Raw response:', data);

		// Upstash XREVRANGE returns: {"result": [["entry-id", ["data", "{json}"]]]}
		if (data.result && data.result.length > 0) {
			const entry = data.result[0];
			if (entry && entry[1] && entry[1][1]) {
				const jsonData = JSON.parse(entry[1][1]);
				console.log('Parsed hardware data:', jsonData);
				return jsonData;
			}
		}

		console.log('No data found in response');
		return null;
	} catch (error) {
		console.error('Failed to fetch metrics:', error);
		return null;
	}
}
