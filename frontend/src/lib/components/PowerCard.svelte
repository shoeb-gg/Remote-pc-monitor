<script lang="ts">
	import { onMount } from 'svelte';
	import RangeSettingsModal from './RangeSettingsModal.svelte';

	interface Props {
		title: string;
		power: number | undefined;
		type: 'cpu' | 'gpu';
		loading?: boolean;
		error?: boolean;
	}

	let { title, power, type, loading = false, error = false }: Props = $props();

	const gradientClass = type === 'cpu'
		? 'from-purple-500 to-pink-600'
		: 'from-amber-500 to-orange-600';

	// Default ranges
	const defaultMin = 0;
	const defaultMax = type === 'cpu' ? 120 : 250; // CPU max ~120W, GPU max ~250W

	// Reactive state for ranges
	let minPower = $state(defaultMin);
	let maxPower = $state(defaultMax);
	let isModalOpen = $state(false);

	// Load saved ranges from localStorage on mount
	onMount(() => {
		const saved = localStorage.getItem(`powerRange_${type}`);
		if (saved) {
			const { min, max } = JSON.parse(saved);
			minPower = min;
			maxPower = max;
		}
	});

	const getPowerColor = (watts: number | undefined) => {
		if (!watts) return 'text-gray-400';
		if (watts < 50) return 'text-green-400';
		if (watts < 100) return 'text-yellow-400';
		if (watts < 150) return 'text-orange-400';
		return 'text-red-400';
	};

	const getStatusDotColor = () => {
		if (loading) return 'bg-yellow-400';
		if (error || power === undefined) return 'bg-red-400';
		return 'bg-green-400';
	};

	const openModal = () => {
		isModalOpen = true;
	};

	const closeModal = () => {
		isModalOpen = false;
	};

	const saveRange = (min: number, max: number) => {
		minPower = min;
		maxPower = max;
		localStorage.setItem(`powerRange_${type}`, JSON.stringify({ min, max }));
	};

	const getProgressPercentage = () => {
		if (power === undefined) return 0;
		const range = maxPower - minPower;
		const value = power - minPower;
		return Math.min(Math.max((value / range) * 100, 0), 100);
	};
</script>

<div class="relative overflow-hidden rounded-2xl bg-gradient-to-br {gradientClass} p-1">
	<div class="rounded-xl bg-gray-900 p-6 backdrop-blur-sm">
		<div class="flex items-center justify-between mb-4">
			<h2 class="text-lg font-semibold text-white">{title}</h2>
			<div class="flex items-center space-x-2">
				<button
					onclick={openModal}
					class="p-2 hover:bg-white/10 rounded-lg transition-colors duration-200"
					aria-label="Settings"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white/70 hover:text-white" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
					</svg>
				</button>
				<div class="{getStatusDotColor()} w-3 h-3 rounded-full animate-pulse"></div>
			</div>
		</div>

		<div class="flex items-baseline space-x-2">
			<span class="text-5xl font-bold {getPowerColor(power)} transition-colors duration-300">
				{power !== undefined ? power.toFixed(1) : '--'}
			</span>
			<span class="text-2xl text-white/60">W</span>
		</div>

		{#if power !== undefined}
			<div class="mt-4 h-2 bg-white/10 rounded-full overflow-hidden">
				<div
					class="h-full bg-gradient-to-r from-white/50 to-white/80 transition-all duration-500 ease-out"
					style="width: {getProgressPercentage()}%"
				></div>
			</div>
		{/if}
	</div>
</div>

<RangeSettingsModal
	isOpen={isModalOpen}
	title={title}
	currentMin={minPower}
	currentMax={maxPower}
	onClose={closeModal}
	onSave={saveRange}
/>
