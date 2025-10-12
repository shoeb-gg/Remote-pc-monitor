<script lang="ts">
	interface Props {
		isOpen: boolean;
		currentInterval: number;
		onClose: () => void;
		onSave: (interval: number) => void;
	}

	let { isOpen, currentInterval, onClose, onSave }: Props = $props();

	let intervalValue = $state(currentInterval);

	// Update local state when props change
	$effect(() => {
		intervalValue = currentInterval;
	});

	const handleSave = () => {
		// Validate interval is at least 1 second
		if (intervalValue < 1) {
			alert('Refresh interval must be at least 1 second');
			return;
		}
		onSave(intervalValue);
		onClose();
	};

	const handleBackdropClick = (e: MouseEvent) => {
		if (e.target === e.currentTarget) {
			onClose();
		}
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			handleSave();
		}
	};

	// Preset intervals
	const presets = [
		{ label: '5s', value: 5 },
		{ label: '10s', value: 10 },
		{ label: '15s', value: 15 },
		{ label: '30s', value: 30 },
		{ label: '60s', value: 60 }
	];

	const setPreset = (value: number) => {
		intervalValue = value;
	};
</script>

{#if isOpen}
	<div
		class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4"
		onclick={handleBackdropClick}
		role="button"
		tabindex="-1"
	>
		<div class="bg-gray-800 rounded-2xl p-6 w-full max-w-md shadow-2xl border border-gray-700">
			<h2 class="text-2xl font-bold text-white mb-6">Refresh Interval Settings</h2>

			<div class="space-y-4">
				<div>
					<label for="interval-value" class="block text-sm font-medium text-white/80 mb-2">
						Refresh Interval (seconds)
					</label>
					<input
						id="interval-value"
						type="number"
						min="1"
						bind:value={intervalValue}
						onkeydown={handleKeyDown}
						class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>

				<div>
					<p class="text-sm font-medium text-white/80 mb-2">Quick Presets</p>
					<div class="flex flex-wrap gap-2">
						{#each presets as preset}
							<button
								onclick={() => setPreset(preset.value)}
								class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors duration-200 {intervalValue === preset.value ? 'ring-2 ring-blue-500' : ''}"
							>
								{preset.label}
							</button>
						{/each}
					</div>
				</div>

				<div class="text-sm text-white/60 bg-gray-900 rounded-lg p-3">
					Data will refresh every {intervalValue} second{intervalValue !== 1 ? 's' : ''}
				</div>

				<div class="text-xs text-yellow-400/80 bg-yellow-900/20 rounded-lg p-3 border border-yellow-700/30">
					<strong>Note:</strong> Shorter intervals use more API calls. Keep it at 30s or higher to stay within free tier limits.
				</div>
			</div>

			<div class="flex space-x-3 mt-6">
				<button
					onclick={handleSave}
					class="flex-1 bg-gradient-to-r from-blue-500 to-cyan-600 hover:from-blue-600 hover:to-cyan-700 text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200"
				>
					Save
				</button>
				<button
					onclick={onClose}
					class="flex-1 bg-gray-700 hover:bg-gray-600 text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200"
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}
