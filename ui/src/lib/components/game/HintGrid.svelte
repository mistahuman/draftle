<script lang="ts">
	import type { Hint } from '$lib/game/logic';

	// Only the unlocked hints: a locked box would leak its label (the lane opponent, the
	// red team's name, which the board hides until the end). The header squares already
	// say how many are left.
	let { hints }: { hints: Hint[] } = $props();
</script>

<div class="space-y-2">
	<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">Hints</p>
	<ol class="grid grid-cols-2 gap-2 sm:grid-cols-5">
		{#each hints as h, i (i)}
			<li
				class="flex flex-col gap-1 rounded-container border border-surface-400-600 bg-surface-100-900/60 p-2.5 {i ===
				4
					? 'col-span-2 sm:col-span-1'
					: ''}"
			>
				<span class="flex items-center gap-1.5 font-mono text-[10px] tracking-wider uppercase">
					<span class="opacity-50">{i + 1}</span>
					<span class="truncate opacity-70" title={h.label}>{h.label}</span>
				</span>
				<span class="text-lg leading-tight font-semibold">{h.value}</span>
				<span class="text-xs leading-snug opacity-60">{h.note}</span>
			</li>
		{/each}
	</ol>
</div>
