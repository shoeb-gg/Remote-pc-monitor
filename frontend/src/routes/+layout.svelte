<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="manifest" href="/manifest.json" />
	<meta name="theme-color" content="#111827" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{@render children?.()}
