<script lang="ts">
	import { onMount } from 'svelte';
	import { asset } from '$app/paths';
	import DraftBoard from '$lib/components/game/DraftBoard.svelte';
	import GuessInput from '$lib/components/game/GuessInput.svelte';
	import ChampIcon from '$lib/components/game/ChampIcon.svelte';
	import {
		MAX_GUESSES,
		classMatch,
		dailyIndex,
		dayNumber,
		picksInRole,
		shareText,
		takenInDraft,
		unlockedHints
	} from '$lib/game/logic';
	import type { Champion, PuzzleData } from '$lib/game/types';

	let data = $state<PuzzleData | null>(null);
	let loadError = $state('');

	const day = dayNumber();
	let mode = $state<'daily' | 'practice'>('daily');
	let index = $state(0);
	let guesses = $state<string[]>([]);
	let copied = $state(false);

	const storageKey = (d: number) => `draftle:day:${d}`;

	onMount(async () => {
		try {
			const res = await fetch(asset('/data/puzzles.json'));
			if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
			data = await res.json();
			index = dailyIndex(day, data!.puzzles.length);
			try {
				const saved = localStorage.getItem(storageKey(day));
				if (saved) guesses = JSON.parse(saved);
			} catch {
				// storage unavailable: play without saving
			}
		} catch (e) {
			loadError = String(e);
		}
	});

	let byName = $derived(new Map((data?.champions ?? []).map((c) => [c.name, c])));
	let puzzle = $derived(data?.puzzles[index]);
	let answer = $derived(puzzle ? byName.get(puzzle.answer) : undefined);
	let guessed = $derived(guesses.map((g) => byName.get(g)).filter((c): c is Champion => !!c));
	let won = $derived(!!puzzle && guesses.includes(puzzle.answer));
	let over = $derived(won || guesses.length >= MAX_GUESSES);
	let hintList = $derived(puzzle && answer ? unlockedHints(puzzle, answer, guessed) : []);
	let blocked = $derived(new Set([...(puzzle ? takenInDraft(puzzle) : []), ...guesses]));

	function guess(c: Champion) {
		if (over) return;
		guesses = [...guesses, c.name];
		if (mode === 'daily') {
			try {
				localStorage.setItem(storageKey(day), JSON.stringify(guesses));
			} catch {
				// storage unavailable
			}
		}
	}

	function practice() {
		if (!data) return;
		mode = 'practice';
		index = Math.floor(Math.random() * data.puzzles.length);
		guesses = [];
		copied = false;
	}

	async function share() {
		if (!answer) return;
		try {
			await navigator.clipboard.writeText(shareText(day, guessed, answer, won));
			copied = true;
		} catch {
			copied = false;
		}
	}

	const matchLabel = { same: 'same', partial: 'partial', none: 'different' };
	// Stamped labels: outline only, like a rubber stamp on the scorecard.
	const matchClass = {
		same: 'preset-outlined-primary-500 font-bold',
		partial: 'preset-outlined-surface-500',
		none: 'opacity-40'
	};
</script>

