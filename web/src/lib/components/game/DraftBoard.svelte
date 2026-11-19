<script lang="ts">
	import ChampIcon from './ChampIcon.svelte';
	import type { Champion, Puzzle } from '$lib/game/types';

	let {
		puzzle,
		version,
		byName,
		reveal
	}: {
		puzzle: Puzzle;
		version: string;
		byName: Map<string, Champion>;
		reveal: boolean;
	} = $props();

	const roleShort: Record<string, string> = {
		Top: 'TOP',
		Jungle: 'JGL',
		Mid: 'MID',
		Bot: 'BOT',
		Support: 'SUP'
	};
</script>

{#snippet bans(list: string[])}
	<div class="flex flex-wrap gap-1">
		{#each list as b (b)}
			<ChampIcon {version} id={byName.get(b)?.id} name={b} size={28} dim />
		{/each}
	</div>
{/snippet}

{#snippet pickRow(champ: string, role: string, label: string, highlight = false)}
	<div
		class="flex items-center gap-3 rounded-lg p-1.5 {highlight
			? 'preset-tonal-primary ring-2 ring-primary-500'
			: ''}"
	>
		<span class="w-6 text-[10px] opacity-50">{label}</span>
		<ChampIcon {version} id={byName.get(champ)?.id} name={champ} size={40} />
		<div class="min-w-0">
			<div class="truncate text-sm font-semibold">{champ}</div>
			<div class="text-[10px] tracking-wider opacity-60">{roleShort[role] ?? role}</div>
		</div>
	</div>
{/snippet}

<div class="grid grid-cols-2 gap-3 sm:gap-6">
	<div class="space-y-2">
		<h3 class="text-sm font-bold text-[#3b82f6]">Blue side</h3>
		{@render bans(puzzle.blue.bans)}
		{#each puzzle.blue.picks as p, i (p.champ)}
			{@render pickRow(p.champ, p.role, `B${i + 1}`)}
		{/each}
	</div>
	<div class="space-y-2">
		<h3 class="text-right text-sm font-bold text-[#ef4444]">Red side</h3>
		<div class="flex justify-end">{@render bans(puzzle.red.bans)}</div>
		{#each puzzle.red.picks as p, i (p.champ)}
			{@render pickRow(p.champ, p.role, `R${i + 1}`)}
		{/each}
		{#if reveal}
			{@render pickRow(puzzle.answer, puzzle.role, 'R5', true)}
		{:else}
			<div class="flex items-center gap-3 rounded-lg p-1.5 ring-2 ring-primary-500">
				<span class="w-6 text-[10px] opacity-50">R5</span>
				<div
					class="flex h-10 w-10 items-center justify-center rounded-md bg-surface-200-800 text-xl font-bold"
				>
					?
				</div>
				<div>
					<div class="text-sm font-semibold">Your pick</div>
					<div class="text-[10px] tracking-wider opacity-60">
						{roleShort[puzzle.role] ?? puzzle.role} · vs {puzzle.opponent}
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
