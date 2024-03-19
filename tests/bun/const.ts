export const PRISME_URL = 'http://prisme.localhost'

export const PRISME_API_URL = PRISME_URL + '/api/v1'
export const PRISME_PAGEVIEWS_URL = PRISME_API_URL + '/events/pageviews'
export const PRISME_CUSTOM_EVENTS_URL = PRISME_API_URL + '/events/customs'

export const TIMESTAMP_REGEX = /\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}/
export const UUID_V4_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
export const COUNTRY_CODE_REGEX = /^[A-Z]{2}$/
