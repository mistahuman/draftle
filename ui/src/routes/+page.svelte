<script lang="ts">
	import { onMount } from 'svelte';
	import { asset } from '$app/paths';
	import DraftBoard from '$lib/components/game/DraftBoard.svelte';
	import GuessInput from '$lib/components/game/GuessInput.svelte';
	import GuessList from '$lib/components/game/GuessList.svelte';
	import HintGrid from '$lib/components/game/HintGrid.svelte';
	import {
		HINTS,
		hints,
		dailyIndex,
		dayNumber,
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
	// The "who won?" bonus after the puzzle: 1 blue, 2 red, null not answered yet.
	let winnerCall = $state<1 | 2 | null>(null);
	let gaveUp = $state(false);

	const storageKey = (d: number) => `draftle:day:${d}`;
	const winnerKey = (d: number) => `draftle:day:${d}:winner`;
	const gaveUpKey = (d: number) => `draftle:day:${d}:gaveup`;

	onMount(async () => {
		try {
			const res = await fetch(asset('/data/puzzles.json'));
			if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
			data = await res.json();
			index = dailyIndex(day, data!.puzzles.length);
			try {
				const saved = localStorage.getItem(storageKey(day));
				if (saved) guesses = JSON.parse(saved);
				const call = localStorage.getItem(winnerKey(day));
				if (call === '1' || call === '2') winnerCall = Number(call) as 1 | 2;
				gaveUp = localStorage.getItem(gaveUpKey(day)) === '1';
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
	let over = $derived(won || gaveUp);
	let misses = $derived(guesses.filter((g) => g !== puzzle?.answer).length);
	let hintList = $derived(puzzle && answer ? hints(puzzle, answer) : []);
	let unlocked = $derived(puzzle ? unlockedHints(guesses, puzzle.answer) : 0);
	let blueName = $derived(puzzle?.blueTeam || 'Blue side');
	let redName = $derived(puzzle?.redTeam || 'Red side');
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

	function giveUp() {
		if (over) return;
		gaveUp = true;
		if (mode === 'daily') {
			try {
				localStorage.setItem(gaveUpKey(day), '1');
			} catch {
				// storage unavailable
			}
		}
	}

	function callWinner(side: 1 | 2) {
		if (winnerCall !== null) return;
		winnerCall = side;
		copied = false;
		if (mode === 'daily') {
			try {
				localStorage.setItem(winnerKey(day), String(side));
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
		winnerCall = null;
		gaveUp = false;
		copied = false;
	}

	async function share() {
		if (!answer) return;
		try {
			await navigator.clipboard.writeText(
				shareText(
					day,
					guessed,
					answer,
					won,
					winnerCall === null ? undefined : winnerCall === puzzle?.winner
				)
			);
			copied = true;
		} catch {
			copied = false;
		}
	}

	const fmtDate = (d: string) =>
		new Date(d + 'T00:00:00Z').toLocaleDateString('en-GB', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
			timeZone: 'UTC'
		});
</script>

<section class="container mx-auto max-w-2xl space-y-6 px-4 py-8">
	{#if loadError}
		<p class="card preset-outlined-error-500 p-4">Could not load the puzzles: {loadError}</p>
	{:else if !data || !puzzle || !answer}
		<p class="text-center opacity-60">Loading…</p>
	{:else}
		<!-- Just the five hint squares: the board says what to do, and the patch lives on it. -->
		<div class="flex items-center justify-center gap-3">
			{#if mode === 'practice'}
				<span class="badge preset-outlined-surface-500 font-mono text-xs">practice</span>
			{/if}
			<div class="flex gap-1.5" title="{unlocked} of {HINTS} hints unlocked">
				{#each { length: HINTS }, i (i)}
					<span
						class="h-4 w-4 rounded-[3px] {i < misses
							? 'bg-(--side-red)'
							: won && i === misses
								? 'bg-success-500'
								: 'bg-surface-300-700'}"
					></span>
				{/each}
			</div>
		</div>

		<div
			class="relative rounded-container border border-surface-300-700 bg-surface-100-900/40 p-4 shadow-sm"
		>
			<!-- the patch sits on the frame, like a label on a scorecard -->
			<span
				class="absolute -top-2.5 left-1/2 -translate-x-1/2 rounded-base border border-surface-400-600 bg-surface-50-950 px-2 font-mono text-[11px] leading-[1.15rem]"
				>patch {puzzle.patch}</span
			>
			<DraftBoard {puzzle} version={data.version} {byName} reveal={over} />
		</div>

		{#if !over}
			<GuessInput champions={data.champions} version={data.version} {blocked} onguess={guess} />
			{#if misses >= HINTS}
				<div class="flex justify-end">
					<button class="btn preset-outlined-surface-500 opacity-70 btn-sm" onclick={giveUp}>
						Give up and show the pick
					</button>
				</div>
			{/if}
		{/if}

		{#if over}
			<div
				class="space-y-3 rounded-container border-[3px] border-double p-5 text-center {won
					? 'border-success-500'
					: 'border-(--side-red)'}"
			>
				<p class="font-display text-2xl font-semibold tracking-tight">
					{won ? `Got it in ${guesses.length}.` : `The pro picked ${puzzle.answer}.`}
				</p>
				<p class="text-sm opacity-80">
					{#if puzzle.metaRank > 0}
						{puzzle.answer} was the #{puzzle.metaRank} most picked {puzzle.role} of the patch.
					{:else}
						Nobody else picked {puzzle.answer} as {puzzle.role} that patch.
					{/if}
				</p>

				<div class="space-y-2 border-t border-dotted border-surface-400-600 pt-3">
					{#if winnerCall === null}
						<p class="font-mono text-[11px] tracking-[0.2em] uppercase opacity-60">
							Bonus · who won the game?
						</p>
						<div class="flex items-center justify-center gap-2">
							<button
								class="btn border border-(--side-blue) text-(--side-blue) hover:bg-(--side-blue)/10"
								onclick={() => callWinner(1)}>{blueName}</button
							>
							<span class="font-mono text-xs opacity-50">vs</span>
							<button
								class="btn border border-(--side-red) text-(--side-red) hover:bg-(--side-red)/10"
								onclick={() => callWinner(2)}>{redName}</button
							>
						</div>
					{:else}
						<p>
							<span
								class="font-semibold {puzzle.winner === 2
									? 'text-(--side-red)'
									: 'text-(--side-blue)'}">{puzzle.winner === 2 ? redName : blueName}</span
							>
							won.
							<span
								class="badge font-mono text-[11px] uppercase {winnerCall === puzzle.winner
									? 'preset-outlined-success-500'
									: 'preset-outlined-surface-500 opacity-60'}"
								>{winnerCall === puzzle.winner ? 'called it' : 'missed'}</span
							>
						</p>
					{/if}
					<p class="font-mono text-xs tracking-wide uppercase opacity-60">
						{puzzle.event.replaceAll('/', ' · ')}{puzzle.date ? ` · ${fmtDate(puzzle.date)}` : ''}
					</p>
				</div>
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

		{#if over || unlocked > 0}
			<HintGrid hints={hintList.slice(0, over ? HINTS : unlocked)} />
		{/if}

		{#if guessed.length > 0}
			<GuessList guesses={guessed} {answer} {puzzle} version={data.version} />
		{/if}
	{/if}
</section>
