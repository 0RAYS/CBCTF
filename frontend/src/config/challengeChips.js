const fallbackChipClass = 'bg-neutral-400/20 text-neutral-400';

const categoryChipClasses = {
  web: 'bg-blue-400/20 text-blue-400',
  crypto: 'bg-purple-400/20 text-purple-400',
  pwn: 'bg-red-400/20 text-red-400',
  reverse: 'bg-green-400/20 text-green-400',
  misc: 'bg-yellow-400/20 text-yellow-400',
};

const typeChipClasses = {
  static: 'bg-geek-400/20 text-geek-400',
  question: 'bg-green-400/20 text-green-400',
  dynamic: 'bg-orange-400/20 text-orange-400',
  pods: 'bg-cyan-400/20 text-cyan-400',
};

export function getChallengeCategoryChipClass(category) {
  if (!category) return fallbackChipClass;
  const key = String(category).toLowerCase();
  return categoryChipClasses[key] || fallbackChipClass;
}

export function getChallengeTypeChipClass(type) {
  if (!type) return fallbackChipClass;
  const key = String(type).toLowerCase();
  return typeChipClasses[key] || fallbackChipClass;
}

