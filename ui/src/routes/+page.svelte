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
	const matchClass = {
		same: 'preset-filled-warning-500',
		partial: 'preset-tonal-warning',
		none: 'preset-tonal-surface'
	};
</script>

<section class="container mx-auto max-w-2xl space-y-5 px-4 py-6">
	{#if loadError}
		<p class="card preset-tonal-error p-4">Could not load the puzzles: {loadError}</p>
	{:else if !data || !puzzle || !answer}
		<p class="text-center opacity-60">Loading…</p>
	{:else}
		<header class="flex items-end justify-between gap-4">
			<div>
				<h1 class="text-lg font-bold">
					{mode === 'daily' ? `Daily #${day}` : 'Practice'}
				</h1>
				<p class="text-xs opacity-60">
					Patch {puzzle.patch} · guess the last red pick
				</p>
			</div>
			<div class="flex gap-1" title="{MAX_GUESSES - guesses.length} tries left">
				{#each { length: MAX_GUESSES }, i (i)}
					{@const g = guesses[i]}
					<span
						class="h-2.5 w-2.5 rounded-full {g === undefined
							? 'bg-surface-300-700'
							: g === puzzle.answer
								? 'bg-success-500'
								: 'bg-error-500'}"
					></span>
				{/each}
			</div>
		</header>

		<div class="card preset-tonal-surface p-4">
			<DraftBoard {puzzle} version={data.version} {byName} reveal={over} />
		</div>

		{#if !over}
			<GuessInput champions={data.champions} version={data.version} {blocked} onguess={guess} />
		{/if}

		{#if hintList.length > 0 && !over}
			<ol class="space-y-1 card border-l-4 border-primary-500 preset-tonal-surface p-3 text-sm">
				{#each hintList as h, i (h)}
					<li class="flex gap-2">
						<span class="font-bold opacity-60">{i + 1}.</span>
						<span>{h}</span>
					</li>
				{/each}
			</ol>
		{/if}

		{#if over}
			<div
				class="space-y-3 card p-4 text-center {won ? 'preset-tonal-success' : 'preset-tonal-error'}"
			>
				<p class="text-lg font-bold">
					{won ? `Got it in ${guesses.length}.` : `The pro picked ${puzzle.answer}.`}
				</p>
				<p class="text-sm opacity-80">
					{puzzle.event.replaceAll('/', ' · ')} — {puzzle.winner === 2 ? 'red' : 'blue'} side won.
					{#if puzzle.metaRank > 0}
						{puzzle.answer} was the #{puzzle.metaRank} most picked {puzzle.role} of the patch.
					{:else}
						Nobody else picked {puzzle.answer} as {puzzle.role} that patch.
					{/if}
				</p>
				<div class="flex justify-center gap-2">
					{#if mode === 'daily'}
						<button class="btn preset-filled-primary-500" onclick={share}>
							{copied ? 'Copied' : 'Share'}
						</button>
					{/if}
					<button class="btn preset-tonal" onclick={practice}>Practice puzzle</button>
				</div>
			</div>
		{/if}

		{#if guessed.length > 0}
			{@const target = picksInRole(puzzle, answer.name)}
			<div class="space-y-1">
				<div
					class="grid grid-cols-[1fr_6rem_6rem] gap-2 px-2 text-[10px] tracking-wider uppercase opacity-50"
				>
					<span>Guess</span>
					<span class="text-center">Class</span>
					<span class="text-center" title="picks as {puzzle.role} in patch {puzzle.patch}"
						>Picks {puzzle.patch}</span
					>
				</div>
				<ol class="space-y-1">
					{#each guessed as g (g.name)}
						{@const m = classMatch(g, answer)}
						{@const mine = picksInRole(puzzle, g.name)}
						<li
							class="grid grid-cols-[1fr_6rem_6rem] items-center gap-2 rounded-md bg-surface-100-900 p-1.5"
						>
							<span class="flex min-w-0 items-center gap-2">
								<ChampIcon version={data.version} id={g.id} name={g.name} size={28} />
								<span class="truncate text-sm font-semibold">{g.name}</span>
							</span>
							{#if g.name === answer.name}
								<span class="col-span-2 badge preset-filled-success-500">correct</span>
							{:else}
								<span class="badge {matchClass[m]}">{matchLabel[m]}</span>
								<span
									class="badge preset-tonal-surface"
									title={mine === target
										? 'same popularity'
										: target > mine
											? 'the answer is picked more'
											: 'the answer is picked less'}
								>
									{mine}
									{mine === target ? '=' : target > mine ? '↑' : '↓'}
								</span>
							{/if}
						</li>
					{/each}
				</ol>
			</div>
		{/if}
	{/if}
</section>
