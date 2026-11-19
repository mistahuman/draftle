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

	// Rows are lanes, so each pick faces its lane opponent. The B1…R5 label keeps the
	// draft order visible.
	const roles = ['Top', 'Jungle', 'Mid', 'Bot', 'Support'];
	const roleShort: Record<string, string> = {
		Top: 'TOP',
		Jungle: 'JGL',
		Mid: 'MID',
		Bot: 'BOT',
		Support: 'SUP'
	};

	type Slot = { champ: string; order: string } | undefined;

	function slot(side: 'blue' | 'red', role: string): Slot {
		const picks = puzzle[side].picks;
		const i = picks.findIndex((p) => p.role === role);
		if (i >= 0) return { champ: picks[i].champ, order: `${side === 'blue' ? 'B' : 'R'}${i + 1}` };
		if (side === 'red' && role === puzzle.role && reveal)
			return { champ: puzzle.answer, order: 'R5' };
		return undefined;
	}
</script>

{#snippet bans(list: string[], end = false)}
	<div class="flex items-center gap-1 {end ? 'flex-row-reverse' : ''}">
		<span
			class="hidden font-mono text-[10px] tracking-widest uppercase opacity-50 sm:inline {end
				? 'ml-1'
				: 'mr-1'}">bans</span
		>
		{#each list as b (b)}
			<ChampIcon {version} id={byName.get(b)?.id} name={b} size={22} dim />
		{/each}
	</div>
{/snippet}

{#snippet pick(s: Slot, end: boolean)}
	{#if s}
		<div class="flex min-w-0 items-center gap-2.5 {end ? 'flex-row-reverse text-right' : ''}">
			<ChampIcon {version} id={byName.get(s.champ)?.id} name={s.champ} size={40} />
			<div class="min-w-0 leading-tight">
				<div class="truncate font-semibold">{s.champ}</div>
				<div class="font-mono text-[11px] opacity-50">{s.order}</div>
			</div>
		</div>
	{/if}
{/snippet}

<div class="space-y-3">
	<div class="grid grid-cols-2 gap-4">
		<div class="space-y-1.5">
			<div class="text-sm font-semibold text-(--side-blue)">Blue side</div>
			{@render bans(puzzle.blue.bans)}
		</div>
		<div class="space-y-1.5 text-right">
			<div class="text-sm font-semibold text-(--side-red)">Red side</div>
			{@render bans(puzzle.red.bans, true)}
		</div>
	</div>

	<ul class="divide-y divide-dotted divide-surface-400-600 border-y border-surface-300-700">
		{#each roles as role (role)}
			{@const open = role === puzzle.role && !reveal}
			{@const target = role === puzzle.role}
			<li class="grid grid-cols-[1fr_auto_1fr] items-center gap-2 py-2">
				{@render pick(slot('blue', role), false)}
				<span class="w-10 text-center font-mono text-[11px] tracking-widest opacity-50">
					{roleShort[role]}
				</span>
				{#if open}
					<div class="flex items-center justify-end gap-2.5">
						<span class="text-sm opacity-60">your pick</span>
						<div
							class="flex h-10 w-10 items-center justify-center rounded-base border-2 border-dashed border-(--side-red) text-lg font-semibold text-(--side-red)"
						>
							?
						</div>
					</div>
				{:else}
					<div
						class={target
							? '-m-1 rounded-base p-1 outline-2 outline-(--side-red) outline-dashed'
							: ''}
					>
						{@render pick(slot('red', role), true)}
					</div>
				{/if}
			</li>
		{/each}
	</ul>
</div>
