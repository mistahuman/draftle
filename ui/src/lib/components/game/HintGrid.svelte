<script lang="ts">
	import { Lock } from 'lucide-svelte';
	import type { Hint } from '$lib/game/logic';

	// Five boxes, always visible: the label says what is coming, the value shows up when
	// the matching square in the header is spent. Locked boxes are dashed like the R5 slot.
	let { hints, unlocked }: { hints: Hint[]; unlocked: number } = $props();
</script>

<div class="space-y-2">
	<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">Hints</p>
	<ol class="grid grid-cols-2 gap-2 sm:grid-cols-5">
		{#each hints as h, i (i)}
			{@const open = i < unlocked}
			{@const left = i + 1 - unlocked}
			<li
				class="flex min-h-24 flex-col gap-1 rounded-container p-2.5 {i === hints.length - 1
					? 'col-span-2 sm:col-span-1'
					: ''} {open
					? 'border border-surface-400-600 bg-surface-100-900/60'
					: 'border border-dashed border-surface-400-600 opacity-60'}"
			>
				<span class="flex items-center gap-1.5 font-mono text-[10px] tracking-wider uppercase">
					<span class="opacity-50">{i + 1}</span>
					<span class="truncate opacity-70" title={h.label}>{h.label}</span>
				</span>
				{#if open}
					<span class="text-lg leading-tight font-semibold">{h.value}</span>
					<span class="text-xs leading-snug opacity-60">{h.note}</span>
				{:else}
					<span class="mt-auto flex items-center gap-1.5 text-xs opacity-70">
						<Lock size={12} />
						{left} more {left === 1 ? 'miss' : 'misses'}
					</span>
				{/if}
			</li>
		{/each}
	</ol>
</div>
