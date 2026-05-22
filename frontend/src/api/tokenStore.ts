const tokenKey = 'lms.session.token'

export function getSessionToken() {
  return window.sessionStorage.getItem(tokenKey)
}

export function setSessionToken(token: string) {
  window.sessionStorage.setItem(tokenKey, token)
}

export function clearSessionToken() {
  window.sessionStorage.removeItem(tokenKey)
}
