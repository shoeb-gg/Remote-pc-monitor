<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { hardwareStore } from '$lib/stores/hardware.svelte';
	import TemperatureCard from '$lib/components/TemperatureCard.svelte';
	import PowerCard from '$lib/components/PowerCard.svelte';
	import StatusIndicator from '$lib/components/StatusIndicator.svelte';
	import RefreshSettingsModal from '$lib/components/RefreshSettingsModal.svelte';

	let refreshTimeout: ReturnType<typeof setTimeout>;
	let refreshIntervalSeconds = $state(30); // Default 30 seconds
	let isRefreshSettingsOpen = $state(false);
	let isAutoRefreshActive = $state(true);

	// Load refresh interval from localStorage on mount
	onMount(async () => {
		const saved = localStorage.getItem('refreshInterval');
		if (saved) {
			refreshIntervalSeconds = parseInt(saved, 10);
		}

		// Initial fetch and start the refresh cycle
		await refreshAndScheduleNext();
	});

	onDestroy(() => {
		isAutoRefreshActive = false;
		if (refreshTimeout) {
			clearTimeout(refreshTimeout);
		}
	});

	const refreshAndScheduleNext = async () => {
		// Fetch data
		await hardwareStore.refresh();

		// Schedule next refresh only if auto-refresh is still active
		if (isAutoRefreshActive) {
			console.log(`Next refresh scheduled in ${refreshIntervalSeconds} seconds`);
			refreshTimeout = setTimeout(refreshAndScheduleNext, refreshIntervalSeconds * 1000);
		}
	};

	const restartRefreshCycle = () => {
		// Clear existing timeout
		if (refreshTimeout) {
			clearTimeout(refreshTimeout);
		}
		// Start new cycle
		refreshAndScheduleNext();
	};

	const openRefreshSettings = () => {
		isRefreshSettingsOpen = true;
	};

	const closeRefreshSettings = () => {
		isRefreshSettingsOpen = false;
	};

	const saveRefreshInterval = (interval: number) => {
		refreshIntervalSeconds = interval;
		localStorage.setItem('refreshInterval', interval.toString());
		restartRefreshCycle();
	};

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
			<div class="flex items-center justify-between mb-2">
				<h1 class="text-4xl font-bold text-white">PC Hardware Monitor</h1>
				<div class="flex items-center space-x-2">
					<button
						onclick={restartRefreshCycle}
						class="p-3 hover:bg-white/10 rounded-lg transition-colors duration-200"
						aria-label="Refresh Now"
						disabled={hardwareStore.loading}
					>
						<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white/70 hover:text-white {hardwareStore.loading ? 'animate-spin' : ''}" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
						</svg>
					</button>
					<button
						onclick={openRefreshSettings}
						class="p-3 hover:bg-white/10 rounded-lg transition-colors duration-200"
						aria-label="Refresh Settings"
					>
						<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white/70 hover:text-white" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
						</svg>
					</button>
				</div>
			</div>
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

<RefreshSettingsModal
	isOpen={isRefreshSettingsOpen}
	currentInterval={refreshIntervalSeconds}
	onClose={closeRefreshSettings}
	onSave={saveRefreshInterval}
/>
