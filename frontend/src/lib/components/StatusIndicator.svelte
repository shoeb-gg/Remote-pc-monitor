<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		loading: boolean;
		error: string | null;
		lastUpdate: Date | null;
		pcName?: string;
	}

	let { loading, error, lastUpdate, pcName }: Props = $props();

	let currentTime = $state(new Date());
	let updateTimer: ReturnType<typeof setInterval>;

	onMount(() => {
		// Update the current time every second to refresh the "ago" display
		updateTimer = setInterval(() => {
			currentTime = new Date();
		}, 1000);
	});

	onDestroy(() => {
		if (updateTimer) {
			clearInterval(updateTimer);
		}
	});

	const getStatusColor = () => {
		if (loading) return 'bg-yellow-400';
		if (error) return 'bg-red-400';
		return 'bg-green-400';
	};

	const getStatusText = () => {
		if (loading) return 'Updating...';
		if (error) return `Error: ${error}`;
		return pcName ? `Connected to ${pcName}` : 'Connected';
	};

	const formatLastUpdate = (date: Date | null) => {
		if (!date) return 'Never';
		// Use currentTime to trigger reactivity
		const diff = Math.floor((currentTime.getTime() - date.getTime()) / 1000);

		// Ensure we never show negative values
		const seconds = Math.max(0, diff);

		if (seconds < 60) return `${seconds}s ago`;
		if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
		return date.toLocaleTimeString();
	};
</script>

<div class="flex items-center space-x-3 text-base">
	<div class="flex items-center space-x-2">
		<div class="{getStatusColor()} w-3 h-3 rounded-full animate-pulse"></div>
		<span class="text-white/80">{getStatusText()}</span>
	</div>

	{#if lastUpdate && !error}
		<span class="text-white/50">•</span>
		<span class="text-white/50">Updated {formatLastUpdate(lastUpdate)}</span>
	{/if}
</div>
