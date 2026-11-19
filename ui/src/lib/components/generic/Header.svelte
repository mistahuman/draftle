<script lang="ts">
	import { resolve } from '$app/paths';
	import GitHubIcon from '$lib/components/icons/GitHubIcon.svelte';
	import Logo from '$lib/components/icons/Logo.svelte';
	import Lightswitch from './Lightswitch.svelte';
	import ThemeSwitch from './ThemeSwitch.svelte';
	import { THEME_PICKER, DEFAULT_THEME } from '$lib/config';
</script>

<header
	class="sticky top-0 z-50 flex h-14 w-full items-center border-b-[1px] border-surface-500/20 bg-surface-50/90 backdrop-blur dark:bg-surface-950/90"
>
	<div class="container mx-auto flex max-w-3xl items-center justify-between px-4">
		<a href={resolve('/')} class="flex items-center gap-2" title="draftle">
			<Logo width="28" height="28" />
			<span class="text-xl font-bold tracking-tight">draftle</span>
		</a>
		<div class="flex items-center gap-1">
			<Lightswitch />
			{#if THEME_PICKER}
				<ThemeSwitch defaultTheme={DEFAULT_THEME} />
			{/if}
			<a
				class="btn-icon hover:preset-tonal"
				href="https://github.com/mistahuman"
				rel="external"
				title="GitHub"
			>
				<GitHubIcon size={20} />
			</a>
		</div>
	</div>
</header>
