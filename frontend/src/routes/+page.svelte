<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { hardwareStore } from '$lib/stores/hardware.svelte';
	import TemperatureCard from '$lib/components/TemperatureCard.svelte';
	import PowerCard from '$lib/components/PowerCard.svelte';
	import StatusIndicator from '$lib/components/StatusIndicator.svelte';

	let refreshInterval: ReturnType<typeof setInterval>;

	onMount(async () => {
		// Initial fetch
		await hardwareStore.refresh();

		// Set up auto-refresh every 5 seconds
		refreshInterval = setInterval(() => {
			hardwareStore.refresh();
		}, 5000);
	});

	onDestroy(() => {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	});

	// Extract CPU and GPU temps from the metrics
	function getCPUTemp(metrics: any): number | undefined {
		if (!metrics) return undefined;

		// The root Children array contains the computer data
		const rootChildren = metrics.Children;
		if (!rootChildren || !Array.isArray(rootChildren)) return undefined;

		// Find the computer (SHOEBRIG)
		const computer = rootChildren[0];
		if (!computer?.Children) return undefined;

		// Find AMD Ryzen CPU
		const cpuData = computer.Children.find((c: any) =>
			c.Text?.includes('AMD Ryzen') || c.Text?.includes('Ryzen')
		);

		if (cpuData?.Children) {
			// Find Temperatures section
			const tempSection = cpuData.Children.find((c: any) =>
				c.Text === 'Temperatures'
			);

			if (tempSection?.Children) {
				// Get Core (Tctl/Tdie) temperature
				const coreTemp = tempSection.Children.find((c: any) =>
					c.Text?.includes('Core') && c.Text?.includes('Tctl')
				);

				if (coreTemp?.Value) {
					// Parse "72.3 °C" to 72.3
					return parseFloat(coreTemp.Value.replace(' °C', ''));
				}
			}
		}
		return undefined;
	}

	function getGPUTemp(metrics: any): number | undefined {
		if (!metrics) return undefined;

		// The root Children array contains the computer data
		const rootChildren = metrics.Children;
		if (!rootChildren || !Array.isArray(rootChildren)) return undefined;

		// Find the computer (SHOEBRIG)
		const computer = rootChildren[0];
		if (!computer?.Children) return undefined;

		// Find NVIDIA GPU
		const gpuData = computer.Children.find((c: any) =>
			c.Text?.includes('NVIDIA') || c.Text?.includes('GeForce')
		);

		if (gpuData?.Children) {
			// Find Temperatures section
			const tempSection = gpuData.Children.find((c: any) =>
				c.Text === 'Temperatures'
			);

			if (tempSection?.Children) {
				// Get GPU Core temperature
				const coreTemp = tempSection.Children.find((c: any) =>
					c.Text === 'GPU Core'
				);

				if (coreTemp?.Value) {
					// Parse "51.0 °C" to 51.0
					return parseFloat(coreTemp.Value.replace(' °C', ''));
				}
			}
		}
		return undefined;
	}

	function getPCName(metrics: any): string | undefined {
		if (!metrics?.Children || !Array.isArray(metrics.Children)) return undefined;
		const computer = metrics.Children[0];
		return computer?.Text;
	}

	function getCPUPower(metrics: any): number | undefined {
		if (!metrics) return undefined;

		const rootChildren = metrics.Children;
		if (!rootChildren || !Array.isArray(rootChildren)) return undefined;

		const computer = rootChildren[0];
		if (!computer?.Children) return undefined;

		// Find AMD Ryzen CPU
		const cpuData = computer.Children.find((c: any) =>
			c.Text?.includes('AMD Ryzen') || c.Text?.includes('Ryzen')
		);

		if (cpuData?.Children) {
			// Find Powers section
			const powerSection = cpuData.Children.find((c: any) =>
				c.Text === 'Powers'
			);

			if (powerSection?.Children) {
				// Get Package power
				const packagePower = powerSection.Children.find((c: any) =>
					c.Text === 'Package'
				);

				if (packagePower?.Value) {
					// Parse "65.7 W" to 65.7
					return parseFloat(packagePower.Value.replace(' W', ''));
				}
			}
		}
		return undefined;
	}

	function getGPUPower(metrics: any): number | undefined {
		if (!metrics) return undefined;

		const rootChildren = metrics.Children;
		if (!rootChildren || !Array.isArray(rootChildren)) return undefined;

		const computer = rootChildren[0];
		if (!computer?.Children) return undefined;

		// Find NVIDIA GPU
		const gpuData = computer.Children.find((c: any) =>
			c.Text?.includes('NVIDIA') || c.Text?.includes('GeForce')
		);

		if (gpuData?.Children) {
			// Find Powers section
			const powerSection = gpuData.Children.find((c: any) =>
				c.Text === 'Powers'
			);

			if (powerSection?.Children) {
				// Get GPU Package power
				const gpuPower = powerSection.Children.find((c: any) =>
					c.Text === 'GPU Package'
				);

				if (gpuPower?.Value) {
					// Parse "179.0 W" to 179.0
					return parseFloat(gpuPower.Value.replace(' W', ''));
				}
			}
		}
		return undefined;
	}
</script>

<svelte:head>
	<title>PC Hardware Monitor</title>
	<meta name="description" content="Real-time PC hardware monitoring" />
</svelte:head>

<div class="min-h-screen bg-gray-900 p-4 md:p-8">
	<div class="max-w-4xl mx-auto">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-4xl font-bold text-white mb-2">PC Hardware Monitor</h1>
			<StatusIndicator
				loading={hardwareStore.loading}
				error={hardwareStore.error}
				lastUpdate={hardwareStore.lastUpdate}
				pcName={getPCName(hardwareStore.metrics)}
			/>
		</div>

		<!-- Temperature Cards -->
		<div class="mb-6">
			<h2 class="text-xl font-semibold text-white mb-4">Temperature</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<TemperatureCard
					title="CPU Temperature"
					temperature={getCPUTemp(hardwareStore.metrics)}
					type="cpu"
					loading={hardwareStore.loading}
					error={!!hardwareStore.error}
				/>

				<TemperatureCard
					title="GPU Temperature"
					temperature={getGPUTemp(hardwareStore.metrics)}
					type="gpu"
					loading={hardwareStore.loading}
					error={!!hardwareStore.error}
				/>
			</div>
		</div>

		<!-- Power Cards -->
		<div class="mb-6">
			<h2 class="text-xl font-semibold text-white mb-4">Power Consumption</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<PowerCard
					title="CPU Power"
					power={getCPUPower(hardwareStore.metrics)}
					type="cpu"
					loading={hardwareStore.loading}
					error={!!hardwareStore.error}
				/>

				<PowerCard
					title="GPU Power"
					power={getGPUPower(hardwareStore.metrics)}
					type="gpu"
					loading={hardwareStore.loading}
					error={!!hardwareStore.error}
				/>
			</div>
		</div>

		<!-- Debug Info (optional - remove in production) -->
		{#if hardwareStore.metrics}
			<details class="mt-8 p-4 bg-gray-800 rounded-lg">
				<summary class="text-white cursor-pointer">Raw Data (Debug)</summary>
				<pre class="mt-2 text-xs text-gray-300 overflow-auto">{JSON.stringify(hardwareStore.metrics, null, 2)}</pre>
			</details>
		{/if}
	</div>
</div>
