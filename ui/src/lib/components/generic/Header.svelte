<script lang="ts">
	import { resolve } from '$app/paths';
	import GitHubIcon from '$lib/components/icons/GitHubIcon.svelte';
	import Logo from '$lib/components/icons/Logo.svelte';
	import HelpDialog from './HelpDialog.svelte';
	import Lightswitch from './Lightswitch.svelte';
</script>

<!-- Same frame as sveltekit-skeleton-starter (full-width container), without the nav. -->
<header
	class="sticky top-0 z-50 flex h-[70px] w-full items-center border-b-[1px] border-surface-500/20 bg-surface-50/90 backdrop-blur dark:bg-surface-950/90"
>
	<div
		class="container mx-auto flex max-w-screen-2xl items-center justify-between gap-4 px-4 xl:px-10"
	>
		<a href={resolve('/')} class="flex items-center gap-2" title="draftle">
			<Logo width="32" height="32" />
			<span class="font-display text-xl font-semibold tracking-tight">draftle</span>
		</a>
		<div class="flex items-stretch justify-end gap-2">
			<HelpDialog />
			<Lightswitch />
			<span class="border-r border-surface-200-800"></span>
			<nav class="flex items-center gap-1">
				<a
					class="btn-icon hover:preset-tonal"
					href="https://github.com/mistahuman"
					rel="external"
					title="GitHub"
				>
					<GitHubIcon size={20} />
				</a>
			</nav>
		</div>
	</div>
</header>
