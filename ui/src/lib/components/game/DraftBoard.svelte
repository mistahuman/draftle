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
	<div class="flex gap-1 {end ? 'justify-end' : ''}">
		{#each list as b (b)}
			<ChampIcon {version} id={byName.get(b)?.id} name={b} size={24} dim />
		{/each}
	</div>
{/snippet}

{#snippet pick(s: Slot, end: boolean)}
	{#if s}
		<div class="flex min-w-0 items-center gap-2 {end ? 'flex-row-reverse text-right' : ''}">
			<ChampIcon {version} id={byName.get(s.champ)?.id} name={s.champ} size={40} />
			<div class="min-w-0">
				<div class="truncate text-sm font-semibold">{s.champ}</div>
				<div class="text-[10px] opacity-50">{s.order}</div>
			</div>
		</div>
	{/if}
{/snippet}

<div class="space-y-3">
	<div class="grid grid-cols-2 items-end gap-4 text-xs font-bold tracking-wider uppercase">
		<div class="space-y-2">
			<div class="text-(--side-blue)">Blue</div>
			{@render bans(puzzle.blue.bans)}
		</div>
		<div class="space-y-2 text-right">
			<div class="text-(--side-red)">Red</div>
			{@render bans(puzzle.red.bans, true)}
		</div>
	</div>

	<ul class="divide-y divide-surface-500/15">
		{#each roles as role (role)}
			{@const open = role === puzzle.role && !reveal}
			{@const target = role === puzzle.role}
			<li class="grid grid-cols-[1fr_auto_1fr] items-center gap-2 py-1.5">
				{@render pick(slot('blue', role), false)}
				<span class="w-9 text-center text-[10px] font-bold tracking-wider opacity-50">
					{roleShort[role]}
				</span>
				{#if open}
					<div class="flex items-center justify-end gap-2">
						<span class="text-xs opacity-70">your pick</span>
						<div
							class="flex h-10 w-10 items-center justify-center rounded-md bg-(--side-red) text-lg font-bold text-white"
						>
							?
						</div>
					</div>
				{:else}
					<div class={target ? 'rounded-md ring-2 ring-(--side-red)' : ''}>
						{@render pick(slot('red', role), true)}
					</div>
				{/if}
			</li>
		{/each}
	</ul>
</div>
