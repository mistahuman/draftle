<script lang="ts">
	let {
		version,
		id,
		name,
		size = 48,
		dim = false
	}: { version: string; id?: string; name: string; size?: number; dim?: boolean } = $props();
</script>

{#if id}
	<img
		src="https://ddragon.leagueoflegends.com/cdn/{version}/img/champion/{id}.png"
		alt={name}
		title={name}
		width={size}
		height={size}
		loading="lazy"
		class="rounded-md {dim ? 'opacity-40 grayscale' : ''}"
	/>
{:else}
	<div
		class="flex items-center justify-center rounded-md bg-surface-200-800 text-xs"
		style="width:{size}px;height:{size}px"
		title={name}
	>
		{name.slice(0, 3)}
	</div>
{/if}
