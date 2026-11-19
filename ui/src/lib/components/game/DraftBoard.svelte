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

	// Laid out like a broadcast draft: blue left, red right, picks top to bottom in pick
	// order, bans in their two phases (1-3, 4-5). The gap after the third pick is where the
	// second ban phase happens: B1 R1 R2 B2 B3 R3 | bans | R4 B4 B5 R5.
	const roleShort: Record<string, string> = {
		Top: 'TOP',
		Jungle: 'JGL',
		Mid: 'MID',
		Bot: 'BOT',
		Support: 'SUP'
	};

	let redPicks = $derived([
		...puzzle.red.picks,
		{ champ: reveal ? puzzle.answer : '', role: puzzle.role }
	]);
</script>

{#snippet banSlot(b: string)}
	{#if b}
		<ChampIcon {version} id={byName.get(b)?.id} name={b} size={24} dim />
	{:else}
		<div
			class="h-6 w-6 rounded-base border border-dashed border-surface-400-600"
			title="no ban"
		></div>
	{/if}
{/snippet}

{#snippet bans(list: string[], red: boolean)}
	<div class="flex items-center gap-2.5 {red ? 'flex-row-reverse' : ''}">
		<div class="flex gap-1 {red ? 'flex-row-reverse' : ''}">
			{#each list.slice(0, 3) as b, i (i)}{@render banSlot(b)}{/each}
		</div>
		<div class="flex gap-1 {red ? 'flex-row-reverse' : ''}">
			{#each list.slice(3) as b, i (i)}{@render banSlot(b)}{/each}
		</div>
	</div>
{/snippet}

{#snippet pick(champ: string, role: string, n: number, red: boolean)}
	{@const last = red && n === 5}
	<li
		class="flex items-center gap-2.5 {red ? 'flex-row-reverse text-right' : ''} {n === 4
			? 'mt-5'
			: ''}"
	>
		{#if champ}
			<div
				class="shrink-0 {last
					? 'rounded-base outline-2 outline-offset-2 outline-(--side-red) outline-dashed'
					: ''}"
			>
				<ChampIcon {version} id={byName.get(champ)?.id} name={champ} size={44} />
			</div>
			<div class="min-w-0 leading-tight">
				<div class="truncate font-semibold">{champ}</div>
				<div class="font-mono text-[11px] opacity-50">{roleShort[role] ?? role}</div>
			</div>
		{:else}
			<div
				class="flex h-11 w-11 shrink-0 items-center justify-center rounded-base border-2 border-dashed border-(--side-red) text-lg font-semibold text-(--side-red)"
			>
				?
			</div>
			<div class="min-w-0 leading-tight">
				<div class="font-semibold text-(--side-red)">Your pick</div>
				<div class="font-mono text-[11px] opacity-50">
					{roleShort[role] ?? role} vs {puzzle.opponent}
				</div>
			</div>
		{/if}
	</li>
{/snippet}

<div class="space-y-4">
	<div class="grid grid-cols-2 gap-4">
		<div class="space-y-2">
			<div class="text-sm font-semibold text-(--side-blue)">Blue side</div>
			{@render bans(puzzle.blue.bans, false)}
		</div>
		<div class="space-y-2 text-right">
			<div class="text-sm font-semibold text-(--side-red)">Red side</div>
			{@render bans(puzzle.red.bans, true)}
		</div>
	</div>

	<div class="grid grid-cols-2 gap-4 border-t border-surface-300-700 pt-4">
		<ol class="space-y-2 border-l-2 border-(--side-blue) pl-3">
			{#each puzzle.blue.picks as p, i (i)}
				{@render pick(p.champ, p.role, i + 1, false)}
			{/each}
		</ol>
		<ol class="space-y-2 border-r-2 border-(--side-red) pr-3">
			{#each redPicks as p, i (i)}
				{@render pick(p.champ, p.role, i + 1, true)}
			{/each}
		</ol>
	</div>
</div>
