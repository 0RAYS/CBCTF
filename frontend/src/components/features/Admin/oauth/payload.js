export function oauthForm(provider) {
  const form = {
    provider: '',
    uri: '',
    protocol: 'oauth2',
    scopes: '',
    auth_url: '',
    token_url: '',
    user_info_url: '',
    callback_url: '',
    client_id: '',
    client_secret: '',
    picture: '',
    picture_claim: '',
    name_claim: '',
    email_claim: '',
    description_claim: '',
    id_claim: '',
    groups_claim: '',
    admin_group: '',
    default_group: 0,
    on: false,
  };
  if (!provider) return form;
  for (const key of Object.keys(form)) {
    if (provider[key] != null) form[key] = provider[key];
  }
  form.protocol = provider.protocol || 'oauth2';
  form.scopes = Array.isArray(provider.scopes) ? provider.scopes.join(', ') : '';
  form.client_secret = '';
  return form;
}

export function parseScopes(value) {
  return value
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean);
}

export function buildOAuthPayload(form) {
  const payload = {};
  for (const key of [
    'protocol',
    'auth_url',
    'user_info_url',
    'callback_url',
    'provider',
    'uri',
    'id_claim',
    'name_claim',
    'email_claim',
    'picture_claim',
    'description_claim',
    'groups_claim',
    'admin_group',
    'default_group',
    'on',
    'picture',
  ])
    payload[key] = form[key];
  if (form.protocol !== 'cas') {
    payload.token_url = form.token_url;
    payload.client_id = form.client_id;
    if (form.client_secret) payload.client_secret = form.client_secret;
    payload.scopes = parseScopes(form.scopes);
  }
  return payload;
}
