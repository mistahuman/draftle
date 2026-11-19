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

{#snippet heading(text: string)}
	<h3
		class="border-b border-dotted border-surface-400-600 pb-1 font-mono text-xs tracking-wider uppercase opacity-60"
	>
		{text}
	</h3>
{/snippet}

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

				<Dialog.Description class="space-y-5 text-sm leading-relaxed">
					<section class="space-y-2">
						<p>
							Every day, a real pro game stopped one pick before the end. Nine champions are locked
							in and so are all ten bans. <strong>Name the last red pick</strong>: you know its role
							and who it plays against in lane.
						</p>
						<ul class="space-y-1.5">
							<li class="flex items-center gap-2">
								<span class="flex shrink-0 gap-1">
									<span class="h-3 w-3 rounded-[2px] bg-surface-300-700"></span>
									<span class="h-3 w-3 rounded-[2px] bg-(--side-red)"></span>
									<span class="h-3 w-3 rounded-[2px] bg-success-500"></span>
								</span>
								<span>Six tries. Each square is one: free, missed, found.</span>
							</li>
						</ul>
					</section>

					<section class="space-y-2">
						{@render heading('Each guess tells you')}
						<ul class="space-y-1.5">
							<li>
								<span class="badge preset-outlined-primary-500 font-mono text-[11px] font-bold"
									>same</span
								>
								<span class="badge preset-outlined-surface-500 font-mono text-[11px]">partial</span>
								<span class="badge font-mono text-[11px] opacity-40">different</span>
								— its class compared with the answer's.
							</li>
							<li>
								<span class="font-mono">12 ↑</span> — your guess was picked 12 times in that role
								this patch; the answer was picked <strong>more</strong> (↓ less, = the same).
							</li>
						</ul>
					</section>

					<section class="space-y-2">
						{@render heading('Clues')}
						<p>
							Every miss unlocks a clue, weakest first: class, picks against the lane opponent, meta
							rank, first letter, then one of three. Every number is from that game's patch.
						</p>
					</section>

					<section class="space-y-2">
						{@render heading('Where it comes from')}
						<p>
							Drafts from pro games of the last six patches, via
							<a class="anchor" href="https://lol.fandom.com" rel="external">Leaguepedia</a>. One
							puzzle a day, new at midnight UTC, the same for everyone. After it you can keep going
							with practice puzzles.
						</p>
						<p class="opacity-70">
							A fan project, not endorsed by Riot Games. Code on
							<a class="anchor" href="https://github.com/mistahuman/draftle" rel="external"
								>GitHub</a
							>.
						</p>
					</section>
				</Dialog.Description>

				<footer class="flex justify-end">
					<Dialog.CloseTrigger class="btn preset-filled-primary-500">Play</Dialog.CloseTrigger>
				</footer>
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
