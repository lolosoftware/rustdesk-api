export const connectByClient = (id) => {
  let a = document.createElement('a')
  a.href = `rustdesk://${id}`
  a.target = '_self'
  a.rel = 'noopener'
  a.click()
}
