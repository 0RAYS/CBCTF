export const DEFAULT_CHALLENGE_CATEGORIES = ['Web', 'Misc', 'Pwn', 'Crypto', 'Reverse'];

export function mergeChallengeCategories(categories = []) {
  const resolved = Array.isArray(categories) ? categories : [];
  return [...new Set([...DEFAULT_CHALLENGE_CATEGORIES, ...resolved])];
}

