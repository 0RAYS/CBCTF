export function smtpForm(smtp) {
  return smtp
    ? { address: smtp.address, host: smtp.host, port: smtp.port, pwd: '', on: smtp.on }
    : { address: '', host: '', port: 587, pwd: '', on: false };
}

export function buildSmtpPayload(form, original) {
  if (!original) return { ...form };
  const changes = {};
  for (const key of ['address', 'host', 'port', 'on']) {
    if (form[key] !== original[key]) changes[key] = form[key];
  }
  if (form.pwd) changes.pwd = form.pwd;
  return changes;
}
