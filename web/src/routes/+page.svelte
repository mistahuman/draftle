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
		hints,
		picksInRole,
		shareText,
		takenInDraft
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
	let misses = $derived(guesses.filter((g) => g !== puzzle?.answer).length);
	let hintList = $derived(puzzle && answer ? hints(puzzle, answer) : []);
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

	const matchLabel = { same: 'same class', partial: 'shares a class', none: 'different class' };
	const matchClass = {
		same: 'preset-filled-warning-500',
		partial: 'preset-tonal-warning',
		none: 'preset-tonal-surface'
	};
</script>

<section class="container mx-auto max-w-3xl space-y-6 px-4 py-8">
	<header class="space-y-1 text-center">
		<h1 class="h2">draftle</h1>
		<p class="text-sm opacity-70">
			A real pro draft, nine picks in. Guess the last red pick: {MAX_GUESSES} tries.
		</p>
	</header>

	{#if loadError}
		<p class="card preset-tonal-error p-4">Could not load the puzzles: {loadError}</p>
	{:else if !data || !puzzle || !answer}
		<p class="text-center opacity-60">Loading…</p>
	{:else}
		<div class="flex items-center justify-between text-xs opacity-70">
			<span>
				{mode === 'daily' ? `Daily #${day}` : 'Practice'} · patch {puzzle.patch}
			</span>
			<span>{over ? '' : `${MAX_GUESSES - guesses.length} tries left`}</span>
		</div>

		<div class="card preset-tonal-surface p-4">
			<DraftBoard {puzzle} version={data.version} {byName} reveal={over} />
		</div>

		{#if !over}
			<GuessInput champions={data.champions} version={data.version} {blocked} onguess={guess} />
		{/if}

		{#if guessed.length > 0}
			<ol class="space-y-2">
				{#each guessed as g (g.name)}
					{@const m = classMatch(g, answer)}
					{@const mine = picksInRole(puzzle, g.name)}
					{@const target = picksInRole(puzzle, answer.name)}
					<li class="flex items-center gap-3 card preset-tonal-surface p-2">
						<ChampIcon version={data.version} id={g.id} name={g.name} size={36} />
						<span class="font-semibold">{g.name}</span>
						{#if g.name === answer.name}
							<span class="ml-auto badge preset-filled-success-500">correct</span>
						{:else}
							<span class="ml-auto badge {matchClass[m]}">{matchLabel[m]}</span>
							<span class="badge preset-tonal-surface" title="picks as {puzzle.role} this patch">
								{mine}× {mine === target
									? '='
									: target > mine
										? '↑ more popular'
										: '↓ less popular'}
							</span>
						{/if}
					</li>
				{/each}
			</ol>
		{/if}

		{#if misses > 0 && !over}
			<div class="space-y-1 card preset-tonal-primary p-4">
				<h3 class="text-sm font-bold">Hints</h3>
				<ul class="list-inside list-disc text-sm">
					{#each hintList.slice(0, misses) as h (h)}
						<li>{h}</li>
					{/each}
				</ul>
			</div>
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
	{/if}
</section>
