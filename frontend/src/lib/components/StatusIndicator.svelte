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
	let updateTimer: ReturnType<typeof setTimeout>;

	const scheduleNextUpdate = () => {
		if (updateTimer) {
			clearTimeout(updateTimer);
		}

		if (!lastUpdate) {
			// No data yet, check again in 1 second
			updateTimer = setTimeout(() => {
				currentTime = new Date();
				scheduleNextUpdate();
			}, 1000);
			return;
		}

		const diff = Math.floor((Date.now() - lastUpdate.getTime()) / 1000);
		let nextUpdateIn: number;

		if (diff < 60) {
			// Under 1 minute: update every 1 second (shows "5s ago", "6s ago"...)
			nextUpdateIn = 1000;
		} else if (diff < 3600) {
			// 1-60 minutes: update every 60 seconds (shows "5m ago", "6m ago"...)
			nextUpdateIn = 60000;
		} else {
			// Over 1 hour: update every hour (shows "2h ago", "3h ago"...)
			nextUpdateIn = 3600000;
		}

		updateTimer = setTimeout(() => {
			currentTime = new Date();
			scheduleNextUpdate();
		}, nextUpdateIn);
	};

	onMount(() => {
		scheduleNextUpdate();
	});

	onDestroy(() => {
		if (updateTimer) {
			clearTimeout(updateTimer);
		}
	});

	// Reschedule when lastUpdate changes
	$effect(() => {
		if (lastUpdate) {
			scheduleNextUpdate();
		}
	});

	const isDataStale = () => {
		if (!lastUpdate) return false;
		const diff = Math.floor((currentTime.getTime() - lastUpdate.getTime()) / 1000);
		// Consider data stale if older than 2 minutes (120 seconds)
		return diff > 120;
	};

	const getStatusColor = () => {
		if (loading) return 'bg-yellow-400';
		if (error) return 'bg-red-400';
		if (isDataStale()) return 'bg-orange-400';
		return 'bg-green-400';
	};

	const getStatusText = () => {
		if (loading) return 'Updating...';
		if (error) return `Error: ${error}`;
		if (isDataStale()) return 'Data may be stale';
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
		const hours = Math.floor(seconds / 3600);
		if (hours < 24) return `${hours}h ago`;
		return date.toLocaleString();
	};
</script>

<div class="flex items-center space-x-3 text-base">
	<div class="flex items-center space-x-2">
		<div class="{getStatusColor()} w-3 h-3 rounded-full animate-pulse"></div>
		<span class="text-white/80">{getStatusText()}</span>
	</div>

	{#if lastUpdate}
		<span class="text-white/50">•</span>
		<span class="{isDataStale() ? 'text-orange-400' : 'text-white/50'}">
			Last data: {formatLastUpdate(lastUpdate)}
		</span>
	{/if}
</div>
