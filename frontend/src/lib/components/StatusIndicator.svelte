<script lang="ts">
	interface Props {
		loading: boolean;
		error: string | null;
		lastUpdate: Date | null;
		pcName?: string;
	}

	let { loading, error, lastUpdate, pcName }: Props = $props();

	const getStatusColor = () => {
		if (loading) return 'bg-yellow-400';
		if (error) return 'bg-red-400';
		return 'bg-green-400';
	};

	const getStatusText = () => {
		if (loading) return 'Loading...';
		if (error) return `Error: ${error}`;
		return pcName ? `Connected to ${pcName}` : 'Connected';
	};

	const formatLastUpdate = (date: Date | null) => {
		if (!date) return 'Never';
		const now = new Date();
		const diff = Math.floor((now.getTime() - date.getTime()) / 1000);

		if (diff < 60) return `${diff}s ago`;
		if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
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
