<script lang="ts">
	let { width = '32', height = '32' }: { width?: string; height?: string } = $props();
</script>

<!-- A minimap: three lanes, the blue base filled, the red base still empty (the last pick).
     Same mark as static/favicon.svg, which hard-codes the colours. -->
<svg {width} {height} viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
	<rect x="1" y="1" width="30" height="30" rx="7" fill="currentColor" opacity="0.1" />
	<g fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" opacity="0.75">
		<path d="M6 26 V8.5 a2.5 2.5 0 0 1 2.5-2.5 H26" />
		<path d="M6 26 H23.5 a2.5 2.5 0 0 0 2.5-2.5 V6" />
		<path d="M6 26 L26 6" />
	</g>
	<circle cx="6.5" cy="25.5" r="4.6" fill="var(--side-blue)" />
	<circle
		cx="25.5"
		cy="6.5"
		r="3.8"
		fill="none"
		stroke="var(--side-red)"
		stroke-width="2"
		stroke-dasharray="2.8 1.6"
	/>
</svg>
