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

// Hints unlock one per miss, weakest first.
export function hints(p: Puzzle, answer: Champion): string[] {
	const rank =
		p.metaRank > 0
			? `It is the #${p.metaRank} most picked ${p.role} of patch ${p.patch} (out of ${p.poolSize}).`
			: `Nobody else picked it as ${p.role} in patch ${p.patch}: a pocket pick.`;
	const vs = p.stats[p.answer]?.[1] ?? 0;
	const decoys = Object.entries(p.stats)
		.filter(([c]) => c !== p.answer)
		.sort((a, b) => b[1][0] - a[1][0])
		.slice(0, 2)
		.map(([c]) => c);
	const three = [...decoys, p.answer].sort();
	return [
		`Class: ${answer.tags.join(' / ')}.`,
		`Picked as ${p.role} against ${p.opponent} ${vs} other time${vs === 1 ? '' : 's'} in pro play this season.`,
		rank,
		`It starts with "${p.answer[0]}".`,
		`It is one of: ${three.join(', ')}.`
	];
}

export function shareText(
	day: number,
	guesses: Champion[],
	answer: Champion,
	won: boolean
): string {
	const squares = guesses
		.map((g) => {
			if (g.name === answer.name) return '🟩';
			return { same: '🟨', partial: '🟧', none: '⬛' }[classMatch(g, answer)];
		})
		.join('');
	return `draftle #${day} ${won ? guesses.length : 'X'}/${MAX_GUESSES}\n${squares}`;
}
