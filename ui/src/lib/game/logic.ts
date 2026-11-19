import type { Champion, Puzzle } from './types';

export const MAX_GUESSES = 6;

// Day 1 is the launch day. Daily puzzles cycle through the exported set.
const EPOCH = Date.UTC(2026, 10, 19);

export function dayNumber(now = new Date()): number {
	const today = Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate());
	return Math.floor((today - EPOCH) / 86_400_000) + 1;
}

export function dailyIndex(day: number, count: number): number {
	return (((day - 1) % count) + count) % count;
}

export type ClassMatch = 'same' | 'partial' | 'none';

export function classMatch(guess: Champion, answer: Champion): ClassMatch {
	const shared = guess.tags.filter((t) => answer.tags.includes(t)).length;
	if (shared === 0) return 'none';
	return shared === guess.tags.length && shared === answer.tags.length ? 'same' : 'partial';
}

export function picksInRole(p: Puzzle, champ: string): number {
	return p.stats[champ]?.[0] ?? 0;
}

export function takenInDraft(p: Puzzle): Set<string> {
	return new Set([
		...p.blue.picks.map((x) => x.champ),
		...p.red.picks.map((x) => x.champ),
		...p.blue.bans,
		...p.red.bans
	]);
}

// One hint unlocks per miss, weakest first. Every number is scoped to the puzzle's patch,
// the same scope as the per-guess popularity. The class hint is skipped when a guess has
// already matched the answer's class exactly. The check runs miss by miss, so a hint
// already shown never disappears.
export function unlockedHints(p: Puzzle, answer: Champion, guesses: Champion[]): string[] {
	const vs = p.stats[p.answer]?.[1] ?? 0;
	const decoys = Object.entries(p.stats)
		.filter(([c]) => c !== p.answer)
		.sort((a, b) => b[1][0] - a[1][0])
		.slice(0, 2)
		.map(([c]) => c);
	const queue: { kind: 'class' | 'other'; text: string }[] = [
		{ kind: 'class', text: `Class: ${answer.tags.join(' / ')}.` },
		{
			kind: 'other',
			text:
				vs > 0
					? `In patch ${p.patch} it was picked as ${p.role} against ${p.opponent} ${vs} other time${vs === 1 ? '' : 's'}.`
					: `In patch ${p.patch} nobody else picked it as ${p.role} against ${p.opponent}.`
		},
		{
			kind: 'other',
			text:
				p.metaRank > 0
					? `It is the #${p.metaRank} most picked ${p.role} of patch ${p.patch} (out of ${p.poolSize}).`
					: `Nobody else picked it as ${p.role} in patch ${p.patch}: a pocket pick.`
		},
		{ kind: 'other', text: `It starts with "${p.answer[0]}".` },
		{ kind: 'other', text: `It is one of: ${[...decoys, p.answer].sort().join(', ')}.` }
	];

	const wrong = guesses.filter((g) => g.name !== answer.name);
	const shown: string[] = [];
	let next = 0;
	wrong.forEach((_, k) => {
		const classKnown = wrong.slice(0, k + 1).some((g) => classMatch(g, answer) === 'same');
		while (next < queue.length && queue[next].kind === 'class' && classKnown) next++;
		if (next < queue.length) shown.push(queue[next++].text);
	});
	return shown;
}

export function shareText(
	day: number,
	guesses: Champion[],
	answer: Champion,
	won: boolean,
	calledWinner?: boolean // the "who won?" bonus; undefined if not answered
): string {
	const squares = guesses
		.map((g) => {
			if (g.name === answer.name) return '🟩';
			return { same: '🟨', partial: '🟧', none: '⬛' }[classMatch(g, answer)];
		})
		.join('');
	const bonus = calledWinner === undefined ? '' : `\nwinner ${calledWinner ? '✅' : '❌'}`;
	return `draftle #${day} ${won ? guesses.length : 'X'}/${MAX_GUESSES}\n${squares}${bonus}`;
}