<section class="container mx-auto max-w-2xl space-y-6 px-4 py-8">
	{#if loadError}
		<p class="card preset-outlined-error-500 p-4">Could not load the puzzles: {loadError}</p>
	{:else if !data || !puzzle || !answer}
		<p class="text-center font-display italic opacity-60">Loading…</p>
	{:else}
		<header
			class="flex items-end justify-between gap-4 border-b-[3px] border-double border-surface-300-700 pb-3"
		>
			<div>
				<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">
					Patch {puzzle.patch}
				</p>
				<h1 class="font-display text-3xl font-semibold">
					{#if mode === 'daily'}Daily <span class="italic">No.</span> {day}{:else}Practice{/if}
				</h1>
				<p class="text-sm opacity-70">Nine picks in. Name the last red one.</p>
			</div>
			<div class="flex gap-1 pb-1" title="{MAX_GUESSES - guesses.length} tries left">
				{#each { length: MAX_GUESSES }, i (i)}
					{@const g = guesses[i]}
					<span
						class="flex h-5 w-5 items-center justify-center rounded-base border border-surface-400-600 font-mono text-xs font-bold {g ===
						undefined
							? ''
							: g === puzzle.answer
								? 'text-success-600-400'
								: 'text-(--side-red)'}"
					>
						{g === undefined ? '' : g === puzzle.answer ? '✓' : '✕'}
					</span>
				{/each}
			</div>
		</header>

		<div
			class="rounded-container border border-surface-300-700 bg-surface-100-900/40 p-4 shadow-sm"
		>
			<DraftBoard {puzzle} version={data.version} {byName} reveal={over} />
		</div>

		{#if !over}
			<GuessInput champions={data.champions} version={data.version} {blocked} onguess={guess} />
		{/if}

		{#if hintList.length > 0 && !over}
			<div class="space-y-1.5 border-l-2 border-dashed border-primary-500 pl-4">
				<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">Clues</p>
				<ol class="space-y-1">
					{#each hintList as h, i (h)}
						<li class="flex gap-2">
							<span class="font-mono text-sm opacity-50">{i + 1}.</span>
							<span>{h}</span>
						</li>
					{/each}
				</ol>
			</div>
		{/if}

		{#if over}
			<div
				class="space-y-3 rounded-container border-[3px] border-double p-5 text-center {won
					? 'border-success-500'
					: 'border-(--side-red)'}"
			>
				<p class="font-display text-2xl font-semibold italic">
					{won ? `Got it in ${guesses.length}.` : `The pro picked ${puzzle.answer}.`}
				</p>
				<p class="font-mono text-xs tracking-wide uppercase opacity-60">
					{puzzle.event.replaceAll('/', ' · ')}
				</p>
				<p class="text-sm opacity-80">
					{puzzle.winner === 2 ? 'Red' : 'Blue'} side won.
					{#if puzzle.metaRank > 0}
						{puzzle.answer} was the #{puzzle.metaRank} most picked {puzzle.role} of the patch.
					{:else}
						Nobody else picked {puzzle.answer} as {puzzle.role} that patch.
					{/if}
				</p>
				<div class="flex justify-center gap-2 pt-1">
					{#if mode === 'daily'}
						<button class="btn preset-filled-primary-500" onclick={share}>
							{copied ? 'Copied' : 'Share'}
						</button>
					{/if}
					<button class="btn preset-outlined-surface-500" onclick={practice}>Practice puzzle</button
					>
				</div>
			</div>
		{/if}

		{#if guessed.length > 0}
			{@const target = picksInRole(puzzle, answer.name)}
			<div>
				<div
					class="grid grid-cols-[1fr_6rem_5rem] gap-2 border-b border-surface-300-700 px-1 pb-1 font-mono text-[11px] tracking-widest uppercase opacity-60"
				>
					<span>Guess</span>
					<span class="text-center">Class</span>
					<span class="text-right" title="picks as {puzzle.role} in patch {puzzle.patch}"
						>Picks</span
					>
				</div>
				<ol class="divide-y divide-dotted divide-surface-400-600">
					{#each guessed as g, i (g.name)}
						{@const m = classMatch(g, answer)}
						{@const mine = picksInRole(puzzle, g.name)}
						<li class="grid grid-cols-[1fr_6rem_5rem] items-center gap-2 px-1 py-1.5">
							<span class="flex min-w-0 items-center gap-2">
								<span class="w-4 font-mono text-xs opacity-40">{i + 1}</span>
								<ChampIcon version={data.version} id={g.id} name={g.name} size={28} />
								<span class="truncate font-semibold">{g.name}</span>
							</span>
							{#if g.name === answer.name}
								<span
									class="col-span-2 badge justify-self-end preset-outlined-success-500 font-mono font-bold uppercase"
									>correct</span
								>
							{:else}
								<span
									class="badge justify-self-center font-mono text-[11px] uppercase {matchClass[m]}"
									>{matchLabel[m]}</span
								>
								<span
									class="text-right font-mono text-sm"
									title={mine === target
										? 'same popularity'
										: target > mine
											? 'the answer is picked more'
											: 'the answer is picked less'}
								>
									{mine}
									<span class="opacity-60">{mine === target ? '=' : target > mine ? '↑' : '↓'}</span
									>
								</span>
							{/if}
						</li>
					{/each}
				</ol>
			</div>
		{/if}
	{/if}
</section>
