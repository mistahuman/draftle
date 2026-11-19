<script lang="ts">
	import { CircleHelp, X } from 'lucide-svelte';
	import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte';
	import { onMount } from 'svelte';

	// Opens by itself on the first visit only; afterwards from the "?" in the header.
	const SEEN = 'draftle:help-seen';
	let open = $state(false);

	onMount(() => {
		try {
			if (!localStorage.getItem(SEEN)) open = true;
		} catch {
			/* storage blocked: just don't auto-open */
		}
	});

	function onOpenChange(e: { open: boolean }) {
		open = e.open;
		if (!open) {
			try {
				localStorage.setItem(SEEN, '1');
			} catch {
				/* ignore */
			}
		}
	}

	const animation =
		'transition transition-discrete opacity-0 translate-y-4 starting:data-[state=open]:opacity-0 starting:data-[state=open]:translate-y-4 data-[state=open]:opacity-100 data-[state=open]:translate-y-0';
</script>

<Dialog {open} {onOpenChange}>
	<Dialog.Trigger class="btn-icon hover:preset-tonal" title="How to play">
		<CircleHelp size={20} />
	</Dialog.Trigger>
	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-950/40" />
		<Dialog.Positioner class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<Dialog.Content
				class="max-h-[85dvh] w-full max-w-lg space-y-5 overflow-y-auto rounded-container border-[3px] border-double border-surface-400-600 bg-surface-50-950 p-5 shadow-xl {animation}"
			>
				<header class="flex items-center justify-between">
					<Dialog.Title class="font-display text-2xl font-semibold tracking-tight">
						How to play
					</Dialog.Title>
					<Dialog.CloseTrigger class="btn-icon hover:preset-tonal" title="Close">
						<X size={18} />
					</Dialog.CloseTrigger>
				</header>

				<Dialog.Description class="space-y-4 text-sm leading-relaxed">
					<p>
						A real pro draft, one pick from the end. <strong>Name the last red pick.</strong>
					</p>
					<ul class="space-y-2">
						<li class="flex items-center gap-2">
							<span class="flex shrink-0 gap-1">
								<span class="h-3 w-3 rounded-[2px] bg-surface-300-700"></span>
								<span class="h-3 w-3 rounded-[2px] bg-(--side-red)"></span>
								<span class="h-3 w-3 rounded-[2px] bg-success-500"></span>
							</span>
							<span>Six tries. Every miss unlocks a clue.</span>
						</li>
						<li>
							<span class="badge preset-outlined-primary-500 font-mono text-[11px] font-bold"
								>same</span
							>
							<span class="badge preset-outlined-surface-500 font-mono text-[11px]">partial</span>
							class as the answer.
						</li>
						<li>
							<span class="font-mono">12 ↑</span> the answer is picked more than your guess this patch.
						</li>
					</ul>
					<p class="border-t border-dotted border-surface-400-600 pt-3 text-xs opacity-70">
						New puzzle daily at midnight UTC. Drafts from
						<a class="anchor" href="https://lol.fandom.com" rel="external">Leaguepedia</a>. Fan
						project, not endorsed by Riot Games.
					</p>
				</Dialog.Description>

				<footer class="flex justify-end">
					<Dialog.CloseTrigger class="btn preset-filled-primary-500">Play</Dialog.CloseTrigger>
				</footer>
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
