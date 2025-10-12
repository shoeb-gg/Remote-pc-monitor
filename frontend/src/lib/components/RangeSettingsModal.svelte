<script lang="ts">
	interface Props {
		isOpen: boolean;
		title: string;
		currentMin: number;
		currentMax: number;
		onClose: () => void;
		onSave: (min: number, max: number) => void;
	}

	let { isOpen, title, currentMin, currentMax, onClose, onSave }: Props = $props();

	let minValue = $state(currentMin);
	let maxValue = $state(currentMax);

	// Update local state when props change
	$effect(() => {
		minValue = currentMin;
		maxValue = currentMax;
	});

	const handleSave = () => {
		// Validate that max is greater than min
		if (maxValue <= minValue) {
			alert('Maximum value must be greater than minimum value');
			return;
		}
		onSave(minValue, maxValue);
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
</script>

{#if isOpen}
	<div
		class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4"
		onclick={handleBackdropClick}
		role="button"
		tabindex="-1"
	>
		<div class="bg-gray-800 rounded-2xl p-6 w-full max-w-md shadow-2xl border border-gray-700">
			<h2 class="text-2xl font-bold text-white mb-6">{title} Range Settings</h2>

			<div class="space-y-4">
				<div>
					<label for="min-value" class="block text-sm font-medium text-white/80 mb-2">
						Minimum Value
					</label>
					<input
						id="min-value"
						type="number"
						bind:value={minValue}
						onkeydown={handleKeyDown}
						class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>

				<div>
					<label for="max-value" class="block text-sm font-medium text-white/80 mb-2">
						Maximum Value
					</label>
					<input
						id="max-value"
						type="number"
						bind:value={maxValue}
						onkeydown={handleKeyDown}
						class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>

				<div class="text-sm text-white/60 bg-gray-900 rounded-lg p-3">
					Current range: {minValue} - {maxValue}
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
