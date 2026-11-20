<script lang="ts">
	import { onMount } from 'svelte';
	import { asset } from '$app/paths';
	import DraftBoard from '$lib/components/game/DraftBoard.svelte';
	import GuessInput from '$lib/components/game/GuessInput.svelte';
	import ChampIcon from '$lib/components/game/ChampIcon.svelte';
	import {
		CLUES,
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
	let hintList = $derived(puzzle && answer ? unlockedHints(puzzle, answer, guessed) : []);
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
		<p class="text-center opacity-60">Loading…</p>
	{:else}
		<header
			class="flex items-end justify-between gap-4 border-b-[3px] border-double border-surface-300-700 pb-3"
		>
			<div>
				<div class="flex items-center gap-3">
					<h1 class="font-display text-3xl font-semibold tracking-tight">
						{mode === 'daily' ? `Daily #${day}` : 'Practice'}
					</h1>
					<span class="badge preset-outlined-surface-500 font-mono text-xs">
						patch {puzzle.patch}
					</span>
				</div>
				<p class="text-sm opacity-70">Nine picks in. Name the last red one.</p>
			</div>
			<div class="flex items-center gap-2 pb-1">
				{#if guesses.length > 0}
					<span class="font-mono text-xs opacity-60">
						{guesses.length}
						{guesses.length === 1 ? 'try' : 'tries'}
					</span>
				{/if}
				<div class="flex gap-1.5" title="{Math.min(misses, CLUES)} of {CLUES} clues unlocked">
					{#each { length: CLUES }, i (i)}
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
		</header>

		<div
			class="rounded-container border border-surface-300-700 bg-surface-100-900/40 p-4 shadow-sm"
		>
			<DraftBoard {puzzle} version={data.version} {byName} reveal={over} />
		</div>

		{#if !over}
			<GuessInput champions={data.champions} version={data.version} {blocked} onguess={guess} />
			{#if misses >= CLUES}
				<div class="flex justify-end">
					<button class="btn preset-outlined-surface-500 opacity-70 btn-sm" onclick={giveUp}>
						Give up and show the pick
					</button>
				</div>
			{/if}
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
