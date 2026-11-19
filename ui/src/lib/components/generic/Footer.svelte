<script lang="ts">
	let { classList = 'p-4' }: { classList?: string } = $props();
</script>

<footer class="border-t-[1px] border-surface-500/20 {classList}">
	<div class="container mx-auto flex flex-col items-center gap-2 text-center text-xs opacity-60">
		<p>
			built by <a class="underline" href="https://www.github.com/mistahuman" target="_blank"
				>mistahuman</a
			>. Draft data from
			<a class="underline" href="https://lol.fandom.com" target="_blank">Leaguepedia</a>
			(CC BY-SA).
		</p>
		<p>
			draftle isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games
			or anyone officially involved in producing or managing Riot Games properties. Riot Games, and
			all associated properties are trademarks or registered trademarks of Riot Games, Inc.
		</p>
	</div>
</footer>
