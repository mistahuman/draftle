<script lang="ts">
	let { width = '32', height = '32' }: { width?: string; height?: string } = $props();

	const rows = [0, 1, 2, 3, 4];
</script>

<!-- The whole draft: five blue picks, four red, and the last red slot empty.
     Same mark as static/favicon.svg, which hard-codes the colours. -->
<svg {width} {height} viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
	{#each rows as i (i)}
		<rect x="4" y={2 + i * 6} width="10" height="4.5" rx="1.5" fill="var(--side-blue)" />
	{/each}
	{#each rows.slice(0, 4) as i (i)}
		<rect x="18" y={2 + i * 6} width="10" height="4.5" rx="1.5" fill="var(--side-red)" />
	{/each}
	<rect
		x="18.75"
		y="26.75"
		width="8.5"
		height="3"
		rx="1"
		fill="none"
		stroke="var(--side-red)"
		stroke-width="1.5"
	/>
</svg>
