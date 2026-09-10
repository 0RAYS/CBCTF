export function cheatReviewForm(cheat) {
  return {
    reason: cheat.reason || '',
    type: cheat.type || '',
    checked: cheat.checked || false,
    comment: cheat.comment || '',
  };
}

export function cheatUpdatePayload(form, cheat) {
  const original = cheatReviewForm(cheat);
  return Object.fromEntries(
    Object.keys(original)
      .filter((key) => form[key] !== original[key])
      .map((key) => [key, form[key]])
  );
}

export function cheatListParams(page, pageSize, type, reasonType) {
  return {
    limit: pageSize,
    offset: (page - 1) * pageSize,
    ...(type ? { type } : {}),
    ...(reasonType ? { reason_type: reasonType } : {}),
  };
}
