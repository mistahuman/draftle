<script lang="ts">
	import { ArrowDown, ArrowUp, Equal } from 'lucide-svelte';
	import ChampIcon from './ChampIcon.svelte';
	import { classMatch, picksInRole } from '$lib/game/logic';
	import type { Champion, Puzzle } from '$lib/game/types';

	// Newest guess on top: with unlimited tries the list grows, and the last answer is the
	// one you read. Each guess gets two small tiles in the same language as the hint boxes.
	let {
		guesses,
		answer,
		puzzle,
		version
	}: { guesses: Champion[]; answer: Champion; puzzle: Puzzle; version: string } = $props();

	let target = $derived(picksInRole(puzzle, answer.name));
	let rows = $derived(guesses.map((g, i) => ({ g, n: i + 1 })).reverse());

	const classLabel = { same: 'Same', partial: 'Partial', none: 'Different' };
	const classTile = {
		same: 'border-primary-500 bg-primary-500/15 text-primary-700-300',
		partial: 'border-surface-500 bg-surface-500/10',
		none: 'border-dashed border-surface-400-600 opacity-50'
	};
</script>

<div class="space-y-2">
	<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">Guesses</p>
	<ol class="space-y-1.5">
		{#each rows as { g, n } (g.name)}
			{@const hit = g.name === answer.name}
			{@const m = classMatch(g, answer)}
			{@const mine = picksInRole(puzzle, g.name)}
			<li
				class="flex items-center gap-2 rounded-container border px-2 py-2 sm:gap-3 sm:px-2.5 {hit
					? 'border-success-500 bg-success-500/10'
					: 'border-surface-300-700 bg-surface-100-900/40'}"
			>
				<span class="hidden w-5 text-right font-mono text-xs opacity-40 sm:inline">{n}</span>
				<ChampIcon {version} id={g.id} name={g.name} size={32} />
				<span class="min-w-0 flex-1 truncate font-semibold">{g.name}</span>
				{#if hit}
					<span
						class="flex items-center rounded-base border border-success-500 px-3 py-2 font-mono text-xs font-bold tracking-wider text-success-700-300 uppercase"
						>Correct</span
					>
				{:else}
					<span
						class="flex w-[4.25rem] shrink-0 flex-col items-center rounded-base border px-1 py-1 sm:w-24 {classTile[
							m
						]}"
						title="class compared with the answer's"
					>
						<span class="font-mono text-[9px] tracking-wider uppercase opacity-60">Class</span>
						<span class="text-sm font-semibold">{classLabel[m]}</span>
					</span>
					<span
						class="flex w-[4.25rem] shrink-0 flex-col items-center rounded-base border border-surface-400-600 px-1 py-1 sm:w-24"
						title="{g.name}: {mine} picks as {puzzle.role} in patch {puzzle.patch}; the answer is picked {mine ===
						target
							? 'as much'
							: target > mine
								? 'more'
								: 'less'}"
					>
						<span class="font-mono text-[9px] tracking-wider uppercase opacity-60">Picks</span>
						<span class="flex items-center gap-1 text-sm font-semibold">
							{mine}
							{#if mine === target}
								<Equal size={14} class="opacity-60" />
							{:else if target > mine}
								<ArrowUp size={14} class="text-primary-600-400" />
							{:else}
								<ArrowDown size={14} class="text-primary-600-400" />
							{/if}
						</span>
					</span>
				{/if}
			</li>
		{/each}
	</ol>
</div>
