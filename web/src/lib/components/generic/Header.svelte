<script lang="ts">
	import { resolve } from '$app/paths';
	import GitHubIcon from '$lib/components/icons/GitHubIcon.svelte';
	import Lightswitch from './Lightswitch.svelte';
</script>

<header
	class="sticky top-0 z-50 flex h-14 w-full items-center border-b-[1px] border-surface-500/20 bg-surface-50/90 backdrop-blur dark:bg-surface-950/90"
>
	<div class="container mx-auto flex max-w-3xl items-center justify-between px-4">
		<a href={resolve('/')} class="text-xl font-bold tracking-tight">draftle</a>
		<div class="flex items-center gap-1">
			<Lightswitch />
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
