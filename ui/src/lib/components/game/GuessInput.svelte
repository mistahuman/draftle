<script lang="ts">
	import ChampIcon from './ChampIcon.svelte';
	import type { Champion } from '$lib/game/types';

	let {
		champions,
		version,
		blocked,
		onguess
	}: {
		champions: Champion[];
		version: string;
		blocked: Set<string>; // already in the draft or already guessed
		onguess: (c: Champion) => void;
	} = $props();

	let query = $state('');

	const norm = (s: string) => s.toLowerCase().replace(/[^a-z]/g, '');

	let matches = $derived.by(() => {
		const q = norm(query);
		if (!q) return [];
		return champions
			.filter((c) => !blocked.has(c.name) && norm(c.name).includes(q))
			.sort((a, b) => Number(!norm(a.name).startsWith(q)) - Number(!norm(b.name).startsWith(q)))
			.slice(0, 6);
	});

	function choose(c: Champion) {
		onguess(c);
		query = '';
	}

	function onkeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && matches.length > 0) {
			e.preventDefault();
			choose(matches[0]);
		}
	}
</script>

<div class="relative">
	<input
		class="input font-display text-lg placeholder:italic placeholder:opacity-60"
		placeholder="Write a champion…"
		bind:value={query}
		{onkeydown}
		autocomplete="off"
		spellcheck="false"
	/>
	{#if matches.length > 0}
		<ul
			class="absolute z-10 mt-1 w-full overflow-hidden rounded-container border border-surface-300-700 bg-surface-50-950 p-1 shadow-lg"
		>
			{#each matches as c (c.id)}
				<li>
					<button
						type="button"
						class="flex w-full items-center gap-3 rounded-md p-1.5 text-left hover:preset-tonal"
						onclick={() => choose(c)}
					>
						<ChampIcon {version} id={c.id} name={c.name} size={28} />
						<span class="font-semibold">{c.name}</span>
						<span class="ml-auto font-mono text-[11px] opacity-50">{c.tags.join(' / ')}</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>
