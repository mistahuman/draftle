import type { Champion, Puzzle } from './types';

// Unlimited tries; each of the first five misses unlocks one clue.
export const CLUES = 5;

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

// One clue per miss, weakest first, five in all. Every number is scoped to the puzzle's
// patch except the team clue, which looks at the whole year. Nothing names or narrows
// the answer to a short list: the old "first letter" and "one of three" clues are gone.
export function unlockedHints(p: Puzzle, answer: Champion, guesses: Champion[]): string[] {
	const vs = p.stats[p.answer]?.[1] ?? 0;
	const r = p.metaRank;
	const role = p.role;
	const clues = [
		`Class: ${answer.tags.join(' / ')}.`,
		`It is a ${answer.range} champion.`,
		vs > 0
			? `In patch ${p.patch} it was picked as ${role} against ${p.opponent} ${vs} other time${vs === 1 ? '' : 's'}.`
			: `In patch ${p.patch} nobody else picked it as ${role} against ${p.opponent}.`,
		r === 0
			? `Nobody else picked it as ${role} in patch ${p.patch}: a pocket pick.`
			: r <= 5
				? `It is one of the 5 most picked ${role}s of patch ${p.patch}.`
				: r <= 15
					? `It is not a top 5 ${role} of patch ${p.patch}, but it is in the top 15.`
					: `It is outside the top 15 ${role}s of patch ${p.patch}.`,
		p.teamPicks > 0
			? `${p.redTeam} played it as ${role} in ${p.teamPicks} other game${p.teamPicks === 1 ? '' : 's'} this year.`
			: `${p.redTeam} never played it as ${role} in another game this year.`
	];
	const misses = guesses.filter((g) => g.name !== answer.name).length;
	return clues.slice(0, Math.min(misses, CLUES));
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
	const misses = guesses.filter((g) => g.name !== answer.name).length;
	const clues = Math.min(misses, CLUES);
	const head = won
		? `${guesses.length} ${guesses.length === 1 ? 'try' : 'tries'} · ${clues}/${CLUES} clues`
		: `gave up after ${guesses.length}`;
	const bonus = calledWinner === undefined ? '' : `\nwinner ${calledWinner ? '✅' : '❌'}`;
	return `draftle #${day} · ${head}\n${squares}${bonus}`;
}
